package presentation

import "pesca/internal/domain"

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
	Distance                  int
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
	Status      StatusView
	PlayerMove  domain.Move
	FishMove    domain.Move
	PlayerLabel string
	FishLabel   string
	Outcome     string
}

type SummaryView struct {
	TotalRounds int
	Distance    int
	Outcome     string
	EndReason   string
	PlayerWins  int
	FishWins    int
	Draws       int
}
