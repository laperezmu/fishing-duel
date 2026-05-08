package playerprofiles

import (
	"fmt"
	"pesca/internal/cards"
	"pesca/internal/domain"
)

type CardRef string

type PresetCard struct {
	Ref  CardRef
	Move domain.Move
	Card cards.PlayerCard
}

func ListPresetCards(presetID string) ([]PresetCard, error) {
	preset, err := ResolveDefaultPreset(presetID)
	if err != nil {
		return nil, err
	}

	orderedMoves := []domain.Move{domain.Blue, domain.Red, domain.Yellow}
	listed := make([]PresetCard, 0)
	for _, move := range orderedMoves {
		cardsForMove := preset.Config.InitialDecks[move]
		for index, card := range cardsForMove {
			listed = append(listed, PresetCard{
				Ref:  buildCardRef(move, index),
				Move: move,
				Card: cards.ClonePlayerCard(card),
			})
		}
	}

	return listed, nil
}

func ResolvePresetCard(presetID string, ref CardRef) (cards.PlayerCard, error) {
	listed, err := ListPresetCards(presetID)
	if err != nil {
		return cards.PlayerCard{}, err
	}
	for _, listedCard := range listed {
		if listedCard.Ref == ref {
			return cards.ClonePlayerCard(listedCard.Card), nil
		}
	}

	return cards.PlayerCard{}, fmt.Errorf("unknown player card ref %q for preset %q", ref, presetID)
}

func buildCardRef(move domain.Move, index int) CardRef {
	return CardRef(fmt.Sprintf("%s-%d", move.String(), index+1))
}
