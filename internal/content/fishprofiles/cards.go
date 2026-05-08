package fishprofiles

import (
	"fmt"
	"pesca/internal/cards"
	"pesca/internal/domain"
)

type CardRef string

type PresetCard struct {
	Ref  CardRef
	Move domain.Move
	Card cards.FishCard
}

func ListPresetCards(presetID ProfileID) ([]PresetCard, error) {
	profile, err := DefaultCatalog().ProfileByID(presetID)
	if err != nil {
		return nil, err
	}
	builtCards := profile.BuildCards()
	listed := make([]PresetCard, 0, len(builtCards))
	moveCounts := make(map[domain.Move]int)
	for _, card := range builtCards {
		moveCounts[card.Move]++
		listed = append(listed, PresetCard{
			Ref:  buildCardRef(card.Move, moveCounts[card.Move]-1),
			Move: card.Move,
			Card: cards.CloneFishCard(card),
		})
	}

	return listed, nil
}

func ResolvePresetCard(presetID ProfileID, ref CardRef) (cards.FishCard, error) {
	listed, err := ListPresetCards(presetID)
	if err != nil {
		return cards.FishCard{}, err
	}
	for _, listedCard := range listed {
		if listedCard.Ref == ref {
			return cards.CloneFishCard(listedCard.Card), nil
		}
	}

	return cards.FishCard{}, fmt.Errorf("unknown fish card ref %q for preset %q", ref, presetID)
}

func buildCardRef(move domain.Move, index int) CardRef {
	return CardRef(fmt.Sprintf("%s-%d", move.String(), index+1))
}
