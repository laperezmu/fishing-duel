package playerrig

import "fmt"

type Config struct {
	MaxDistance int
	MaxDepth    int
}

type State struct {
	MaxDistance int
	MaxDepth    int
}

func DefaultConfig() Config {
	return Config{
		MaxDistance: 5,
		MaxDepth:    3,
	}
}

func (config Config) Validate() error {
	if config.MaxDistance <= 0 {
		return fmt.Errorf("max distance must be greater than 0")
	}
	if config.MaxDepth < 0 {
		return fmt.Errorf("max depth must be greater than or equal to 0")
	}

	return nil
}

func NewState(config Config) (State, error) {
	if err := config.Validate(); err != nil {
		return State{}, err
	}

	return State{
		MaxDistance: config.MaxDistance,
		MaxDepth:    config.MaxDepth,
	}, nil
}
