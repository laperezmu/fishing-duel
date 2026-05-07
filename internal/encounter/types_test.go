package encounter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewState(t *testing.T) {
	validCases := []struct {
		title     string
		config    Config
		wantState State
	}{
		{
			title:  "returns an ongoing state when the default config is valid",
			config: DefaultConfig(),
			wantState: State{
				Config:    DefaultConfig(),
				Distance:  3,
				Depth:     1,
				Status:    StatusOngoing,
				EndReason: EndReasonNone,
			},
		},
		{
			title: "allows an encounter to start close to shore",
			config: Config{
				InitialDistance:           1,
				InitialDepth:              0,
				SurfaceDepth:              0,
				CaptureDistance:           0,
				ExhaustionCaptureDistance: 2,
				PlayerWinStep:             1,
				FishWinStep:               1,
				SplashProfile:             DefaultSplashProfile(),
			},
			wantState: State{
				Config: Config{
					InitialDistance:           1,
					InitialDepth:              0,
					SurfaceDepth:              0,
					CaptureDistance:           0,
					ExhaustionCaptureDistance: 2,
					PlayerWinStep:             1,
					FishWinStep:               1,
					SplashProfile:             DefaultSplashProfile(),
				},
				Distance:  1,
				Depth:     0,
				Status:    StatusOngoing,
				EndReason: EndReasonNone,
			},
		},
	}

	for _, test := range validCases {
		t.Run(test.title, func(t *testing.T) {
			state, err := NewState(test.config)

			require.NoError(t, err)
			assert.Equal(t, test.wantState, state)
		})
	}

	invalidCases := []struct {
		title       string
		config      Config
		wantErrText string
	}{
		{
			title: "returns an error when initial distance is negative",
			config: Config{
				InitialDistance:           -1,
				InitialDepth:              1,
				SurfaceDepth:              0,
				CaptureDistance:           0,
				ExhaustionCaptureDistance: 2,
				PlayerWinStep:             1,
				FishWinStep:               1,
				SplashProfile:             DefaultSplashProfile(),
			},
			wantErrText: "initial distance must be greater than or equal to 0",
		},
		{
			title: "returns an error when initial depth is above the surface",
			config: Config{
				InitialDistance:           3,
				InitialDepth:              -1,
				SurfaceDepth:              0,
				CaptureDistance:           0,
				ExhaustionCaptureDistance: 2,
				PlayerWinStep:             1,
				FishWinStep:               1,
				SplashProfile:             DefaultSplashProfile(),
			},
			wantErrText: "initial depth must be greater than or equal to surface depth",
		},
		{
			title: "returns an error when player win step is not positive",
			config: Config{
				InitialDistance:           3,
				InitialDepth:              1,
				SurfaceDepth:              0,
				CaptureDistance:           0,
				ExhaustionCaptureDistance: 2,
				PlayerWinStep:             0,
				FishWinStep:               1,
				SplashProfile:             DefaultSplashProfile(),
			},
			wantErrText: "player win step must be greater than 0",
		},
		{
			title: "returns an error when fish win step is not positive",
			config: Config{
				InitialDistance:           3,
				InitialDepth:              1,
				SurfaceDepth:              0,
				CaptureDistance:           0,
				ExhaustionCaptureDistance: 2,
				PlayerWinStep:             1,
				FishWinStep:               0,
				SplashProfile:             DefaultSplashProfile(),
			},
			wantErrText: "fish win step must be greater than 0",
		},
		{
			title: "returns an error when exhaustion capture distance is below capture distance",
			config: Config{
				InitialDistance:           3,
				InitialDepth:              1,
				SurfaceDepth:              0,
				CaptureDistance:           1,
				ExhaustionCaptureDistance: 0,
				PlayerWinStep:             1,
				FishWinStep:               1,
				SplashProfile:             DefaultSplashProfile(),
			},
			wantErrText: "exhaustion capture distance must be at least the capture distance",
		},
		{
			title: "returns an error when splash jump count falls outside the supported range",
			config: Config{
				InitialDistance:           3,
				InitialDepth:              1,
				SurfaceDepth:              0,
				CaptureDistance:           0,
				ExhaustionCaptureDistance: 2,
				PlayerWinStep:             1,
				FishWinStep:               1,
				SplashProfile:             SplashProfile{JumpCount: 6},
			},
			wantErrText: "splash jump count must be between 1 and 5",
		},
	}

	for _, test := range invalidCases {
		t.Run(test.title, func(t *testing.T) {
			state, err := NewState(test.config)

			require.Error(t, err)
			assert.ErrorContains(t, err, test.wantErrText)
			assert.Equal(t, State{}, state)
		})
	}
}
