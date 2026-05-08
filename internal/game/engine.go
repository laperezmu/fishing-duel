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
	drawEffects := cards.OrderOwnedEffects(cards.FilterOwnedEffects(playerCard.Effects, cards.EffectContext{
		Owner:    cards.OwnerPlayer,
		Phase:    cards.PhaseDraw,
		CardMove: playerCard.Move,
	}))
	drawEffects = append(drawEffects, cards.FilterOwnedEffects(fishCard.Effects, cards.EffectContext{
		Owner:    cards.OwnerFish,
		Phase:    cards.PhaseDraw,
		CardMove: fishCard.Move,
	})...)
	orderedDrawEffects := cards.FlattenOwnedEffects(cards.OrderOwnedEffects(drawEffects))
	encounter.ApplyThresholdEffects(&engine.state.Round.Thresholds, orderedDrawEffects)

	roundOutcome := engine.roundEvaluator.Evaluate(playerMove, fishCard.Move)
	playerMoveRuntime = engine.state.PlayerMoveRuntime()
	engine.playerMoves.ConsumeMove(&playerMoveRuntime, playerMove)
	preProgressionEncounter := engine.state.Encounter
	preProgressionDeck := engine.state.Deck
	outcomeEffects := cards.FilterOwnedEffects(playerCard.Effects, cards.EffectContext{
		Owner:          cards.OwnerPlayer,
		Phase:          cards.PhaseOutcome,
		Outcome:        roundOutcome,
		CardMove:       playerCard.Move,
		TriggerMove:    playerMove,
		HasTriggerMove: true,
		FishDidSplash:  preProgressionEncounter.LastEvent.Kind == encounter.EventKindSplash || preProgressionEncounter.Splash != nil,
		FishReshuffled: preProgressionDeck.RecycleCount > 0,
		FishExhausted:  preProgressionDeck.Exhausted,
	})
	outcomeEffects = append(outcomeEffects, cards.FilterOwnedEffects(fishCard.Effects, cards.EffectContext{
		Owner:          cards.OwnerFish,
		Phase:          cards.PhaseOutcome,
		Outcome:        roundOutcome,
		CardMove:       fishCard.Move,
		TriggerMove:    fishCard.Move,
		HasTriggerMove: true,
		FishDidSplash:  preProgressionEncounter.LastEvent.Kind == encounter.EventKindSplash || preProgressionEncounter.Splash != nil,
		FishReshuffled: preProgressionDeck.RecycleCount > 0,
		FishExhausted:  preProgressionDeck.Exhausted,
	})...)
	orderedOutcomeEffects := cards.FlattenOwnedEffects(cards.OrderOwnedEffects(outcomeEffects))
	encounter.ApplyThresholdEffects(&engine.state.Round.Thresholds, orderedOutcomeEffects)
	orderedDrawOwned := cards.OrderOwnedEffects(drawEffects)
	orderedOutcomeOwned := cards.OrderOwnedEffects(outcomeEffects)
	progressionState := engine.state.ProgressionState()
	engine.progressionPolicy.Apply(&progressionState, match.ResolvedRound{
		PlayerMove:     playerMove,
		PlayerCard:     playerCard,
		FishCard:       fishCard,
		DrawOwned:      append([]cards.OwnedEffect(nil), orderedDrawOwned...),
		OutcomeOwned:   append([]cards.OwnedEffect(nil), orderedOutcomeOwned...),
		DrawEffects:    append([]cards.CardEffect(nil), orderedDrawEffects...),
		OutcomeEffects: append([]cards.CardEffect(nil), orderedOutcomeEffects...),
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
		Round:           engine.state.Round.Number,
		PlayerMove:      playerMove,
		PlayerCard:      playerCard,
		FishMove:        fishCard.Move,
		FishCard:        fishCard,
		Outcome:         roundOutcome,
		ResolvedEffects: append(buildResolvedEffectStates(orderedDrawOwned), buildResolvedEffectStates(orderedOutcomeOwned)...),
		Status:          match.NewStatusSnapshot(engine.state),
		Encounter:       match.NewEncounterEventSnapshot(engine.state.Encounter),
	}, nil
}

func buildResolvedEffectStates(effects []cards.OwnedEffect) []match.ResolvedEffectState {
	resolved := make([]match.ResolvedEffectState, 0, len(effects))
	for _, effect := range effects {
		normalized := effect.Effect.Normalize()
		resolved = append(resolved, match.ResolvedEffectState{
			Owner:    effect.Owner,
			Trigger:  normalized.Trigger,
			Type:     normalized.Type,
			Priority: normalized.Priority,
		})
	}

	return resolved
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
