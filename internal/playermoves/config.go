package playermoves

import (
	"fmt"
	"pesca/internal/domain"
)

type Config struct {
	InitialUsesPerMove  map[domain.Move]int
	RecoveryDelayRounds int
}

func DefaultConfig() Config {
	return Config{
		InitialUsesPerMove: map[domain.Move]int{
			domain.Blue:   3,
			domain.Red:    3,
			domain.Yellow: 3,
		},
		RecoveryDelayRounds: 1,
	}
}

func (config Config) Validate() error {
	if config.RecoveryDelayRounds < 0 {
		return fmt.Errorf("recovery delay rounds must be greater than or equal to 0")
	}

	for _, move := range supportedMoves() {
		initialUses, ok := config.InitialUsesPerMove[move]
		if !ok {
			return fmt.Errorf("initial uses for move %s are required", move)
		}
		if initialUses < 0 {
			return fmt.Errorf("initial uses for move %s must be greater than or equal to 0", move)
		}
	}

	return nil
}

func (config Config) initialUsesFor(move domain.Move) int {
	return config.InitialUsesPerMove[move]
}

func supportedMoves() []domain.Move {
	return []domain.Move{domain.Blue, domain.Red, domain.Yellow}
}
