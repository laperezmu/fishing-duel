package deck

import (
	"pesca/internal/cards"
	"pesca/internal/domain"
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

func TestDeckVisibilitySnapshot(t *testing.T) {
	t.Run("tracks visible entries for the current cycle", func(t *testing.T) {
		fishDeck := New([]cards.FishCard{
			cards.NewNamedFishCard("Oleaje abierto", "Empuja hacia mar abierto.", domain.Blue),
			func() cards.FishCard {
				card := cards.NewNamedFishCard("Marea calma", "Solo deja visible el movimiento.", domain.Red)
				card.DiscardVisibility = cards.DiscardVisibilityMoveOnly
				return card
			}(),
			func() cards.FishCard {
				card := cards.NewNamedFishCard("Rastro oculto", "No deja rastro en el historial.", domain.Yellow)
				card.DiscardVisibility = cards.DiscardVisibilityHidden
				return card
			}(),
			func() cards.FishCard {
				card := cards.NewNamedFishCard("Silueta turbia", "Se registra como carta desconocida.", domain.Red)
				card.DiscardVisibility = cards.DiscardVisibilityMasked
				return card
			}(),
		}, nil, RemoveCardsRecyclePolicy{CardsToRemove: 0})

		for drawIndex := 0; drawIndex < 4; drawIndex++ {
			_, err := fishDeck.Draw()
			require.NoError(t, err)
		}

		snapshot := fishDeck.VisibilitySnapshot()

		assert.Equal(t, 1, snapshot.CurrentCycle.Number)
		assert.Equal(t, 4, snapshot.CurrentCycle.TotalCards)
		require.Len(t, snapshot.CurrentCycle.Entries, 3)
		assert.Equal(t, cards.DiscardVisibilityMasked, snapshot.CurrentCycle.Entries[0].Visibility)
		assert.Equal(t, cards.DiscardVisibilityMoveOnly, snapshot.CurrentCycle.Entries[1].Visibility)
		assert.Equal(t, domain.Red, snapshot.CurrentCycle.Entries[1].Move)
		assert.Equal(t, cards.DiscardVisibilityFull, snapshot.CurrentCycle.Entries[2].Visibility)
		assert.Equal(t, "Oleaje abierto", snapshot.CurrentCycle.Entries[2].Name)
		assert.Equal(t, 0, snapshot.CardsToRemove)
		assert.False(t, snapshot.ShufflesOnRecycle)
	})

	t.Run("moves completed cycles into the previous cycle summaries after recycling", func(t *testing.T) {
		fishDeck := New([]cards.FishCard{
			cards.NewFishCard(domain.Blue),
			func() cards.FishCard {
				card := cards.NewFishCard(domain.Red)
				card.DiscardVisibility = cards.DiscardVisibilityHidden
				return card
			}(),
		}, nil, RemoveCardsRecyclePolicy{CardsToRemove: 1})

		for drawIndex := 0; drawIndex < 2; drawIndex++ {
			_, err := fishDeck.Draw()
			require.NoError(t, err)
		}

		fishDeck.PrepareNextRound()

		snapshot := fishDeck.VisibilitySnapshot()

		assert.Equal(t, 1, snapshot.RecycleCount)
		assert.Equal(t, 1, snapshot.CardsToRemove)
		assert.Equal(t, 2, snapshot.CurrentCycle.Number)
		assert.Equal(t, 0, snapshot.CurrentCycle.TotalCards)
		require.Len(t, snapshot.PreviousCycles, 1)
		assert.Equal(t, 1, snapshot.PreviousCycles[0].Number)
		assert.Equal(t, 2, snapshot.PreviousCycles[0].TotalCards)
		assert.Equal(t, 1, snapshot.PreviousCycles[0].VisibleCards)
		assert.Equal(t, 1, snapshot.PreviousCycles[0].HiddenCards)
	})
}
