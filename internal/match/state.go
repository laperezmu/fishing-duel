package match

import (
	"pesca/internal/domain"
	"pesca/internal/encounter"
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

type State struct {
	Round     int
	Deck      DeckState
	Encounter encounter.State
	Stats     Stats
	Finished  bool
}

type RoundResult struct {
	Round      int
	PlayerMove domain.Move
	FishMove   domain.Move
	Outcome    domain.RoundOutcome
	State      State
}
