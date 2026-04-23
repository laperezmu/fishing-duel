package game

import (
	"pesca/internal/cards"
	"pesca/internal/match"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyRoundScopedEffects(t *testing.T) {
	state := match.State{}

	applyRoundScopedEffects(&state, []cards.CardEffect{{
		CaptureDistanceBonus: 1,
	}, {
		ExhaustionCaptureDistanceBonus: -1,
		SurfaceDepthBonus:              2,
	}})

	assert.Equal(t, 1, state.RoundState.Thresholds.CaptureDistanceBonus)
	assert.Equal(t, -1, state.RoundState.Thresholds.ExhaustionCaptureDistanceBonus)
	assert.Equal(t, 2, state.RoundState.Thresholds.SurfaceDepthBonus)
}
