package endings

import (
	"pesca/internal/encounter"
	"pesca/internal/match"
	"pesca/internal/player/playerrig"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncounterEndConditionApply(t *testing.T) {
	tests := []struct {
		title         string
		state         match.State
		wantFinished  bool
		wantStatus    encounter.Status
		wantEndReason encounter.EndReason
	}{
		{
			title:         "captures the fish when distance reaches the capture threshold at the surface",
			state:         newMatchState(t, 0, 0),
			wantFinished:  true,
			wantStatus:    encounter.StatusCaptured,
			wantEndReason: encounter.EndReasonTrackCapture,
		},
		{
			title:         "keeps the encounter open when the fish is close enough but still below the surface",
			state:         newMatchState(t, 0, 1),
			wantFinished:  false,
			wantStatus:    encounter.StatusOngoing,
			wantEndReason: encounter.EndReasonNone,
		},
		{
			title: "captures when a temporary round bonus extends the capture distance",
			state: func() match.State {
				state := newMatchState(t, 1, 0)
				state.RoundState.Thresholds.CaptureDistanceBonus = 1
				return state
			}(),
			wantFinished:  true,
			wantStatus:    encounter.StatusCaptured,
			wantEndReason: encounter.EndReasonTrackCapture,
		},
		{
			title: "captures when a temporary round bonus raises the effective surface threshold",
			state: func() match.State {
				state := newMatchState(t, 0, 1)
				state.RoundState.Thresholds.SurfaceDepthBonus = 1
				return state
			}(),
			wantFinished:  true,
			wantStatus:    encounter.StatusCaptured,
			wantEndReason: encounter.EndReasonTrackCapture,
		},
		{
			title:         "escapes when the fish exceeds the player's max distance",
			state:         newMatchState(t, 6, 1),
			wantFinished:  true,
			wantStatus:    encounter.StatusEscaped,
			wantEndReason: encounter.EndReasonTrackEscape,
		},
		{
			title:         "escapes when the fish exceeds the player's max depth",
			state:         newMatchState(t, 3, 4),
			wantFinished:  true,
			wantStatus:    encounter.StatusEscaped,
			wantEndReason: encounter.EndReasonDepthEscape,
		},
		{
			title: "preserves splash escape set by encounter progression",
			state: func() match.State {
				state := newMatchState(t, 3, 0)
				state.Encounter.Status = encounter.StatusEscaped
				state.Encounter.EndReason = encounter.EndReasonSplashEscape
				return state
			}(),
			wantFinished:  true,
			wantStatus:    encounter.StatusEscaped,
			wantEndReason: encounter.EndReasonSplashEscape,
		},
		{
			title: "captures when the deck is exhausted near the player and close to the surface",
			state: func() match.State {
				state := newMatchState(t, 2, 1)
				state.Deck.Exhausted = true
				return state
			}(),
			wantFinished:  true,
			wantStatus:    encounter.StatusCaptured,
			wantEndReason: encounter.EndReasonDeckCapture,
		},
		{
			title: "captures on deck exhaustion when a temporary round bonus extends the exhaustion threshold",
			state: func() match.State {
				state := newMatchState(t, 3, 1)
				state.Deck.Exhausted = true
				state.RoundState.Thresholds.ExhaustionCaptureDistanceBonus = 1
				return state
			}(),
			wantFinished:  true,
			wantStatus:    encounter.StatusCaptured,
			wantEndReason: encounter.EndReasonDeckCapture,
		},
		{
			title: "escapes when the deck is exhausted near the player but the fish is too deep",
			state: func() match.State {
				state := newMatchState(t, 2, 2)
				state.Deck.Exhausted = true
				return state
			}(),
			wantFinished:  true,
			wantStatus:    encounter.StatusEscaped,
			wantEndReason: encounter.EndReasonDeckEscape,
		},
		{
			title: "escapes when the deck is exhausted far from the player",
			state: func() match.State {
				state := newMatchState(t, 4, 1)
				state.Deck.Exhausted = true
				return state
			}(),
			wantFinished:  true,
			wantStatus:    encounter.StatusEscaped,
			wantEndReason: encounter.EndReasonDeckEscape,
		},
		{
			title:         "remains ongoing when no terminal condition is met",
			state:         newMatchState(t, 3, 1),
			wantFinished:  false,
			wantStatus:    encounter.StatusOngoing,
			wantEndReason: encounter.EndReasonNone,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			state := test.state

			EncounterEndCondition{}.Apply(&state)

			assert.Equal(t, test.wantFinished, state.Finished)
			assert.Equal(t, test.wantStatus, state.Encounter.Status)
			assert.Equal(t, test.wantEndReason, state.Encounter.EndReason)
		})
	}
}

func newMatchState(t *testing.T, distance, depth int) match.State {
	t.Helper()

	encounterState, err := encounter.NewState(encounter.DefaultConfig())
	require.NoError(t, err)
	encounterState.Distance = distance
	encounterState.Depth = depth
	playerRigState, err := playerrig.NewState(playerrig.DefaultConfig())
	require.NoError(t, err)

	return match.State{
		Encounter: encounterState,
		PlayerRig: playerRigState,
	}
}
