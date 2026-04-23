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

	preset.Config.InitialDecks[domain.Blue][0].Effects[0].CaptureDistanceBonus = 99
	assert.Equal(t, 1, config.InitialDecks[domain.Blue][0].Effects[0].CaptureDistanceBonus)

	config.DeckShuffler(config.InitialDecks[domain.Blue])
	assert.Equal(t, 1, shuffleCalls)
}
