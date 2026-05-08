package match_test

import (
	"testing"

	"pesca/internal/deck"
	"pesca/internal/domain"
	"pesca/internal/match"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProgressionState(t *testing.T) {
	t.Run("returns progression state with references to round, encounter, and lifecycle", func(t *testing.T) {
		state := &match.State{}

		progression := state.ProgressionState()

		assert.NotNil(t, progression.Round)
		assert.NotNil(t, progression.Encounter)
		assert.NotNil(t, progression.Lifecycle)
	})
}

func TestEndingState(t *testing.T) {
	t.Run("returns ending state with references to all subsystems", func(t *testing.T) {
		state := &match.State{}

		ending := state.EndingState()

		assert.NotNil(t, ending.Round)
		assert.NotNil(t, ending.Deck)
		assert.NotNil(t, ending.Encounter)
		assert.NotNil(t, ending.Player)
		assert.NotNil(t, ending.Lifecycle)
	})
}

func TestPlayerMoveRuntime(t *testing.T) {
	t.Run("returns player move runtime with round and moves", func(t *testing.T) {
		state := &match.State{}

		runtime := state.PlayerMoveRuntime()

		assert.NotNil(t, runtime.Round)
		assert.NotNil(t, runtime.Moves)
	})
}

func TestNewDeckState(t *testing.T) {
	t.Run("creates deck state with provided values", func(t *testing.T) {
		visibility := deck.VisibilitySnapshot{
			ShufflesOnRecycle: true,
			CardsToRemove:     2,
			CurrentCycle: deck.VisibleDiscardCycle{
				Number:     1,
				TotalCards: 10,
				Entries:    []deck.VisibleDiscardEntry{},
			},
			PreviousCycles: []deck.VisibleDiscardCycleSummary{},
		}

		deckState := match.NewDeckState(15, 5, 2, false, visibility)

		assert.Equal(t, 15, deckState.ActiveCards)
		assert.Equal(t, 5, deckState.DiscardCards)
		assert.Equal(t, 2, deckState.RecycleCount)
		assert.False(t, deckState.Exhausted)
		assert.True(t, deckState.ShufflesOnRecycle)
		assert.Equal(t, 2, deckState.CardsToRemove)
	})

	t.Run("maps current cycle entries", func(t *testing.T) {
		visibility := deck.VisibilitySnapshot{
			CurrentCycle: deck.VisibleDiscardCycle{
				Number:     1,
				TotalCards: 5,
				Entries: []deck.VisibleDiscardEntry{
					{
						Visibility: "visible",
						Move:       domain.Blue,
						Name:       "Test Card",
						Summary:    "Test summary",
					},
				},
			},
		}

		deckState := match.NewDeckState(10, 0, 0, false, visibility)

		require.Len(t, deckState.CurrentCycle.Entries, 1)
		assert.Equal(t, 1, deckState.CurrentCycle.Number)
		assert.Equal(t, 5, deckState.CurrentCycle.TotalCards)
		assert.Equal(t, domain.Blue, deckState.CurrentCycle.Entries[0].Move)
		assert.Equal(t, "Test Card", deckState.CurrentCycle.Entries[0].Name)
	})

	t.Run("maps previous cycle summaries", func(t *testing.T) {
		visibility := deck.VisibilitySnapshot{
			PreviousCycles: []deck.VisibleDiscardCycleSummary{
				{
					Number:       1,
					TotalCards:   10,
					VisibleCards: 6,
					HiddenCards:  4,
				},
			},
		}

		deckState := match.NewDeckState(10, 0, 0, false, visibility)

		require.Len(t, deckState.PreviousCycleStats, 1)
		assert.Equal(t, 1, deckState.PreviousCycleStats[0].Number)
		assert.Equal(t, 10, deckState.PreviousCycleStats[0].TotalCards)
		assert.Equal(t, 6, deckState.PreviousCycleStats[0].VisibleCards)
		assert.Equal(t, 4, deckState.PreviousCycleStats[0].HiddenCards)
	})
}
