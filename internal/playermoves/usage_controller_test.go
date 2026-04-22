package playermoves

import (
	"pesca/internal/domain"
	"pesca/internal/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsageControllerInitialize(t *testing.T) {
	controller, err := NewUsageController(Config{
		InitialUsesPerMove: map[domain.Move]int{
			domain.Blue:   2,
			domain.Red:    4,
			domain.Yellow: 1,
		},
		RecoveryDelayRounds: 1,
	})
	require.NoError(t, err)

	state := match.State{}
	controller.Initialize(&state)

	require.Len(t, state.PlayerMoves.Moves, 3)
	assert.Equal(t, domain.Blue, state.PlayerMoves.Moves[0].Move)
	assert.Equal(t, 2, state.PlayerMoves.Moves[0].MaxUses)
	assert.Equal(t, 2, state.PlayerMoves.Moves[0].RemainingUses)
	assert.Equal(t, domain.Red, state.PlayerMoves.Moves[1].Move)
	assert.Equal(t, 4, state.PlayerMoves.Moves[1].MaxUses)
	assert.Equal(t, 4, state.PlayerMoves.Moves[1].RemainingUses)
	assert.Equal(t, domain.Yellow, state.PlayerMoves.Moves[2].Move)
	assert.Equal(t, 1, state.PlayerMoves.Moves[2].MaxUses)
	assert.Equal(t, 1, state.PlayerMoves.Moves[2].RemainingUses)
}

func TestUsageControllerPrepareRound(t *testing.T) {
	controller, err := NewUsageController(DefaultConfig())
	require.NoError(t, err)

	tests := []struct {
		title                 string
		state                 match.State
		wantRemainingBlueUses int
		wantBlueRestoreRound  int
	}{
		{
			title: "restores a move when its recharge round has arrived",
			state: match.State{
				Round: 2,
				PlayerMoves: match.PlayerMoveResources{Moves: []match.PlayerMoveState{
					{Move: domain.Blue, MaxUses: 3, RemainingUses: 0, RestoresOnRound: 3},
					{Move: domain.Red, MaxUses: 3, RemainingUses: 2},
					{Move: domain.Yellow, MaxUses: 3, RemainingUses: 1},
				}},
			},
			wantRemainingBlueUses: 3,
			wantBlueRestoreRound:  0,
		},
		{
			title: "keeps a move blocked when its recharge round has not arrived yet",
			state: match.State{
				Round: 1,
				PlayerMoves: match.PlayerMoveResources{Moves: []match.PlayerMoveState{
					{Move: domain.Blue, MaxUses: 3, RemainingUses: 0, RestoresOnRound: 3},
					{Move: domain.Red, MaxUses: 3, RemainingUses: 2},
					{Move: domain.Yellow, MaxUses: 3, RemainingUses: 1},
				}},
			},
			wantRemainingBlueUses: 0,
			wantBlueRestoreRound:  3,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			state := test.state

			controller.PrepareRound(&state)

			assert.Equal(t, test.wantRemainingBlueUses, state.PlayerMoves.Moves[0].RemainingUses)
			assert.Equal(t, test.wantBlueRestoreRound, state.PlayerMoves.Moves[0].RestoresOnRound)
		})
	}
}

func TestUsageControllerValidateMove(t *testing.T) {
	controller, err := NewUsageController(DefaultConfig())
	require.NoError(t, err)

	tests := []struct {
		title   string
		state   match.State
		move    domain.Move
		wantErr error
	}{
		{
			title: "returns nil when the move still has uses remaining",
			state: match.State{PlayerMoves: match.PlayerMoveResources{Moves: []match.PlayerMoveState{
				{Move: domain.Blue, MaxUses: 3, RemainingUses: 1},
				{Move: domain.Red, MaxUses: 3, RemainingUses: 3},
				{Move: domain.Yellow, MaxUses: 3, RemainingUses: 3},
			}}},
			move: domain.Blue,
		},
		{
			title: "returns an unavailable move error when the move is recharging",
			state: match.State{PlayerMoves: match.PlayerMoveResources{Moves: []match.PlayerMoveState{
				{Move: domain.Blue, MaxUses: 3, RemainingUses: 0, RestoresOnRound: 4},
				{Move: domain.Red, MaxUses: 3, RemainingUses: 3},
				{Move: domain.Yellow, MaxUses: 3, RemainingUses: 3},
			}}},
			move:    domain.Blue,
			wantErr: ErrMoveUnavailable,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			err := controller.ValidateMove(test.state, test.move)

			if test.wantErr == nil {
				assert.NoError(t, err)
				return
			}

			assert.ErrorIs(t, err, test.wantErr)
		})
	}
}

func TestUsageControllerConsumeMove(t *testing.T) {
	controller, err := NewUsageController(Config{
		InitialUsesPerMove: map[domain.Move]int{
			domain.Blue:   2,
			domain.Red:    3,
			domain.Yellow: 3,
		},
		RecoveryDelayRounds: 1,
	})
	require.NoError(t, err)

	tests := []struct {
		title                 string
		state                 match.State
		move                  domain.Move
		wantRemainingBlueUses int
		wantBlueRestoreRound  int
	}{
		{
			title: "decrements remaining uses when the move still has charges left",
			state: match.State{
				Round: 1,
				PlayerMoves: match.PlayerMoveResources{Moves: []match.PlayerMoveState{
					{Move: domain.Blue, MaxUses: 2, RemainingUses: 2},
					{Move: domain.Red, MaxUses: 3, RemainingUses: 3},
					{Move: domain.Yellow, MaxUses: 3, RemainingUses: 3},
				}},
			},
			move:                  domain.Blue,
			wantRemainingBlueUses: 1,
			wantBlueRestoreRound:  0,
		},
		{
			title: "schedules a recharge round when the final use is consumed",
			state: match.State{
				Round: 3,
				PlayerMoves: match.PlayerMoveResources{Moves: []match.PlayerMoveState{
					{Move: domain.Blue, MaxUses: 2, RemainingUses: 1},
					{Move: domain.Red, MaxUses: 3, RemainingUses: 3},
					{Move: domain.Yellow, MaxUses: 3, RemainingUses: 3},
				}},
			},
			move:                  domain.Blue,
			wantRemainingBlueUses: 0,
			wantBlueRestoreRound:  5,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			state := test.state

			controller.ConsumeMove(&state, test.move)

			assert.Equal(t, test.wantRemainingBlueUses, state.PlayerMoves.Moves[0].RemainingUses)
			assert.Equal(t, test.wantBlueRestoreRound, state.PlayerMoves.Moves[0].RestoresOnRound)
		})
	}
}
