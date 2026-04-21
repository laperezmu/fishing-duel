package deck

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pesca/internal/domain"
)

func TestManagerRecyclesAndRemovesThreeCards(t *testing.T) {
	manager := NewManager(
		NewStandardFishDeck(),
		func([]domain.Move) {},
		RemoveCardsRecyclePolicy{CardsToRemove: 3},
	)

	for i := 0; i < 9; i++ {
		_, err := manager.Draw()
		require.NoError(t, err, "draw %d failed", i+1)
	}
	assert.Equal(t, 0, manager.ActiveCount())
	assert.Equal(t, 9, manager.DiscardCount())

	manager.PrepareNextRound()
	assert.Equal(t, 6, manager.ActiveCount())
	assert.Equal(t, 0, manager.DiscardCount())
	assert.Equal(t, 1, manager.RecycleCount())

	for i := 0; i < 6; i++ {
		_, err := manager.Draw()
		require.NoError(t, err, "second cycle draw %d failed", i+1)
	}
	manager.PrepareNextRound()
	assert.Equal(t, 3, manager.ActiveCount())
	assert.Equal(t, 2, manager.RecycleCount())

	for i := 0; i < 3; i++ {
		_, err := manager.Draw()
		require.NoError(t, err, "third cycle draw %d failed", i+1)
	}
	manager.PrepareNextRound()
	assert.Equal(t, 0, manager.ActiveCount())
	assert.True(t, manager.Exhausted())
	assert.Equal(t, 3, manager.RecycleCount())

	_, err := manager.Draw()
	assert.ErrorIs(t, err, ErrNoCardsAvailable)
}
