package progression

import (
	"pesca/internal/domain"
	"pesca/internal/game"
)

type TrackPolicy struct{}

func (TrackPolicy) Apply(state *game.State, outcome domain.RoundOutcome) {
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
