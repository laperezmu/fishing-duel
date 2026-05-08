package app

import (
	"fmt"
	"pesca/internal/content/fishprofiles"
	"pesca/internal/content/playerprofiles"
	"pesca/internal/content/watercontexts"
)

type SandboxMode string

const (
	SandboxModeGuided   SandboxMode = "guided"
	SandboxModeManual   SandboxMode = "manual"
	SandboxModeScenario SandboxMode = "scenario"
)

type SandboxCardOwner string

const (
	SandboxCardOwnerPlayer SandboxCardOwner = "player"
	SandboxCardOwnerFish   SandboxCardOwner = "fish"
)

type SandboxCardSelectionMode string

const (
	SandboxCardSelectionKeepBase SandboxCardSelectionMode = "keep_base"
	SandboxCardSelectionOverlay  SandboxCardSelectionMode = "overlay"
	SandboxCardSelectionScenario SandboxCardSelectionMode = "scenario"
)

type SandboxCardOrigin string

const (
	SandboxCardOriginPresetBase        SandboxCardOrigin = "preset_base"
	SandboxCardOriginManualReplacement SandboxCardOrigin = "manual_replacement"
	SandboxCardOriginScenarioDefined   SandboxCardOrigin = "scenario_defined"
)

type SandboxCardSelectionEntry struct {
	CardRef string
	Origin  SandboxCardOrigin
}

type SandboxCardSelection struct {
	OwnerScope     SandboxCardOwner
	SourcePresetID string
	SelectionMode  SandboxCardSelectionMode
	SelectedCards  []SandboxCardSelectionEntry
}

type SandboxStateOverrides struct {
	InitialDistance *int
	InitialDepth    *int
	CaptureDistance *int
	RecycleCount    *int
}

type SandboxScenario struct {
	ID                 string
	Name               string
	Description        string
	Seed               *int64
	PlayerDeckPresetID string
	FishPresetID       fishprofiles.ProfileID
	WaterContextID     watercontexts.ID
	CardSelections     []SandboxCardSelection
	StateOverrides     SandboxStateOverrides
}

type SandboxConfig struct {
	Mode               SandboxMode
	PlayerDeckPresetID string
	FishPresetID       fishprofiles.ProfileID
	WaterContextID     watercontexts.ID
	Seed               *int64
	ScenarioID         string
	CardSelections     []SandboxCardSelection
	StateOverrides     SandboxStateOverrides
}

func (config SandboxConfig) Validate() error {
	if config.Mode == "" {
		return fmt.Errorf("sandbox mode is required")
	}
	if config.Mode != SandboxModeGuided && config.Mode != SandboxModeManual && config.Mode != SandboxModeScenario {
		return fmt.Errorf("unsupported sandbox mode %q", config.Mode)
	}
	if config.Mode == SandboxModeScenario && config.ScenarioID == "" {
		return fmt.Errorf("scenario id is required in scenario mode")
	}
	for index, selection := range config.CardSelections {
		if err := selection.Validate(); err != nil {
			return fmt.Errorf("card selection %d: %w", index, err)
		}
	}

	return config.StateOverrides.Validate()
}

func (selection SandboxCardSelection) Validate() error {
	if selection.OwnerScope != SandboxCardOwnerPlayer && selection.OwnerScope != SandboxCardOwnerFish {
		return fmt.Errorf("unsupported card owner %q", selection.OwnerScope)
	}
	if selection.SourcePresetID == "" {
		return fmt.Errorf("source preset id is required")
	}
	if selection.SelectionMode == "" {
		return fmt.Errorf("selection mode is required")
	}
	if selection.SelectionMode != SandboxCardSelectionKeepBase && selection.SelectionMode != SandboxCardSelectionOverlay && selection.SelectionMode != SandboxCardSelectionScenario {
		return fmt.Errorf("unsupported selection mode %q", selection.SelectionMode)
	}
	seen := make(map[string]struct{}, len(selection.SelectedCards))
	for index, card := range selection.SelectedCards {
		if card.CardRef == "" {
			return fmt.Errorf("selected card %d: card ref is required", index)
		}
		if card.Origin == "" {
			return fmt.Errorf("selected card %d: card origin is required", index)
		}
		if _, ok := seen[card.CardRef]; ok {
			return fmt.Errorf("selected card %d: duplicated card ref %q", index, card.CardRef)
		}
		seen[card.CardRef] = struct{}{}
	}

	return nil
}

func (overrides SandboxStateOverrides) Validate() error {
	if overrides.InitialDistance != nil && *overrides.InitialDistance < 0 {
		return fmt.Errorf("initial distance override must be greater than or equal to 0")
	}
	if overrides.InitialDepth != nil && *overrides.InitialDepth < 0 {
		return fmt.Errorf("initial depth override must be greater than or equal to 0")
	}
	if overrides.CaptureDistance != nil && *overrides.CaptureDistance < 0 {
		return fmt.Errorf("capture distance override must be greater than or equal to 0")
	}
	if overrides.RecycleCount != nil && *overrides.RecycleCount < 0 {
		return fmt.Errorf("recycle count override must be greater than or equal to 0")
	}

	return nil
}

func (scenario SandboxScenario) Validate() error {
	if scenario.ID == "" {
		return fmt.Errorf("scenario id is required")
	}
	if scenario.Name == "" {
		return fmt.Errorf("scenario name is required")
	}
	config := SandboxConfig{
		Mode:               SandboxModeScenario,
		PlayerDeckPresetID: scenario.PlayerDeckPresetID,
		FishPresetID:       scenario.FishPresetID,
		WaterContextID:     scenario.WaterContextID,
		Seed:               scenario.Seed,
		ScenarioID:         scenario.ID,
		CardSelections:     scenario.CardSelections,
		StateOverrides:     scenario.StateOverrides,
	}

	return config.Validate()
}

func (config SandboxConfig) ResolvePlayerPreset() (playerprofiles.DeckPreset, error) {
	if config.PlayerDeckPresetID == "" {
		return playerprofiles.DeckPreset{}, fmt.Errorf("player deck preset id is required")
	}

	return playerprofiles.ResolveDefaultPreset(config.PlayerDeckPresetID)
}
