package app

import (
	"fmt"
	"math/rand"
	"pesca/internal/cards"
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

type EncounterBootstrapUI interface {
	SetupUI
	OpeningUI
	SpawnUI
}

type Randomizer interface {
	Float64() float64
	Shuffle(n int, swap func(i, j int))
}

func BootstrapEncounter(title string, rng Randomizer, ui EncounterBootstrapUI) (*game.Engine, error) {
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
	opening, spawn, err := resolveEncounterBootstrap(title, playerLoadout, ui, presenter)
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

func resolveEncounterBootstrap(title string, playerLoadout loadout.State, ui EncounterBootstrapUI, presenter presentation.Presenter) (encounter.Opening, fishprofiles.Spawn, error) {
	opening, err := ResolveEncounterOpening(title, encounter.DefaultConfig(), playerLoadout, watercontexts.DefaultPresets(), ui, presenter)
	if err != nil {
		return encounter.Opening{}, fishprofiles.Spawn{}, fmt.Errorf("resolve encounter opening: %w", err)
	}
	spawn, err := ResolveFishSpawn(title, opening, playerLoadout, fishprofiles.DefaultProfiles(), ui, presenter)
	if err != nil {
		return encounter.Opening{}, fishprofiles.Spawn{}, fmt.Errorf("resolve fish spawn: %w", err)
	}

	return opening, spawn, nil
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
		progression.TrackPolicy{SplashEscapeDecider: progression.SplashEscapeDeciderFunc(func(chance float64) bool {
			return rng.Float64() < chance
		})},
		endings.EncounterEndCondition{},
		match.State{Encounter: encounterState, Player: match.PlayerState{Loadout: playerLoadout}},
	)
	if err != nil {
		return nil, fmt.Errorf("initialize game engine: %w", err)
	}

	return engine, nil
}

type seededRandom struct{ *rand.Rand }

func NewSeededRandom(source int64) Randomizer {
	//nolint:gosec // Gameplay randomness is non-cryptographic and intentionally reproducible from a seed.
	return seededRandom{Rand: rand.New(rand.NewSource(source))}
}

func DefaultRandomizer() Randomizer {
	return NewSeededRandom(time.Now().UnixNano())
}
