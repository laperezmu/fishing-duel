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
				"Nueve cartas lisas sin efectos: tres rojas, tres azules y tres amarillas.",
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
			Description: "Baraja del pez pensada para abrir el round con ventajas temporales y cierres tempranos.",
			Details: []string{
				"Rojo - Tiron de apertura: al revelarse permite capturar desde un paso mas lejos ese round.",
				"Azul - Salto de espuma: al revelarse el pez cuenta como un nivel mas cerca de la superficie ese round.",
				"Amarillo - Ultima ventana: al revelarse amplia el margen de captura cuando se agota la baraja ese round.",
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
				"Azul - Tiron al fondo: si gana, el pez baja un nivel mas profundo.",
				"Rojo - Respiro corto: si pierde, el pez sube un nivel hacia la superficie.",
				"Amarillo - Caida larga: si gana, el pez baja un nivel mas profundo.",
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
			Description: "Baraja del pez que mezcla ventajas al revelarse con respuestas segun el resultado.",
			Details: []string{
				"Rojo - Corriente cerrada: al revelarse amplia la captura y, si pierde, sube un nivel hacia la superficie.",
				"Azul - Oleaje abierto: al revelarse se acerca a la superficie y, si gana, empuja un paso hacia mar abierto.",
				"Amarillo - Deriva neutra: en empate gana un paso hacia el escape horizontal.",
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
		clonedFishCards = append(clonedFishCards, cards.CloneFishCard(fishCard))
	}

	return clonedFishCards
}
