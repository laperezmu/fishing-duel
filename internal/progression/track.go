package progression

import (
	"pesca/internal/cards"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/match"
)

type TrackPolicy struct{}

func (policy TrackPolicy) Apply(state *match.ProgressionState, round match.ResolvedRound) {
	captureDistance := encounter.EffectiveCaptureDistance(state.Encounter.Config, state.Round.Thresholds)
	surfaceDepth := encounter.EffectiveSurfaceDepth(state.Encounter.Config, state.Round.Thresholds)
	baseDelta := encounter.Delta{}

	switch round.Outcome {
	case domain.PlayerWin:
		state.Lifecycle.Stats.PlayerWins++
		if state.Encounter.Distance <= captureDistance && state.Encounter.Depth > surfaceDepth {
			baseDelta.DepthShift--
		} else {
			baseDelta.DistanceShift -= state.Encounter.Config.PlayerWinStep
		}
	case domain.FishWin:
		state.Lifecycle.Stats.FishWins++
		baseDelta.DistanceShift += state.Encounter.Config.FishWinStep
	default:
		state.Lifecycle.Stats.Draws++
	}

	encounter.ApplyDelta(state.Encounter, baseDelta)
	encounter.ApplyMovementEffects(state.Encounter, round.OutcomeEffects)
}

func accumulateCardEffects(delta encounter.Delta, effects []cards.CardEffect) encounter.Delta {
	for _, effect := range effects {
		effect = effect.Normalize()
		delta.DistanceShift += effect.DistanceShift
		delta.DepthShift += effect.DepthShift
	}

	return delta
}
