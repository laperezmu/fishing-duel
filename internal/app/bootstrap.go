package app

import (
	"fmt"
	"math/rand"
	"pesca/internal/cards"
	"pesca/internal/content/anglerprofiles"
	"pesca/internal/content/attachmentpresets"
	"pesca/internal/content/fishprofiles"
	"pesca/internal/content/playerprofiles"
	"pesca/internal/content/rodpresets"
	"pesca/internal/content/watercontexts"
	"pesca/internal/encounter"
	"pesca/internal/endings"
	"pesca/internal/game"
	"pesca/internal/match"
	"pesca/internal/player/loadout"
	"pesca/internal/player/playermoves"
	"pesca/internal/player/rod"
	"pesca/internal/presentation"
	"pesca/internal/progression"
	"pesca/internal/rules"
	"time"
)

type SetupUI interface {
	ChoosePlayerDeckPreset(title string, presets []playerprofiles.DeckPreset) (playerprofiles.DeckPreset, error)
	ChooseRodPreset(title string, presets []rodpresets.Preset) (rodpresets.Preset, error)
	ChooseAttachmentPreset(title string, baseRod rod.State, presets []attachmentpresets.Preset) (attachmentpresets.Preset, error)
}

type RunSetupUI interface {
	ChooseAnglerProfile(title string, profiles []anglerprofiles.Profile) (anglerprofiles.Profile, error)
}

type EncounterBootstrapUI interface {
	SetupUI
	OpeningUI
	SpawnUI
}

type Randomizer interface {
	Float64() float64
	Intn(n int) int
	Shuffle(n int, swap func(i, j int))
}

type EncounterBootstrapConfig struct {
	FishCatalog fishprofiles.Catalog
	FishPoolID  fishprofiles.PoolID
}

func BootstrapEncounter(title string, rng Randomizer, ui EncounterBootstrapUI) (*game.Engine, error) {
	return BootstrapEncounterWithConfig(title, rng, ui, EncounterBootstrapConfig{})
}

func BootstrapEncounterWithConfig(title string, rng Randomizer, ui EncounterBootstrapUI, config EncounterBootstrapConfig) (*game.Engine, error) {
	if rng == nil {
		return nil, fmt.Errorf("randomizer is required")
	}
	if ui == nil {
		return nil, fmt.Errorf("encounter bootstrap ui is required")
	}

	playerDeckPreset, playerLoadout, err := resolvePlayerSetup(title, ui)
	if err != nil {
		return nil, err
	}
	presenter := presentation.NewPresenter(presentation.DefaultCatalog())
	opening, spawn, err := resolveEncounterBootstrap(title, playerLoadout, ui, presenter, rng, config)
	if err != nil {
		return nil, err
	}

	return buildEncounterEngine(rng, playerDeckPreset, playerLoadout, opening, spawn)
}

func resolvePlayerSetup(title string, ui SetupUI) (playerprofiles.DeckPreset, loadout.State, error) {
	playerDeckPreset, err := ui.ChoosePlayerDeckPreset(title, playerprofiles.DefaultPresets())
	if err != nil {
		return playerprofiles.DeckPreset{}, loadout.State{}, fmt.Errorf("choose player deck preset: %w", err)
	}
	rodPreset, err := ui.ChooseRodPreset(title, rodpresets.DefaultPresets())
	if err != nil {
		return playerprofiles.DeckPreset{}, loadout.State{}, fmt.Errorf("choose player rod preset: %w", err)
	}
	playerRod, err := rodPreset.BuildRod()
	if err != nil {
		return playerprofiles.DeckPreset{}, loadout.State{}, fmt.Errorf("build player rod: %w", err)
	}
	attachmentPreset, err := ui.ChooseAttachmentPreset(title, playerRod, attachmentpresets.DefaultPresets())
	if err != nil {
		return playerprofiles.DeckPreset{}, loadout.State{}, fmt.Errorf("choose attachment preset: %w", err)
	}
	playerLoadout, err := rodPreset.BuildLoadoutWithAttachments(attachmentPreset.BuildAttachments())
	if err != nil {
		return playerprofiles.DeckPreset{}, loadout.State{}, fmt.Errorf("build player loadout: %w", err)
	}

	return playerDeckPreset, playerLoadout, nil
}

