package attachmentpresets

import (
	"fmt"
	"pesca/internal/content/habitats"
	"pesca/internal/player/loadout"
)

var defaultPresets = []Preset{
	{
		ID:          "no-attachments",
		Name:        "Sin aditamentos",
		Description: "Usa solo la rod base, sin sesgos extra en apertura o track.",
		Details: []string{
			"No modifica distancia ni profundidad.",
			"Sirve como referencia limpia para comparar rods.",
		},
	},
	{
		ID:          "bottom-kit",
		Name:        "Kit de fondo",
		Description: "Empuja la pesca hacia capas profundas a cambio de recortar algo de alcance horizontal.",
		Details: []string{
			"Apertura: -1 distancia, +1 profundidad.",
			"Track: +1 profundidad.",
			"Habitats: fondo, canal.",
		},
		Attachments: []loadout.Attachment{{
			ID:                      "sinker-heavy",
			Name:                    "Plomada pesada",
			Description:             "Ayuda a abrir la pesca mas abajo y sostener profundidad.",
			OpeningDistanceModifier: -1,
			OpeningDepthModifier:    1,
			TrackDepthModifier:      1,
			HabitatTags:             []habitats.Tag{habitats.Bottom, habitats.Channel},
		}},
	},
	{
		ID:          "long-cast-kit",
		Name:        "Kit de lance largo",
		Description: "Gana margen horizontal para abrir y sostener distancia, pero pierde acceso vertical temprano.",
		Details: []string{
			"Apertura: +1 distancia, -1 profundidad.",
			"Track: +1 distancia.",
			"Habitats: costa abierta, superficie.",
		},
		Attachments: []loadout.Attachment{{
			ID:                      "spool-long-cast",
			Name:                    "Bobina de lance",
			Description:             "Premia aguas abiertas con mas recorrido horizontal.",
			OpeningDistanceModifier: 1,
			OpeningDepthModifier:    -1,
			TrackDistanceModifier:   1,
			HabitatTags:             []habitats.Tag{habitats.OpenWater, habitats.Surface},
		}},
	},
	{
		ID:          "stability-kit",
		Name:        "Kit de estabilidad",
		Description: "No abre mas lejos, pero amplia el margen defensivo del track en ambos ejes.",
		Details: []string{
			"Apertura: sin cambios.",
			"Track: +1 distancia, +1 profundidad.",
			"Habitats: maleza, roca.",
		},
		Attachments: []loadout.Attachment{{
			ID:                    "line-reinforced",
			Name:                  "Linea reforzada",
			Description:           "Aguanta mejor la presion del pez en el tablero.",
			TrackDistanceModifier: 1,
			TrackDepthModifier:    1,
			HabitatTags:           []habitats.Tag{habitats.Weed, habitats.Rock},
		}},
	},
}

var defaultPresetByID = buildDefaultPresetIndex(defaultPresets)

type Preset struct {
	ID          string
	Name        string
	Description string
	Details     []string
	Attachments []loadout.Attachment
}

func (preset Preset) BuildAttachments() []loadout.Attachment {
	attachments := make([]loadout.Attachment, 0, len(preset.Attachments))
	for _, attachment := range preset.Attachments {
		clonedAttachment := attachment
		clonedAttachment.HabitatTags = append([]habitats.Tag(nil), attachment.HabitatTags...)
		attachments = append(attachments, clonedAttachment)
	}

	return attachments
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
		return Preset{}, fmt.Errorf("unknown attachment preset %q", id)
	}

	return clonePreset(preset), nil
}

func clonePreset(preset Preset) Preset {
	clonedPreset := preset
	clonedPreset.Details = append([]string(nil), preset.Details...)
	clonedPreset.Attachments = preset.BuildAttachments()

	return clonedPreset
}

func buildDefaultPresetIndex(presets []Preset) map[string]Preset {
	index := make(map[string]Preset, len(presets))
	for _, preset := range presets {
		index[preset.ID] = clonePreset(preset)
	}

	return index
}
