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

	captureDistance := state.Encounter.Config.CaptureDistance + state.RoundState.Thresholds.CaptureDistanceBonus
	surfaceDepth := state.Encounter.Config.SurfaceDepth + state.RoundState.Thresholds.SurfaceDepthBonus
	exhaustionCaptureDistance := state.Encounter.Config.ExhaustionCaptureDistance + state.RoundState.Thresholds.ExhaustionCaptureDistanceBonus

	state.Finished = false
	state.Encounter.Status = encounter.StatusOngoing
	state.Encounter.EndReason = encounter.EndReasonNone

	switch {
	case state.Encounter.Distance <= captureDistance && state.Encounter.Depth <= surfaceDepth:
		state.Encounter.Status = encounter.StatusCaptured
		state.Encounter.EndReason = encounter.EndReasonTrackCapture
	case state.Encounter.Distance > state.PlayerRig.MaxDistance:
		state.Encounter.Status = encounter.StatusEscaped
		state.Encounter.EndReason = encounter.EndReasonTrackEscape
	case state.Encounter.Depth > state.PlayerRig.MaxDepth:
		state.Encounter.Status = encounter.StatusEscaped
		state.Encounter.EndReason = encounter.EndReasonDepthEscape
	case state.Deck.Exhausted:
		if state.Encounter.Distance <= exhaustionCaptureDistance && state.Encounter.Depth <= surfaceDepth+1 {
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
