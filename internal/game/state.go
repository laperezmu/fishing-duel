package game

import "pesca/internal/encounter"

type DeckState struct {
	ActiveCards  int
	DiscardCards int
	RecycleCount int
	Exhausted    bool
}

type RoundStats struct {
	PlayerWins int
	FishWins   int
	Draws      int
}

type State struct {
	Round     int
	Deck      DeckState
	Encounter encounter.State
	Stats     RoundStats
	Finished  bool
}
