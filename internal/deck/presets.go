package deck

import (
	"pesca/internal/cards"
	"pesca/internal/domain"
)

type CustomFishDeck struct {
	Name          string
	Description   string
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
	return []CustomFishDeck{
		{
			Name:          "Clasico",
			Description:   "Baraja base de 9 cartas sin efectos, mezclada y con reciclado clasico.",
			FishCards:     NewStandardFishCards(),
			CardsToRemove: 3,
			Shuffle:       true,
		},
		{
			Name:        "Apertura con anzuelo",
			Description: "Prueba efectos `on_draw` que alteran distancia de captura, superficie y cierre por agotamiento.",
			FishCards: []cards.FishCard{
				cards.NewFishCard(domain.Red, cards.CardEffect{Trigger: cards.TriggerOnDraw, CaptureDistanceBonus: 1}),
				cards.NewFishCard(domain.Blue, cards.CardEffect{Trigger: cards.TriggerOnDraw, SurfaceDepthBonus: 1}),
				cards.NewFishCard(domain.Yellow, cards.CardEffect{Trigger: cards.TriggerOnDraw, ExhaustionCaptureDistanceBonus: 1}),
				cards.NewFishCard(domain.Red, cards.CardEffect{Trigger: cards.TriggerOnDraw, CaptureDistanceBonus: 1}),
			},
			CardsToRemove: 0,
			Shuffle:       false,
		},
		{
			Name:        "Presion vertical",
			Description: "Prueba hundimiento al ganar y subida al perder para validar la capa de profundidad.",
			FishCards: []cards.FishCard{
				cards.NewFishCard(domain.Blue, cards.CardEffect{Trigger: cards.TriggerOnOwnerWin, DepthShift: 1}),
				cards.NewFishCard(domain.Red, cards.CardEffect{Trigger: cards.TriggerOnOwnerLose, DepthShift: -1}),
				cards.NewFishCard(domain.Yellow, cards.CardEffect{Trigger: cards.TriggerOnOwnerWin, DepthShift: 1}),
				cards.NewFishCard(domain.Blue, cards.CardEffect{Trigger: cards.TriggerOnOwnerLose, DepthShift: -1}),
			},
			CardsToRemove: 0,
			Shuffle:       false,
		},
		{
			Name:        "Corriente mixta",
			Description: "Combina `on_draw` y efectos post-outcome para revisar el pipeline completo en una sola partida.",
			FishCards: []cards.FishCard{
				cards.NewFishCard(domain.Red,
					cards.CardEffect{Trigger: cards.TriggerOnDraw, CaptureDistanceBonus: 1},
					cards.CardEffect{Trigger: cards.TriggerOnOwnerLose, DepthShift: -1},
				),
				cards.NewFishCard(domain.Blue,
					cards.CardEffect{Trigger: cards.TriggerOnDraw, SurfaceDepthBonus: 1},
					cards.CardEffect{Trigger: cards.TriggerOnOwnerWin, DistanceShift: 1},
				),
				cards.NewFishCard(domain.Yellow,
					cards.CardEffect{Trigger: cards.TriggerOnRoundDraw, DistanceShift: 1},
				),
				cards.NewFishCard(domain.Red,
					cards.CardEffect{Trigger: cards.TriggerOnOwnerLose, DepthShift: -1},
				),
			},
			CardsToRemove: 0,
			Shuffle:       false,
		},
	}
}

func cloneFishCards(fishCards []cards.FishCard) []cards.FishCard {
	clonedFishCards := make([]cards.FishCard, 0, len(fishCards))
	for _, fishCard := range fishCards {
		clonedFishCards = append(clonedFishCards, cards.NewFishCard(fishCard.Move, fishCard.Effects...))
	}

	return clonedFishCards
}
