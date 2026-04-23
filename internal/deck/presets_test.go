package deck

import (
	"pesca/internal/cards"
	"pesca/internal/domain"
	"pesca/internal/fishprofiles"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomFishDeckBuildCopiesCardsAndUsesDeckConfig(t *testing.T) {
	customFishDeck := CustomFishDeck{
		Name: "Prueba",
		FishCards: []cards.FishCard{
			cards.NewNamedFishCard("Tiron de apertura", "Permite capturar desde un paso mas lejos este round.", domain.Red, cards.CardEffect{Trigger: cards.TriggerOnDraw, CaptureDistanceBonus: 1}),
		},
		CardsToRemove: 0,
		Shuffle:       false,
	}

	builtDeck := customFishDeck.Build(func([]cards.FishCard) {
		t.Fatal("the shuffler should not be called for fixed-order presets")
	})

	require.Len(t, builtDeck.activeCards, 1)
	assert.Equal(t, "Tiron de apertura", builtDeck.activeCards[0].Name)
	assert.Equal(t, "Permite capturar desde un paso mas lejos este round.", builtDeck.activeCards[0].Summary)
	assert.Equal(t, domain.Red, builtDeck.activeCards[0].Move)
	require.Len(t, builtDeck.activeCards[0].Effects, 1)
	assert.Equal(t, 1, builtDeck.activeCards[0].Effects[0].CaptureDistanceBonus)

	customFishDeck.FishCards[0].Effects[0].CaptureDistanceBonus = 99
	assert.Equal(t, 1, builtDeck.activeCards[0].Effects[0].CaptureDistanceBonus)
	assert.Equal(t, 0, builtDeck.recyclePolicy.(RemoveCardsRecyclePolicy).CardsToRemove)
}

func TestDefaultCustomFishDecks(t *testing.T) {
	customFishDecks := DefaultCustomFishDecks()

	require.Len(t, customFishDecks, 7)
	assert.Equal(t, "classic", customFishDecks[0].ID)
	assert.Equal(t, fishprofiles.ArchetypeBaselineCycle, customFishDecks[0].ArchetypeID)
	assert.Equal(t, "Clasico", customFishDecks[0].Name)
	assert.NotEmpty(t, customFishDecks[0].Details)
	assert.True(t, customFishDecks[0].Shuffle)
	assert.Equal(t, "mixed-current", customFishDecks[6].ID)
	assert.False(t, customFishDecks[1].Shuffle)
	assert.Equal(t, fishprofiles.ArchetypeDrawTempo, customFishDecks[1].ArchetypeID)
	assert.Contains(t, customFishDecks[1].Description, "tempo")
	assert.Equal(t, fishprofiles.ArchetypeHybridPressure, customFishDecks[6].ArchetypeID)
	assert.Contains(t, customFishDecks[6].Description, "respuestas segun el resultado")

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
