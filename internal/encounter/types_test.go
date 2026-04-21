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
			title: "returns an error when initial distance is two or lower",
			config: Config{
				InitialDistance:           2,
				CaptureDistance:           0,
				EscapeDistance:            5,
				ExhaustionCaptureDistance: 2,
				PlayerWinStep:             1,
				FishWinStep:               1,
			},
			wantErrText: "initial distance must be greater than 2",
		},
		{
			title: "returns an error when escape distance is lower than initial distance",
			config: Config{
				InitialDistance:           3,
				CaptureDistance:           0,
				EscapeDistance:            2,
				ExhaustionCaptureDistance: 2,
				PlayerWinStep:             1,
				FishWinStep:               1,
			},
			wantErrText: "escape distance must be at least the initial distance",
		},
		{
			title: "returns an error when player win step is not positive",
			config: Config{
				InitialDistance:           3,
				CaptureDistance:           0,
				EscapeDistance:            5,
				ExhaustionCaptureDistance: 2,
				PlayerWinStep:             0,
				FishWinStep:               1,
			},
			wantErrText: "player win step must be greater than 0",
		},
		{
			title: "returns an error when fish win step is not positive",
			config: Config{
				InitialDistance:           3,
				CaptureDistance:           0,
				EscapeDistance:            5,
				ExhaustionCaptureDistance: 2,
				PlayerWinStep:             1,
				FishWinStep:               0,
			},
			wantErrText: "fish win step must be greater than 0",
		},
		{
			title: "returns an error when exhaustion capture distance is below capture distance",
			config: Config{
				InitialDistance:           3,
				CaptureDistance:           1,
				EscapeDistance:            5,
				ExhaustionCaptureDistance: 0,
				PlayerWinStep:             1,
				FishWinStep:               1,
			},
			wantErrText: "exhaustion capture distance must be at least the capture distance",
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
