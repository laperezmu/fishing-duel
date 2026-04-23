package match

import (
	"pesca/internal/cards"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/playerrig"
)

type DeckState struct {
	ActiveCards  int
	DiscardCards int
	RecycleCount int
	Exhausted    bool
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

type RoundThresholdState struct {
	CaptureDistanceBonus           int
	ExhaustionCaptureDistanceBonus int
	SurfaceDepthBonus              int
}

type RoundState struct {
	Thresholds RoundThresholdState
}

type State struct {
	Round       int
	RoundState  RoundState
	Deck        DeckState
	Encounter   encounter.State
	PlayerRig   playerrig.State
	PlayerMoves PlayerMoveResources
	Stats       Stats
	Finished    bool
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
