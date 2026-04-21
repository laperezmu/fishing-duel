package game

import (
	"errors"
	"fmt"
	"pesca/internal/deck"
	"pesca/internal/domain"
	"pesca/internal/match"
	"pesca/internal/playermoves"
)

var ErrGameFinished = errors.New("game already finished")

type RoundEvaluator interface {
	Evaluate(playerMove, fishMove domain.Move) domain.RoundOutcome
}

type MatchProgressionPolicy interface {
	Apply(state *match.State, outcome domain.RoundOutcome)
}

type MatchEndCondition interface {
	Apply(state *match.State)
}

type Engine struct {
	fishDeck          *deck.Deck
	playerMoves       *playermoves.UsageController
	roundEvaluator    RoundEvaluator
	progressionPolicy MatchProgressionPolicy
	endCondition      MatchEndCondition
	state             match.State
}

func NewEngine(fishDeck *deck.Deck, playerMoves *playermoves.UsageController, roundEvaluator RoundEvaluator, progressionPolicy MatchProgressionPolicy, endCondition MatchEndCondition, initialState match.State) (*Engine, error) {
	switch {
	case fishDeck == nil:
		return nil, fmt.Errorf("fish deck is required")
	case playerMoves == nil:
		return nil, fmt.Errorf("player moves are required")
	case roundEvaluator == nil:
		return nil, fmt.Errorf("round evaluator is required")
	case progressionPolicy == nil:
		return nil, fmt.Errorf("progression policy is required")
	case endCondition == nil:
		return nil, fmt.Errorf("end condition is required")
	}

	engine := &Engine{
		fishDeck:          fishDeck,
		playerMoves:       playerMoves,
		roundEvaluator:    roundEvaluator,
		progressionPolicy: progressionPolicy,
		endCondition:      endCondition,
		state:             initialState,
	}
	engine.playerMoves.Initialize(&engine.state)
	engine.fishDeck.PrepareNextRound()
	engine.refreshState()
	engine.endCondition.Apply(&engine.state)

	return engine, nil
}

func (engine *Engine) State() match.State {
	return engine.state
}

func (engine *Engine) PlayRound(playerMove domain.Move) (match.RoundResult, error) {
	if engine.state.Finished {
		return match.RoundResult{}, ErrGameFinished
	}

	engine.playerMoves.PrepareRound(&engine.state)
	if err := engine.playerMoves.ValidateMove(engine.state, playerMove); err != nil {
		return match.RoundResult{}, err
	}

	fishMove, err := engine.fishDeck.Draw()
	if err != nil {
		engine.refreshState()
		engine.endCondition.Apply(&engine.state)
		if engine.state.Finished {
			return match.RoundResult{}, ErrGameFinished
		}
		return match.RoundResult{}, err
	}

	roundOutcome := engine.roundEvaluator.Evaluate(playerMove, fishMove)
	engine.state.Round++
	engine.playerMoves.ConsumeMove(&engine.state, playerMove)
	engine.progressionPolicy.Apply(&engine.state, roundOutcome)

	engine.fishDeck.PrepareNextRound()
	engine.refreshState()
	engine.playerMoves.PrepareRound(&engine.state)
	engine.endCondition.Apply(&engine.state)

	return match.RoundResult{
		Round:      engine.state.Round,
		PlayerMove: playerMove,
		FishMove:   fishMove,
		Outcome:    roundOutcome,
		State:      engine.state,
	}, nil
}

func (engine *Engine) refreshState() {
	engine.state.Deck.ActiveCards = engine.fishDeck.ActiveCount()
	engine.state.Deck.DiscardCards = engine.fishDeck.DiscardCount()
	engine.state.Deck.RecycleCount = engine.fishDeck.RecycleCount()
	engine.state.Deck.Exhausted = engine.fishDeck.Exhausted()
}
