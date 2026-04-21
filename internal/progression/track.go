package progression

import (
	"pesca/internal/domain"
	"pesca/internal/match"
)

type TrackPolicy struct{}

func (TrackPolicy) Apply(state *match.State, outcome domain.RoundOutcome) {
	switch outcome {
	case domain.PlayerWin:
		state.Stats.PlayerWins++
		state.Encounter.Distance -= state.Encounter.Config.PlayerWinStep
	case domain.FishWin:
		state.Stats.FishWins++
		state.Encounter.Distance += state.Encounter.Config.FishWinStep
	default:
		state.Stats.Draws++
	}
}
