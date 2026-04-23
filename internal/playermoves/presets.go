package playermoves

import (
	"pesca/internal/cards"
	"pesca/internal/domain"
)

type PlayerDeckPreset struct {
	ID          string
	Name        string
	Description string
	Details     []string
	Config      Config
}

func (preset PlayerDeckPreset) BuildConfig(shuffler func([]cards.PlayerCard)) Config {
	config := Config{
		InitialDecks:        cloneInitialDecks(preset.Config.InitialDecks),
		DeckShuffler:        shuffler,
		RecoveryDelayRounds: preset.Config.RecoveryDelayRounds,
	}

	return config
}

func DefaultPlayerDeckPresets() []PlayerDeckPreset {
	return []PlayerDeckPreset{
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
			Config: DefaultConfig(),
		},
		{
			ID:          "hooked-opening",
			Name:        "Apertura preparada",
			Description: "Cartas del jugador con `on_draw` para manipular thresholds temporales del round.",
			Details: []string{
				"Azul: primera carta `on_draw capt +1`, luego 2 cartas lisas.",
				"Rojo: primera carta `on_draw sup +1`, luego 2 cartas lisas.",
				"Amarillo: primera carta `on_draw baraja +1`, luego 2 cartas lisas.",
				"Objetivo: validar aperturas tacticas del jugador sin alterar la UX actual.",
			},
			Config: Config{
				InitialDecks: map[domain.Move][]cards.PlayerCard{
					domain.Blue: {
						cards.NewPlayerCard(domain.Blue, cards.CardEffect{Trigger: cards.TriggerOnDraw, CaptureDistanceBonus: 1}),
						cards.NewPlayerCard(domain.Blue),
						cards.NewPlayerCard(domain.Blue),
					},
					domain.Red: {
						cards.NewPlayerCard(domain.Red, cards.CardEffect{Trigger: cards.TriggerOnDraw, SurfaceDepthBonus: 1}),
						cards.NewPlayerCard(domain.Red),
						cards.NewPlayerCard(domain.Red),
					},
					domain.Yellow: {
						cards.NewPlayerCard(domain.Yellow, cards.CardEffect{Trigger: cards.TriggerOnDraw, ExhaustionCaptureDistanceBonus: 1}),
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
				"Azul: primera carta `si gana prof -1`, luego 2 cartas lisas.",
				"Rojo: primera carta `si pierde prof -1`, luego 2 cartas lisas.",
				"Amarillo: primera carta `empate prof -1`, luego 2 cartas lisas.",
				"Objetivo: validar respuestas verticales del jugador segun el outcome.",
			},
			Config: Config{
				InitialDecks: map[domain.Move][]cards.PlayerCard{
					domain.Blue: {
						cards.NewPlayerCard(domain.Blue, cards.CardEffect{Trigger: cards.TriggerOnOwnerWin, DepthShift: -1}),
						cards.NewPlayerCard(domain.Blue),
						cards.NewPlayerCard(domain.Blue),
					},
					domain.Red: {
						cards.NewPlayerCard(domain.Red, cards.CardEffect{Trigger: cards.TriggerOnOwnerLose, DepthShift: -1}),
						cards.NewPlayerCard(domain.Red),
						cards.NewPlayerCard(domain.Red),
					},
					domain.Yellow: {
						cards.NewPlayerCard(domain.Yellow, cards.CardEffect{Trigger: cards.TriggerOnRoundDraw, DepthShift: -1}),
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
			Description: "Barajas del jugador con efectos en draw y post-outcome para probar el pipeline completo.",
			Details: []string{
				"Azul: primera carta `draw capt +1` y `si gana dist -1`, luego 2 cartas lisas.",
				"Rojo: primera carta `si pierde prof -1`, luego 2 cartas lisas.",
				"Amarillo: primera carta `empate sup +1`, luego 2 cartas lisas.",
				"Objetivo: validar convivencia entre fases y triggers del jugador.",
			},
			Config: Config{
				InitialDecks: map[domain.Move][]cards.PlayerCard{
					domain.Blue: {
						cards.NewPlayerCard(domain.Blue,
							cards.CardEffect{Trigger: cards.TriggerOnDraw, CaptureDistanceBonus: 1},
							cards.CardEffect{Trigger: cards.TriggerOnOwnerWin, DistanceShift: -1},
						),
						cards.NewPlayerCard(domain.Blue),
						cards.NewPlayerCard(domain.Blue),
					},
					domain.Red: {
						cards.NewPlayerCard(domain.Red,
							cards.CardEffect{Trigger: cards.TriggerOnOwnerLose, DepthShift: -1},
						),
						cards.NewPlayerCard(domain.Red),
						cards.NewPlayerCard(domain.Red),
					},
					domain.Yellow: {
						cards.NewPlayerCard(domain.Yellow,
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
			clonedDeck = append(clonedDeck, cards.NewPlayerCard(playerCard.Move, playerCard.Effects...))
		}
		clonedDecks[move] = clonedDeck
	}

	return clonedDecks
}
