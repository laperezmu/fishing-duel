package game

import (
	"errors"
	"fmt"
	"pesca/internal/cards"
	"pesca/internal/domain"
	"pesca/internal/match"
)

var ErrGameFinished = errors.New("game already finished")

type RoundEvaluator interface {
	Evaluate(playerMove, fishMove domain.Move) domain.RoundOutcome
}

type MatchProgressionPolicy interface {
	Apply(state *match.State, round match.ResolvedRound)
}

type MatchEndCondition interface {
	Apply(state *match.State)
}

type FishDeck interface {
	Draw() (cards.FishCard, error)
	PrepareNextRound()
	ActiveCount() int
	DiscardCount() int
	RecycleCount() int
	Exhausted() bool
}

type PlayerMoveController interface {
	Initialize(state *match.State)
	PrepareRound(state *match.State)
	ValidateMove(state match.State, playerMove domain.Move) error
	ConsumeMove(state *match.State, playerMove domain.Move)
}

type Engine struct {
	fishDeck          FishDeck
	playerMoves       PlayerMoveController
	roundEvaluator    RoundEvaluator
	progressionPolicy MatchProgressionPolicy
	endCondition      MatchEndCondition
	state             match.State
}

func NewEngine(fishDeck FishDeck, playerMoves PlayerMoveController, roundEvaluator RoundEvaluator, progressionPolicy MatchProgressionPolicy, endCondition MatchEndCondition, initialState match.State) (*Engine, error) {
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

	engine.resetRoundState()
	engine.playerMoves.PrepareRound(&engine.state)
	if err := engine.playerMoves.ValidateMove(engine.state, playerMove); err != nil {
		return match.RoundResult{}, err
	}

	fishCard, err := engine.fishDeck.Draw()
	if err != nil {
		engine.refreshState()
		engine.endCondition.Apply(&engine.state)
		if engine.state.Finished {
			return match.RoundResult{}, ErrGameFinished
		}
		return match.RoundResult{}, err
	}

	roundOutcome := engine.roundEvaluator.Evaluate(playerMove, fishCard.Move)
	engine.state.Round++
	engine.playerMoves.ConsumeMove(&engine.state, playerMove)
	cardEffects := cards.FilterEffects(fishCard.Effects, cards.EffectContext{
		Owner:   cards.OwnerFish,
		Outcome: roundOutcome,
	})
	engine.progressionPolicy.Apply(&engine.state, match.ResolvedRound{
		PlayerMove:  playerMove,
		FishCard:    fishCard,
		CardEffects: append([]cards.CardEffect(nil), cardEffects...),
		Outcome:     roundOutcome,
	})

	engine.fishDeck.PrepareNextRound()
	engine.refreshState()
	engine.playerMoves.PrepareRound(&engine.state)
	engine.endCondition.Apply(&engine.state)
	engine.resetRoundState()

	return match.RoundResult{
		Round:      engine.state.Round,
		PlayerMove: playerMove,
		FishMove:   fishCard.Move,
		FishCard:   fishCard,
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

func (engine *Engine) resetRoundState() {
	engine.state.RoundState = match.RoundState{}
}
