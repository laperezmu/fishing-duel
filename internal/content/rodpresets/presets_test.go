package rodpresets

import (
	"pesca/internal/player/loadout"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPresetBuildLoadout(t *testing.T) {
	preset := Preset{
		ID:          "test",
		Name:        "Test",
		Description: "Preset de prueba.",
		Config:      DefaultPresets()[0].Config,
	}

	loadoutState, err := preset.BuildLoadout()

	require.NoError(t, err)
	assert.Equal(t, preset.Config.OpeningMaxDistance, loadoutState.Rod.OpeningMaxDistance)
	assert.Equal(t, preset.Config.OpeningMaxDepth, loadoutState.Rod.OpeningMaxDepth)
	assert.Equal(t, preset.Config.TrackMaxDistance, loadoutState.Rod.TrackMaxDistance)
	assert.Equal(t, preset.Config.TrackMaxDepth, loadoutState.Rod.TrackMaxDepth)
	assert.Empty(t, loadoutState.Attachments)
}

func TestDefaultPresets(t *testing.T) {
	presets := DefaultPresets()

	require.Len(t, presets, 3)
	assert.Equal(t, "Control costero", presets[0].Name)
	assert.Equal(t, 3, presets[0].Config.OpeningMaxDistance)
	assert.Equal(t, 5, presets[0].Config.TrackMaxDistance)
	assert.Equal(t, "Presion de fondo", presets[2].Name)
	assert.Equal(t, 4, presets[2].Config.OpeningMaxDepth)
	assert.Equal(t, 5, presets[2].Config.TrackMaxDepth)
}

func TestBuildLoadoutWithAttachments(t *testing.T) {
	preset := DefaultPresets()[1]
	loadoutState, err := preset.BuildLoadoutWithAttachments([]loadout.Attachment{{
		ID:                    "line-reinforced",
		Name:                  "Linea reforzada",
		TrackDistanceModifier: 1,
	}})

	require.NoError(t, err)
	assert.Equal(t, 6, loadoutState.TrackMaxDistance())
	assert.Len(t, loadoutState.Attachments, 1)
}

func TestResolveDefaultPreset(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		id        string
		wantFound bool
	}{
		{"coastal control", "coastal-control", true},
		{"versatile standard", "versatile-standard", true},
		{"bottom pressure", "bottom-pressure", true},
		{"invalid id", "invalid-id", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			preset, err := ResolveDefaultPreset(tt.id)
			if tt.wantFound {
				require.NoError(t, err)
				assert.Equal(t, tt.id, preset.ID)
			} else {
				require.Error(t, err)
			}
		})
	}
}
