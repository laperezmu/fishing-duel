package endings

import (
	"pesca/internal/encounter"
	"pesca/internal/game"
)

type EncounterCondition struct{}

func (EncounterCondition) Apply(state *game.State) {
	state.Finished = false
	state.Encounter.Status = encounter.StatusOngoing
	state.Encounter.EndReason = encounter.EndReasonNone

	switch {
	case state.Encounter.Distance <= state.Encounter.Config.CaptureDistance:
		state.Encounter.Status = encounter.StatusCaptured
		state.Encounter.EndReason = encounter.EndReasonTrackCapture
	case state.Encounter.Distance > state.Encounter.Config.EscapeDistance:
		state.Encounter.Status = encounter.StatusEscaped
		state.Encounter.EndReason = encounter.EndReasonTrackEscape
	case state.Deck.Exhausted:
		if state.Encounter.Distance <= state.Encounter.Config.ExhaustionCaptureDistance {
			state.Encounter.Status = encounter.StatusCaptured
			state.Encounter.EndReason = encounter.EndReasonDeckCapture
		} else {
			state.Encounter.Status = encounter.StatusEscaped
			state.Encounter.EndReason = encounter.EndReasonDeckEscape
		}
	default:
		return
	}

	state.Finished = true
}
