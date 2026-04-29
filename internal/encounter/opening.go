package encounter

import "fmt"

type CastBand string

const (
	CastBandVeryShort CastBand = "very_short"
	CastBandShort     CastBand = "short"
	CastBandMedium    CastBand = "medium"
	CastBandLong      CastBand = "long"
	CastBandVeryLong  CastBand = "very_long"
)

type WaterContext struct {
	ID                  string
	Name                string
	Description         string
	VisibleSignals      []string
	PoolTag             string
	BandInitialDistance map[CastBand]int
	BaseInitialDepth    int
}

type CastResult struct {
	Band CastBand
}

type Opening struct {
	WaterContext    WaterContext
	CastResult      CastResult
	InitialDistance int
	InitialDepth    int
	Config          Config
}

func OrderedCastBands() []CastBand {
	return []CastBand{
		CastBandVeryShort,
		CastBandShort,
		CastBandMedium,
		CastBandLong,
		CastBandVeryLong,
	}
}

func (band CastBand) Label() string {
	switch band {
	case CastBandVeryShort:
		return "muy corto"
	case CastBandShort:
		return "corto"
	case CastBandMedium:
		return "medio"
	case CastBandLong:
		return "largo"
	case CastBandVeryLong:
		return "muy largo"
	default:
		return string(band)
	}
}

func (context WaterContext) Validate() error {
	if context.ID == "" {
		return fmt.Errorf("water context id is required")
	}
	if context.Name == "" {
		return fmt.Errorf("water context name is required")
	}
	if len(context.BandInitialDistance) != len(OrderedCastBands()) {
		return fmt.Errorf("water context must define an initial distance for every cast band")
	}

	for _, band := range OrderedCastBands() {
		initialDistance, ok := context.BandInitialDistance[band]
		if !ok {
			return fmt.Errorf("water context is missing an initial distance for cast band %s", band)
		}
		if initialDistance < 0 {
			return fmt.Errorf("water context initial distance for cast band %s must be greater than or equal to 0", band)
		}
	}
	if context.BaseInitialDepth < 0 {
		return fmt.Errorf("water context base initial depth must be greater than or equal to 0")
	}

	return nil
}

func (context WaterContext) InitialDistanceForBand(band CastBand) (int, error) {
	if err := context.Validate(); err != nil {
		return 0, err
	}

	initialDistance, ok := context.BandInitialDistance[band]
	if !ok {
		return 0, fmt.Errorf("water context is missing an initial distance for cast band %s", band)
	}

	return initialDistance, nil
}

func ResolveOpening(baseConfig Config, context WaterContext, castResult CastResult) (Opening, error) {
	if err := context.Validate(); err != nil {
		return Opening{}, err
	}

	initialDistance, err := context.InitialDistanceForBand(castResult.Band)
	if err != nil {
		return Opening{}, err
	}

	resolvedConfig := baseConfig
	resolvedConfig.InitialDistance = initialDistance
	resolvedConfig.InitialDepth = context.BaseInitialDepth
	if err := resolvedConfig.Validate(); err != nil {
		return Opening{}, err
	}

	return Opening{
		WaterContext:    context,
		CastResult:      castResult,
		InitialDistance: resolvedConfig.InitialDistance,
		InitialDepth:    resolvedConfig.InitialDepth,
		Config:          resolvedConfig,
	}, nil
}
