package app

import (
	"fmt"
	"pesca/internal/content/watercontexts"
	"pesca/internal/encounter"
	"pesca/internal/player/loadout"
	"pesca/internal/presentation"
)

type OpeningUI interface {
	ChooseWaterContext(title string, presets []watercontexts.Preset) (watercontexts.Preset, error)
	ResolveCast(title string, context encounter.WaterContext, presenter CastPresenter) (encounter.CastResult, error)
	ShowEncounterOpening(title string, opening presentation.OpeningView) error
}

type OpeningRuntimeUI interface {
	ResolveCast(title string, context encounter.WaterContext, presenter CastPresenter) (encounter.CastResult, error)
	ShowEncounterOpening(title string, opening presentation.OpeningView) error
}

type OpeningPresenter interface {
	CastPresenter
	Opening(opening encounter.Opening) presentation.OpeningView
}

type CastPresenter interface {
	Cast(context encounter.WaterContext, position, totalSlots, sectionWidth int) presentation.CastView
}

func ResolveEncounterOpening(title string, baseConfig encounter.Config, playerLoadout loadout.State, presets []watercontexts.Preset, ui OpeningUI, presenter OpeningPresenter) (encounter.Opening, error) {
	if ui == nil {
		return encounter.Opening{}, fmt.Errorf("opening ui is required")
	}
	if presenter == nil {
		return encounter.Opening{}, fmt.Errorf("opening presenter is required")
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

	return ResolveEncounterOpeningWithPreset(title, baseConfig, playerLoadout, selectedPreset, ui, presenter)
}

func ResolveEncounterOpeningWithPreset(title string, baseConfig encounter.Config, playerLoadout loadout.State, preset watercontexts.Preset, ui OpeningRuntimeUI, presenter OpeningPresenter) (encounter.Opening, error) {
	if ui == nil {
		return encounter.Opening{}, fmt.Errorf("opening ui is required")
	}
	if presenter == nil {
		return encounter.Opening{}, fmt.Errorf("opening presenter is required")
	}
	if err := playerLoadout.Validate(); err != nil {
		return encounter.Opening{}, fmt.Errorf("player loadout: %w", err)
	}

	waterContext := preset.BuildContext()
	castResult, err := ui.ResolveCast(title, waterContext, presenter)
	if err != nil {
		return encounter.Opening{}, fmt.Errorf("resolve cast: %w", err)
	}

	opening, err := encounter.ResolveOpening(baseConfig, waterContext, castResult, playerLoadout.OpeningLimits())
	if err != nil {
		return encounter.Opening{}, fmt.Errorf("resolve encounter opening: %w", err)
	}

	if err := ui.ShowEncounterOpening(title, presenter.Opening(opening)); err != nil {
		return encounter.Opening{}, fmt.Errorf("show encounter opening: %w", err)
	}

	return opening, nil
}
