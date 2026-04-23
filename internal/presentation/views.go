package presentation

import (
	"pesca/internal/domain"
	"pesca/internal/encounter"
)

type MoveOption struct {
	Index           int
	Move            domain.Move
	Label           string
	CardHint        string
	RemainingUses   int
	MaxUses         int
	Available       bool
	RestoresOnRound int
}

type IntroView struct {
	Title   string
	Options []MoveOption
}

type StatusView struct {
	RoundNumber               int
	FishDistance              int
	FishDepth                 int
	SurfaceDepth              int
	MaxDistance               int
	MaxDepth                  int
	CaptureDistance           int
	ExhaustionCaptureDistance int
	ActiveCards               int
	DiscardCards              int
	RecycleCount              int
	PlayerWins                int
	FishWins                  int
	Draws                     int
	MoveOptions               []MoveOption
}

type RoundView struct {
	Status       StatusView
	PlayerMove   domain.Move
	FishMove     domain.Move
	PlayerLabel  string
	FishLabel    string
	Outcome      domain.RoundOutcome
	OutcomeLabel string
	EventLabel   string
}

type SummaryView struct {
	TotalRounds     int
	FishDistance    int
	FishDepth       int
	EncounterStatus encounter.Status
	OutcomeLabel    string
	EndReasonLabel  string
	PlayerWins      int
	FishWins        int
	Draws           int
}
