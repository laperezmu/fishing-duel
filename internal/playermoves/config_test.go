package playermoves

import (
	"pesca/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUsageController(t *testing.T) {
	tests := []struct {
		title   string
		config  Config
		wantErr string
	}{
		{
			title:  "returns a controller when the config is valid",
			config: DefaultConfig(),
		},
		{
			title: "returns an error when recovery delay is negative",
			config: Config{
				InitialUsesPerMove: map[domain.Move]int{
					domain.Blue:   3,
					domain.Red:    3,
					domain.Yellow: 3,
				},
				RecoveryDelayRounds: -1,
			},
			wantErr: "recovery delay rounds must be greater than or equal to 0",
		},
		{
			title: "returns an error when a move is missing from the config",
			config: Config{
				InitialUsesPerMove: map[domain.Move]int{
					domain.Blue: 3,
					domain.Red:  3,
				},
				RecoveryDelayRounds: 1,
			},
			wantErr: "initial uses for move yellow are required",
		},
		{
			title: "returns an error when a move starts with negative uses",
			config: Config{
				InitialUsesPerMove: map[domain.Move]int{
					domain.Blue:   3,
					domain.Red:    -1,
					domain.Yellow: 3,
				},
				RecoveryDelayRounds: 1,
			},
			wantErr: "initial uses for move red must be greater than or equal to 0",
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			controller, err := NewUsageController(test.config)

			if test.wantErr != "" {
				require.Error(t, err)
				assert.EqualError(t, err, test.wantErr)
				assert.Nil(t, controller)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, controller)
		})
	}
}
