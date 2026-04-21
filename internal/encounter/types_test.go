package encounter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewState(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		wantState   State
		wantErrText string
	}{
		{
			name:   "default config is valid",
			config: DefaultConfig(),
			wantState: State{
				Config:    DefaultConfig(),
				Distance:  3,
				Status:    StatusOngoing,
				EndReason: EndReasonNone,
			},
		},
		{
			name: "rejects initial distance at or below two",
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
			name: "rejects escape distance below initial distance",
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
			name: "rejects non positive player win step",
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
			name: "rejects non positive fish win step",
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
			name: "rejects exhaustion capture distance below capture distance",
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state, err := NewState(test.config)

			if test.wantErrText != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.wantErrText)
				assert.Equal(t, State{}, state)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.wantState, state)
		})
	}
}
