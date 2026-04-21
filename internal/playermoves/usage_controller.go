package playermoves

import (
	"errors"
	"fmt"
	"pesca/internal/domain"
	"pesca/internal/match"
)

var ErrMoveUnavailable = errors.New("player move unavailable")

type MoveUnavailableError struct {
	Move            domain.Move
	RestoresOnRound int
}

func (err MoveUnavailableError) Error() string {
	if err.RestoresOnRound > 0 {
		return fmt.Sprintf("player move %s is recharging until round %d", err.Move, err.RestoresOnRound)
	}

	return fmt.Sprintf("player move %s has no remaining uses", err.Move)
}

func (err MoveUnavailableError) Unwrap() error {
	return ErrMoveUnavailable
}

type UsageController struct {
	config Config
}

func NewUsageController(config Config) (*UsageController, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &UsageController{config: config}, nil
}

func (controller *UsageController) Initialize(state *match.State) {
	state.PlayerMoves = match.PlayerMoveResources{Moves: make([]match.PlayerMoveState, 0, len(supportedMoves()))}
	for _, move := range supportedMoves() {
		initialUses := controller.config.initialUsesFor(move)
		state.PlayerMoves.Moves = append(state.PlayerMoves.Moves, match.PlayerMoveState{
			Move:          move,
			MaxUses:       initialUses,
			RemainingUses: initialUses,
		})
	}

	controller.PrepareRound(state)
}

func (controller *UsageController) PrepareRound(state *match.State) {
	currentSelectionRound := state.Round + 1
	for moveIndex := range state.PlayerMoves.Moves {
		moveState := &state.PlayerMoves.Moves[moveIndex]
		if moveState.RestoresOnRound == 0 {
			continue
		}
		if currentSelectionRound < moveState.RestoresOnRound {
			continue
		}

		moveState.RemainingUses = moveState.MaxUses
		moveState.RestoresOnRound = 0
	}
}

func (controller *UsageController) ValidateMove(state match.State, playerMove domain.Move) error {
	moveState, ok := findMoveState(state.PlayerMoves, playerMove)
	if !ok {
		return MoveUnavailableError{Move: playerMove}
	}
	if moveState.RemainingUses > 0 {
		return nil
	}

	return MoveUnavailableError{
		Move:            playerMove,
		RestoresOnRound: moveState.RestoresOnRound,
	}
}

func (controller *UsageController) ConsumeMove(state *match.State, playerMove domain.Move) {
	moveState, ok := findMoveStatePointer(&state.PlayerMoves, playerMove)
	if !ok || moveState.RemainingUses == 0 {
		return
	}

	moveState.RemainingUses--
	if moveState.RemainingUses > 0 {
		return
	}

	moveState.RestoresOnRound = state.Round + controller.config.RecoveryDelayRounds + 1
}

func findMoveState(resources match.PlayerMoveResources, playerMove domain.Move) (match.PlayerMoveState, bool) {
	for _, moveState := range resources.Moves {
		if moveState.Move == playerMove {
			return moveState, true
		}
	}

	return match.PlayerMoveState{}, false
}

func findMoveStatePointer(resources *match.PlayerMoveResources, playerMove domain.Move) (*match.PlayerMoveState, bool) {
	for moveIndex := range resources.Moves {
		if resources.Moves[moveIndex].Move == playerMove {
			return &resources.Moves[moveIndex], true
		}
	}

	return nil, false
}
