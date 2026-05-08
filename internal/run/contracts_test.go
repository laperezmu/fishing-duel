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

func TestNodeKindValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		kind    NodeKind
		wantErr bool
	}{
		{"valid start", NodeKindStart, false},
		{"valid fishing", NodeKindFishing, false},
		{"valid service", NodeKindService, false},
		{"valid checkpoint", NodeKindCheckpoint, false},
		{"valid boss", NodeKindBoss, false},
		{"valid end", NodeKindEnd, false},
		{"invalid kind", NodeKind("invalid"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.kind.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestStatusValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		status  Status
		wantErr bool
	}{
		{"valid in progress", StatusInProgress, false},
		{"valid victory", StatusVictory, false},
		{"valid defeat", StatusDefeat, false},
		{"valid retired", StatusRetired, false},
		{"invalid status", Status("invalid"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.status.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestEncounterOutcomeValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		outcome EncounterOutcome
		wantErr bool
	}{
		{"valid captured", EncounterOutcomeCaptured, false},
		{"valid escaped", EncounterOutcomeEscaped, false},
		{"valid defeated", EncounterOutcomeDefeated, false},
		{"valid retired", EncounterOutcomeRetired, false},
		{"invalid outcome", EncounterOutcome("invalid"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.outcome.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCaptureRecordValidate(t *testing.T) {
	t.Parallel()

	t.Run("valid capture record", func(t *testing.T) {
		t.Parallel()

		record := CaptureRecord{
			FishID:      "bass-1",
			FishName:    "Lubina",
			EncounterID: "enc-1",
		}

		require.NoError(t, record.Validate())
	})

	t.Run("missing fish id", func(t *testing.T) {
		t.Parallel()

		record := CaptureRecord{
			FishName: "Lubina",
		}

		err := record.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "fish id is required")
	})

	t.Run("missing fish name", func(t *testing.T) {
		t.Parallel()

		record := CaptureRecord{
			FishID: "bass-1",
		}

		err := record.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "fish name is required")
	})
}

func TestNodeStateValidate(t *testing.T) {
	t.Parallel()

	t.Run("valid node state without preset", func(t *testing.T) {
		t.Parallel()

		node := NodeState{
			ZoneID: "coast",
			NodeID: "fish-1",
			Kind:   NodeKindFishing,
		}

		require.NoError(t, node.Validate())
	})

	t.Run("missing zone id", func(t *testing.T) {
		t.Parallel()

		node := NodeState{
			NodeID: "fish-1",
			Kind:   NodeKindFishing,
		}

		err := node.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "zone id is required")
	})

	t.Run("missing node id", func(t *testing.T) {
		t.Parallel()

		node := NodeState{
			ZoneID: "coast",
			Kind:   NodeKindFishing,
		}

		err := node.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "node id is required")
	})

	t.Run("invalid water preset", func(t *testing.T) {
		t.Parallel()

		node := NodeState{
			ZoneID:        "coast",
			NodeID:        "fish-1",
			Kind:          NodeKindFishing,
			WaterPresetID: "invalid-preset",
		}

		err := node.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown water context preset")
	})
}

func TestEncounterResultValidate(t *testing.T) {
	t.Parallel()

	t.Run("valid captured result", func(t *testing.T) {
		t.Parallel()

		result := EncounterResult{
			Outcome:       EncounterOutcomeCaptured,
			Status:        encounter.StatusCaptured,
			EndReason:     encounter.EndReasonTrackCapture,
			Capture:       &CaptureRecord{FishID: "bass-1", FishName: "Lubina"},
			NodeResolved:  true,
			FinishedMatch: true,
		}

		require.NoError(t, result.Validate())
	})

	t.Run("negative thread damage", func(t *testing.T) {
		t.Parallel()

		result := EncounterResult{
			Outcome:      EncounterOutcomeEscaped,
			ThreadDamage: -1,
		}

		err := result.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "thread damage must be greater than or equal to 0")
	})

	t.Run("ongoing status not allowed", func(t *testing.T) {
		t.Parallel()

		result := EncounterResult{
			Outcome: EncounterOutcomeEscaped,
			Status:  encounter.StatusOngoing,
		}

		err := result.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot use ongoing status")
	})

	t.Run("invalid capture data", func(t *testing.T) {
		t.Parallel()

		result := EncounterResult{
			Outcome:       EncounterOutcomeCaptured,
			Status:        encounter.StatusCaptured,
			Capture:       &CaptureRecord{FishID: "", FishName: ""},
			NodeResolved:  true,
			FinishedMatch: true,
		}

		err := result.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "capture:")
	})
}

func TestStateValidateCurrency(t *testing.T) {
	t.Parallel()

	playerRod, err := rod.NewState(rod.DefaultConfig())
	require.NoError(t, err)
	playerLoadout, err := loadout.NewState(playerRod, nil)
	require.NoError(t, err)

	t.Run("negative currency is invalid", func(t *testing.T) {
		t.Parallel()

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
			Currency: -1,
		}

		err := state.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "currency must be greater than or equal to 0")
	})
}
