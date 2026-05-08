package match

import (
	"pesca/internal/cards"
	"pesca/internal/domain"
	"pesca/internal/encounter"
)

type TrackSnapshot struct {
	Distance                  int
	Depth                     int
	SurfaceDepth              int
	CaptureDistance           int
	ExhaustionCaptureDistance int
	MaxDistance               int
	MaxDepth                  int
}

type EncounterSummarySnapshot struct {
	Distance  int
	Depth     int
	Status    encounter.Status
	EndReason encounter.EndReason
}

type EncounterEventSnapshot struct {
	LastEvent encounter.Event
	Splash    *SplashSnapshot
}

type SplashSnapshot struct {
	TotalJumps    int
	ResolvedJumps int
	CurrentJump   int
	TimeLimit     int64
}

type MoveResourceSnapshot struct {
	Move            domain.Move
	MaxUses         int
	RemainingUses   int
	RestoresOnRound int
	HasTopCard      bool
	TopCard         cards.PlayerCard
}

type PlayerOptionsSnapshot struct {
	Moves []MoveResourceSnapshot
}

type FishDiscardSnapshot struct {
	ActiveCards        int
	DiscardCards       int
	RecycleCount       int
	ShufflesOnRecycle  bool
	CardsToRemove      int
	CurrentCycle       FishDiscardCycleState
	PreviousCycleStats []FishDiscardCycleSummaryState
}

type StatusSnapshot struct {
	RoundNumber int
	Track       TrackSnapshot
	FishDiscard FishDiscardSnapshot
	Player      PlayerOptionsSnapshot
	Stats       Stats
}

type RoundSnapshot struct {
	Status          StatusSnapshot
	Encounter       EncounterEventSnapshot
	PlayerMove      domain.Move
	FishMove        domain.Move
	Outcome         domain.RoundOutcome
	ResolvedEffects []ResolvedEffectState
}

type SummarySnapshot struct {
	TotalRounds int
	Encounter   EncounterSummarySnapshot
	Stats       Stats
}

func NewStatusSnapshot(state State) StatusSnapshot {
	return StatusSnapshot{
		RoundNumber: state.Round.Number + 1,
		Track: TrackSnapshot{
			Distance:                  state.Encounter.Distance,
			Depth:                     state.Encounter.Depth,
			SurfaceDepth:              state.Encounter.Config.SurfaceDepth,
			CaptureDistance:           state.Encounter.Config.CaptureDistance,
			ExhaustionCaptureDistance: state.Encounter.Config.ExhaustionCaptureDistance,
			MaxDistance:               state.Player.Loadout.TrackMaxDistance(),
			MaxDepth:                  state.Player.Loadout.TrackMaxDepth(),
		},
		FishDiscard: FishDiscardSnapshot{
			ActiveCards:        state.Deck.ActiveCards,
			DiscardCards:       state.Deck.DiscardCards,
			RecycleCount:       state.Deck.RecycleCount,
			ShufflesOnRecycle:  state.Deck.ShufflesOnRecycle,
			CardsToRemove:      state.Deck.CardsToRemove,
			CurrentCycle:       cloneFishDiscardCycleState(state.Deck.CurrentCycle),
			PreviousCycleStats: append([]FishDiscardCycleSummaryState(nil), state.Deck.PreviousCycleStats...),
		},
		Player: PlayerOptionsSnapshot{
			Moves: buildMoveResourceSnapshots(state.Player.Moves.Moves),
		},
		Stats: state.Lifecycle.Stats,
	}
}

func NewRoundSnapshot(result RoundResult) RoundSnapshot {
	return RoundSnapshot{
		Status:          result.Status,
		Encounter:       result.Encounter,
		PlayerMove:      result.PlayerMove,
		FishMove:        result.FishMove,
		Outcome:         result.Outcome,
		ResolvedEffects: append([]ResolvedEffectState(nil), result.ResolvedEffects...),
	}
}

func NewEncounterEventSnapshot(state encounter.State) EncounterEventSnapshot {
	return EncounterEventSnapshot{
		LastEvent: state.LastEvent,
		Splash:    newSplashSnapshot(state.Splash),
	}
}

func newSplashSnapshot(state *encounter.SplashState) *SplashSnapshot {
	if state == nil {
		return nil
	}

	return &SplashSnapshot{
		TotalJumps:    state.TotalJumps,
		ResolvedJumps: state.ResolvedJumps,
		CurrentJump:   state.CurrentJump(),
		TimeLimit:     int64(state.TimeLimit),
	}
}

func NewSummarySnapshot(state State) SummarySnapshot {
	return SummarySnapshot{
		TotalRounds: state.Round.Number,
		Encounter: EncounterSummarySnapshot{
			Distance:  state.Encounter.Distance,
			Depth:     state.Encounter.Depth,
			Status:    state.Encounter.Status,
			EndReason: state.Encounter.EndReason,
		},
		Stats: state.Lifecycle.Stats,
	}
}

func buildMoveResourceSnapshots(moves []PlayerMoveState) []MoveResourceSnapshot {
	clonedMoves := make([]MoveResourceSnapshot, 0, len(moves))
	for _, moveState := range moves {
		clonedState := MoveResourceSnapshot{
			Move:            moveState.Move,
			MaxUses:         moveState.MaxUses,
			RemainingUses:   moveState.RemainingUses,
			RestoresOnRound: moveState.RestoresOnRound,
		}
		if len(moveState.ActiveCards) > 0 {
			clonedState.HasTopCard = true
			clonedState.TopCard = moveState.ActiveCards[0]
		}
		clonedMoves = append(clonedMoves, clonedState)
	}

	return clonedMoves
}

func cloneFishDiscardCycleState(state FishDiscardCycleState) FishDiscardCycleState {
	clonedState := state
	clonedState.Entries = append([]FishDiscardEntryState(nil), state.Entries...)

	return clonedState
}
