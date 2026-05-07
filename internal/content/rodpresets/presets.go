package rodpresets

import (
	"fmt"
	"pesca/internal/player/loadout"
	"pesca/internal/player/rod"
)

var defaultPresets = []Preset{
	{
		ID:          "coastal-control",
		Name:        "Control costero",
		Description: "Abre cerca y con poca profundidad, pero deja margen para sostener el duelo en el tablero.",
		Details: []string{
			"Apertura: distancia 3, profundidad 1.",
			"Track: distancia 5, profundidad 2.",
			"Pensada para aguas cerradas, lectura cercana y cierres rapidos.",
		},
		Config: rod.Config{
			OpeningMaxDistance: 3,
			OpeningMaxDepth:    1,
			TrackMaxDistance:   5,
			TrackMaxDepth:      2,
		},
	},
	{
		ID:          "versatile-standard",
		Name:        "Versatil estandar",
		Description: "Preset equilibrado para probar el loop base sin sesgo fuerte hacia costa o fondo.",
		Details: []string{
			"Apertura: distancia 5, profundidad 3.",
			"Track: distancia 5, profundidad 3.",
			"Sirve como referencia general del sistema.",
		},
		Config: rod.DefaultConfig(),
	},
	{
		ID:          "bottom-pressure",
		Name:        "Presion de fondo",
		Description: "Pierde algo de margen horizontal, pero abre y sostiene mejor pescas profundas.",
		Details: []string{
			"Apertura: distancia 4, profundidad 4.",
			"Track: distancia 4, profundidad 5.",
			"Pensada para priorizar capas profundas antes que mar abierto lejano.",
		},
		Config: rod.Config{
			OpeningMaxDistance: 4,
			OpeningMaxDepth:    4,
			TrackMaxDistance:   4,
			TrackMaxDepth:      5,
		},
	},
}

var defaultPresetByID = buildDefaultPresetIndex(defaultPresets)

type Preset struct {
	ID          string
	Name        string
	Description string
	Details     []string
	Config      rod.Config
}

func (preset Preset) BuildRod() (rod.State, error) {
	return rod.NewState(preset.Config)
}

func (preset Preset) BuildLoadout() (loadout.State, error) {
	return preset.BuildLoadoutWithAttachments(nil)
}

func (preset Preset) BuildLoadoutWithAttachments(attachments []loadout.Attachment) (loadout.State, error) {
	playerRod, err := preset.BuildRod()
	if err != nil {
		return loadout.State{}, err
	}

	return loadout.NewState(playerRod, attachments)
}

func DefaultPresets() []Preset {
	presets := make([]Preset, 0, len(defaultPresets))
	for _, preset := range defaultPresets {
		presets = append(presets, clonePreset(preset))
	}

	return presets
}

func ResolveDefaultPreset(id string) (Preset, error) {
	preset, ok := defaultPresetByID[id]
	if !ok {
		return Preset{}, fmt.Errorf("unknown rod preset %q", id)
	}

	return clonePreset(preset), nil
}

func clonePreset(preset Preset) Preset {
	clonedPreset := preset
	clonedPreset.Details = append([]string(nil), preset.Details...)

	return clonedPreset
}

func buildDefaultPresetIndex(presets []Preset) map[string]Preset {
	index := make(map[string]Preset, len(presets))
	for _, preset := range presets {
		index[preset.ID] = clonePreset(preset)
	}

	return index
}
