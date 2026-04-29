package watercontexts

import (
	"pesca/internal/encounter"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultPresetsBuildValidContexts(t *testing.T) {
	presets := DefaultPresets()
	require.Len(t, presets, 3)

	context := presets[0].BuildContext()
	assert.Equal(t, presets[0].Name, context.Name)
	assert.Equal(t, presets[0].PoolTag, context.PoolTag)
	assert.Equal(t, 0, context.BandInitialDistance[encounter.CastBandVeryShort])
	assert.Equal(t, 1, context.BandInitialDistance[encounter.CastBandShort])
	assert.Equal(t, 2, context.BandInitialDistance[encounter.CastBandMedium])
	assert.Equal(t, 3, context.BandInitialDistance[encounter.CastBandLong])
	assert.Equal(t, 4, context.BandInitialDistance[encounter.CastBandVeryLong])
	require.NoError(t, context.Validate())
}
