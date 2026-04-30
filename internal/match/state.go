package match

import (
	"pesca/internal/cards"
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
	State      State
}
