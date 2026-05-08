package watercontexts

import (
	"fmt"
	"pesca/internal/content/waterpools"
	"pesca/internal/encounter"
)

type ID string

const (
	ShorelineCove ID = "shoreline-cove"
	OpenChannel   ID = "open-channel"
	BrokenCurrent ID = "broken-current"
	ReefShadow    ID = "reef-shadow"
	TidalGate     ID = "tidal-gate"
	WeedPocket    ID = "weed-pocket"
	StoneDrop     ID = "stone-drop"
	WindLane      ID = "wind-lane"
	DeepLedge     ID = "deep-ledge"
)

type Preset struct {
	ID           ID
	Name         string
	Description  string
	Signals      []string
	PoolTag      waterpools.ID
	Distances    map[encounter.CastBand]int
	InitialDepth int
}

func (preset Preset) BuildContext() encounter.WaterContext {
	bandDistances := make(map[encounter.CastBand]int, len(preset.Distances))
	for band, initialDistance := range preset.Distances {
		bandDistances[band] = initialDistance
	}

	return encounter.WaterContext{
		ID:                  string(preset.ID),
		Name:                preset.Name,
		Description:         preset.Description,
		VisibleSignals:      append([]string(nil), preset.Signals...),
		PoolTag:             preset.PoolTag,
		BandInitialDistance: bandDistances,
		BaseInitialDepth:    preset.InitialDepth,
	}
}

func DefaultPresets() []Preset {
	presets := make([]Preset, 0, len(defaultPresets))
	for _, preset := range defaultPresets {
		presets = append(presets, clonePreset(preset))
	}

	return presets
}

func ResolveDefaultPreset(id ID) (Preset, error) {
	preset, ok := defaultPresetByID[id]
	if !ok {
		return Preset{}, fmt.Errorf("unknown water context preset %q", id)
	}

	return clonePreset(preset), nil
}

func DefaultPhaseLabel(id ID) string {
	preset, ok := defaultPresetByID[id]
	if !ok {
		return string(id)
	}
	phaseIndex, ok := defaultPhaseIndex[id]
	if !ok {
		return preset.Name
	}

	return fmt.Sprintf("Fase %d - %s", phaseIndex+1, preset.Name)
}

func consecutiveDistances(start int) map[encounter.CastBand]int {
	bands := encounter.OrderedCastBands()
	distances := make(map[encounter.CastBand]int, len(bands))
	for offset, band := range bands {
		distances[band] = start + offset
	}

	return distances
}

func clonePreset(preset Preset) Preset {
	clonedPreset := preset
	clonedPreset.Signals = append([]string(nil), preset.Signals...)
	clonedPreset.Distances = make(map[encounter.CastBand]int, len(preset.Distances))
	for band, distance := range preset.Distances {
		clonedPreset.Distances[band] = distance
	}

	return clonedPreset
}

func buildDefaultPresetIndex(presets []Preset) map[ID]Preset {
	index := make(map[ID]Preset, len(presets))
	for _, preset := range presets {
		index[preset.ID] = clonePreset(preset)
	}

	return index
}

func buildDefaultPhaseIndex(order []ID) map[ID]int {
	index := make(map[ID]int, len(order))
	for i, id := range order {
		index[id] = i
	}

	return index
}
