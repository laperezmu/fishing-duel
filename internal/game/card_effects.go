package game

import (
	"pesca/internal/cards"
	"pesca/internal/match"
)

func applyRoundScopedEffects(state *match.State, effects []cards.CardEffect) {
	for _, effect := range effects {
		state.RoundState.Thresholds.CaptureDistanceBonus += effect.CaptureDistanceBonus
		state.RoundState.Thresholds.ExhaustionCaptureDistanceBonus += effect.ExhaustionCaptureDistanceBonus
		state.RoundState.Thresholds.SurfaceDepthBonus += effect.SurfaceDepthBonus
	}
}
