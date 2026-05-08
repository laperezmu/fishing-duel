package app

import (
	"testing"

	"pesca/internal/content/fishprofiles"
	"pesca/internal/content/watercontexts"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSandboxConfigValidate(t *testing.T) {
	t.Run("accepts scenario config", func(t *testing.T) {
		seed := int64(42)
		distance := 2
		config := SandboxConfig{
			Mode:               SandboxModeScenario,
			PlayerDeckPresetID: "classic",
			FishPresetID:       fishprofiles.ProfileID("classic"),
			WaterContextID:     watercontexts.ShorelineCove,
			Seed:               &seed,
			ScenarioID:         "qa-surface-control",
			CardSelections: []SandboxCardSelection{{
				OwnerScope:     SandboxCardOwnerPlayer,
				SourcePresetID: "classic",
				SelectionMode:  SandboxCardSelectionOverlay,
				SelectedCards: []SandboxCardSelectionEntry{{
					CardRef: "blue-1",
					Origin:  SandboxCardOriginManualReplacement,
				}},
			}},
			StateOverrides: SandboxStateOverrides{InitialDistance: &distance},
		}

		require.NoError(t, config.Validate())
	})

	t.Run("rejects unknown mode", func(t *testing.T) {
		config := SandboxConfig{Mode: SandboxMode("broken")}

		err := config.Validate()
		require.Error(t, err)
		assert.EqualError(t, err, "unsupported sandbox mode \"broken\"")
	})

	t.Run("requires scenario id in scenario mode", func(t *testing.T) {
		config := SandboxConfig{Mode: SandboxModeScenario}

		err := config.Validate()
		require.Error(t, err)
		assert.EqualError(t, err, "scenario id is required in scenario mode")
	})
}

func TestSandboxCardSelectionValidate(t *testing.T) {
	t.Run("rejects duplicate refs", func(t *testing.T) {
		selection := SandboxCardSelection{
			OwnerScope:     SandboxCardOwnerFish,
			SourcePresetID: "classic",
			SelectionMode:  SandboxCardSelectionScenario,
			SelectedCards:  []SandboxCardSelectionEntry{{CardRef: "red-1", Origin: SandboxCardOriginScenarioDefined}, {CardRef: "red-1", Origin: SandboxCardOriginScenarioDefined}},
		}

		err := selection.Validate()
		require.Error(t, err)
		assert.EqualError(t, err, "selected card 1: duplicated card ref \"red-1\"")
	})
}

func TestSandboxScenarioValidate(t *testing.T) {
	scenario := SandboxScenario{ID: "surface-control", Name: "Surface control"}

	require.NoError(t, scenario.Validate())
}

func TestSandboxConfigResolvePlayerPreset(t *testing.T) {
	preset, err := (SandboxConfig{PlayerDeckPresetID: "classic"}).ResolvePlayerPreset()

	require.NoError(t, err)
	assert.Equal(t, "classic", preset.ID)
}
