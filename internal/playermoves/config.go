package playermoves

import (
	"fmt"
	"pesca/internal/cards"
	"pesca/internal/domain"
)

type Config struct {
	InitialDecks        map[domain.Move][]cards.PlayerCard
	DeckShuffler        func([]cards.PlayerCard)
	RecoveryDelayRounds int
}

func DefaultConfig() Config {
	return Config{
		InitialDecks: map[domain.Move][]cards.PlayerCard{
			domain.Blue: {
				cards.NewPlayerCard(domain.Blue),
				cards.NewPlayerCard(domain.Blue),
				cards.NewPlayerCard(domain.Blue),
			},
			domain.Red: {
				cards.NewPlayerCard(domain.Red),
				cards.NewPlayerCard(domain.Red),
				cards.NewPlayerCard(domain.Red),
			},
			domain.Yellow: {
				cards.NewPlayerCard(domain.Yellow),
				cards.NewPlayerCard(domain.Yellow),
				cards.NewPlayerCard(domain.Yellow),
			},
		},
		RecoveryDelayRounds: 1,
	}
}

func (config Config) Validate() error {
	if config.RecoveryDelayRounds < 0 {
		return fmt.Errorf("recovery delay rounds must be greater than or equal to 0")
	}

	for _, move := range supportedMoves() {
		initialDeck, ok := config.InitialDecks[move]
		if !ok {
			return fmt.Errorf("initial deck for move %s is required", move)
		}
		for _, playerCard := range initialDeck {
			if playerCard.Move != move {
				return fmt.Errorf("initial deck for move %s contains a card for move %s", move, playerCard.Move)
			}
		}
		if len(initialDeck) == 0 {
			return fmt.Errorf("initial deck for move %s must contain at least one card", move)
		}
	}

	return nil
}

func (config Config) initialDeckFor(move domain.Move) []cards.PlayerCard {
	configuredDeck := config.InitialDecks[move]
	clonedDeck := make([]cards.PlayerCard, 0, len(configuredDeck))
	for _, playerCard := range configuredDeck {
		clonedDeck = append(clonedDeck, cards.NewPlayerCard(playerCard.Move, playerCard.Effects...))
	}

	return clonedDeck
}

func supportedMoves() []domain.Move {
	return []domain.Move{domain.Blue, domain.Red, domain.Yellow}
}
