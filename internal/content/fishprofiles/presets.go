package fishprofiles

import (
	"pesca/internal/cards"
	"pesca/internal/deck"
)

type FishDeckPreset struct {
	ID            ProfileID
	ArchetypeID   ArchetypeID
	Name          string
	Description   string
	Details       []string
	FishCards     []cards.FishCard
	CardsToRemove int
	Shuffle       bool
}

func (preset FishDeckPreset) BuildDeck(shuffler deck.Shuffler) *deck.Deck {
	configuredShuffler := shuffler
	if !preset.Shuffle {
		configuredShuffler = nil
	}

	return deck.New(
		cloneFishCards(preset.FishCards),
		configuredShuffler,
		deck.RemoveCardsRecyclePolicy{CardsToRemove: preset.CardsToRemove},
	)
}

func DefaultPresets() []FishDeckPreset {
	profiles := DefaultProfiles()
	presets := make([]FishDeckPreset, 0, len(profiles))
	for _, profile := range profiles {
		presets = append(presets, profile.BuildPreset())
	}

	return presets
}

func cloneFishCards(fishCards []cards.FishCard) []cards.FishCard {
	clonedFishCards := make([]cards.FishCard, 0, len(fishCards))
	for _, fishCard := range fishCards {
		clonedFishCards = append(clonedFishCards, cards.CloneFishCard(fishCard))
	}

	return clonedFishCards
}
