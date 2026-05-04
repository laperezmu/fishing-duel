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
