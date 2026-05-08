package game

import (
	"errors"
	"fmt"
	"pesca/internal/cards"
	"pesca/internal/deck"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/match"
)

var ErrGameFinished = errors.New("game already finished")

type RoundEvaluator interface {
	Evaluate(playerMove, fishMove domain.Move) domain.RoundOutcome
}

type MatchProgressionPolicy interface {
	Apply(state *match.ProgressionState, round match.ResolvedRound)
}

type MatchEndCondition interface {
	Apply(state *match.EndingState)
}

type FishDeck interface {
	Draw() (cards.FishCard, error)
	PrepareNextRound()
	ActiveCount() int
	DiscardCount() int
	RecycleCount() int
	Exhausted() bool
	VisibilitySnapshot() deck.VisibilitySnapshot
}

type PlayerMoveController interface {
	Initialize(state *match.PlayerMoveRuntime)
	PrepareRound(state *match.PlayerMoveRuntime)
	ValidateMove(state match.PlayerMoveRuntime, playerMove domain.Move) error
	PeekMoveCard(state match.PlayerMoveRuntime, playerMove domain.Move) (cards.PlayerCard, error)
	ConsumeMove(state *match.PlayerMoveRuntime, playerMove domain.Move) cards.PlayerCard
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
	playerMoveRuntime := engine.state.PlayerMoveRuntime()
	engine.playerMoves.Initialize(&playerMoveRuntime)
	engine.fishDeck.PrepareNextRound()
	engine.refreshState()
	endingState := engine.state.EndingState()
	engine.endCondition.Apply(&endingState)

	return engine, nil
}

func (engine *Engine) State() match.State {
	return engine.state
}

func (engine *Engine) PlayRound(playerMove domain.Move) (match.RoundResult, error) {
	if engine.state.Lifecycle.Finished {
		return match.RoundResult{}, ErrGameFinished
	}

	engine.resetRoundState()
	playerMoveRuntime := engine.state.PlayerMoveRuntime()
	engine.playerMoves.PrepareRound(&playerMoveRuntime)
	if err := engine.playerMoves.ValidateMove(playerMoveRuntime, playerMove); err != nil {
		return match.RoundResult{}, err
	}
	playerCard, err := engine.playerMoves.PeekMoveCard(playerMoveRuntime, playerMove)
	if err != nil {
		return match.RoundResult{}, err
	}

	fishCard, err := engine.fishDeck.Draw()
	if err != nil {
		engine.refreshState()
		endingState := engine.state.EndingState()
		engine.endCondition.Apply(&endingState)
		if engine.state.Lifecycle.Finished {
			return match.RoundResult{}, ErrGameFinished
		}
		return match.RoundResult{}, err
	}

	engine.state.Round.Number++
	drawEffects := cards.FilterEffects(playerCard.Effects, cards.EffectContext{
		Owner: cards.OwnerPlayer,
		Phase: cards.PhaseDraw,
	})
	drawEffects = append(drawEffects, cards.FilterEffects(fishCard.Effects, cards.EffectContext{
		Owner: cards.OwnerFish,
		Phase: cards.PhaseDraw,
	})...)
	encounter.ApplyThresholdEffects(&engine.state.Round.Thresholds, drawEffects)

	roundOutcome := engine.roundEvaluator.Evaluate(playerMove, fishCard.Move)
	playerMoveRuntime = engine.state.PlayerMoveRuntime()
	engine.playerMoves.ConsumeMove(&playerMoveRuntime, playerMove)
	outcomeEffects := cards.FilterEffects(playerCard.Effects, cards.EffectContext{
		Owner:   cards.OwnerPlayer,
		Phase:   cards.PhaseOutcome,
		Outcome: roundOutcome,
	})
	outcomeEffects = append(outcomeEffects, cards.FilterEffects(fishCard.Effects, cards.EffectContext{
		Owner:   cards.OwnerFish,
		Phase:   cards.PhaseOutcome,
		Outcome: roundOutcome,
	})...)
	encounter.ApplyThresholdEffects(&engine.state.Round.Thresholds, outcomeEffects)
	progressionState := engine.state.ProgressionState()
	engine.progressionPolicy.Apply(&progressionState, match.ResolvedRound{
		PlayerMove:     playerMove,
		PlayerCard:     playerCard,
		FishCard:       fishCard,
		DrawEffects:    append([]cards.CardEffect(nil), drawEffects...),
		OutcomeEffects: append([]cards.CardEffect(nil), outcomeEffects...),
		Outcome:        roundOutcome,
	})

	engine.fishDeck.PrepareNextRound()
	engine.refreshState()
	playerMoveRuntime = engine.state.PlayerMoveRuntime()
	engine.playerMoves.PrepareRound(&playerMoveRuntime)
	endingState := engine.state.EndingState()
	engine.endCondition.Apply(&endingState)
	if engine.state.Encounter.Splash == nil {
		engine.resetRoundState()
	}

	return match.RoundResult{
		Round:      engine.state.Round.Number,
		PlayerMove: playerMove,
		PlayerCard: playerCard,
		FishMove:   fishCard.Move,
		FishCard:   fishCard,
		Outcome:    roundOutcome,
		Status:     match.NewStatusSnapshot(engine.state),
		Encounter:  match.NewEncounterEventSnapshot(engine.state.Encounter),
	}, nil
}

func (engine *Engine) ResolveSplash(resolution encounter.SplashResolution) error {
	if engine.state.Lifecycle.Finished {
		return ErrGameFinished
	}

	encounter.ApplySplashResolution(&engine.state.Encounter, resolution)
	engine.refreshState()
	endingState := engine.state.EndingState()
	engine.endCondition.Apply(&endingState)
	if engine.state.Lifecycle.Finished {
		return nil
	}
	if engine.state.Encounter.Splash == nil {
		engine.resetRoundState()
	}

	return nil
}

func (engine *Engine) refreshState() {
	visibilitySnapshot := engine.fishDeck.VisibilitySnapshot()
	engine.state.Deck = match.NewDeckState(
		engine.fishDeck.ActiveCount(),
		engine.fishDeck.DiscardCount(),
		engine.fishDeck.RecycleCount(),
		engine.fishDeck.Exhausted(),
		visibilitySnapshot,
	)
}

func (engine *Engine) resetRoundState() {
	engine.state.Round = match.RoundState{Number: engine.state.Round.Number}
}
