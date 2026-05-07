package app

import (
	"fmt"
	"pesca/internal/content/fishprofiles"
	"pesca/internal/encounter"
	"pesca/internal/match"
	"pesca/internal/run"
)

func ResolveEncounterResult(state match.State, spawn fishprofiles.Spawn) (run.EncounterResult, error) {
	if !state.Lifecycle.Finished {
		return run.EncounterResult{}, fmt.Errorf("encounter result requires a finished match")
	}

	result := run.EncounterResult{
		Status:        state.Encounter.Status,
		EndReason:     state.Encounter.EndReason,
		FinishedMatch: state.Lifecycle.Finished,
	}

	switch state.Encounter.Status {
	case encounter.StatusCaptured:
		result.Outcome = run.EncounterOutcomeCaptured
		result.NodeResolved = true
		result.Capture = &run.CaptureRecord{
			FishID:   string(spawn.Profile.ID),
			FishName: spawn.Profile.Name,
		}
	case encounter.StatusEscaped:
		result.Outcome = run.EncounterOutcomeEscaped
		result.NodeResolved = true
		result.Retryable = false
		result.ThreadDamage = resolveEncounterThreadDamage(state)
	default:
		return run.EncounterResult{}, fmt.Errorf("unsupported encounter status %q", state.Encounter.Status)
	}

	if err := result.Validate(); err != nil {
		return run.EncounterResult{}, err
	}

	return result, nil
}

func resolveEncounterThreadDamage(state match.State) int {
	switch state.Encounter.EndReason {
	case encounter.EndReasonTrackEscape, encounter.EndReasonDepthEscape:
		return 1
	case encounter.EndReasonSplashEscape:
		return 0
	default:
		return 0
	}
}
