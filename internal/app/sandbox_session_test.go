package app_test

import (
	"errors"
	"testing"
	"testing/iotest"

	"pesca/internal/app"
	"pesca/internal/content/attachmentpresets"
	"pesca/internal/content/fishprofiles"
	"pesca/internal/content/playerprofiles"
	"pesca/internal/content/rodpresets"
	"pesca/internal/content/watercontexts"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/match"
	"pesca/internal/player/loadout"
	"pesca/internal/player/rod"
	"pesca/internal/presentation"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSandboxSessionRunModeManual(t *testing.T) {
	ui := &mockSandboxUI{}
	presenter := presentation.NewPresenter(presentation.DefaultCatalog())
	session, err := app.NewSandboxSession("Sandbox", app.NewSeededRandom(7), ui, presenter)
	require.NoError(t, err)

	deckPreset := playerprofiles.DefaultPresets()[0]
	rodPreset := rodpresets.DefaultPresets()[0]
	attachmentPreset := attachmentpresets.DefaultPresets()[0]
	waterPreset := watercontexts.DefaultPresets()[0]
	fishPreset := fishprofiles.DefaultPresets()[0]

	ui.On("ChooseFishDeckPreset", "Sandbox", mock.Anything).Return(fishPreset, nil).Once()
	ui.On("ChoosePlayerDeckPreset", "Sandbox", mock.Anything).Return(deckPreset, nil).Once()
	ui.On("ChooseRodPreset", "Sandbox", mock.Anything).Return(rodPreset, nil).Once()
	ui.On("ChooseAttachmentPreset", "Sandbox", mock.Anything, mock.Anything).Return(attachmentPreset, nil).Once()
	ui.On("ChooseWaterContext", "Sandbox", mock.Anything).Return(waterPreset, nil).Once()
	ui.On("ResolveCast", "Sandbox", mock.Anything, mock.Anything).Return(encounter.CastResult{Band: encounter.CastBandShort}, nil).Once()
	ui.On("ShowEncounterOpening", "Sandbox", mock.Anything).Return(nil).Once()
	ui.On("ShowFishSpawn", "Sandbox", mock.Anything).Return(nil).Once()
	ui.On("ShowIntro", mock.Anything).Return(nil).Once()
	ui.On("ChooseMove", mock.Anything, mock.Anything).Return(domain.Blue, iotest.ErrTimeout).Maybe()

	err = session.RunMode(app.SandboxModeManual)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "choose move")
	ui.AssertExpectations(t)
}

func TestSandboxSessionRunScenarioByID(t *testing.T) {
	ui := &mockSandboxUI{}
	presenter := presentation.NewPresenter(presentation.DefaultCatalog())
	session, err := app.NewSandboxSession("Sandbox", app.NewSeededRandom(7), ui, presenter)
	require.NoError(t, err)

	ui.On("ResolveCast", "Sandbox - Control de superficie en ensenada", mock.Anything, mock.Anything).Return(encounter.CastResult{Band: encounter.CastBandShort}, nil).Once()
	ui.On("ShowEncounterOpening", "Sandbox - Control de superficie en ensenada", mock.Anything).Return(nil).Once()
	ui.On("ShowFishSpawn", "Sandbox - Control de superficie en ensenada", mock.Anything).Return(nil).Once()
	ui.On("ShowIntro", mock.Anything).Return(nil).Once()
	ui.On("ChooseMove", mock.Anything, mock.Anything).Return(domain.Blue, errors.New("stop after boot")).Once()

	err = session.RunScenarioByID("surface-control-shoreline")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "choose move")
	ui.AssertExpectations(t)
}

type mockSandboxUI struct{ mock.Mock }

func (ui *mockSandboxUI) ChooseSandboxMode(title string, modes []app.SandboxModeOption) (app.SandboxModeOption, error) {
	args := ui.Called(title, modes)
	return args.Get(0).(app.SandboxModeOption), args.Error(1)
}

func (ui *mockSandboxUI) ChooseSandboxScenario(title string, scenarios []app.SandboxScenario) (app.SandboxScenario, error) {
	args := ui.Called(title, scenarios)
	return args.Get(0).(app.SandboxScenario), args.Error(1)
}

func (ui *mockSandboxUI) ChooseFishDeckPreset(title string, presets []fishprofiles.FishDeckPreset) (fishprofiles.FishDeckPreset, error) {
	args := ui.Called(title, presets)
	return args.Get(0).(fishprofiles.FishDeckPreset), args.Error(1)
}

func (ui *mockSandboxUI) ChoosePlayerDeckPreset(title string, presets []playerprofiles.DeckPreset) (playerprofiles.DeckPreset, error) {
	args := ui.Called(title, presets)
	return args.Get(0).(playerprofiles.DeckPreset), args.Error(1)
}

func (ui *mockSandboxUI) ChooseRodPreset(title string, presets []rodpresets.Preset) (rodpresets.Preset, error) {
	args := ui.Called(title, presets)
	return args.Get(0).(rodpresets.Preset), args.Error(1)
}

func (ui *mockSandboxUI) ChooseAttachmentPreset(title string, baseRod rod.State, presets []attachmentpresets.Preset) (attachmentpresets.Preset, error) {
	args := ui.Called(title, baseRod, presets)
	return args.Get(0).(attachmentpresets.Preset), args.Error(1)
}

func (ui *mockSandboxUI) ChooseWaterContext(title string, presets []watercontexts.Preset) (watercontexts.Preset, error) {
	args := ui.Called(title, presets)
	return args.Get(0).(watercontexts.Preset), args.Error(1)
}

func (ui *mockSandboxUI) ResolveCast(title string, context encounter.WaterContext, presenter app.CastPresenter) (encounter.CastResult, error) {
	args := ui.Called(title, context, presenter)
	return args.Get(0).(encounter.CastResult), args.Error(1)
}

func (ui *mockSandboxUI) ShowEncounterOpening(title string, opening presentation.OpeningView) error {
	return ui.Called(title, opening).Error(0)
}

func (ui *mockSandboxUI) ShowFishSpawn(title string, spawn presentation.SpawnView) error {
	return ui.Called(title, spawn).Error(0)
}

func (ui *mockSandboxUI) ShowIntro(view presentation.IntroView) error {
	return ui.Called(view).Error(0)
}

func (ui *mockSandboxUI) ChooseMove(status presentation.StatusView, options []presentation.MoveOption) (domain.Move, error) {
	args := ui.Called(status, options)
	return args.Get(0).(domain.Move), args.Error(1)
}

func (ui *mockSandboxUI) ResolveSplash(view presentation.SplashView) (encounter.SplashResolution, error) {
	args := ui.Called(view)
	return args.Get(0).(encounter.SplashResolution), args.Error(1)
}

func (ui *mockSandboxUI) ShowRound(_ presentation.RoundView) error      { return nil }
func (ui *mockSandboxUI) ShowGameOver(_ presentation.SummaryView) error { return nil }
func (ui *mockSandboxUI) ShowNotice(message string) error             { return nil }

var (
	_ = match.StatusSnapshot{}
	_ = loadout.State{}
)
