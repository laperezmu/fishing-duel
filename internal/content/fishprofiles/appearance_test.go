package fishprofiles

import (
	"pesca/internal/content/habitats"
	"pesca/internal/content/waterpools"
	"pesca/internal/encounter"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveSpawnFallsBackToGenericProfileWithoutHabitats(t *testing.T) {
	opening := encounter.Opening{
		WaterContext:    encounter.WaterContext{PoolTag: waterpools.Shoreline},
		InitialDistance: 0,
		InitialDepth:    1,
	}
	context, err := NewSpawnContext(opening, nil)
	require.NoError(t, err)

	spawn, err := ResolveSpawn(DefaultProfiles(), context)

	require.NoError(t, err)
	assert.Equal(t, "classic", spawn.Profile.ID)
	assert.Equal(t, 1, spawn.CandidateCount)
	assert.Equal(t, waterpools.Shoreline, spawn.Context.WaterPoolTag)
}

func TestResolveSpawnPrefersSpecificHabitatMatch(t *testing.T) {
	opening := encounter.Opening{
		WaterContext:    encounter.WaterContext{PoolTag: waterpools.Shoreline},
		InitialDistance: 1,
		InitialDepth:    1,
	}
	context, err := NewSpawnContext(opening, []habitats.Tag{habitats.Surface, habitats.OpenWater})
	require.NoError(t, err)

	spawn, err := ResolveSpawn(DefaultProfiles(), context)

	require.NoError(t, err)
	assert.Equal(t, "surface-control", spawn.Profile.ID)
	assert.GreaterOrEqual(t, spawn.CandidateCount, 2)
}

func TestResolveSpawnWithRandomizerVariesWithinTopScoreTies(t *testing.T) {
	context := SpawnContext{
		WaterPoolTag:    waterpools.Shoreline,
		InitialDistance: 1,
		InitialDepth:    1,
	}
	profiles := []Profile{
		{
			ID:          "tie-first",
			ArchetypeID: ArchetypeBaselineCycle,
			Name:        "Tie First",
			Appearance: Appearance{
				WaterPoolTags:      []waterpools.ID{waterpools.Shoreline},
				MinInitialDistance: 0,
				MaxInitialDistance: 2,
				MinInitialDepth:    0,
				MaxInitialDepth:    2,
			},
		},
		{
			ID:          "tie-second",
			ArchetypeID: ArchetypeBaselineCycle,
			Name:        "Tie Second",
			Appearance: Appearance{
				WaterPoolTags:      []waterpools.ID{waterpools.Shoreline},
				MinInitialDistance: 0,
				MaxInitialDistance: 2,
				MinInitialDepth:    0,
				MaxInitialDepth:    2,
			},
		},
	}

	spawn, err := ResolveSpawnWithRandomizer(profiles, context, fixedSpawnRandomizer{value: 1})

	require.NoError(t, err)
	assert.Equal(t, "tie-second", spawn.Profile.ID)
	assert.Equal(t, 2, spawn.CandidateCount)
}

func TestResolveSpawnWithoutRandomizerRemainsStableForTies(t *testing.T) {
	context := SpawnContext{
		WaterPoolTag:    waterpools.Shoreline,
		InitialDistance: 1,
		InitialDepth:    1,
	}
	profiles := []Profile{
		{
			ID:          "tie-first",
			ArchetypeID: ArchetypeBaselineCycle,
			Name:        "Tie First",
			Appearance: Appearance{
				WaterPoolTags:      []waterpools.ID{waterpools.Shoreline},
				MinInitialDistance: 0,
				MaxInitialDistance: 2,
				MinInitialDepth:    0,
				MaxInitialDepth:    2,
			},
		},
		{
			ID:          "tie-second",
			ArchetypeID: ArchetypeBaselineCycle,
			Name:        "Tie Second",
			Appearance: Appearance{
				WaterPoolTags:      []waterpools.ID{waterpools.Shoreline},
				MinInitialDistance: 0,
				MaxInitialDistance: 2,
				MinInitialDepth:    0,
				MaxInitialDepth:    2,
			},
		},
	}

	spawn, err := ResolveSpawn(profiles, context)

	require.NoError(t, err)
	assert.Equal(t, "tie-first", spawn.Profile.ID)
}

func TestResolveSpawnReturnsErrorWhenNoProfileMatches(t *testing.T) {
	context := SpawnContext{
		WaterPoolTag:    waterpools.Offshore,
		InitialDistance: 5,
		InitialDepth:    4,
	}

	_, err := ResolveSpawn(DefaultProfiles(), context)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no fish profile matches water offshore")
}

func TestProfileBuildPreset(t *testing.T) {
	profile := DefaultProfiles()[0]
	preset := profile.BuildPreset()

	assert.Equal(t, profile.ID, preset.ID)
	assert.Equal(t, profile.Name, preset.Name)
	assert.Equal(t, profile.Shuffle, preset.Shuffle)
	assert.Equal(t, profile.CardsToRemove, preset.CardsToRemove)
	require.Len(t, preset.FishCards, len(profile.Cards))
}

type fixedSpawnRandomizer struct {
	value int
}

func (randomizer fixedSpawnRandomizer) Intn(n int) int {
	if n <= 0 {
		return 0
	}
	return randomizer.value % n
}
