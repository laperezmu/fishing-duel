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
