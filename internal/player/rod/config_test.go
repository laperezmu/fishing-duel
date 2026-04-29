package rod

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
			title:  "returns a rod state when the config is valid",
			config: DefaultConfig(),
			wantState: State{
				OpeningMaxDistance: 5,
				OpeningMaxDepth:    3,
				TrackMaxDistance:   5,
				TrackMaxDepth:      3,
			},
		},
		{
			title: "returns an error when opening max distance is negative",
			config: Config{
				OpeningMaxDistance: -1,
				OpeningMaxDepth:    3,
				TrackMaxDistance:   5,
				TrackMaxDepth:      3,
			},
			wantErrText: "opening max distance must be greater than or equal to 0",
		},
		{
			title: "returns an error when opening max depth is negative",
			config: Config{
				OpeningMaxDistance: 5,
				OpeningMaxDepth:    -1,
				TrackMaxDistance:   5,
				TrackMaxDepth:      3,
			},
			wantErrText: "opening max depth must be greater than or equal to 0",
		},
		{
			title: "returns an error when track max distance is not positive",
			config: Config{
				OpeningMaxDistance: 5,
				OpeningMaxDepth:    3,
				TrackMaxDistance:   0,
				TrackMaxDepth:      3,
			},
			wantErrText: "track max distance must be greater than 0",
		},
		{
			title: "returns an error when track max depth is negative",
			config: Config{
				OpeningMaxDistance: 5,
				OpeningMaxDepth:    3,
				TrackMaxDistance:   5,
				TrackMaxDepth:      -1,
			},
			wantErrText: "track max depth must be greater than or equal to 0",
		},
		{
			title: "returns an error when opening distance exceeds track distance",
			config: Config{
				OpeningMaxDistance: 6,
				OpeningMaxDepth:    3,
				TrackMaxDistance:   5,
				TrackMaxDepth:      3,
			},
			wantErrText: "opening max distance must be less than or equal to track max distance",
		},
		{
			title: "returns an error when opening depth exceeds track depth",
			config: Config{
				OpeningMaxDistance: 5,
				OpeningMaxDepth:    4,
				TrackMaxDistance:   5,
				TrackMaxDepth:      3,
			},
			wantErrText: "opening max depth must be less than or equal to track max depth",
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
