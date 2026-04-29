package playerprofiles

import (
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
	return []DeckPreset{
		{
			ID:          "classic",
			Name:        "Clasico",
			Description: "Tres cartas lisas por color, sin efectos, para reproducir la experiencia base.",
			Details: []string{
				"Azul: 3 cartas lisas sin efectos.",
				"Rojo: 3 cartas lisas sin efectos.",
				"Amarillo: 3 cartas lisas sin efectos.",
				"Recuperacion: 1 ronda tras vaciar una baraja de color.",
			},
			Config: playermoves.DefaultConfig(),
		},
		{
			ID:          "hooked-opening",
			Name:        "Apertura preparada",
			Description: "Barajas del jugador pensadas para abrir el round con ventajas temporales al revelar la carta.",
			Details: []string{
				"Azul - Anzuelo tenso: al revelar la carta permite capturar desde un paso mas lejos ese round.",
				"Rojo - Giro de superficie: al revelar la carta deja al pez contar como un nivel mas cerca de la superficie ese round.",
				"Amarillo - Reserva final: al revelar la carta amplia el margen de captura cuando se agota la baraja ese round.",
				"Objetivo: validar aperturas iniciales del jugador sin alterar la UX actual.",
			},
			Config: playermoves.Config{
				InitialDecks: map[domain.Move][]cards.PlayerCard{
					domain.Blue: {
						cards.NewNamedPlayerCard("Anzuelo tenso", "Capturas desde un paso mas lejos este round.", domain.Blue, cards.CardEffect{Trigger: cards.TriggerOnDraw, CaptureDistanceBonus: 1}),
						cards.NewPlayerCard(domain.Blue),
						cards.NewPlayerCard(domain.Blue),
					},
					domain.Red: {
						cards.NewNamedPlayerCard("Giro de superficie", "El pez cuenta como un nivel mas cerca de la superficie este round.", domain.Red, cards.CardEffect{Trigger: cards.TriggerOnDraw, SurfaceDepthBonus: 1}),
						cards.NewPlayerCard(domain.Red),
						cards.NewPlayerCard(domain.Red),
					},
					domain.Yellow: {
						cards.NewNamedPlayerCard("Reserva final", "Amplia el margen de captura por agotamiento durante este round.", domain.Yellow, cards.CardEffect{Trigger: cards.TriggerOnDraw, ExhaustionCaptureDistanceBonus: 1}),
						cards.NewPlayerCard(domain.Yellow),
						cards.NewPlayerCard(domain.Yellow),
					},
				},
				RecoveryDelayRounds: 1,
			},
		},
		{
			ID:          "vertical-pressure",
			Name:        "Respuesta vertical",
			Description: "Cartas del jugador que castigan resultados con cambios de profundidad.",
			Details: []string{
				"Azul - Tiron certero: si ganas, subes al pez un nivel hacia la superficie.",
				"Rojo - Recobro paciente: si pierdes, aun asi haces subir al pez un nivel.",
				"Amarillo - Tregua corta: en empate, fuerzas al pez a subir un nivel.",
				"Objetivo: validar respuestas verticales del jugador segun el outcome.",
			},
			Config: playermoves.Config{
				InitialDecks: map[domain.Move][]cards.PlayerCard{
					domain.Blue: {
						cards.NewNamedPlayerCard("Tiron certero", "Si ganas, subes al pez un nivel hacia la superficie.", domain.Blue, cards.CardEffect{Trigger: cards.TriggerOnOwnerWin, DepthShift: -1}),
						cards.NewPlayerCard(domain.Blue),
						cards.NewPlayerCard(domain.Blue),
					},
					domain.Red: {
						cards.NewNamedPlayerCard("Recobro paciente", "Si pierdes, haces subir al pez un nivel de todos modos.", domain.Red, cards.CardEffect{Trigger: cards.TriggerOnOwnerLose, DepthShift: -1}),
						cards.NewPlayerCard(domain.Red),
						cards.NewPlayerCard(domain.Red),
					},
					domain.Yellow: {
						cards.NewNamedPlayerCard("Tregua corta", "Si empatas, haces subir al pez un nivel.", domain.Yellow, cards.CardEffect{Trigger: cards.TriggerOnRoundDraw, DepthShift: -1}),
						cards.NewPlayerCard(domain.Yellow),
						cards.NewPlayerCard(domain.Yellow),
					},
				},
				RecoveryDelayRounds: 1,
			},
		},
		{
			ID:          "mixed-current",
			Name:        "Corriente mixta",
			Description: "Barajas del jugador que mezclan ventajas al revelarse con respuestas segun el resultado.",
			Details: []string{
				"Azul - Carrera corta: al revelarse acerca la captura y, si ganas, arrastra al pez un paso mas hacia la orilla.",
				"Rojo - Rescate profundo: si pierdes, haces subir al pez un nivel hacia la superficie.",
				"Amarillo - Calma tensa: en empate el pez cuenta como un nivel mas cerca de la superficie ese round.",
				"Objetivo: validar convivencia entre fases y triggers del jugador.",
			},
			Config: playermoves.Config{
				InitialDecks: map[domain.Move][]cards.PlayerCard{
					domain.Blue: {
						cards.NewNamedPlayerCard("Carrera corta", "Acerca la captura al revelarse y tira del pez si ganas.", domain.Blue,
							cards.CardEffect{Trigger: cards.TriggerOnDraw, CaptureDistanceBonus: 1},
							cards.CardEffect{Trigger: cards.TriggerOnOwnerWin, DistanceShift: -1},
						),
						cards.NewPlayerCard(domain.Blue),
						cards.NewPlayerCard(domain.Blue),
					},
					domain.Red: {
						cards.NewNamedPlayerCard("Rescate profundo", "Si pierdes, haces subir al pez un nivel.", domain.Red,
							cards.CardEffect{Trigger: cards.TriggerOnOwnerLose, DepthShift: -1},
						),
						cards.NewPlayerCard(domain.Red),
						cards.NewPlayerCard(domain.Red),
					},
					domain.Yellow: {
						cards.NewNamedPlayerCard("Calma tensa", "En empate el pez cuenta como un nivel mas cerca de la superficie.", domain.Yellow,
							cards.CardEffect{Trigger: cards.TriggerOnRoundDraw, SurfaceDepthBonus: 1},
						),
						cards.NewPlayerCard(domain.Yellow),
						cards.NewPlayerCard(domain.Yellow),
					},
				},
				RecoveryDelayRounds: 1,
			},
		},
	}
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
