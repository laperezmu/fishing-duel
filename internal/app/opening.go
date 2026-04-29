package app

import (
	"fmt"
	"pesca/internal/content/watercontexts"
	"pesca/internal/encounter"
	"pesca/internal/player/loadout"
)

type OpeningUI interface {
	ChooseWaterContext(title string, presets []watercontexts.Preset) (watercontexts.Preset, error)
	ResolveCast(title string, context encounter.WaterContext) (encounter.CastResult, error)
	ShowEncounterOpening(title string, opening encounter.Opening) error
}

func ResolveEncounterOpening(title string, baseConfig encounter.Config, playerLoadout loadout.State, presets []watercontexts.Preset, ui OpeningUI) (encounter.Opening, error) {
	if ui == nil {
		return encounter.Opening{}, fmt.Errorf("opening ui is required")
	}
	if err := playerLoadout.Validate(); err != nil {
		return encounter.Opening{}, fmt.Errorf("player loadout: %w", err)
	}
	if len(presets) == 0 {
		return encounter.Opening{}, fmt.Errorf("at least one water context preset is required")
	}

	selectedPreset, err := ui.ChooseWaterContext(title, presets)
	if err != nil {
		return encounter.Opening{}, fmt.Errorf("choose water context: %w", err)
	}

	waterContext := selectedPreset.BuildContext()
	castResult, err := ui.ResolveCast(title, waterContext)
	if err != nil {
		return encounter.Opening{}, fmt.Errorf("resolve cast: %w", err)
	}

	opening, err := encounter.ResolveOpening(baseConfig, waterContext, castResult, playerLoadout.OpeningLimits())
	if err != nil {
		return encounter.Opening{}, fmt.Errorf("resolve encounter opening: %w", err)
	}

	if err := ui.ShowEncounterOpening(title, opening); err != nil {
		return encounter.Opening{}, fmt.Errorf("show encounter opening: %w", err)
	}

	return opening, nil
}
