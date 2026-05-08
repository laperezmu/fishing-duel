package attachmentpresets

import (
	"fmt"
	"pesca/internal/content/habitats"
	"pesca/internal/player/loadout"
)

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
