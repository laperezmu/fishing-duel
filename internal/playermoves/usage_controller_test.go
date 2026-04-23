package playermoves

import (
	"pesca/internal/cards"
	"pesca/internal/domain"
	"pesca/internal/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsageControllerInitialize(t *testing.T) {
	controller, err := NewUsageController(Config{
		InitialDecks: map[domain.Move][]cards.PlayerCard{
			domain.Blue:   {cards.NewPlayerCard(domain.Blue), cards.NewPlayerCard(domain.Blue)},
			domain.Red:    {cards.NewPlayerCard(domain.Red), cards.NewPlayerCard(domain.Red), cards.NewPlayerCard(domain.Red), cards.NewPlayerCard(domain.Red)},
			domain.Yellow: {cards.NewPlayerCard(domain.Yellow)},
		},
		RecoveryDelayRounds: 1,
	})
	require.NoError(t, err)

	state := match.State{}
	controller.Initialize(&state)

	require.Len(t, state.PlayerMoves.Moves, 3)
	assert.Equal(t, domain.Blue, state.PlayerMoves.Moves[0].Move)
	assert.Equal(t, 2, state.PlayerMoves.Moves[0].MaxUses)
	assert.Equal(t, 2, state.PlayerMoves.Moves[0].RemainingUses)
	require.Len(t, state.PlayerMoves.Moves[0].ActiveCards, 2)
	assert.Equal(t, domain.Blue, state.PlayerMoves.Moves[0].ActiveCards[0].Move)
	assert.Equal(t, domain.Red, state.PlayerMoves.Moves[1].Move)
	assert.Equal(t, 4, state.PlayerMoves.Moves[1].MaxUses)
	assert.Equal(t, 4, state.PlayerMoves.Moves[1].RemainingUses)
	assert.Equal(t, domain.Yellow, state.PlayerMoves.Moves[2].Move)
	assert.Equal(t, 1, state.PlayerMoves.Moves[2].MaxUses)
	assert.Equal(t, 1, state.PlayerMoves.Moves[2].RemainingUses)
}

func TestUsageControllerPrepareRound(t *testing.T) {
	controller, err := NewUsageController(Config{
		InitialDecks: map[domain.Move][]cards.PlayerCard{
			domain.Blue:   {cards.NewPlayerCard(domain.Blue)},
			domain.Red:    {cards.NewPlayerCard(domain.Red)},
			domain.Yellow: {cards.NewPlayerCard(domain.Yellow)},
		},
		RecoveryDelayRounds: 1,
		DeckShuffler: func(playerCards []cards.PlayerCard) {
			if len(playerCards) < 2 {
				return
			}
			playerCards[0], playerCards[1] = playerCards[1], playerCards[0]
		},
	})
	require.NoError(t, err)

	t.Run("restores a deck when its recharge round has arrived and shuffles it", func(t *testing.T) {
		firstCard := cards.NewPlayerCard(domain.Blue, cards.CardEffect{Trigger: cards.TriggerOnDraw, CaptureDistanceBonus: 1})
		secondCard := cards.NewPlayerCard(domain.Blue, cards.CardEffect{Trigger: cards.TriggerOnDraw, SurfaceDepthBonus: 1})
		state := match.State{
			Round: 2,
			PlayerMoves: match.PlayerMoveResources{Moves: []match.PlayerMoveState{
				{Move: domain.Blue, MaxUses: 2, RemainingUses: 0, RestoresOnRound: 3, DiscardedCards: []cards.PlayerCard{firstCard, secondCard}},
				{Move: domain.Red, MaxUses: 1, RemainingUses: 1, ActiveCards: []cards.PlayerCard{cards.NewPlayerCard(domain.Red)}},
				{Move: domain.Yellow, MaxUses: 1, RemainingUses: 1, ActiveCards: []cards.PlayerCard{cards.NewPlayerCard(domain.Yellow)}},
			}},
		}

		controller.PrepareRound(&state)

		assert.Equal(t, 2, state.PlayerMoves.Moves[0].RemainingUses)
		assert.Equal(t, 2, state.PlayerMoves.Moves[0].MaxUses)
		assert.Equal(t, 0, state.PlayerMoves.Moves[0].RestoresOnRound)
		require.Len(t, state.PlayerMoves.Moves[0].ActiveCards, 2)
		assert.Equal(t, 1, state.PlayerMoves.Moves[0].ActiveCards[0].Effects[0].SurfaceDepthBonus)
		assert.Empty(t, state.PlayerMoves.Moves[0].DiscardedCards)
	})

	t.Run("keeps a deck blocked when its recharge round has not arrived yet", func(t *testing.T) {
		state := match.State{
			Round: 1,
			PlayerMoves: match.PlayerMoveResources{Moves: []match.PlayerMoveState{
				{Move: domain.Blue, MaxUses: 2, RemainingUses: 0, RestoresOnRound: 3, DiscardedCards: []cards.PlayerCard{cards.NewPlayerCard(domain.Blue), cards.NewPlayerCard(domain.Blue)}},
				{Move: domain.Red, MaxUses: 1, RemainingUses: 1, ActiveCards: []cards.PlayerCard{cards.NewPlayerCard(domain.Red)}},
				{Move: domain.Yellow, MaxUses: 1, RemainingUses: 1, ActiveCards: []cards.PlayerCard{cards.NewPlayerCard(domain.Yellow)}},
			}},
		}

		controller.PrepareRound(&state)

		assert.Equal(t, 0, state.PlayerMoves.Moves[0].RemainingUses)
		assert.Equal(t, 3, state.PlayerMoves.Moves[0].RestoresOnRound)
		assert.Empty(t, state.PlayerMoves.Moves[0].ActiveCards)
		require.Len(t, state.PlayerMoves.Moves[0].DiscardedCards, 2)
	})
}

