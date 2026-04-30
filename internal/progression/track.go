package progression

import (
	"pesca/internal/cards"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/match"
)

type TrackPolicy struct {
	SplashEscapeDecider encounter.SplashEscapeDecider
}

func (policy TrackPolicy) Apply(state *match.State, round match.ResolvedRound) {
	delta := encounter.Delta{}
	captureDistance := encounter.EffectiveCaptureDistance(state.Encounter.Config, state.Round.Thresholds)
	surfaceDepth := encounter.EffectiveSurfaceDepth(state.Encounter.Config, state.Round.Thresholds)

	switch round.Outcome {
	case domain.PlayerWin:
		state.Lifecycle.Stats.PlayerWins++
		if state.Encounter.Distance <= captureDistance && state.Encounter.Depth > surfaceDepth {
			delta.DepthShift--
		} else {
			delta.DistanceShift -= state.Encounter.Config.PlayerWinStep
		}
	case domain.FishWin:
		state.Lifecycle.Stats.FishWins++
		delta.DistanceShift += state.Encounter.Config.FishWinStep
	default:
		state.Lifecycle.Stats.Draws++
	}

	delta = accumulateCardEffects(delta, round.OutcomeEffects)

	encounter.ApplyDelta(&state.Encounter, delta, policy.SplashEscapeDecider)
}

func accumulateCardEffects(delta encounter.Delta, effects []cards.CardEffect) encounter.Delta {
	for _, effect := range effects {
		delta.DistanceShift += effect.DistanceShift
		delta.DepthShift += effect.DepthShift
	}

	return delta
}
