package fishprofiles

import (
	"pesca/internal/cards"
	"pesca/internal/content/habitats"
	"pesca/internal/content/waterpools"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultCatalogResolvesClosedPool(t *testing.T) {
	catalog := DefaultCatalog()

	profiles, err := catalog.ResolvePool(PoolID("shoreline-basics"))

	require.NoError(t, err)
	require.Len(t, profiles, 4)
	assert.Equal(t, ProfileID("classic"), profiles[0].ID)
	assert.Equal(t, ProfileID("deck-exhaustion"), profiles[3].ID)
}

func TestDefaultCatalogResolvesWeightedPoolEntries(t *testing.T) {
	catalog := DefaultCatalog()

	profiles, err := catalog.ResolvePool(PoolID("shoreline-weighted"))

	require.NoError(t, err)
	require.Len(t, profiles, 6)
	assert.Equal(t, ProfileID("classic"), profiles[0].ID)
	assert.Equal(t, ProfileID("classic"), profiles[1].ID)
	assert.Equal(t, ProfileID("classic"), profiles[2].ID)
	assert.Equal(t, ProfileID("surface-control"), profiles[4].ID)
	assert.Equal(t, ProfileID("deck-exhaustion"), profiles[5].ID)
}

func TestLoadCatalogRejectsUnknownPoolProfileReference(t *testing.T) {
	profileData := []byte(`{
		"profiles": [
			{
				"id": "classic",
				"archetype_id": "baseline_cycle",
				"name": "Clasico",
				"description": "desc",
				"details": [],
				"appearance": {
					"water_pool_tags": ["shoreline"],
					"min_initial_distance": 0,
					"max_initial_distance": 1,
					"min_initial_depth": 0,
					"max_initial_depth": 1,
					"required_habitat_tags": []
				},
				"cards": [{"move": "blue"}],
				"cards_to_remove": 0,
				"shuffle": false
			}
		]
	}`)
	poolData := []byte(`{
		"pools": [
			{
				"id": "broken-pool",
				"name": "Broken Pool",
				"description": "desc",
				"profile_ids": ["missing-profile"]
			}
		]
	}`)

	_, err := LoadCatalog(profileData, poolData)

	require.Error(t, err)
	assert.EqualError(t, err, "fish pool broken-pool references unknown profile id missing-profile")
}

func TestLoadCatalogRejectsInvalidPoolWeight(t *testing.T) {
	profileData := []byte(`{
		"profiles": [
			{
				"id": "classic",
				"archetype_id": "baseline_cycle",
				"name": "Clasico",
				"description": "desc",
				"details": [],
				"appearance": {
					"water_pool_tags": ["shoreline"],
					"min_initial_distance": 0,
					"max_initial_distance": 1,
					"min_initial_depth": 0,
					"max_initial_depth": 1,
					"required_habitat_tags": []
				},
				"cards": [{"move": "blue"}],
				"cards_to_remove": 0,
				"shuffle": false
			}
		]
	}`)
	poolData := []byte(`{
		"pools": [
			{
				"id": "broken-pool",
				"name": "Broken Pool",
				"description": "desc",
				"entries": [{"profile_id": "classic", "weight": 0}]
			}
		]
	}`)

	_, err := LoadCatalog(profileData, poolData)

	require.Error(t, err)
	assert.EqualError(t, err, "pool broken-pool has invalid weight 0 for profile id classic")
}

func TestResolveSpawnCanUseClosedPoolSubset(t *testing.T) {
	catalog := DefaultCatalog()
	profiles, err := catalog.ResolvePool(PoolID("shoreline-basics"))
	require.NoError(t, err)

	spawn, err := ResolveSpawn(profiles, SpawnContext{
		WaterPoolTag:    waterpools.Shoreline,
		InitialDistance: 1,
		InitialDepth:    1,
		HabitatTags:     []habitats.Tag{habitats.Surface},
	})

	require.NoError(t, err)
	assert.Equal(t, ProfileID("surface-control"), spawn.Profile.ID)
	assert.Equal(t, 3, spawn.CandidateCount)
}

func TestDefaultCatalogUsesNormalizedEffects(t *testing.T) {
	for _, profile := range DefaultProfiles() {
		for _, pattern := range profile.Cards {
			for _, effect := range pattern.Effects {
				assert.NotEqual(t, cards.EffectTypeUnknown, effect.Type)
				assert.Positive(t, effect.Priority)
			}
		}
	}
}

func TestDefaultCatalogCoversLegacyMappings(t *testing.T) {
	hasCaptureWindow := false
	hasSurfaceWindow := false
	hasExhaustionWindow := false
	hasHorizontal := false
	hasVertical := false

	for _, profile := range DefaultProfiles() {
		for _, pattern := range profile.Cards {
			for _, effect := range pattern.Effects {
				switch effect.Type {
				case cards.EffectTypeLegacyCaptureWindow:
					hasCaptureWindow = true
				case cards.EffectTypeLegacySurfaceWindow:
					hasSurfaceWindow = true
				case cards.EffectTypeLegacyExhaustionWindow:
					hasExhaustionWindow = true
				case cards.EffectTypeAdvanceHorizontal:
					hasHorizontal = true
				case cards.EffectTypeAdvanceVertical:
					hasVertical = true
				}
			}
		}
	}

	assert.True(t, hasCaptureWindow)
	assert.True(t, hasSurfaceWindow)
	assert.True(t, hasExhaustionWindow)
	assert.True(t, hasHorizontal)
	assert.True(t, hasVertical)
}

func TestDefaultCatalogJSONUsesExplicitEffectMetadata(t *testing.T) {
	profiles := DefaultProfiles()
	for _, profile := range profiles {
		for _, card := range profile.Cards {
			for _, effect := range card.Effects {
				assert.NotZero(t, effect.Priority)
				assert.NotEqual(t, cards.EffectTypeUnknown, effect.Type)
			}
		}
	}
}
