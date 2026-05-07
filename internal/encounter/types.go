package encounter

import (
	"fmt"
	"time"
)

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

type SplashProfile struct {
	JumpCount int
	TimeLimit time.Duration
}

func DefaultSplashProfile() SplashProfile {
	return SplashProfile{
		JumpCount: 1,
		TimeLimit: time.Second,
	}
}

func (profile SplashProfile) WithDefaults() SplashProfile {
	resolved := profile
	if resolved.JumpCount == 0 {
		resolved.JumpCount = DefaultSplashProfile().JumpCount
	}
	if resolved.TimeLimit == 0 {
		resolved.TimeLimit = DefaultSplashProfile().TimeLimit
	}

	return resolved
}

func (profile SplashProfile) Validate() error {
	resolved := profile.WithDefaults()
	if resolved.JumpCount < 1 || resolved.JumpCount > 5 {
		return fmt.Errorf("splash jump count must be between 1 and 5")
	}
	if resolved.TimeLimit <= 0 {
		return fmt.Errorf("splash time limit must be greater than 0")
	}

	return nil
}

type SplashState struct {
	TotalJumps    int
	ResolvedJumps int
	TimeLimit     time.Duration
}

func (state SplashState) Pending() bool {
	return state.ResolvedJumps < state.TotalJumps
}

func (state SplashState) CurrentJump() int {
	if !state.Pending() {
		return state.TotalJumps
	}

	return state.ResolvedJumps + 1
}

type SplashResolution struct {
	Escaped               bool
	SuccessfulJumps       int
	DistanceRewardApplied int
}

type Config struct {
	InitialDistance           int
	InitialDepth              int
	SurfaceDepth              int
	CaptureDistance           int
	ExhaustionCaptureDistance int
	PlayerWinStep             int
	FishWinStep               int
	SplashProfile             SplashProfile
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
		SplashProfile:             DefaultSplashProfile(),
	}
}

func (c Config) Validate() error {
	if c.InitialDistance < 0 {
		return fmt.Errorf("initial distance must be greater than or equal to 0")
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
	return c.SplashProfile.Validate()
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
	Splash    *SplashState
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
