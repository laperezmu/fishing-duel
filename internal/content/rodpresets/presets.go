package rodpresets

import (
	"fmt"
	"pesca/internal/player/loadout"
	"pesca/internal/player/rod"
)

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
