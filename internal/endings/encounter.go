package endings

import (
	"pesca/internal/encounter"
	"pesca/internal/match"
)

type EncounterEndCondition struct{}

func (EncounterEndCondition) Apply(state *match.State) {
	if state.Encounter.EndReason == encounter.EndReasonSplashEscape {
		state.Finished = true
		state.Encounter.Status = encounter.StatusEscaped
		return
	}

	state.Finished = false
	state.Encounter.Status = encounter.StatusOngoing
	state.Encounter.EndReason = encounter.EndReasonNone

	switch {
	case state.Encounter.Distance <= state.Encounter.Config.CaptureDistance:
		state.Encounter.Status = encounter.StatusCaptured
		state.Encounter.EndReason = encounter.EndReasonTrackCapture
	case state.Encounter.Distance > state.PlayerRig.MaxDistance:
		state.Encounter.Status = encounter.StatusEscaped
		state.Encounter.EndReason = encounter.EndReasonTrackEscape
	case state.Encounter.Depth > state.PlayerRig.MaxDepth:
		state.Encounter.Status = encounter.StatusEscaped
		state.Encounter.EndReason = encounter.EndReasonDepthEscape
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
