package app

import (
	"fmt"
	"pesca/internal/content/fishprofiles"
	"pesca/internal/content/playerprofiles"
	"pesca/internal/game"
	"pesca/internal/player/loadout"
	"pesca/internal/presentation"
)

type EncounterBootstrapRuntime struct {
	Engine *game.Engine
	Spawn  fishprofiles.Spawn
}

func BootstrapEncounterForRun(title string, rng Randomizer, ui EncounterBootstrapUI, presenter presentation.Presenter, playerDeckPreset playerprofiles.DeckPreset, playerLoadout loadout.State, config EncounterBootstrapConfig) (EncounterBootstrapRuntime, error) {
	if rng == nil {
		return EncounterBootstrapRuntime{}, fmt.Errorf("randomizer is required")
	}
	if ui == nil {
		return EncounterBootstrapRuntime{}, fmt.Errorf("encounter bootstrap ui is required")
	}

	opening, spawn, err := resolveEncounterBootstrap(title, playerLoadout, ui, presenter, rng, config)
	if err != nil {
		return EncounterBootstrapRuntime{}, err
	}
	engine, err := buildEncounterEngine(rng, playerDeckPreset, playerLoadout, opening, spawn)
	if err != nil {
		return EncounterBootstrapRuntime{}, err
	}

	return EncounterBootstrapRuntime{Engine: engine, Spawn: spawn}, nil
}
