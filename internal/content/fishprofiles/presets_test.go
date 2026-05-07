package fishprofiles

import (
	"pesca/internal/cards"
	"pesca/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFishDeckPresetBuildDeckCopiesCardsAndUsesDeckConfig(t *testing.T) {
	preset := FishDeckPreset{
		Name: "Prueba",
		FishCards: []cards.FishCard{
			func() cards.FishCard {
				card := cards.NewNamedFishCard("Tiron de apertura", "Permite capturar desde un paso mas lejos este round.", domain.Red, cards.CardEffect{Trigger: cards.TriggerOnDraw, CaptureDistanceBonus: 1})
				card.DiscardVisibility = cards.DiscardVisibilityMoveOnly
				return card
			}(),
		},
		CardsToRemove: 0,
		Shuffle:       false,
	}

	builtDeck := preset.BuildDeck(func([]cards.FishCard) {
		t.Fatal("the shuffler should not be called for fixed-order presets")
	})
	drawnCard, err := builtDeck.Draw()
	require.NoError(t, err)

	assert.Equal(t, "Tiron de apertura", drawnCard.Name)
	assert.Equal(t, "Permite capturar desde un paso mas lejos este round.", drawnCard.Summary)
	assert.Equal(t, domain.Red, drawnCard.Move)
	assert.Equal(t, cards.DiscardVisibilityMoveOnly, drawnCard.DiscardVisibility)
	require.Len(t, drawnCard.Effects, 1)
	assert.Equal(t, 1, drawnCard.Effects[0].CaptureDistanceBonus)

	preset.FishCards[0].Effects[0].CaptureDistanceBonus = 99
	assert.Equal(t, 1, drawnCard.Effects[0].CaptureDistanceBonus)
}

func TestDefaultPresets(t *testing.T) {
	presets := DefaultPresets()

	require.Len(t, presets, 7)
	assert.Equal(t, ProfileID("classic"), presets[0].ID)
	assert.Equal(t, ArchetypeBaselineCycle, presets[0].ArchetypeID)
	assert.Equal(t, "Clasico", presets[0].Name)
	assert.NotEmpty(t, presets[0].Details)
	assert.True(t, presets[0].Shuffle)
	assert.Equal(t, ProfileID("mixed-current"), presets[6].ID)
	assert.False(t, presets[1].Shuffle)
	assert.Equal(t, ArchetypeDrawTempo, presets[1].ArchetypeID)
	assert.Contains(t, presets[1].Description, "tempo")
	assert.Equal(t, ArchetypeHybridPressure, presets[6].ArchetypeID)
	assert.Contains(t, presets[6].Description, "respuestas segun el resultado")

	hasDrawEffect := false
	for _, fishCard := range presets[1].FishCards {
		for _, effect := range fishCard.Effects {
			if effect.Trigger == cards.TriggerOnDraw {
				hasDrawEffect = true
			}
		}
	}

	assert.True(t, hasDrawEffect)
}
