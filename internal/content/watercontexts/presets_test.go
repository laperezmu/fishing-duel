package watercontexts

import (
	"pesca/internal/encounter"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultPresetsBuildValidContexts(t *testing.T) {
	presets := DefaultPresets()
	require.Len(t, presets, 9)

	for _, preset := range presets {
		context := preset.BuildContext()
		assert.Equal(t, preset.Name, context.Name)
		assert.Equal(t, preset.PoolTag, context.PoolTag)
		require.NoError(t, context.Validate())
	}

	assert.Equal(t, 0, presets[0].BuildContext().BandInitialDistance[encounter.CastBandVeryShort])
	assert.Equal(t, 1, presets[0].BuildContext().BandInitialDistance[encounter.CastBandShort])
	assert.Equal(t, 2, presets[0].BuildContext().BandInitialDistance[encounter.CastBandMedium])
	assert.Equal(t, 3, presets[0].BuildContext().BandInitialDistance[encounter.CastBandLong])
	assert.Equal(t, 4, presets[0].BuildContext().BandInitialDistance[encounter.CastBandVeryLong])
}

func TestResolveDefaultPreset(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		id        ID
		wantFound bool
	}{
		{"shoreline cove", ShorelineCove, true},
		{"open channel", OpenChannel, true},
		{"broken current", BrokenCurrent, true},
		{"reef shadow", ReefShadow, true},
		{"tidal gate", TidalGate, true},
		{"weed pocket", WeedPocket, true},
		{"stone drop", StoneDrop, true},
		{"wind lane", WindLane, true},
		{"deep ledge", DeepLedge, true},
		{"invalid id", ID("invalid-id"), false},
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

func TestDefaultPhaseLabel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		zoneID ID
		want   string
	}{
		{"shoreline cove", ShorelineCove, "Fase 1 - Ensenada cercana"},
		{"open channel", OpenChannel, "Fase 2 - Canal abierto"},
		{"broken current", BrokenCurrent, "Fase 3 - Corriente irregular"},
		{"reef shadow", ReefShadow, "Fase 4 - Sombra de arrecife"},
		{"tidal gate", TidalGate, "Fase 5 - Paso de marea"},
		{"weed pocket", WeedPocket, "Fase 6 - Bolsillo de maleza"},
		{"stone drop", StoneDrop, "Fase 7 - Caida de piedra"},
		{"wind lane", WindLane, "Calle de viento"},
		{"deep ledge", DeepLedge, "Fase 8 - Cornisa profunda"},
		{"unknown zone", ID("unknown"), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			label := DefaultPhaseLabel(tt.zoneID)
			assert.Equal(t, tt.want, label)
		})
	}
}
