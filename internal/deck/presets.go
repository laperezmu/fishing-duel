package deck

import (
	"pesca/internal/cards"
	"pesca/internal/fishprofiles"
)

type CustomFishDeck struct {
	ID            string
	ArchetypeID   fishprofiles.ArchetypeID
	Name          string
	Description   string
	Details       []string
	FishCards     []cards.FishCard
	CardsToRemove int
	Shuffle       bool
}

func (customFishDeck CustomFishDeck) Build(shuffler Shuffler) *Deck {
	configuredShuffler := shuffler
	if !customFishDeck.Shuffle {
		configuredShuffler = nil
	}

	return New(
		cloneFishCards(customFishDeck.FishCards),
		configuredShuffler,
		RemoveCardsRecyclePolicy{CardsToRemove: customFishDeck.CardsToRemove},
	)
}

func DefaultCustomFishDecks() []CustomFishDeck {
	profiles := fishprofiles.DefaultProfiles()
	customFishDecks := make([]CustomFishDeck, 0, len(profiles))
	for _, profile := range profiles {
		customFishDecks = append(customFishDecks, CustomFishDeck{
			ID:            profile.ID,
			ArchetypeID:   profile.ArchetypeID,
			Name:          profile.Name,
			Description:   profile.Description,
			Details:       append([]string(nil), profile.Details...),
			FishCards:     profile.BuildCards(),
			CardsToRemove: profile.CardsToRemove,
			Shuffle:       profile.Shuffle,
		})
	}

	return customFishDecks
}

func cloneFishCards(fishCards []cards.FishCard) []cards.FishCard {
	clonedFishCards := make([]cards.FishCard, 0, len(fishCards))
	for _, fishCard := range fishCards {
		clonedFishCards = append(clonedFishCards, cards.CloneFishCard(fishCard))
	}

	return clonedFishCards
}
