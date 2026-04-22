package encounter

import "fmt"

type Status string

const (
	StatusOngoing  Status = "ongoing"
	StatusCaptured Status = "captured"
	StatusEscaped  Status = "escaped"
)

type EndReason string

const (
	EndReasonNone         EndReason = ""
	EndReasonTrackCapture EndReason = "track_capture"
	EndReasonTrackEscape  EndReason = "track_escape"
	EndReasonDepthEscape  EndReason = "depth_escape"
	EndReasonSplashEscape EndReason = "splash_escape"
	EndReasonDeckCapture  EndReason = "deck_capture"
	EndReasonDeckEscape   EndReason = "deck_escape"
)

type EventKind string

const (
	EventKindNone   EventKind = ""
	EventKindSplash EventKind = "splash"
)

type Event struct {
	Kind    EventKind
	Escaped bool
}

type Config struct {
	InitialDistance           int
	InitialDepth              int
	SurfaceDepth              int
	CaptureDistance           int
	ExhaustionCaptureDistance int
	PlayerWinStep             int
	FishWinStep               int
	SplashEscapeChance        float64
}

func DefaultConfig() Config {
	return Config{
		InitialDistance:           3,
		InitialDepth:              1,
		SurfaceDepth:              0,
		CaptureDistance:           0,
		ExhaustionCaptureDistance: 2,
		PlayerWinStep:             1,
		FishWinStep:               1,
		SplashEscapeChance:        0.5,
	}
}

func (c Config) Validate() error {
	if c.InitialDistance <= 2 {
		return fmt.Errorf("initial distance must be greater than 2")
	}
	if c.InitialDepth < c.SurfaceDepth {
		return fmt.Errorf("initial depth must be greater than or equal to surface depth")
	}
	if c.PlayerWinStep <= 0 {
		return fmt.Errorf("player win step must be greater than 0")
	}
	if c.FishWinStep <= 0 {
		return fmt.Errorf("fish win step must be greater than 0")
	}
	if c.ExhaustionCaptureDistance < c.CaptureDistance {
		return fmt.Errorf("exhaustion capture distance must be at least the capture distance")
	}
	if c.SplashEscapeChance < 0 || c.SplashEscapeChance > 1 {
		return fmt.Errorf("splash escape chance must be between 0 and 1")
	}

	return nil
}

type Delta struct {
	DistanceShift int
	DepthShift    int
}

type State struct {
	Config    Config
	Distance  int
	Depth     int
	Status    Status
	EndReason EndReason
	LastEvent Event
}

func NewState(config Config) (State, error) {
	if err := config.Validate(); err != nil {
		return State{}, err
	}

	return State{
		Config:    config,
		Distance:  config.InitialDistance,
		Depth:     config.InitialDepth,
		Status:    StatusOngoing,
		EndReason: EndReasonNone,
	}, nil
}
