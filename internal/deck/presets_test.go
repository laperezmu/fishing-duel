package deck

import (
	"pesca/internal/cards"
	"pesca/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomFishDeckBuildCopiesCardsAndUsesDeckConfig(t *testing.T) {
	customFishDeck := CustomFishDeck{
		Name: "Prueba",
		FishCards: []cards.FishCard{
			cards.NewFishCard(domain.Red, cards.CardEffect{Trigger: cards.TriggerOnDraw, CaptureDistanceBonus: 1}),
		},
		CardsToRemove: 0,
		Shuffle:       false,
	}

	builtDeck := customFishDeck.Build(func([]cards.FishCard) {
		t.Fatal("the shuffler should not be called for fixed-order presets")
	})

	require.Len(t, builtDeck.activeCards, 1)
	assert.Equal(t, domain.Red, builtDeck.activeCards[0].Move)
	require.Len(t, builtDeck.activeCards[0].Effects, 1)
	assert.Equal(t, 1, builtDeck.activeCards[0].Effects[0].CaptureDistanceBonus)

	customFishDeck.FishCards[0].Effects[0].CaptureDistanceBonus = 99
	assert.Equal(t, 1, builtDeck.activeCards[0].Effects[0].CaptureDistanceBonus)
	assert.Equal(t, 0, builtDeck.recyclePolicy.(RemoveCardsRecyclePolicy).CardsToRemove)
}

func TestDefaultCustomFishDecks(t *testing.T) {
	customFishDecks := DefaultCustomFishDecks()

	require.Len(t, customFishDecks, 4)
	assert.Equal(t, "classic", customFishDecks[0].ID)
	assert.Equal(t, "Clasico", customFishDecks[0].Name)
	assert.NotEmpty(t, customFishDecks[0].Details)
	assert.True(t, customFishDecks[0].Shuffle)
	assert.Equal(t, "mixed-current", customFishDecks[3].ID)
	assert.False(t, customFishDecks[1].Shuffle)
	assert.Contains(t, customFishDecks[1].Description, "on_draw")
	assert.Contains(t, customFishDecks[3].Description, "post-outcome")

	hasDrawEffect := false
	for _, fishCard := range customFishDecks[1].FishCards {
		for _, effect := range fishCard.Effects {
			if effect.Trigger == cards.TriggerOnDraw {
				hasDrawEffect = true
			}
		}
	}

	assert.True(t, hasDrawEffect)
}
