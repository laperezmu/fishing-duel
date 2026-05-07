package fishprofiles

import (
	"fmt"
	"pesca/internal/cards"
	"pesca/internal/domain"
)

type ProfileID string

type CardPattern struct {
	Name              string
	Summary           string
	Move              domain.Move
	Effects           []cards.CardEffect
	DiscardVisibility cards.DiscardVisibility
}

func (pattern CardPattern) BuildCard() cards.FishCard {
	var card cards.FishCard
	if pattern.Name != "" || pattern.Summary != "" {
		card = cards.NewNamedFishCard(pattern.Name, pattern.Summary, pattern.Move, pattern.Effects...)
	} else {
		card = cards.NewFishCard(pattern.Move, pattern.Effects...)
	}

	if pattern.DiscardVisibility != "" {
		card.DiscardVisibility = pattern.DiscardVisibility
	}

	return card
}

type Profile struct {
	ID            ProfileID
	ArchetypeID   ArchetypeID
	Name          string
	Description   string
	Details       []string
	Appearance    Appearance
	Cards         []CardPattern
	CardsToRemove int
	Shuffle       bool
}

func (id ProfileID) Validate() error {
	if id == "" {
		return fmt.Errorf("profile id is required")
	}

	return nil
}

func (profile Profile) Validate() error {
	if profile.ID == "" {
		return profile.ID.Validate()
	}
	if profile.Name == "" {
		return fmt.Errorf("profile name is required")
	}
	if err := profile.ArchetypeID.Validate(); err != nil {
		return err
	}
	if err := profile.Appearance.Validate(); err != nil {
		return fmt.Errorf("appearance: %w", err)
	}

	return nil
}

func (profile Profile) BuildPreset() FishDeckPreset {
	return FishDeckPreset{
		ID:            profile.ID,
		ArchetypeID:   profile.ArchetypeID,
		Name:          profile.Name,
		Description:   profile.Description,
		Details:       append([]string(nil), profile.Details...),
		FishCards:     profile.BuildCards(),
		CardsToRemove: profile.CardsToRemove,
		Shuffle:       profile.Shuffle,
	}
}

func (profile Profile) BuildCards() []cards.FishCard {
	builtCards := make([]cards.FishCard, 0, len(profile.Cards))
	for _, pattern := range profile.Cards {
		builtCards = append(builtCards, pattern.BuildCard())
	}

	return builtCards
}

func DefaultProfiles() []Profile {
	return DefaultCatalog().Profiles()
}
