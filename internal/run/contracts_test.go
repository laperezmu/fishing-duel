package run

import (
	"pesca/internal/encounter"
	"pesca/internal/player/loadout"
	"pesca/internal/player/rod"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStateValidate(t *testing.T) {
	t.Parallel()

	playerRod, err := rod.NewState(rod.DefaultConfig())
	require.NoError(t, err)
	playerLoadout, err := loadout.NewState(playerRod, nil)
	require.NoError(t, err)

	state := State{
		Status: StatusInProgress,
		Progress: ProgressState{
			ZoneIndex: 0,
			NodeIndex: 1,
			Current:   NodeState{ZoneID: "coast", NodeID: "fish-2", Kind: NodeKindFishing},
			Next:      &NodeState{ZoneID: "coast", NodeID: "service-1", Kind: NodeKindService},
		},
		Thread:   ThreadState{Current: 5, Maximum: 7},
		Loadout:  playerLoadout,
		Currency: 3,
		Captures: []CaptureRecord{{FishID: "mackerel", FishName: "Caballa"}},
	}

	assert.NoError(t, state.Validate())
}

func TestStateValidateRejectsInvalidProgress(t *testing.T) {
	t.Parallel()

	playerRod, err := rod.NewState(rod.DefaultConfig())
	require.NoError(t, err)
	playerLoadout, err := loadout.NewState(playerRod, nil)
	require.NoError(t, err)

	state := State{
		Status: StatusInProgress,
		Progress: ProgressState{
			ZoneIndex: 0,
			NodeIndex: 0,
			Current:   NodeState{},
		},
		Thread:  ThreadState{Current: 5, Maximum: 7},
		Loadout: playerLoadout,
	}

	err = state.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "current node")
}

func TestEncounterResultValidateRequiresCaptureWhenCaptured(t *testing.T) {
	t.Parallel()

	err := EncounterResult{
		Outcome:       EncounterOutcomeCaptured,
		Status:        encounter.StatusCaptured,
		EndReason:     encounter.EndReasonTrackCapture,
		NodeResolved:  true,
		FinishedMatch: true,
	}.Validate()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires capture data")
}
