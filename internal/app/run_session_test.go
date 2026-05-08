package app_test

import (
	"pesca/internal/app"
	"pesca/internal/content/anglerprofiles"
	"pesca/internal/content/attachmentpresets"
	"pesca/internal/content/fishprofiles"
	"pesca/internal/content/playerprofiles"
	"pesca/internal/content/rodpresets"
	"pesca/internal/content/watercontexts"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/player/rod"
	"pesca/internal/presentation"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewRunSession(t *testing.T) {
	t.Run("returns a run session when all dependencies are provided", func(t *testing.T) {
		engine := &mockRNG{}
		ui := &mockRunUI{}
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())

		session, err := app.NewRunSession("Test Run", engine, ui, presenter)

		require.NoError(t, err)
		assert.NotNil(t, session)
	})

	t.Run("returns an error when title is missing", func(t *testing.T) {
		engine := &mockRNG{}
		ui := &mockRunUI{}
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())

		session, err := app.NewRunSession("", engine, ui, presenter)

		require.Error(t, err)
		assert.EqualError(t, err, "title is required")
		assert.Nil(t, session)
	})

	t.Run("returns an error when rng is missing", func(t *testing.T) {
		ui := &mockRunUI{}
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())

		session, err := app.NewRunSession("Test Run", nil, ui, presenter)

		require.Error(t, err)
		assert.EqualError(t, err, "randomizer is required")
		assert.Nil(t, session)
	})

	t.Run("returns an error when ui is missing", func(t *testing.T) {
		engine := &mockRNG{}
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())

		session, err := app.NewRunSession("Test Run", engine, nil, presenter)

		require.Error(t, err)
		assert.EqualError(t, err, "run ui is required")
		assert.Nil(t, session)
	})
}

type mockRNG struct {
	mock.Mock
}

func (rng *mockRNG) Float64() float64 {
	return rng.Called().Get(0).(float64)
}

func (rng *mockRNG) Intn(n int) int {
	return rng.Called(n).Int(0)
}

func (rng *mockRNG) Shuffle(n int, swap func(i, j int)) {
	rng.Called(n, swap)
}

type mockRunUI struct {
	mock.Mock
}

func (ui *mockRunUI) ChooseAnglerProfile(name string, profiles []anglerprofiles.Profile) (anglerprofiles.Profile, error) {
	args := ui.Called(name, profiles)
	return args.Get(0).(anglerprofiles.Profile), args.Error(1)
}

func (ui *mockRunUI) ChoosePlayerDeckPreset(name string, presets []playerprofiles.DeckPreset) (playerprofiles.DeckPreset, error) {
	args := ui.Called(name, presets)
	return args.Get(0).(playerprofiles.DeckPreset), args.Error(1)
}

func (ui *mockRunUI) ChooseRodPreset(name string, presets []rodpresets.Preset) (rodpresets.Preset, error) {
	args := ui.Called(name, presets)
	return args.Get(0).(rodpresets.Preset), args.Error(1)
}

func (ui *mockRunUI) ChooseAttachmentPreset(name string, baseRod rod.State, presets []attachmentpresets.Preset) (attachmentpresets.Preset, error) {
	args := ui.Called(name, baseRod, presets)
	return args.Get(0).(attachmentpresets.Preset), args.Error(1)
}

func (ui *mockRunUI) ShowOpening(view presentation.OpeningView) error {
	return ui.Called(view).Error(0)
}

func (ui *mockRunUI) ResolveCast(name string, context encounter.WaterContext, presenter app.CastPresenter) (encounter.CastResult, error) {
	args := ui.Called(name, context, presenter)
	return args.Get(0).(encounter.CastResult), args.Error(1)
}

func (ui *mockRunUI) ShowEncounterOpening(title string, opening presentation.OpeningView) error {
	return ui.Called(title, opening).Error(0)
}

func (ui *mockRunUI) ShowFishSpawn(title string, spawn presentation.SpawnView) error {
	return ui.Called(title, spawn).Error(0)
}

func (ui *mockRunUI) ChooseWaterContext(name string, contexts []watercontexts.Preset) (watercontexts.Preset, error) {
	args := ui.Called(name, contexts)
	return args.Get(0).(watercontexts.Preset), args.Error(1)
}

func (ui *mockRunUI) ChooseFishDeckPreset(name string, presets []fishprofiles.Profile) (fishprofiles.Profile, error) {
	args := ui.Called(name, presets)
	return args.Get(0).(fishprofiles.Profile), args.Error(1)
}

func (ui *mockRunUI) ShowRunIntro(view presentation.RunIntroView) error {
	return ui.Called(view).Error(0)
}

func (ui *mockRunUI) ShowRunNode(view presentation.RunNodeView) error {
	return ui.Called(view).Error(0)
}

func (ui *mockRunUI) ShowRunNodeSummary(view presentation.RunNodeSummaryView) error {
	return ui.Called(view).Error(0)
}

func (ui *mockRunUI) ShowRunSummary(view presentation.RunSummaryView) error {
	return ui.Called(view).Error(0)
}

func (ui *mockRunUI) ShowIntro(view presentation.IntroView) error {
	return ui.Called(view).Error(0)
}

func (ui *mockRunUI) ChooseMove(view presentation.StatusView, options []presentation.MoveOption) (domain.Move, error) {
	args := ui.Called(view, options)
	return args.Get(0).(domain.Move), args.Error(1)
}

func (ui *mockRunUI) ShowRound(view presentation.RoundView) error {
	return ui.Called(view).Error(0)
}

func (ui *mockRunUI) ResolveSplash(view presentation.SplashView) (encounter.SplashResolution, error) {
	args := ui.Called(view)
	return args.Get(0).(encounter.SplashResolution), args.Error(1)
}

func (ui *mockRunUI) ShowGameOver(view presentation.SummaryView) error {
	return ui.Called(view).Error(0)
}
