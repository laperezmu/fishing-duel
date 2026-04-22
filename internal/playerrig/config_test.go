package playerrig

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewState(t *testing.T) {
	tests := []struct {
		title       string
		config      Config
		wantState   State
		wantErrText string
	}{
		{
			title:  "returns a rig state when the config is valid",
			config: DefaultConfig(),
			wantState: State{
				MaxDistance: 5,
				MaxDepth:    3,
			},
		},
		{
			title: "returns an error when max distance is not positive",
			config: Config{
				MaxDistance: 0,
				MaxDepth:    3,
			},
			wantErrText: "max distance must be greater than 0",
		},
		{
			title: "returns an error when max depth is negative",
			config: Config{
				MaxDistance: 5,
				MaxDepth:    -1,
			},
			wantErrText: "max depth must be greater than or equal to 0",
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			state, err := NewState(test.config)

			if test.wantErrText != "" {
				require.Error(t, err)
				assert.EqualError(t, err, test.wantErrText)
				assert.Equal(t, State{}, state)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.wantState, state)
		})
	}
}
