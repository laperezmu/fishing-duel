package fishprofiles

import (
	"pesca/internal/cards"
	"pesca/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProfileBuildCards(t *testing.T) {
	profile := Profile{
		Cards: []CardPattern{{
			Name:              "Tiron de apertura",
			Summary:           "Permite capturar desde un paso mas lejos este round.",
			Move:              domain.Red,
			DiscardVisibility: cards.DiscardVisibilityMoveOnly,
			Effects: []cards.CardEffect{{
				Trigger:              cards.TriggerOnDraw,
				CaptureDistanceBonus: 1,
			}},
		}},
	}

	builtCards := profile.BuildCards()

	require.Len(t, builtCards, 1)
	assert.Equal(t, "Tiron de apertura", builtCards[0].Name)
	assert.Equal(t, domain.Red, builtCards[0].Move)
	assert.Equal(t, cards.DiscardVisibilityMoveOnly, builtCards[0].DiscardVisibility)
	assert.Equal(t, 1, builtCards[0].Effects[0].CaptureDistanceBonus)
	assert.Equal(t, cards.EffectTypeLegacyCaptureWindow, builtCards[0].Effects[0].Type)
	assert.Equal(t, 60, builtCards[0].Effects[0].Priority)
}

func TestDefaultProfiles(t *testing.T) {
	profiles := DefaultProfiles()

	require.Len(t, profiles, 7)
	assert.Equal(t, ArchetypeBaselineCycle, profiles[0].ArchetypeID)
	assert.Equal(t, ArchetypeDrawTempo, profiles[1].ArchetypeID)
	assert.Equal(t, ArchetypeHorizontalPressure, profiles[2].ArchetypeID)
	assert.Equal(t, ArchetypeVerticalEscape, profiles[3].ArchetypeID)
	assert.Equal(t, ArchetypeSurfaceControl, profiles[4].ArchetypeID)
	assert.Equal(t, ArchetypeDeckExhaustion, profiles[5].ArchetypeID)
	assert.Equal(t, ArchetypeHybridPressure, profiles[6].ArchetypeID)
	assert.NotEmpty(t, profiles[1].Details)
	require.Len(t, profiles[1].Cards, 4)
	assert.Equal(t, "Tiron de apertura", profiles[1].Cards[0].Name)
	assert.Equal(t, cards.DiscardVisibilityMoveOnly, profiles[4].Cards[1].DiscardVisibility)
	assert.Equal(t, 1, profiles[0].Splash.BuildEncounterProfile().JumpCount)
}

func TestSplashProfileBuildEncounterProfileUsesDefaultsAndOverrides(t *testing.T) {
	assert.Equal(t, 1, (SplashProfile{}).BuildEncounterProfile().JumpCount)
	assert.Equal(t, 3, (SplashProfile{JumpCount: 3, TimeLimitMillis: 1500}).BuildEncounterProfile().JumpCount)
	assert.Equal(t, 1500, int((SplashProfile{JumpCount: 3, TimeLimitMillis: 1500}).BuildEncounterProfile().TimeLimit.Milliseconds()))
}
