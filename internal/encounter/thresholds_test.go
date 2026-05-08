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

func TestIsTrackCaptureWithDifferentDistances(t *testing.T) {
	state, err := NewState(DefaultConfig())
	require.NoError(t, err)

	t.Run("returns true when distance is at capture distance with bonus", func(t *testing.T) {
		state.Distance = 1
		state.Depth = 0
		thresholds := RoundThresholds{CaptureDistanceBonus: 1}
		assert.True(t, IsTrackCapture(state, thresholds))
	})

	t.Run("returns false when distance is above capture distance", func(t *testing.T) {
		state.Distance = 2
		state.Depth = 0
		thresholds := RoundThresholds{}
		assert.False(t, IsTrackCapture(state, thresholds))
	})

	t.Run("returns false when depth is above surface depth", func(t *testing.T) {
		state.Distance = 0
		state.Depth = 1
		thresholds := RoundThresholds{}
		assert.False(t, IsTrackCapture(state, thresholds))
	})
}

func TestIsDeckExhaustionCaptureWithDifferentStates(t *testing.T) {
	state, err := NewState(DefaultConfig())
	require.NoError(t, err)
	thresholds := RoundThresholds{}

	t.Run("returns true when distance and depth are low", func(t *testing.T) {
		state.Distance = 1
		state.Depth = 0
		assert.True(t, IsDeckExhaustionCapture(state, thresholds))
	})

	t.Run("returns false when distance is too high", func(t *testing.T) {
		state.Distance = 3
		state.Depth = 0
		assert.False(t, IsDeckExhaustionCapture(state, thresholds))
	})

	t.Run("returns false when depth is too high", func(t *testing.T) {
		state.Distance = 1
		state.Depth = 2
		assert.False(t, IsDeckExhaustionCapture(state, thresholds))
	})
}
