package game_test

import (
	"pesca/internal/deck"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/endings"
	"pesca/internal/game"
	"pesca/internal/progression"
	"pesca/internal/rules"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnginePlayRound(t *testing.T) {
	tests := []struct {
		title               string
		cards               []domain.Move
		config              encounter.Config
		plays               []domain.Move
		wantRoundFinished   []bool
		wantStatus          encounter.Status
		wantEndReason       encounter.EndReason
		wantDistance        int
		wantPlayerWins      int
		wantFishWins        int
		wantDraws           int
		wantFollowUpPlayErr error
	}{
		{
			title:               "returns capture when repeated player wins move the fish to the shore",
			cards:               []domain.Move{domain.Red, domain.Red, domain.Red},
			config:              encounter.DefaultConfig(),
			plays:               []domain.Move{domain.Blue, domain.Blue, domain.Blue},
			wantRoundFinished:   []bool{false, false, true},
			wantStatus:          encounter.StatusCaptured,
			wantEndReason:       encounter.EndReasonTrackCapture,
			wantDistance:        0,
			wantPlayerWins:      3,
			wantFollowUpPlayErr: game.ErrGameFinished,
		},
		{
			title:               "returns escape when repeated fish wins move beyond the limit",
			cards:               []domain.Move{domain.Yellow, domain.Yellow, domain.Yellow},
			config:              encounter.DefaultConfig(),
			plays:               []domain.Move{domain.Blue, domain.Blue, domain.Blue},
			wantRoundFinished:   []bool{false, false, true},
			wantStatus:          encounter.StatusEscaped,
			wantEndReason:       encounter.EndReasonTrackEscape,
			wantDistance:        6,
			wantFishWins:        3,
			wantFollowUpPlayErr: game.ErrGameFinished,
		},
		{
			title:               "returns deck capture when the deck ends near the player",
			cards:               []domain.Move{domain.Red},
			config:              encounter.DefaultConfig(),
			plays:               []domain.Move{domain.Blue},
			wantRoundFinished:   []bool{true},
			wantStatus:          encounter.StatusCaptured,
			wantEndReason:       encounter.EndReasonDeckCapture,
			wantDistance:        2,
			wantPlayerWins:      1,
			wantFollowUpPlayErr: game.ErrGameFinished,
		},
		{
			title: "returns deck escape when the deck ends beyond the capture threshold",
			cards: []domain.Move{domain.Yellow},
			config: encounter.Config{
				InitialDistance:           3,
				CaptureDistance:           0,
				EscapeDistance:            10,
				ExhaustionCaptureDistance: 2,
				PlayerWinStep:             1,
				FishWinStep:               1,
			},
			plays:               []domain.Move{domain.Blue},
			wantRoundFinished:   []bool{true},
			wantStatus:          encounter.StatusEscaped,
			wantEndReason:       encounter.EndReasonDeckEscape,
			wantDistance:        4,
			wantFishWins:        1,
			wantFollowUpPlayErr: game.ErrGameFinished,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			engine := newEngineForTest(t, test.cards, test.config)

			for i, move := range test.plays {
				result, err := engine.PlayRound(move)
				require.NoError(t, err)
				assert.Equal(t, test.wantRoundFinished[i], result.State.Finished)
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

			_, err := engine.PlayRound(domain.Blue)
			assert.ErrorIs(t, err, test.wantFollowUpPlayErr)
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
