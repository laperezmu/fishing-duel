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
	EndReasonDeckCapture  EndReason = "deck_capture"
	EndReasonDeckEscape   EndReason = "deck_escape"
)

type Config struct {
	InitialDistance           int
	CaptureDistance           int
	EscapeDistance            int
	ExhaustionCaptureDistance int
	PlayerWinStep             int
	FishWinStep               int
}

func DefaultConfig() Config {
	return Config{
		InitialDistance:           3,
		CaptureDistance:           0,
		EscapeDistance:            5,
		ExhaustionCaptureDistance: 2,
		PlayerWinStep:             1,
		FishWinStep:               1,
	}
}

func (c Config) Validate() error {
	if c.InitialDistance <= 2 {
		return fmt.Errorf("initial distance must be greater than 2")
	}
	if c.EscapeDistance < c.InitialDistance {
		return fmt.Errorf("escape distance must be at least the initial distance")
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

	return nil
}

type State struct {
	Config    Config
	Distance  int
	Status    Status
	EndReason EndReason
}

func NewState(config Config) (State, error) {
	if err := config.Validate(); err != nil {
		return State{}, err
	}

	return State{
		Config:    config,
		Distance:  config.InitialDistance,
		Status:    StatusOngoing,
		EndReason: EndReasonNone,
	}, nil
}
