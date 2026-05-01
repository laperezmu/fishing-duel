package presentation

import (
	"pesca/internal/domain"
	"pesca/internal/encounter"
)

type OpeningView struct {
	WaterLabel      string
	CastLabel       string
	InitialDistance int
	InitialDepth    int
}

type SpawnView struct {
	ProfileLabel    string
	WaterBaseLabel  string
	InitialDistance int
	InitialDepth    int
	CandidateCount  int
	HabitatLabels   []string
}

type CastView struct {
	WaterLabel   string
	Position     int
	TotalSlots   int
	SectionWidth int
}

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

type FishDiscardEntryView struct {
	Label string
}

type FishDiscardCycleSummaryView struct {
	CycleNumber  int
	TotalCards   int
	VisibleCards int
	HiddenCards  int
}

type FishDiscardView struct {
	CurrentCycleNumber     int
	CurrentCycleTotalCards int
	CurrentCycleEntries    []FishDiscardEntryView
	PreviousCycles         []FishDiscardCycleSummaryView
	ShufflesOnRecycle      bool
	CardsToRemove          int
	RecycleCount           int
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
	FishDiscard               FishDiscardView
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
