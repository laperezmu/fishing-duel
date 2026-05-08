package rodpresets

import _ "embed"

var (
	//go:embed data/default_presets.json
	defaultPresetsJSON []byte
	defaultPresets     = loadDefaultPresets()
	defaultPresetByID  = buildDefaultPresetIndex(defaultPresets)
)

func loadDefaultPresets() []Preset {
	presets, err := LoadPresets(defaultPresetsJSON)
	if err != nil {
		panic(err)
	}

	return presets
}
