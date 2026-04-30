package game

import (
	"errors"
	"fmt"
	"pesca/internal/cards"
	"pesca/internal/deck"
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
	VisibilitySnapshot() deck.VisibilitySnapshot
}

type PlayerMoveController interface {
	Initialize(state *match.State)
	PrepareRound(state *match.State)
	ValidateMove(state match.State, playerMove domain.Move) error
	PeekMoveCard(state match.State, playerMove domain.Move) (cards.PlayerCard, error)
	ConsumeMove(state *match.State, playerMove domain.Move) cards.PlayerCard
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
	if engine.state.Lifecycle.Finished {
		return match.RoundResult{}, ErrGameFinished
	}

	engine.resetRoundState()
	engine.playerMoves.PrepareRound(&engine.state)
	if err := engine.playerMoves.ValidateMove(engine.state, playerMove); err != nil {
		return match.RoundResult{}, err
	}
	playerCard, err := engine.playerMoves.PeekMoveCard(engine.state, playerMove)
	if err != nil {
		return match.RoundResult{}, err
	}

	fishCard, err := engine.fishDeck.Draw()
	if err != nil {
		engine.refreshState()
		engine.endCondition.Apply(&engine.state)
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
	applyRoundScopedEffects(&engine.state, drawEffects)

	roundOutcome := engine.roundEvaluator.Evaluate(playerMove, fishCard.Move)
	engine.playerMoves.ConsumeMove(&engine.state, playerMove)
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
	applyRoundScopedEffects(&engine.state, outcomeEffects)
	engine.progressionPolicy.Apply(&engine.state, match.ResolvedRound{
		PlayerMove:     playerMove,
		PlayerCard:     playerCard,
		FishCard:       fishCard,
		DrawEffects:    append([]cards.CardEffect(nil), drawEffects...),
		OutcomeEffects: append([]cards.CardEffect(nil), outcomeEffects...),
		Outcome:        roundOutcome,
	})

	engine.fishDeck.PrepareNextRound()
	engine.refreshState()
	engine.playerMoves.PrepareRound(&engine.state)
	engine.endCondition.Apply(&engine.state)
	engine.resetRoundState()

	return match.RoundResult{
		Round:      engine.state.Round.Number,
		PlayerMove: playerMove,
		PlayerCard: playerCard,
		FishMove:   fishCard.Move,
		FishCard:   fishCard,
		Outcome:    roundOutcome,
		State:      engine.state,
	}, nil
}

func (engine *Engine) refreshState() {
	visibilitySnapshot := engine.fishDeck.VisibilitySnapshot()
	engine.state.Deck.ActiveCards = engine.fishDeck.ActiveCount()
	engine.state.Deck.DiscardCards = engine.fishDeck.DiscardCount()
	engine.state.Deck.RecycleCount = engine.fishDeck.RecycleCount()
	engine.state.Deck.Exhausted = engine.fishDeck.Exhausted()
	engine.state.Deck.ShufflesOnRecycle = visibilitySnapshot.ShufflesOnRecycle
	engine.state.Deck.CardsToRemove = visibilitySnapshot.CardsToRemove
	engine.state.Deck.CurrentCycle = mapVisibleDiscardCycleState(visibilitySnapshot.CurrentCycle)
	engine.state.Deck.PreviousCycleStats = mapVisibleDiscardCycleSummaryStates(visibilitySnapshot.PreviousCycles)
}

func (engine *Engine) resetRoundState() {
	engine.state.Round = match.RoundState{Number: engine.state.Round.Number}
}

func mapVisibleDiscardCycleState(cycle deck.VisibleDiscardCycle) match.FishDiscardCycleState {
	entries := make([]match.FishDiscardEntryState, 0, len(cycle.Entries))
	for _, entry := range cycle.Entries {
		entries = append(entries, match.FishDiscardEntryState{
			Visibility: entry.Visibility,
			Move:       entry.Move,
			Name:       entry.Name,
			Summary:    entry.Summary,
		})
	}

	return match.FishDiscardCycleState{
		Number:     cycle.Number,
		TotalCards: cycle.TotalCards,
		Entries:    entries,
	}
}

func mapVisibleDiscardCycleSummaryStates(summaries []deck.VisibleDiscardCycleSummary) []match.FishDiscardCycleSummaryState {
	mappedSummaries := make([]match.FishDiscardCycleSummaryState, 0, len(summaries))
	for _, summary := range summaries {
		mappedSummaries = append(mappedSummaries, match.FishDiscardCycleSummaryState{
			Number:       summary.Number,
			TotalCards:   summary.TotalCards,
			VisibleCards: summary.VisibleCards,
			HiddenCards:  summary.HiddenCards,
		})
	}

	return mappedSummaries
}
