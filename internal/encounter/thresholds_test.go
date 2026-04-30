package encounter

import (
	"pesca/internal/cards"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyThresholdEffects(t *testing.T) {
	thresholds := RoundThresholds{}

	ApplyThresholdEffects(&thresholds, []cards.CardEffect{{
		CaptureDistanceBonus: 1,
	}, {
		ExhaustionCaptureDistanceBonus: -1,
		SurfaceDepthBonus:              2,
	}})

	assert.Equal(t, 1, thresholds.CaptureDistanceBonus)
	assert.Equal(t, -1, thresholds.ExhaustionCaptureDistanceBonus)
	assert.Equal(t, 2, thresholds.SurfaceDepthBonus)
}

func TestEffectiveThresholdHelpers(t *testing.T) {
	state, err := NewState(DefaultConfig())
	require.NoError(t, err)
	state.Distance = 1
	state.Depth = 0
	thresholds := RoundThresholds{CaptureDistanceBonus: 1, SurfaceDepthBonus: 0, ExhaustionCaptureDistanceBonus: 1}

	assert.Equal(t, 1, EffectiveCaptureDistance(state.Config, thresholds))
	assert.Equal(t, 0, EffectiveSurfaceDepth(state.Config, thresholds))
	assert.Equal(t, 3, EffectiveExhaustionCaptureDistance(state.Config, thresholds))
	assert.True(t, IsTrackCapture(state, thresholds))
	assert.True(t, IsDeckExhaustionCapture(state, thresholds))
}
