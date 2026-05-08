package watercontexts

import _ "embed"

var (
	//go:embed data/default_presets.json
	defaultPresetsJSON []byte
	defaultPresets     = loadDefaultPresets()
	defaultPresetByID  = buildDefaultPresetIndex(defaultPresets)
	defaultPhaseIndex  = loadDefaultPhaseIndex()
)

func loadDefaultPresets() []Preset {
	presets, phaseIndex, err := LoadPresets(defaultPresetsJSON)
	if err != nil {
		panic(err)
	}

	defaultPhaseIndex = phaseIndex
	return presets
}

func loadDefaultPhaseIndex() map[ID]int {
	_, phaseIndex, err := LoadPresets(defaultPresetsJSON)
	if err != nil {
		panic(err)
	}

	return phaseIndex
}
