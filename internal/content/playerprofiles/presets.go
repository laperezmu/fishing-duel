package playerprofiles

import (
	"fmt"
	"pesca/internal/cards"
	"pesca/internal/domain"
	"pesca/internal/player/playermoves"
)

type DeckPreset struct {
	ID          string
	Name        string
	Description string
	Details     []string
	Config      playermoves.Config
}

func (preset DeckPreset) BuildConfig(shuffler func([]cards.PlayerCard)) playermoves.Config {
	config := playermoves.Config{
		InitialDecks:        cloneInitialDecks(preset.Config.InitialDecks),
		DeckShuffler:        shuffler,
		RecoveryDelayRounds: preset.Config.RecoveryDelayRounds,
	}

	return config
}

func DefaultPresets() []DeckPreset {
	presets := make([]DeckPreset, 0, len(defaultPresets))
	for _, preset := range defaultPresets {
		presets = append(presets, clonePreset(preset))
	}

	return presets
}

func ResolveDefaultPreset(id string) (DeckPreset, error) {
	preset, ok := defaultPresetByID[id]
	if !ok {
		return DeckPreset{}, fmt.Errorf("unknown player deck preset %q", id)
	}

	return clonePreset(preset), nil
}

func cloneInitialDecks(initialDecks map[domain.Move][]cards.PlayerCard) map[domain.Move][]cards.PlayerCard {
	clonedDecks := make(map[domain.Move][]cards.PlayerCard, len(initialDecks))
	for move, configuredDeck := range initialDecks {
		clonedDeck := make([]cards.PlayerCard, 0, len(configuredDeck))
		for _, playerCard := range configuredDeck {
			clonedDeck = append(clonedDeck, cards.ClonePlayerCard(playerCard))
		}
		clonedDecks[move] = clonedDeck
	}

	return clonedDecks
}

func clonePreset(preset DeckPreset) DeckPreset {
	clonedPreset := preset
	clonedPreset.Details = append([]string(nil), preset.Details...)
	clonedPreset.Config = preset.BuildConfig(nil)

	return clonedPreset
}

func buildDefaultPresetIndex(presets []DeckPreset) map[string]DeckPreset {
	index := make(map[string]DeckPreset, len(presets))
	for _, preset := range presets {
		index[preset.ID] = clonePreset(preset)
	}

	return index
}