func TestUsageControllerValidateMove(t *testing.T) {
	controller, err := NewUsageController(DefaultConfig())
	require.NoError(t, err)

	t.Run("returns nil when the selected color still has cards remaining", func(t *testing.T) {
		state := match.State{PlayerMoves: match.PlayerMoveResources{Moves: []match.PlayerMoveState{
			{Move: domain.Blue, MaxUses: 3, RemainingUses: 1, ActiveCards: []cards.PlayerCard{cards.NewPlayerCard(domain.Blue)}},
			{Move: domain.Red, MaxUses: 3, RemainingUses: 3, ActiveCards: []cards.PlayerCard{cards.NewPlayerCard(domain.Red), cards.NewPlayerCard(domain.Red), cards.NewPlayerCard(domain.Red)}},
			{Move: domain.Yellow, MaxUses: 3, RemainingUses: 3, ActiveCards: []cards.PlayerCard{cards.NewPlayerCard(domain.Yellow), cards.NewPlayerCard(domain.Yellow), cards.NewPlayerCard(domain.Yellow)}},
		}}}

		assert.NoError(t, controller.ValidateMove(state, domain.Blue))
	})

	t.Run("returns an unavailable error when the selected color is recharging", func(t *testing.T) {
		state := match.State{PlayerMoves: match.PlayerMoveResources{Moves: []match.PlayerMoveState{
			{Move: domain.Blue, MaxUses: 3, RemainingUses: 0, RestoresOnRound: 4},
			{Move: domain.Red, MaxUses: 3, RemainingUses: 3, ActiveCards: []cards.PlayerCard{cards.NewPlayerCard(domain.Red), cards.NewPlayerCard(domain.Red), cards.NewPlayerCard(domain.Red)}},
			{Move: domain.Yellow, MaxUses: 3, RemainingUses: 3, ActiveCards: []cards.PlayerCard{cards.NewPlayerCard(domain.Yellow), cards.NewPlayerCard(domain.Yellow), cards.NewPlayerCard(domain.Yellow)}},
		}}}

		assert.ErrorIs(t, controller.ValidateMove(state, domain.Blue), ErrMoveUnavailable)
	})
}

func TestUsageControllerPeekMoveCard(t *testing.T) {
	controller, err := NewUsageController(DefaultConfig())
	require.NoError(t, err)

	t.Run("returns the visible top card for the selected color", func(t *testing.T) {
		expectedCard := cards.NewPlayerCard(domain.Blue, cards.CardEffect{Trigger: cards.TriggerOnDraw, CaptureDistanceBonus: 1})
		state := match.State{PlayerMoves: match.PlayerMoveResources{Moves: []match.PlayerMoveState{
			{Move: domain.Blue, MaxUses: 2, RemainingUses: 2, ActiveCards: []cards.PlayerCard{expectedCard, cards.NewPlayerCard(domain.Blue)}},
			{Move: domain.Red, MaxUses: 1, RemainingUses: 1, ActiveCards: []cards.PlayerCard{cards.NewPlayerCard(domain.Red)}},
			{Move: domain.Yellow, MaxUses: 1, RemainingUses: 1, ActiveCards: []cards.PlayerCard{cards.NewPlayerCard(domain.Yellow)}},
		}}}

		playerCard, err := controller.PeekMoveCard(state, domain.Blue)

		require.NoError(t, err)
		assert.Equal(t, expectedCard, playerCard)
	})

	t.Run("returns an unavailable error when the deck has no visible card", func(t *testing.T) {
		state := match.State{PlayerMoves: match.PlayerMoveResources{Moves: []match.PlayerMoveState{{Move: domain.Blue, RestoresOnRound: 4}}}}

		_, err := controller.PeekMoveCard(state, domain.Blue)

		assert.ErrorIs(t, err, ErrMoveUnavailable)
	})
}

