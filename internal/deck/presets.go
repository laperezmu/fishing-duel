package deck

import (
	"pesca/internal/cards"
	"pesca/internal/domain"
)

type CustomFishDeck struct {
	ID            string
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
	return []CustomFishDeck{
		{
			ID:          "classic",
			Name:        "Clasico",
			Description: "Baraja base de 9 cartas sin efectos con reciclado clasico.",
			Details: []string{
				"9 cartas lisas sin efectos: 3 rojas, 3 azules y 3 amarillas.",
				"Orden: barajada antes de empezar y en cada reciclado.",
				"Reciclado: retira 3 cartas por ciclo.",
			},
			FishCards:     NewStandardFishCards(),
			CardsToRemove: 3,
			Shuffle:       true,
		},
		{
			ID:          "hooked-opening",
			Name:        "Apertura con anzuelo",
			Description: "Baraja del pez centrada en `on_draw` para thresholds y cierres tempranos.",
			Details: []string{
				"Rojo: `draw capt +1`.",
				"Azul: `draw sup +1`.",
				"Amarillo: `draw baraja +1`.",
				"Orden: fijo para probar cada apertura de forma reproducible.",
			},
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
			ID:          "vertical-pressure",
			Name:        "Presion vertical",
			Description: "Baraja del pez orientada a hundirse al ganar y subir al perder.",
			Details: []string{
				"Azul: `si gana prof +1` o `si pierde prof -1` segun la carta.",
				"Rojo: `si pierde prof -1`.",
				"Amarillo: `si gana prof +1`.",
				"Orden: fijo para validar la capa vertical paso a paso.",
			},
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
			ID:          "mixed-current",
			Name:        "Corriente mixta",
			Description: "Baraja del pez que mezcla `on_draw` y respuestas post-outcome en una sola secuencia.",
			Details: []string{
				"Rojo: `draw capt +1` y `si pierde prof -1`.",
				"Azul: `draw sup +1` y `si gana dist +1`.",
				"Amarillo: `empate dist +1`.",
				"Orden: fijo para revisar el pipeline completo round a round.",
			},
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
