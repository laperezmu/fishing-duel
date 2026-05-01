package match

import (
	"pesca/internal/cards"
	"pesca/internal/deck"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/player/loadout"
)

type DeckState struct {
	ActiveCards        int
	DiscardCards       int
	RecycleCount       int
	Exhausted          bool
	ShufflesOnRecycle  bool
	CardsToRemove      int
	CurrentCycle       FishDiscardCycleState
	PreviousCycleStats []FishDiscardCycleSummaryState
}

type FishDiscardEntryState struct {
	Visibility cards.DiscardVisibility
	Move       domain.Move
	Name       string
	Summary    string
}

type FishDiscardCycleState struct {
	Number     int
	TotalCards int
	Entries    []FishDiscardEntryState
}

type FishDiscardCycleSummaryState struct {
	Number       int
	TotalCards   int
	VisibleCards int
	HiddenCards  int
}

type Stats struct {
	PlayerWins int
	FishWins   int
	Draws      int
}

type PlayerMoveState struct {
	Move            domain.Move
	MaxUses         int
	RemainingUses   int
	RestoresOnRound int
	ActiveCards     []cards.PlayerCard
	DiscardedCards  []cards.PlayerCard
}

type PlayerMoveResources struct {
	Moves []PlayerMoveState
}

type RoundState struct {
	Number     int
	Thresholds encounter.RoundThresholds
}

type PlayerState struct {
	Loadout loadout.State
	Moves   PlayerMoveResources
}

type LifecycleState struct {
	Stats    Stats
	Finished bool
}

type State struct {
	Round     RoundState
	Deck      DeckState
	Encounter encounter.State
	Player    PlayerState
	Lifecycle LifecycleState
}

type ProgressionState struct {
	Round     *RoundState
	Encounter *encounter.State
	Lifecycle *LifecycleState
}

type EndingState struct {
	Round     *RoundState
	Deck      *DeckState
	Encounter *encounter.State
	Player    *PlayerState
	Lifecycle *LifecycleState
}

type ResolvedRound struct {
	PlayerMove     domain.Move
	PlayerCard     cards.PlayerCard
	FishCard       cards.FishCard
	DrawEffects    []cards.CardEffect
	OutcomeEffects []cards.CardEffect
	Outcome        domain.RoundOutcome
}

type RoundResult struct {
	Round      int
	PlayerMove domain.Move
	PlayerCard cards.PlayerCard
	FishMove   domain.Move
	FishCard   cards.FishCard
	Outcome    domain.RoundOutcome
	Status     StatusSnapshot
	Encounter  EncounterEventSnapshot
}

func (state *State) ProgressionState() ProgressionState {
	return ProgressionState{
		Round:     &state.Round,
		Encounter: &state.Encounter,
		Lifecycle: &state.Lifecycle,
	}
}

func (state *State) EndingState() EndingState {
	return EndingState{
		Round:     &state.Round,
		Deck:      &state.Deck,
		Encounter: &state.Encounter,
		Player:    &state.Player,
		Lifecycle: &state.Lifecycle,
	}
}

func NewDeckState(activeCards, discardCards, recycleCount int, exhausted bool, visibility deck.VisibilitySnapshot) DeckState {
	return DeckState{
		ActiveCards:        activeCards,
		DiscardCards:       discardCards,
		RecycleCount:       recycleCount,
		Exhausted:          exhausted,
		ShufflesOnRecycle:  visibility.ShufflesOnRecycle,
		CardsToRemove:      visibility.CardsToRemove,
		CurrentCycle:       mapVisibleDiscardCycleState(visibility.CurrentCycle),
		PreviousCycleStats: mapVisibleDiscardCycleSummaryStates(visibility.PreviousCycles),
	}
}

func mapVisibleDiscardCycleState(cycle deck.VisibleDiscardCycle) FishDiscardCycleState {
	entries := make([]FishDiscardEntryState, 0, len(cycle.Entries))
	for _, entry := range cycle.Entries {
		entries = append(entries, FishDiscardEntryState{
			Visibility: entry.Visibility,
			Move:       entry.Move,
			Name:       entry.Name,
			Summary:    entry.Summary,
		})
	}

	return FishDiscardCycleState{
		Number:     cycle.Number,
		TotalCards: cycle.TotalCards,
		Entries:    entries,
	}
}

func mapVisibleDiscardCycleSummaryStates(summaries []deck.VisibleDiscardCycleSummary) []FishDiscardCycleSummaryState {
	mappedSummaries := make([]FishDiscardCycleSummaryState, 0, len(summaries))
	for _, summary := range summaries {
		mappedSummaries = append(mappedSummaries, FishDiscardCycleSummaryState{
			Number:       summary.Number,
			TotalCards:   summary.TotalCards,
			VisibleCards: summary.VisibleCards,
			HiddenCards:  summary.HiddenCards,
		})
	}

	return mappedSummaries
}
