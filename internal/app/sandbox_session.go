package app

import (
	"fmt"
	"pesca/internal/content/fishprofiles"
	"pesca/internal/presentation"
	"strings"
)

type SandboxMenuUI interface {
	EncounterBootstrapUI
	UI
	ChooseSandboxMode(title string, modes []SandboxModeOption) (SandboxModeOption, error)
	ChooseSandboxScenario(title string, scenarios []SandboxScenario) (SandboxScenario, error)
	ChooseFishDeckPreset(title string, presets []fishprofiles.FishDeckPreset) (fishprofiles.FishDeckPreset, error)
	ShowNotice(message string) error
}

type SandboxModeOption struct {
	Mode        SandboxMode
	Name        string
	Description string
}

type SandboxSession struct {
	title     string
	rng       Randomizer
	ui        SandboxMenuUI
	presenter presentation.Presenter
	scenarios []SandboxScenario
}

func NewSandboxSession(title string, rng Randomizer, ui SandboxMenuUI, presenter presentation.Presenter) (*SandboxSession, error) {
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if rng == nil {
		return nil, fmt.Errorf("randomizer is required")
	}
	if ui == nil {
		return nil, fmt.Errorf("sandbox ui is required")
	}

	return &SandboxSession{
		title:     title,
		rng:       rng,
		ui:        ui,
		presenter: presenter,
		scenarios: DefaultSandboxScenarios(),
	}, nil
}

func DefaultSandboxModeOptions() []SandboxModeOption {
	return []SandboxModeOption{
		{Mode: SandboxModeGuided, Name: "Guiado", Description: "Flujo compatible con el duelo actual, con seleccion asistida y pez derivado por pool."},
		{Mode: SandboxModeManual, Name: "Manual", Description: "Seleccion explicita de presets base y preset de pez para reproducir encuentros concretos."},
		{Mode: SandboxModeScenario, Name: "Escenario", Description: "Replay semi-reproducible con seed fija y setup predefinido."},
	}
}

func (session *SandboxSession) Run() error {
	option, err := session.ui.ChooseSandboxMode(session.title, DefaultSandboxModeOptions())
	if err != nil {
		return fmt.Errorf("choose sandbox mode: %w", err)
	}

	return session.RunMode(option.Mode)
}

func (session *SandboxSession) RunMode(mode SandboxMode) error {
	switch mode {
	case SandboxModeGuided:
		engine, err := BootstrapEncounter(session.title, session.rng, session.ui)
		if err != nil {
			return fmt.Errorf("bootstrap guided sandbox: %w", err)
		}
		return session.runEngine(engine)
	case SandboxModeManual:
		return session.runManualMode()
	case SandboxModeScenario:
		return session.runScenarioSelection()
	default:
		return fmt.Errorf("unsupported sandbox mode %q", mode)
	}
}

func (session *SandboxSession) RunScenarioByID(id string) error {
	for _, scenario := range session.scenarios {
		if scenario.ID == id {
			return session.runScenario(scenario)
		}
	}

	return fmt.Errorf("unknown sandbox scenario %q", id)
}

var noFishProfileMatchErr = fmt.Errorf("no fish profile matches")

func (session *SandboxSession) runManualMode() error {
	fishPreset, err := session.ui.ChooseFishDeckPreset(session.title, fishprofiles.DefaultPresets())
	if err != nil {
		return fmt.Errorf("choose fish deck preset: %w", err)
	}

	for {
		engine, err := session.tryManualBootstrap(fishPreset)
		if err == nil {
			return session.runEngine(engine)
		}
		if !isNoFishProfileMatch(err) {
			return fmt.Errorf("bootstrap manual sandbox: %w", err)
		}
		if showErr := session.ui.ShowNotice("No esta picando nada. Intenta un nuevo lance."); showErr != nil {
			return fmt.Errorf("show notice: %w", showErr)
		}
	}
}

func (session *SandboxSession) tryManualBootstrap(fishPreset fishprofiles.FishDeckPreset) (Engine, error) {
	engine, err := BootstrapEncounterWithConfig(session.title, session.rng, session.ui, EncounterBootstrapConfig{FishPresetID: fishPreset.ID})
	if err != nil {
		if isNoFishProfileMatch(err) {
			return nil, noFishProfileMatchErr
		}
		return nil, err
	}
	return engine, nil
}

func isNoFishProfileMatch(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "no fish profile matches")
}

func (session *SandboxSession) runScenarioSelection() error {
	scenario, err := session.ui.ChooseSandboxScenario(session.title, session.scenarios)
	if err != nil {
		return fmt.Errorf("choose sandbox scenario: %w", err)
	}

	return session.runScenario(scenario)
}

func (session *SandboxSession) runScenario(scenario SandboxScenario) error {
	if err := scenario.Validate(); err != nil {
		return fmt.Errorf("scenario %s: %w", scenario.ID, err)
	}

	rng := session.rng
	if scenario.Seed != nil {
		rng = NewSeededRandom(*scenario.Seed)
	}
	engine, err := BootstrapEncounterWithConfig(session.title+" - "+scenario.Name, rng, session.ui, EncounterBootstrapConfig{
		PlayerDeckPresetID: scenario.PlayerDeckPresetID,
		RodPresetID:        "coastal-control",
		AttachmentPresetID: "no-attachments",
		WaterContextID:     scenario.WaterContextID,
		FishPresetID:       scenario.FishPresetID,
	})
	if err != nil {
		return fmt.Errorf("bootstrap scenario sandbox: %w", err)
	}

	return session.runEngine(engine)
}

func (session *SandboxSession) runEngine(engine Engine) error {
	encounterSession, err := NewSession(engine, session.ui, session.presenter)
	if err != nil {
		return fmt.Errorf("initialize sandbox session: %w", err)
	}

	return encounterSession.Run()
}
