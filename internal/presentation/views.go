package presentation

import (
	"pesca/internal/domain"
	"pesca/internal/encounter"
)

type MoveOption struct {
	Index int
	Move  domain.Move
	Label string
}

type IntroView struct {
	Title   string
	Options []MoveOption
}

type StatusView struct {
	RoundNumber               int
	FishDistance              int
	CaptureDistance           int
	EscapeDistance            int
	ExhaustionCaptureDistance int
	ActiveCards               int
	DiscardCards              int
	RecycleCount              int
	PlayerWins                int
	FishWins                  int
	Draws                     int
}

type RoundView struct {
	Status       StatusView
	PlayerMove   domain.Move
	FishMove     domain.Move
	PlayerLabel  string
	FishLabel    string
	Outcome      domain.RoundOutcome
	OutcomeLabel string
}

type SummaryView struct {
	TotalRounds     int
	FishDistance    int
	EncounterStatus encounter.Status
	OutcomeLabel    string
	EndReasonLabel  string
	PlayerWins      int
	FishWins        int
	Draws           int
}
