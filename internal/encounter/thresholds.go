package encounter

import "pesca/internal/cards"

type RoundThresholds struct {
	CaptureDistanceBonus           int
	ExhaustionCaptureDistanceBonus int
	SurfaceDepthBonus              int
}

func ApplyThresholdEffects(thresholds *RoundThresholds, effects []cards.CardEffect) {
	for _, effect := range effects {
		thresholds.CaptureDistanceBonus += effect.CaptureDistanceBonus
		thresholds.ExhaustionCaptureDistanceBonus += effect.ExhaustionCaptureDistanceBonus
		thresholds.SurfaceDepthBonus += effect.SurfaceDepthBonus
	}
}

func EffectiveCaptureDistance(config Config, thresholds RoundThresholds) int {
	return config.CaptureDistance + thresholds.CaptureDistanceBonus
}

func EffectiveSurfaceDepth(config Config, thresholds RoundThresholds) int {
	return config.SurfaceDepth + thresholds.SurfaceDepthBonus
}

func EffectiveExhaustionCaptureDistance(config Config, thresholds RoundThresholds) int {
	return config.ExhaustionCaptureDistance + thresholds.ExhaustionCaptureDistanceBonus
}

func IsTrackCapture(state State, thresholds RoundThresholds) bool {
	return state.Distance <= EffectiveCaptureDistance(state.Config, thresholds) && state.Depth <= EffectiveSurfaceDepth(state.Config, thresholds)
}

func IsDeckExhaustionCapture(state State, thresholds RoundThresholds) bool {
	return state.Distance <= EffectiveExhaustionCaptureDistance(state.Config, thresholds) && state.Depth <= EffectiveSurfaceDepth(state.Config, thresholds)+1
}
