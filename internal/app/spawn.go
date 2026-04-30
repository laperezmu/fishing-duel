package app

import (
	"fmt"
	"pesca/internal/content/fishprofiles"
	"pesca/internal/encounter"
	"pesca/internal/player/loadout"
)

type SpawnUI interface {
	ShowFishSpawn(title string, spawn fishprofiles.Spawn) error
}

func ResolveFishSpawn(title string, opening encounter.Opening, playerLoadout loadout.State, profiles []fishprofiles.Profile, ui SpawnUI) (fishprofiles.Spawn, error) {
	if ui == nil {
		return fishprofiles.Spawn{}, fmt.Errorf("spawn ui is required")
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

	spawn, err := fishprofiles.ResolveSpawn(profiles, spawnContext)
	if err != nil {
		return fishprofiles.Spawn{}, fmt.Errorf("resolve fish spawn: %w", err)
	}

	if err := ui.ShowFishSpawn(title, spawn); err != nil {
		return fishprofiles.Spawn{}, fmt.Errorf("show fish spawn: %w", err)
	}

	return spawn, nil
}