func resolveRunSetup(title string, ui RunSetupUI) (anglerprofiles.ResolvedStart, error) {
	profile, err := ui.ChooseAnglerProfile(title, anglerprofiles.DefaultUnlockedProfiles())
	if err != nil {
		return anglerprofiles.ResolvedStart{}, fmt.Errorf("choose angler profile: %w", err)
	}
	resolved, err := anglerprofiles.ResolveStart(profile)
	if err != nil {
		return anglerprofiles.ResolvedStart{}, fmt.Errorf("resolve angler profile: %w", err)
	}

	return resolved, nil
}

func resolveEncounterBootstrap(title string, playerLoadout loadout.State, ui EncounterBootstrapUI, presenter presentation.Presenter, rng Randomizer, config EncounterBootstrapConfig) (encounter.Opening, fishprofiles.Spawn, error) {
	opening, err := ResolveEncounterOpening(title, encounter.DefaultConfig(), playerLoadout, watercontexts.DefaultPresets(), ui, presenter)
	if err != nil {
		return encounter.Opening{}, fishprofiles.Spawn{}, fmt.Errorf("resolve encounter opening: %w", err)
	}
	profiles, err := resolveEncounterFishProfiles(config)
	if err != nil {
		return encounter.Opening{}, fishprofiles.Spawn{}, err
	}
	spawn, err := ResolveFishSpawnWithRandomizer(title, opening, playerLoadout, profiles, ui, presenter, rng)
	if err != nil {
		return encounter.Opening{}, fishprofiles.Spawn{}, fmt.Errorf("resolve fish spawn: %w", err)
	}

	return opening, spawn, nil
}

func resolveEncounterFishProfiles(config EncounterBootstrapConfig) ([]fishprofiles.Profile, error) {
	catalog := config.FishCatalog
	if len(catalog.Profiles()) == 0 {
		catalog = fishprofiles.DefaultCatalog()
	}
	poolID := resolveEncounterFishPoolID(config)
	profiles, err := catalog.ResolvePool(poolID)
	if err != nil {
		return nil, fmt.Errorf("resolve encounter fish pool: %w", err)
	}

	return profiles, nil
}

func resolveEncounterFishPoolID(config EncounterBootstrapConfig) fishprofiles.PoolID {
	if config.FishPoolID == "" {
		return fishprofiles.DefaultEncounterFishPoolID
	}

	return config.FishPoolID
}

func buildEncounterEngine(rng Randomizer, playerDeckPreset playerprofiles.DeckPreset, playerLoadout loadout.State, opening encounter.Opening, spawn fishprofiles.Spawn) (*game.Engine, error) {
	shuffler := func(fishCards []cards.FishCard) {
		rng.Shuffle(len(fishCards), func(i, j int) {
			fishCards[i], fishCards[j] = fishCards[j], fishCards[i]
		})
	}
	playerCardShuffler := func(playerCards []cards.PlayerCard) {
		rng.Shuffle(len(playerCards), func(i, j int) {
			playerCards[i], playerCards[j] = playerCards[j], playerCards[i]
		})
	}

	encounterState, err := encounter.NewState(opening.Config)
	if err != nil {
		return nil, fmt.Errorf("initialize encounter: %w", err)
	}
	playerMoveConfig := playerDeckPreset.BuildConfig(playerCardShuffler)
	playerMoveController, err := playermoves.NewUsageController(playerMoveConfig)
	if err != nil {
		return nil, fmt.Errorf("configure player moves: %w", err)
	}

	engine, err := game.NewEngine(
		spawn.Profile.BuildPreset().BuildDeck(shuffler),
		playerMoveController,
		rules.NewClassicEvaluator(rules.NewFishCombatProfile()),
		progression.TrackPolicy{},
		endings.EncounterEndCondition{},
		match.State{Encounter: applySpawnSplashProfile(encounterState, spawn), Player: match.PlayerState{Loadout: playerLoadout}},
	)
	if err != nil {
		return nil, fmt.Errorf("initialize game engine: %w", err)
	}

	return engine, nil
}

func applySpawnSplashProfile(state encounter.State, spawn fishprofiles.Spawn) encounter.State {
	state.Config.SplashProfile = spawn.Profile.Splash.BuildEncounterProfile()
	return state
}

type seededRandom struct{ *rand.Rand }

func NewSeededRandom(source int64) Randomizer {
	//nolint:gosec // Gameplay randomness is non-cryptographic and intentionally reproducible from a seed.
	return seededRandom{Rand: rand.New(rand.NewSource(source))}
}

func DefaultRandomizer() Randomizer {
	return NewSeededRandom(time.Now().UnixNano())
}
