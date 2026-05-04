package run

import (
	"pesca/internal/encounter"
	"pesca/internal/player/loadout"
	"pesca/internal/player/rod"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testThreadMaximum = 3

func TestNewStateUsesFirstRouteNode(t *testing.T) {
	t.Parallel()

	state, err := NewState(mustBuildRouteLoadout(t), DefaultRoute(), testThreadMaximum)
	require.NoError(t, err)
	assert.Equal(t, NodeKindStart, state.Progress.Current.Kind)
	assert.NotNil(t, state.Progress.Next)
	assert.Equal(t, NodeKindFishing, state.Progress.Next.Kind)
	assert.Equal(t, testThreadMaximum, state.Thread.Current)
}

func TestAdvanceMovesToNextNode(t *testing.T) {
	t.Parallel()

	route := DefaultRoute()
	state, err := NewState(mustBuildRouteLoadout(t), route, testThreadMaximum)
	require.NoError(t, err)

	require.NoError(t, Advance(&state, route))
	assert.Equal(t, "fishing-1", state.Progress.Current.NodeID)
	assert.NotNil(t, state.Progress.Next)
	assert.Equal(t, "fishing-2", state.Progress.Next.NodeID)
}

func TestApplyEncounterResultUpdatesThreadAndCaptures(t *testing.T) {
	t.Parallel()

	state, err := NewState(mustBuildRouteLoadout(t), DefaultRoute(), testThreadMaximum)
	require.NoError(t, err)

	err = ApplyEncounterResult(&state, EncounterResult{
		Outcome:       EncounterOutcomeEscaped,
		Status:        encounter.StatusEscaped,
		EndReason:     encounter.EndReasonTrackEscape,
		ThreadDamage:  2,
		NodeResolved:  true,
		FinishedMatch: true,
	})
	require.NoError(t, err)
	assert.Equal(t, 1, state.Thread.Current)

	err = ApplyEncounterResult(&state, EncounterResult{
		Outcome:       EncounterOutcomeCaptured,
		Status:        encounter.StatusCaptured,
		EndReason:     encounter.EndReasonTrackCapture,
		Capture:       &CaptureRecord{FishID: "tuna", FishName: "Atun"},
		NodeResolved:  true,
		FinishedMatch: true,
	})
	require.NoError(t, err)
	assert.Len(t, state.Captures, 1)
	assert.Equal(t, StatusInProgress, state.Status)
}

func TestApplyEncounterResultDefeatsRunWhenThreadReachesZero(t *testing.T) {
	t.Parallel()

	state, err := NewState(mustBuildRouteLoadout(t), DefaultRoute(), testThreadMaximum)
	require.NoError(t, err)

	err = ApplyEncounterResult(&state, EncounterResult{
		Outcome:       EncounterOutcomeEscaped,
		Status:        encounter.StatusEscaped,
		EndReason:     encounter.EndReasonTrackEscape,
		ThreadDamage:  testThreadMaximum,
		NodeResolved:  true,
		FinishedMatch: true,
	})
	require.NoError(t, err)
	assert.Equal(t, StatusDefeat, state.Status)
	assert.Zero(t, state.Thread.Current)
}

func mustBuildRouteLoadout(t *testing.T) loadout.State {
	t.Helper()

	playerRod, err := rod.NewState(rod.DefaultConfig())
	require.NoError(t, err)
	playerLoadout, err := loadout.NewState(playerRod, nil)
	require.NoError(t, err)

	return playerLoadout
}
