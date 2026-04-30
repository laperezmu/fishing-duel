package endings

import (
	"pesca/internal/encounter"
	"pesca/internal/match"
)

type EncounterEndCondition struct{}

func (EncounterEndCondition) Apply(state *match.State) {
	if state.Encounter.EndReason == encounter.EndReasonSplashEscape {
		state.Lifecycle.Finished = true
		state.Encounter.Status = encounter.StatusEscaped
		return
	}

	state.Lifecycle.Finished = false
	state.Encounter.Status = encounter.StatusOngoing
	state.Encounter.EndReason = encounter.EndReasonNone

	switch {
	case encounter.IsTrackCapture(state.Encounter, state.Round.Thresholds):
		state.Encounter.Status = encounter.StatusCaptured
		state.Encounter.EndReason = encounter.EndReasonTrackCapture
	case state.Encounter.Distance > state.Player.Loadout.TrackMaxDistance():
		state.Encounter.Status = encounter.StatusEscaped
		state.Encounter.EndReason = encounter.EndReasonTrackEscape
	case state.Encounter.Depth > state.Player.Loadout.TrackMaxDepth():
		state.Encounter.Status = encounter.StatusEscaped
		state.Encounter.EndReason = encounter.EndReasonDepthEscape
	case state.Deck.Exhausted:
		if encounter.IsDeckExhaustionCapture(state.Encounter, state.Round.Thresholds) {
			state.Encounter.Status = encounter.StatusCaptured
			state.Encounter.EndReason = encounter.EndReasonDeckCapture
		} else {
			state.Encounter.Status = encounter.StatusEscaped
			state.Encounter.EndReason = encounter.EndReasonDeckEscape
		}
	default:
		return
	}

	state.Lifecycle.Finished = true
}
