package game

import (
	"errors"
	"fmt"

	"pesca/internal/deck"
	"pesca/internal/domain"
)

var ErrGameFinished = errors.New("game already finished")

type Evaluator interface {
	Evaluate(player, fish domain.Move) domain.RoundOutcome
}

type ProgressionPolicy interface {
	Apply(state *State, outcome domain.RoundOutcome)
}

type EndCondition interface {
	Apply(state *State)
}

type RoundResult struct {
	Round      int
	PlayerMove domain.Move
	FishMove   domain.Move
	Outcome    domain.RoundOutcome
	State      State
}

type Engine struct {
	deck        *deck.Manager
	evaluator   Evaluator
	progression ProgressionPolicy
	ending      EndCondition
	state       State
}

func NewEngine(deckManager *deck.Manager, evaluator Evaluator, progression ProgressionPolicy, ending EndCondition, initialState State) (*Engine, error) {
	switch {
	case deckManager == nil:
		return nil, fmt.Errorf("deck manager is required")
	case evaluator == nil:
		return nil, fmt.Errorf("evaluator is required")
	case progression == nil:
		return nil, fmt.Errorf("progression policy is required")
	case ending == nil:
		return nil, fmt.Errorf("end condition is required")
	}

	engine := &Engine{
		deck:        deckManager,
		evaluator:   evaluator,
		progression: progression,
		ending:      ending,
		state:       initialState,
	}
	engine.deck.PrepareNextRound()
	engine.refreshState()
	engine.ending.Apply(&engine.state)

	return engine, nil
}

func (e *Engine) State() State {
	return e.state
}

func (e *Engine) PlayRound(playerMove domain.Move) (RoundResult, error) {
	if e.state.Finished {
		return RoundResult{}, ErrGameFinished
	}

	fishMove, err := e.deck.Draw()
	if err != nil {
		e.refreshState()
		e.ending.Apply(&e.state)
		if e.state.Finished {
			return RoundResult{}, ErrGameFinished
		}
		return RoundResult{}, err
	}

	outcome := e.evaluator.Evaluate(playerMove, fishMove)
	e.state.Round++
	e.progression.Apply(&e.state, outcome)

	e.deck.PrepareNextRound()
	e.refreshState()
	e.ending.Apply(&e.state)

	return RoundResult{
		Round:      e.state.Round,
		PlayerMove: playerMove,
		FishMove:   fishMove,
		Outcome:    outcome,
		State:      e.state,
	}, nil
}

func (e *Engine) refreshState() {
	e.state.Deck.ActiveCards = e.deck.ActiveCount()
	e.state.Deck.DiscardCards = e.deck.DiscardCount()
	e.state.Deck.RecycleCount = e.deck.RecycleCount()
	e.state.Deck.Exhausted = e.deck.Exhausted()
}
