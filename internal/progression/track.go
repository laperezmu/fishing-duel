package progression

import (
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/match"
)

type TrackPolicy struct {
	SplashEscapeDecider encounter.SplashEscapeDecider
}

func (policy TrackPolicy) Apply(state *match.State, round match.ResolvedRound) {
	delta := encounter.Delta{}

	switch round.Outcome {
	case domain.PlayerWin:
		state.Stats.PlayerWins++
		if state.Encounter.Distance <= state.Encounter.Config.CaptureDistance && state.Encounter.Depth > state.Encounter.Config.SurfaceDepth {
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

	for _, modifier := range round.EncounterModifiers {
		if !modifier.Applies(round.Outcome) {
			continue
		}

		delta.DistanceShift += modifier.DistanceShift
		delta.DepthShift += modifier.DepthShift
	}

	encounter.ApplyDelta(&state.Encounter, delta, policy.SplashEscapeDecider)
}