func TestUsageControllerConsumeMove(t *testing.T) {
	controller, err := NewUsageController(Config{
		InitialDecks: map[domain.Move][]cards.PlayerCard{
			domain.Blue:   {cards.NewPlayerCard(domain.Blue), cards.NewPlayerCard(domain.Blue)},
			domain.Red:    {cards.NewPlayerCard(domain.Red), cards.NewPlayerCard(domain.Red), cards.NewPlayerCard(domain.Red)},
			domain.Yellow: {cards.NewPlayerCard(domain.Yellow), cards.NewPlayerCard(domain.Yellow), cards.NewPlayerCard(domain.Yellow)},
		},
		RecoveryDelayRounds: 1,
	})
	require.NoError(t, err)

	t.Run("consumes the visible card and reveals the next one", func(t *testing.T) {
		firstCard := cards.NewPlayerCard(domain.Blue, cards.CardEffect{Trigger: cards.TriggerOnDraw, CaptureDistanceBonus: 1})
		secondCard := cards.NewPlayerCard(domain.Blue, cards.CardEffect{Trigger: cards.TriggerOnDraw, SurfaceDepthBonus: 1})
		state := match.State{
			Round: 1,
			PlayerMoves: match.PlayerMoveResources{Moves: []match.PlayerMoveState{
				{Move: domain.Blue, MaxUses: 2, RemainingUses: 2, ActiveCards: []cards.PlayerCard{firstCard, secondCard}},
				{Move: domain.Red, MaxUses: 3, RemainingUses: 3, ActiveCards: []cards.PlayerCard{cards.NewPlayerCard(domain.Red), cards.NewPlayerCard(domain.Red), cards.NewPlayerCard(domain.Red)}},
				{Move: domain.Yellow, MaxUses: 3, RemainingUses: 3, ActiveCards: []cards.PlayerCard{cards.NewPlayerCard(domain.Yellow), cards.NewPlayerCard(domain.Yellow), cards.NewPlayerCard(domain.Yellow)}},
			}},
		}

		usedCard := controller.ConsumeMove(&state, domain.Blue)

		assert.Equal(t, firstCard, usedCard)
		assert.Equal(t, 1, state.PlayerMoves.Moves[0].RemainingUses)
		assert.Equal(t, 2, state.PlayerMoves.Moves[0].MaxUses)
		require.Len(t, state.PlayerMoves.Moves[0].ActiveCards, 1)
		assert.Equal(t, secondCard, state.PlayerMoves.Moves[0].ActiveCards[0])
		require.Len(t, state.PlayerMoves.Moves[0].DiscardedCards, 1)
		assert.Equal(t, firstCard, state.PlayerMoves.Moves[0].DiscardedCards[0])
		assert.Equal(t, 0, state.PlayerMoves.Moves[0].RestoresOnRound)
	})

	t.Run("schedules recovery when the last card in a color deck is consumed", func(t *testing.T) {
		lastCard := cards.NewPlayerCard(domain.Blue)
		state := match.State{
			Round: 3,
			PlayerMoves: match.PlayerMoveResources{Moves: []match.PlayerMoveState{
				{Move: domain.Blue, MaxUses: 1, RemainingUses: 1, ActiveCards: []cards.PlayerCard{lastCard}},
				{Move: domain.Red, MaxUses: 3, RemainingUses: 3, ActiveCards: []cards.PlayerCard{cards.NewPlayerCard(domain.Red), cards.NewPlayerCard(domain.Red), cards.NewPlayerCard(domain.Red)}},
				{Move: domain.Yellow, MaxUses: 3, RemainingUses: 3, ActiveCards: []cards.PlayerCard{cards.NewPlayerCard(domain.Yellow), cards.NewPlayerCard(domain.Yellow), cards.NewPlayerCard(domain.Yellow)}},
			}},
		}

		usedCard := controller.ConsumeMove(&state, domain.Blue)

		assert.Equal(t, lastCard, usedCard)
		assert.Equal(t, 0, state.PlayerMoves.Moves[0].RemainingUses)
		assert.Equal(t, 1, state.PlayerMoves.Moves[0].MaxUses)
		assert.Equal(t, 5, state.PlayerMoves.Moves[0].RestoresOnRound)
		assert.Empty(t, state.PlayerMoves.Moves[0].ActiveCards)
		require.Len(t, state.PlayerMoves.Moves[0].DiscardedCards, 1)
	})
}
