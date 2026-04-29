package encounter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveOpening(t *testing.T) {
	baseConfig := DefaultConfig()
	waterContext := WaterContext{
		ID:               "shoreline-cove",
		Name:             "Ensenada cercana",
		Description:      "Actividad cerca de la orilla.",
		VisibleSignals:   []string{"Espuma corta."},
		PoolTag:          "shoreline",
		BaseInitialDepth: 1,
		BandInitialDistance: map[CastBand]int{
			CastBandVeryShort: 3,
			CastBandShort:     3,
			CastBandMedium:    4,
			CastBandLong:      5,
			CastBandVeryLong:  6,
		},
	}

	opening, err := ResolveOpening(baseConfig, waterContext, CastResult{Band: CastBandShort})

	require.NoError(t, err)
	assert.Equal(t, CastBandShort, opening.CastResult.Band)
	assert.Equal(t, 3, opening.InitialDistance)
	assert.Equal(t, 1, opening.InitialDepth)
	assert.Equal(t, 3, opening.Config.InitialDistance)
	assert.Equal(t, 1, opening.Config.InitialDepth)
}

func TestResolveOpeningAllowsCloseCastBands(t *testing.T) {
	baseConfig := DefaultConfig()
	waterContext := WaterContext{
		ID:               "tight-pocket",
		Name:             "Bolsillo cercano",
		Description:      "Actividad pegada a la orilla.",
		VisibleSignals:   []string{"Ondas cortas junto a la costa."},
		PoolTag:          "pocket",
		BaseInitialDepth: 0,
		BandInitialDistance: map[CastBand]int{
			CastBandVeryShort: 0,
			CastBandShort:     1,
			CastBandMedium:    2,
			CastBandLong:      3,
			CastBandVeryLong:  4,
		},
	}

	opening, err := ResolveOpening(baseConfig, waterContext, CastResult{Band: CastBandShort})

	require.NoError(t, err)
	assert.Equal(t, 1, opening.InitialDistance)
	assert.Equal(t, 0, opening.InitialDepth)
}

func TestWaterContextValidate(t *testing.T) {
	t.Run("requires every cast band to be configured", func(t *testing.T) {
		waterContext := WaterContext{
			ID:               "broken",
			Name:             "Incompleto",
			BaseInitialDepth: 1,
			BandInitialDistance: map[CastBand]int{
				CastBandVeryShort: 3,
			},
		}

		err := waterContext.Validate()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "water context must define an initial distance")
	})

	t.Run("allows water contexts to start on the closest encounter cells", func(t *testing.T) {
		waterContext := WaterContext{
			ID:               "closest-cells",
			Name:             "Celdas cercanas",
			BaseInitialDepth: 0,
			BandInitialDistance: map[CastBand]int{
				CastBandVeryShort: 0,
				CastBandShort:     1,
				CastBandMedium:    2,
				CastBandLong:      3,
				CastBandVeryLong:  4,
			},
		}

		err := waterContext.Validate()

		require.NoError(t, err)
	})
}
