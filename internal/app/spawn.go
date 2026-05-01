package app

import (
	"fmt"
	"pesca/internal/content/fishprofiles"
	"pesca/internal/encounter"
	"pesca/internal/player/loadout"
	"pesca/internal/presentation"
)

type SpawnUI interface {
	ShowFishSpawn(title string, spawn presentation.SpawnView) error
}

type SpawnPresenter interface {
	Spawn(spawn fishprofiles.Spawn) presentation.SpawnView
}

func ResolveFishSpawn(title string, opening encounter.Opening, playerLoadout loadout.State, profiles []fishprofiles.Profile, ui SpawnUI, presenter SpawnPresenter) (fishprofiles.Spawn, error) {
	return resolveFishSpawn(title, opening, playerLoadout, profiles, ui, presenter, nil)
}

func ResolveFishSpawnWithRandomizer(title string, opening encounter.Opening, playerLoadout loadout.State, profiles []fishprofiles.Profile, ui SpawnUI, presenter SpawnPresenter, randomizer fishprofiles.SpawnRandomizer) (fishprofiles.Spawn, error) {
	return resolveFishSpawn(title, opening, playerLoadout, profiles, ui, presenter, randomizer)
}

func resolveFishSpawn(title string, opening encounter.Opening, playerLoadout loadout.State, profiles []fishprofiles.Profile, ui SpawnUI, presenter SpawnPresenter, randomizer fishprofiles.SpawnRandomizer) (fishprofiles.Spawn, error) {
	if ui == nil {
		return fishprofiles.Spawn{}, fmt.Errorf("spawn ui is required")
	}
	if presenter == nil {
		return fishprofiles.Spawn{}, fmt.Errorf("spawn presenter is required")
	}
	if err := playerLoadout.Validate(); err != nil {
		return fishprofiles.Spawn{}, fmt.Errorf("player loadout: %w", err)
	}
	if len(profiles) == 0 {
		return fishprofiles.Spawn{}, fmt.Errorf("at least one fish profile is required")
	}

	spawnContext, err := fishprofiles.NewSpawnContext(opening, playerLoadout.HabitatTags())
	if err != nil {
		return fishprofiles.Spawn{}, fmt.Errorf("spawn context: %w", err)
	}

	spawn, err := fishprofiles.ResolveSpawnWithRandomizer(profiles, spawnContext, randomizer)
	if err != nil {
		return fishprofiles.Spawn{}, fmt.Errorf("resolve fish spawn: %w", err)
	}

	if err := ui.ShowFishSpawn(title, presenter.Spawn(spawn)); err != nil {
		return fishprofiles.Spawn{}, fmt.Errorf("show fish spawn: %w", err)
	}

	return spawn, nil
}
