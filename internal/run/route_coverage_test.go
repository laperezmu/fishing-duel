package run_test

import (
	"testing"

	"pesca/internal/content/watercontexts"
	"pesca/internal/player/loadout"
	"pesca/internal/player/rod"
	"pesca/internal/run"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustLoadout(t *testing.T) loadout.State {
	t.Helper()
	state, err := loadout.NewState(rod.State{OpeningMaxDistance: 4, OpeningMaxDepth: 3, TrackMaxDistance: 6, TrackMaxDepth: 4}, nil)
	require.NoError(t, err)
	return state
}

func TestComplete(t *testing.T) {
	t.Run("sets run status to victory", func(t *testing.T) {
		state, err := run.NewState(mustLoadout(t), run.DefaultRoute(), 10)
		require.NoError(t, err)

		err = run.Complete(&state)

		require.NoError(t, err)
		assert.Equal(t, run.StatusVictory, state.Status)
	})

	t.Run("returns error when state is nil", func(t *testing.T) {
		err := run.Complete(nil)

		assert.EqualError(t, err, "run state is required")
	})
}

func TestResolveWaterPreset(t *testing.T) {
	t.Run("resolves known water preset", func(t *testing.T) {
		preset, err := run.ResolveWaterPreset(watercontexts.ShorelineCove)

		require.NoError(t, err)
		assert.NotEmpty(t, preset)
	})
}

func TestStateValidate(t *testing.T) {
	t.Run("validates thread current cannot exceed maximum", func(t *testing.T) {
		state, err := run.NewState(mustLoadout(t), run.DefaultRoute(), 10)
		require.NoError(t, err)
		state.Thread.Current = 15

		err = state.Validate()

		assert.EqualError(t, err, "thread: run thread current must be less than or equal to maximum")
	})

	t.Run("validates progress zone id is required", func(t *testing.T) {
		state := run.State{
			Status: run.StatusInProgress,
			Progress: run.ProgressState{
				Current: run.NodeState{},
			},
		}

		err := state.Validate()

		assert.EqualError(t, err, "progress: current node: run node zone id is required")
	})
}
