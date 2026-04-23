package playermoves

import (
	"pesca/internal/cards"
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
				InitialDecks: map[domain.Move][]cards.PlayerCard{
					domain.Blue:   {cards.NewPlayerCard(domain.Blue)},
					domain.Red:    {cards.NewPlayerCard(domain.Red)},
					domain.Yellow: {cards.NewPlayerCard(domain.Yellow)},
				},
				RecoveryDelayRounds: -1,
			},
			wantErr: "recovery delay rounds must be greater than or equal to 0",
		},
		{
			title: "returns an error when a move is missing from the config",
			config: Config{
				InitialDecks: map[domain.Move][]cards.PlayerCard{
					domain.Blue: {cards.NewPlayerCard(domain.Blue)},
					domain.Red:  {cards.NewPlayerCard(domain.Red)},
				},
				RecoveryDelayRounds: 1,
			},
			wantErr: "initial deck for move yellow is required",
		},
		{
			title: "returns an error when a move starts with an empty deck",
			config: Config{
				InitialDecks: map[domain.Move][]cards.PlayerCard{
					domain.Blue:   {cards.NewPlayerCard(domain.Blue)},
					domain.Red:    {},
					domain.Yellow: {cards.NewPlayerCard(domain.Yellow)},
				},
				RecoveryDelayRounds: 1,
			},
			wantErr: "initial deck for move red must contain at least one card",
		},
		{
			title: "returns an error when a deck contains a card from another move",
			config: Config{
				InitialDecks: map[domain.Move][]cards.PlayerCard{
					domain.Blue:   {cards.NewPlayerCard(domain.Red)},
					domain.Red:    {cards.NewPlayerCard(domain.Red)},
					domain.Yellow: {cards.NewPlayerCard(domain.Yellow)},
				},
				RecoveryDelayRounds: 1,
			},
			wantErr: "initial deck for move blue contains a card for move red",
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
