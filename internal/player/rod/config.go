package rod

import "fmt"

type Config struct {
	OpeningMaxDistance int
	OpeningMaxDepth    int
	TrackMaxDistance   int
	TrackMaxDepth      int
}

type State struct {
	OpeningMaxDistance int
	OpeningMaxDepth    int
	TrackMaxDistance   int
	TrackMaxDepth      int
}

func (state State) Validate() error {
	return Config{
		OpeningMaxDistance: state.OpeningMaxDistance,
		OpeningMaxDepth:    state.OpeningMaxDepth,
		TrackMaxDistance:   state.TrackMaxDistance,
		TrackMaxDepth:      state.TrackMaxDepth,
	}.Validate()
}

func DefaultConfig() Config {
	return Config{
		OpeningMaxDistance: 5,
		OpeningMaxDepth:    3,
		TrackMaxDistance:   5,
		TrackMaxDepth:      3,
	}
}

func (config Config) Validate() error {
	if config.OpeningMaxDistance < 0 {
		return fmt.Errorf("opening max distance must be greater than or equal to 0")
	}
	if config.OpeningMaxDepth < 0 {
		return fmt.Errorf("opening max depth must be greater than or equal to 0")
	}
	if config.TrackMaxDistance <= 0 {
		return fmt.Errorf("track max distance must be greater than 0")
	}
	if config.TrackMaxDepth < 0 {
		return fmt.Errorf("track max depth must be greater than or equal to 0")
	}
	if config.OpeningMaxDistance > config.TrackMaxDistance {
		return fmt.Errorf("opening max distance must be less than or equal to track max distance")
	}
	if config.OpeningMaxDepth > config.TrackMaxDepth {
		return fmt.Errorf("opening max depth must be less than or equal to track max depth")
	}

	return nil
}

func NewState(config Config) (State, error) {
	if err := config.Validate(); err != nil {
		return State{}, err
	}

	return State{
		OpeningMaxDistance: config.OpeningMaxDistance,
		OpeningMaxDepth:    config.OpeningMaxDepth,
		TrackMaxDistance:   config.TrackMaxDistance,
		TrackMaxDepth:      config.TrackMaxDepth,
	}, nil
}
