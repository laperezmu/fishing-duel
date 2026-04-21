package game_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pesca/internal/deck"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/endings"
	"pesca/internal/game"
	"pesca/internal/progression"
	"pesca/internal/rules"
)

func TestEngineAppliesEncounterEndConditions(t *testing.T) {
	tests := []struct {
		name              string
		cards             []domain.Move
		config            encounter.Config
		plays             []domain.Move
		wantStatus        encounter.Status
		wantEndReason     encounter.EndReason
		wantDistance      int
		wantPlayerWins    int
		wantFishWins      int
		wantDraws         int
		wantFinishedError error
	}{
		{
			name:           "captures when track reaches player",
			cards:          []domain.Move{domain.Red, domain.Red, domain.Red},
			config:         encounter.DefaultConfig(),
			plays:          []domain.Move{domain.Blue, domain.Blue, domain.Blue},
			wantStatus:     encounter.StatusCaptured,
			wantEndReason:  encounter.EndReasonTrackCapture,
			wantDistance:   0,
			wantPlayerWins: 3,
		},
		{
			name:          "escapes when track exceeds limit",
			cards:         []domain.Move{domain.Yellow, domain.Yellow, domain.Yellow},
			config:        encounter.DefaultConfig(),
			plays:         []domain.Move{domain.Blue, domain.Blue, domain.Blue},
			wantStatus:    encounter.StatusEscaped,
			wantEndReason: encounter.EndReasonTrackEscape,
			wantDistance:  6,
			wantFishWins:  3,
		},
		{
			name:           "captures on deck exhaustion near player",
			cards:          []domain.Move{domain.Red},
			config:         encounter.DefaultConfig(),
			plays:          []domain.Move{domain.Blue},
			wantStatus:     encounter.StatusCaptured,
			wantEndReason:  encounter.EndReasonDeckCapture,
			wantDistance:   2,
			wantPlayerWins: 1,
		},
		{
			name:  "escapes on deck exhaustion far from player",
			cards: []domain.Move{domain.Yellow},
			config: encounter.Config{
				InitialDistance:           3,
				CaptureDistance:           0,
				EscapeDistance:            10,
				ExhaustionCaptureDistance: 2,
				PlayerWinStep:             1,
				FishWinStep:               1,
			},
			plays:             []domain.Move{domain.Blue},
			wantStatus:        encounter.StatusEscaped,
			wantEndReason:     encounter.EndReasonDeckEscape,
			wantDistance:      4,
			wantFishWins:      1,
			wantFinishedError: game.ErrGameFinished,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			engine := newEngineForTest(t, test.cards, test.config)

			for i, move := range test.plays {
				result, err := engine.PlayRound(move)
				require.NoError(t, err)
				if i < len(test.plays)-1 {
					assert.False(t, result.State.Finished)
				}
			}

			state := engine.State()
			assert.True(t, state.Finished)
			assert.Equal(t, len(test.plays), state.Round)
			assert.Equal(t, test.wantStatus, state.Encounter.Status)
			assert.Equal(t, test.wantEndReason, state.Encounter.EndReason)
			assert.Equal(t, test.wantDistance, state.Encounter.Distance)
			assert.Equal(t, test.wantPlayerWins, state.Stats.PlayerWins)
			assert.Equal(t, test.wantFishWins, state.Stats.FishWins)
			assert.Equal(t, test.wantDraws, state.Stats.Draws)

			if test.wantFinishedError != nil {
				_, err := engine.PlayRound(domain.Blue)
				assert.ErrorIs(t, err, test.wantFinishedError)
			}
		})
	}
}

func newEngineForTest(t *testing.T, cards []domain.Move, config encounter.Config) *game.Engine {
	t.Helper()

	encounterState, err := encounter.NewState(config)
	require.NoError(t, err)

	engine, err := game.NewEngine(
		deck.NewManager(cards, func([]domain.Move) {}, deck.RemoveCardsRecyclePolicy{CardsToRemove: 3}),
		rules.NewClassicEvaluator(rules.NewFishCombatProfile()),
		progression.TrackPolicy{},
		endings.EncounterCondition{},
		game.State{Encounter: encounterState},
	)
	require.NoError(t, err)

	return engine
}
