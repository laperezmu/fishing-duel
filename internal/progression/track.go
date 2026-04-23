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
	captureDistance := state.Encounter.Config.CaptureDistance + state.RoundState.Thresholds.CaptureDistanceBonus
	surfaceDepth := state.Encounter.Config.SurfaceDepth + state.RoundState.Thresholds.SurfaceDepthBonus

	switch round.Outcome {
	case domain.PlayerWin:
		state.Stats.PlayerWins++
		if state.Encounter.Distance <= captureDistance && state.Encounter.Depth > surfaceDepth {
			delta.DepthShift--
		} else {
			delta.DistanceShift -= state.Encounter.Config.PlayerWinStep
		}
	case domain.FishWin:
		state.Stats.FishWins++
		delta.DistanceShift += state.Encounter.Config.FishWinStep
	default:
		state.Stats.Draws++
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
