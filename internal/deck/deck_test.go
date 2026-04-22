package deck

import (
	"pesca/internal/cards"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeckDrawRecyclesAndRemovesThreeCards(t *testing.T) {
	fishDeck := New(
		NewStandardFishCards(),
		func([]cards.FishCard) {},
		RemoveCardsRecyclePolicy{CardsToRemove: 3},
	)

	for i := 0; i < 9; i++ {
		_, err := fishDeck.Draw()
		require.NoError(t, err, "draw %d failed", i+1)
	}
	assert.Equal(t, 0, fishDeck.ActiveCount())
	assert.Equal(t, 9, fishDeck.DiscardCount())

	fishDeck.PrepareNextRound()
	assert.Equal(t, 6, fishDeck.ActiveCount())
	assert.Equal(t, 0, fishDeck.DiscardCount())
	assert.Equal(t, 1, fishDeck.RecycleCount())

	for i := 0; i < 6; i++ {
		_, err := fishDeck.Draw()
		require.NoError(t, err, "second cycle draw %d failed", i+1)
	}
	fishDeck.PrepareNextRound()
	assert.Equal(t, 3, fishDeck.ActiveCount())
	assert.Equal(t, 2, fishDeck.RecycleCount())

	for i := 0; i < 3; i++ {
		_, err := fishDeck.Draw()
		require.NoError(t, err, "third cycle draw %d failed", i+1)
	}
	fishDeck.PrepareNextRound()
	assert.Equal(t, 0, fishDeck.ActiveCount())
	assert.True(t, fishDeck.Exhausted())
	assert.Equal(t, 3, fishDeck.RecycleCount())

	_, err := fishDeck.Draw()
	assert.ErrorIs(t, err, ErrNoCardsAvailable)
}
