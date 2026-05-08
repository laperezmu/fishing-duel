package presentation

import (
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/run"
	"time"
)

type AnglerProfileView struct {
	ProfileID         string
	Name              string
	Description       string
	Details           []string
	StartingThread    int
	DeckLabel         string
	RodLabel          string
	AttachmentLabel   string
	UnlockedByDefault bool
}

type RunIntroView struct {
	Title       string
	RouteLabels []string
	Thread      int
}

type RunNodeView struct {
	Title        string
	ZoneLabel    string
	NodeLabel    string
	NodeKind     run.NodeKind
	Thread       int
	ThreadMax    int
	CaptureCount int
}

type RunSummaryView struct {
	Title         string
	Status        run.Status
	StatusLabel   string
	Thread        int
	ThreadMax     int
	CaptureCount  int
	LastNodeLabel string
}

type RunNodeSummaryView struct {
	Title               string
	NodeLabel           string
	NodeKind            run.NodeKind
	OutcomeLabel        string
	Thread              int
	ThreadMax           int
	ThreadDelta         int
	CaptureCount        int
	LastCaptureLabel    string
	NextNodeLabel       string
	ContinuePromptLabel string
}

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
	Title        string
	WaterLabel   string
	Position     int
	TotalSlots   int
	SectionWidth int
}

type SplashView struct {
	Title                 string
	EventLabel            string
	CurrentJump           int
	TotalJumps            int
	TimeLimit             time.Duration
	SuccessRewardDistance int
	TotalSlots            int
	TargetStart           int
	TargetWidth           int
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
	BeforeStatus StatusView
	AfterStatus  StatusView
	PlayerMove   domain.Move
	FishMove     domain.Move
	PlayerLabel  string
	FishLabel    string
	Outcome      domain.RoundOutcome
	OutcomeLabel string
	EventLabel   string
	Resolved     []string
	TraceSummary []string
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
