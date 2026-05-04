package app

import (
	"fmt"
	"pesca/internal/content/fishprofiles"
	"pesca/internal/content/playerprofiles"
	"pesca/internal/content/watercontexts"
	"pesca/internal/encounter"
	"pesca/internal/game"
	"pesca/internal/player/loadout"
	"pesca/internal/presentation"
)

type EncounterBootstrapRuntime struct {
	Engine *game.Engine
	Spawn  fishprofiles.Spawn
}

type RunEncounterBootstrapConfig struct {
	Encounter EncounterBootstrapConfig
	Water     watercontexts.Preset
}

func BootstrapEncounterForRun(title string, rng Randomizer, openingUI OpeningRuntimeUI, spawnUI SpawnUI, presenter presentation.Presenter, playerDeckPreset playerprofiles.DeckPreset, playerLoadout loadout.State, config RunEncounterBootstrapConfig) (EncounterBootstrapRuntime, error) {
	if rng == nil {
		return EncounterBootstrapRuntime{}, fmt.Errorf("randomizer is required")
	}
	if openingUI == nil {
		return EncounterBootstrapRuntime{}, fmt.Errorf("encounter bootstrap ui is required")
	}
	if spawnUI == nil {
		return EncounterBootstrapRuntime{}, fmt.Errorf("spawn ui is required")
	}

	opening, err := ResolveEncounterOpeningWithPreset(title, encounter.DefaultConfig(), playerLoadout, config.Water, openingUI, presenter)
	if err != nil {
		return EncounterBootstrapRuntime{}, fmt.Errorf("resolve encounter opening: %w", err)
	}
	profiles, err := resolveEncounterFishProfiles(config.Encounter)
	if err != nil {
		return EncounterBootstrapRuntime{}, err
	}
	spawn, err := ResolveFishSpawnWithRandomizer(title, opening, playerLoadout, profiles, spawnUI, presenter, rng)
	if err != nil {
		return EncounterBootstrapRuntime{}, err
	}
	engine, err := buildEncounterEngine(rng, playerDeckPreset, playerLoadout, opening, spawn)
	if err != nil {
		return EncounterBootstrapRuntime{}, err
	}

	return EncounterBootstrapRuntime{Engine: engine, Spawn: spawn}, nil
}
