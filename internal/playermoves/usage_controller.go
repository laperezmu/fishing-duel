package playermoves

import (
	"errors"
	"fmt"
	"pesca/internal/cards"
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
		moveState := match.PlayerMoveState{
			Move:        move,
			ActiveCards: controller.config.initialDeckFor(move),
		}
		syncMoveStateCounts(&moveState)
		state.PlayerMoves.Moves = append(state.PlayerMoves.Moves, moveState)
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

		moveState.ActiveCards = append([]cards.PlayerCard(nil), moveState.DiscardedCards...)
		moveState.DiscardedCards = nil
		if controller.config.DeckShuffler != nil {
			controller.config.DeckShuffler(moveState.ActiveCards)
		}
		moveState.RestoresOnRound = 0
		syncMoveStateCounts(moveState)
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

func (controller *UsageController) PeekMoveCard(state match.State, playerMove domain.Move) (cards.PlayerCard, error) {
	moveState, ok := findMoveState(state.PlayerMoves, playerMove)
	if !ok || len(moveState.ActiveCards) == 0 {
		return cards.PlayerCard{}, MoveUnavailableError{Move: playerMove, RestoresOnRound: moveState.RestoresOnRound}
	}

	return moveState.ActiveCards[0], nil
}

func (controller *UsageController) ConsumeMove(state *match.State, playerMove domain.Move) cards.PlayerCard {
	moveState, ok := findMoveStatePointer(&state.PlayerMoves, playerMove)
	if !ok || len(moveState.ActiveCards) == 0 {
		return cards.PlayerCard{}
	}

	selectedCard := moveState.ActiveCards[0]
	moveState.ActiveCards = append([]cards.PlayerCard(nil), moveState.ActiveCards[1:]...)
	moveState.DiscardedCards = append(moveState.DiscardedCards, selectedCard)
	syncMoveStateCounts(moveState)
	if moveState.RemainingUses > 0 {
		return selectedCard
	}

	moveState.RestoresOnRound = state.Round + controller.config.RecoveryDelayRounds + 1
	return selectedCard
}

func syncMoveStateCounts(moveState *match.PlayerMoveState) {
	moveState.RemainingUses = len(moveState.ActiveCards)
	moveState.MaxUses = len(moveState.ActiveCards) + len(moveState.DiscardedCards)
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
