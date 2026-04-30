package game

import (
	"pesca/internal/cards"
	"pesca/internal/match"
)

func applyRoundScopedEffects(state *match.State, effects []cards.CardEffect) {
	for _, effect := range effects {
		state.Round.Thresholds.CaptureDistanceBonus += effect.CaptureDistanceBonus
		state.Round.Thresholds.ExhaustionCaptureDistanceBonus += effect.ExhaustionCaptureDistanceBonus
		state.Round.Thresholds.SurfaceDepthBonus += effect.SurfaceDepthBonus
	}
}
