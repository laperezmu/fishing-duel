package encounter

import (
	"pesca/internal/content/waterpools"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveOpening(t *testing.T) {
	baseConfig := DefaultConfig()
	openingLimits := OpeningLimits{MaxInitialDistance: 6, MaxInitialDepth: 3}
	waterContext := WaterContext{
		ID:               "shoreline-cove",
		Name:             "Ensenada cercana",
		Description:      "Actividad cerca de la orilla.",
		VisibleSignals:   []string{"Espuma corta."},
		PoolTag:          waterpools.Shoreline,
		BaseInitialDepth: 1,
		BandInitialDistance: map[CastBand]int{
			CastBandVeryShort: 3,
			CastBandShort:     3,
			CastBandMedium:    4,
			CastBandLong:      5,
			CastBandVeryLong:  6,
		},
	}

	opening, err := ResolveOpening(baseConfig, waterContext, CastResult{Band: CastBandShort}, openingLimits)

	require.NoError(t, err)
	assert.Equal(t, CastBandShort, opening.CastResult.Band)
	assert.Equal(t, 3, opening.InitialDistance)
	assert.Equal(t, 1, opening.InitialDepth)
	assert.Equal(t, 3, opening.Config.InitialDistance)
	assert.Equal(t, 1, opening.Config.InitialDepth)
}

func TestResolveOpeningAllowsCloseCastBands(t *testing.T) {
	baseConfig := DefaultConfig()
	openingLimits := OpeningLimits{MaxInitialDistance: 4, MaxInitialDepth: 3}
	waterContext := WaterContext{
		ID:               "tight-pocket",
		Name:             "Bolsillo cercano",
		Description:      "Actividad pegada a la orilla.",
		VisibleSignals:   []string{"Ondas cortas junto a la costa."},
		PoolTag:          waterpools.Shoreline,
		BaseInitialDepth: 0,
		BandInitialDistance: map[CastBand]int{
			CastBandVeryShort: 0,
			CastBandShort:     1,
			CastBandMedium:    2,
			CastBandLong:      3,
			CastBandVeryLong:  4,
		},
	}

	opening, err := ResolveOpening(baseConfig, waterContext, CastResult{Band: CastBandShort}, openingLimits)

	require.NoError(t, err)
	assert.Equal(t, 1, opening.InitialDistance)
	assert.Equal(t, 0, opening.InitialDepth)
}

func TestResolveOpeningClampsInitialValuesToOpeningLimits(t *testing.T) {
	baseConfig := DefaultConfig()
	openingLimits := OpeningLimits{MaxInitialDistance: 3, MaxInitialDepth: 2}
	waterContext := WaterContext{
		ID:               "deep-channel",
		Name:             "Canal profundo",
		Description:      "Actividad lejos y mas abajo.",
		VisibleSignals:   []string{"Remolinos largos."},
		PoolTag:          waterpools.Offshore,
		BaseInitialDepth: 4,
		BandInitialDistance: map[CastBand]int{
			CastBandVeryShort: 2,
			CastBandShort:     3,
			CastBandMedium:    4,
			CastBandLong:      5,
			CastBandVeryLong:  6,
		},
	}

	opening, err := ResolveOpening(baseConfig, waterContext, CastResult{Band: CastBandVeryLong}, openingLimits)

	require.NoError(t, err)
	assert.Equal(t, 3, opening.InitialDistance)
	assert.Equal(t, 2, opening.InitialDepth)
	assert.Equal(t, 3, opening.Config.InitialDistance)
	assert.Equal(t, 2, opening.Config.InitialDepth)
}

func TestWaterContextValidate(t *testing.T) {
	t.Run("requires every cast band to be configured", func(t *testing.T) {
		waterContext := WaterContext{
			ID:               "broken",
			Name:             "Incompleto",
			PoolTag:          waterpools.Shoreline,
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
			PoolTag:          waterpools.Shoreline,
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

func TestOpeningLimitsValidate(t *testing.T) {
	tests := []struct {
		title       string
		limits      OpeningLimits
		wantErrText string
	}{
		{
			title:  "accepts non negative limits",
			limits: OpeningLimits{MaxInitialDistance: 5, MaxInitialDepth: 3},
		},
		{
			title:       "rejects negative distance",
			limits:      OpeningLimits{MaxInitialDistance: -1, MaxInitialDepth: 3},
			wantErrText: "opening max initial distance must be greater than or equal to 0",
		},
		{
			title:       "rejects negative depth",
			limits:      OpeningLimits{MaxInitialDistance: 5, MaxInitialDepth: -1},
			wantErrText: "opening max initial depth must be greater than or equal to 0",
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			err := test.limits.Validate()

			if test.wantErrText != "" {
				require.Error(t, err)
				assert.EqualError(t, err, test.wantErrText)
				return
			}

			require.NoError(t, err)
		})
	}
}
