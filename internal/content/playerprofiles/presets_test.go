package playerprofiles

import (
	"pesca/internal/cards"
	"pesca/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultPresets(t *testing.T) {
	presets := DefaultPresets()

	require.Len(t, presets, 4)
	assert.Equal(t, "classic", presets[0].ID)
	assert.Equal(t, "Clasico", presets[0].Name)
	assert.Contains(t, presets[1].Description, "ventajas temporales")
	require.NotEmpty(t, presets[1].Details)
	assert.Contains(t, presets[1].Details[0], "Azul")
}

func TestDeckPresetBuildConfig(t *testing.T) {
	preset := DefaultPresets()[3]
	shuffleCalls := 0
	config := preset.BuildConfig(func([]cards.PlayerCard) {
		shuffleCalls++
	})

	require.NoError(t, config.Validate())
	assert.NotNil(t, config.DeckShuffler)
	require.Len(t, config.InitialDecks[domain.Blue], 3)
	require.Len(t, config.InitialDecks[domain.Blue][0].Effects, 2)
	assert.Equal(t, cards.TriggerOnDraw, config.InitialDecks[domain.Blue][0].Effects[0].Trigger)
	assert.Equal(t, cards.EffectTypeLegacyCaptureWindow, config.InitialDecks[domain.Blue][0].Effects[0].Type)
	assert.Equal(t, 60, config.InitialDecks[domain.Blue][0].Effects[0].Priority)

	preset.Config.InitialDecks[domain.Blue][0].Effects[0].CaptureDistanceBonus = 99
	assert.Equal(t, 1, config.InitialDecks[domain.Blue][0].Effects[0].CaptureDistanceBonus)

	config.DeckShuffler(config.InitialDecks[domain.Blue])
	assert.Equal(t, 1, shuffleCalls)
}

func TestDefaultPresetsUseNormalizedEffects(t *testing.T) {
	hasCaptureWindow := false
	hasSurfaceWindow := false
	hasExhaustionWindow := false
	hasHorizontal := false
	hasVertical := false

	for _, preset := range DefaultPresets() {
		for move, cardsForMove := range preset.Config.InitialDecks {
			for _, card := range cardsForMove {
				assert.Equal(t, move, card.Move)
				for _, effect := range card.Effects {
					assert.NotEqual(t, cards.EffectTypeUnknown, effect.Type)
					assert.Positive(t, effect.Priority)
					switch effect.Type {
					case cards.EffectTypeLegacyCaptureWindow:
						hasCaptureWindow = true
					case cards.EffectTypeLegacySurfaceWindow:
						hasSurfaceWindow = true
					case cards.EffectTypeLegacyExhaustionWindow:
						hasExhaustionWindow = true
					case cards.EffectTypeAdvanceHorizontal:
						hasHorizontal = true
					case cards.EffectTypeAdvanceVertical:
						hasVertical = true
					}
				}
			}
		}
	}

	assert.True(t, hasCaptureWindow)
	assert.True(t, hasSurfaceWindow)
	assert.True(t, hasExhaustionWindow)
	assert.True(t, hasHorizontal)
	assert.True(t, hasVertical)
}

func TestListPresetCards(t *testing.T) {
	listed, err := ListPresetCards("hooked-opening")

	require.NoError(t, err)
	require.NotEmpty(t, listed)
	assert.Equal(t, CardRef("blue-1"), listed[0].Ref)
	assert.Equal(t, domain.Blue, listed[0].Move)
	assert.Equal(t, "Anzuelo tenso", listed[0].Card.Name)
}

func TestResolvePresetCard(t *testing.T) {
	card, err := ResolvePresetCard("hooked-opening", CardRef("red-1"))

	require.NoError(t, err)
	assert.Equal(t, domain.Red, card.Move)
	assert.Equal(t, "Giro de superficie", card.Name)

	_, err = ResolvePresetCard("hooked-opening", CardRef("red-9"))
	require.Error(t, err)
	assert.EqualError(t, err, "unknown player card ref \"red-9\" for preset \"hooked-opening\"")
}
