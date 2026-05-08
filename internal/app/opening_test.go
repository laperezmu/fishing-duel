package app_test

import (
	"errors"
	"pesca/internal/app"
	"pesca/internal/content/watercontexts"
	"pesca/internal/encounter"
	"pesca/internal/player/loadout"
	"pesca/internal/player/rod"
	"pesca/internal/presentation"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestResolveEncounterOpening(t *testing.T) {
	baseConfig := encounter.DefaultConfig()
	playerLoadout := sampleLoadout(t)
	preset := watercontexts.DefaultPresets()[0]
	resolvedContext := preset.BuildContext()
	castResult := encounter.CastResult{Band: encounter.CastBandShort}
	ui := &mockOpeningUI{}
	presenter := presentation.NewPresenter(presentation.DefaultCatalog())

	ui.On("ChooseWaterContext", "Pesca", []watercontexts.Preset{preset}).Return(preset, nil).Once()
	ui.On("ResolveCast", "Pesca", resolvedContext, mock.Anything).Return(castResult, nil).Once()
	ui.On("ShowEncounterOpening", "Pesca", mock.MatchedBy(func(opening presentation.OpeningView) bool {
		return opening.InitialDistance == 1 && opening.InitialDepth == preset.InitialDepth && opening.CastLabel == encounter.CastBandShort.Label()
	})).Return(nil).Once()

	opening, err := app.ResolveEncounterOpening("Pesca", baseConfig, playerLoadout, []watercontexts.Preset{preset}, ui, presenter)

	require.NoError(t, err)
	assert.Equal(t, 1, opening.InitialDistance)
	assert.Equal(t, preset.InitialDepth, opening.InitialDepth)
	ui.AssertExpectations(t)
}

func TestResolveEncounterOpeningWrapsErrors(t *testing.T) {
	t.Run("returns an error when ui is missing", func(t *testing.T) {
		playerLoadout := sampleLoadout(t)

		_, err := app.ResolveEncounterOpening("Pesca", encounter.DefaultConfig(), playerLoadout, watercontexts.DefaultPresets(), nil, presentation.NewPresenter(presentation.DefaultCatalog()))

		require.Error(t, err)
		assert.EqualError(t, err, "opening ui is required")
	})

	t.Run("returns an error when presenter is missing", func(t *testing.T) {
		playerLoadout := sampleLoadout(t)

		_, err := app.ResolveEncounterOpening("Pesca", encounter.DefaultConfig(), playerLoadout, watercontexts.DefaultPresets(), &mockOpeningUI{}, nil)

		require.Error(t, err)
		assert.EqualError(t, err, "opening presenter is required")
	})

	t.Run("returns an error when player loadout is invalid", func(t *testing.T) {
		_, err := app.ResolveEncounterOpening("Pesca", encounter.DefaultConfig(), loadout.State{Rod: rod.State{TrackMaxDistance: 0, TrackMaxDepth: 3}}, watercontexts.DefaultPresets(), &mockOpeningUI{}, presentation.NewPresenter(presentation.DefaultCatalog()))

		require.Error(t, err)
		assert.EqualError(t, err, "player loadout: rod: track max distance must be greater than 0")
	})

	t.Run("returns an error when presets are empty", func(t *testing.T) {
		playerLoadout := sampleLoadout(t)

		_, err := app.ResolveEncounterOpening("Pesca", encounter.DefaultConfig(), playerLoadout, nil, &mockOpeningUI{}, presentation.NewPresenter(presentation.DefaultCatalog()))

		require.Error(t, err)
		assert.EqualError(t, err, "at least one water context preset is required")
	})

	t.Run("wraps selection errors", func(t *testing.T) {
		preset := watercontexts.DefaultPresets()[0]
		playerLoadout := sampleLoadout(t)
		ui := &mockOpeningUI{}
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())
		ui.On("ChooseWaterContext", "Pesca", []watercontexts.Preset{preset}).Return(watercontexts.Preset{}, errors.New("selection failed")).Once()

		_, err := app.ResolveEncounterOpening("Pesca", encounter.DefaultConfig(), playerLoadout, []watercontexts.Preset{preset}, ui, presenter)

		require.Error(t, err)
		assert.EqualError(t, err, "choose water context: selection failed")
		ui.AssertExpectations(t)
	})
}

func TestResolveEncounterOpeningWithPreset(t *testing.T) {
	t.Run("returns error when ui is nil", func(t *testing.T) {
		playerLoadout := sampleLoadout(t)
		preset := watercontexts.DefaultPresets()[0]

		_, err := app.ResolveEncounterOpeningWithPreset("Pesca", encounter.DefaultConfig(), playerLoadout, preset, nil, presentation.NewPresenter(presentation.DefaultCatalog()))

		require.Error(t, err)
		assert.EqualError(t, err, "opening ui is required")
	})

	t.Run("returns error when presenter is nil", func(t *testing.T) {
		playerLoadout := sampleLoadout(t)
		preset := watercontexts.DefaultPresets()[0]
		ui := &mockOpeningUI{}

		_, err := app.ResolveEncounterOpeningWithPreset("Pesca", encounter.DefaultConfig(), playerLoadout, preset, ui, nil)

		require.Error(t, err)
		assert.EqualError(t, err, "opening presenter is required")
	})

	t.Run("returns error when loadout is invalid", func(t *testing.T) {
		preset := watercontexts.DefaultPresets()[0]
		ui := &mockOpeningUI{}

		_, err := app.ResolveEncounterOpeningWithPreset("Pesca", encounter.DefaultConfig(), loadout.State{Rod: rod.State{TrackMaxDistance: 0, TrackMaxDepth: 3}}, preset, ui, presentation.NewPresenter(presentation.DefaultCatalog()))

		require.Error(t, err)
		assert.EqualError(t, err, "player loadout: rod: track max distance must be greater than 0")
	})

	t.Run("resolves opening successfully", func(t *testing.T) {
		playerLoadout := sampleLoadout(t)
		preset := watercontexts.DefaultPresets()[0]
		resolvedContext := preset.BuildContext()
		castResult := encounter.CastResult{Band: encounter.CastBandMedium}
		ui := &mockOpeningUI{}
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())

		ui.On("ResolveCast", "Pesca", resolvedContext, mock.Anything).Return(castResult, nil).Once()
		ui.On("ShowEncounterOpening", "Pesca", mock.Anything).Return(nil).Once()

		opening, err := app.ResolveEncounterOpeningWithPreset("Pesca", encounter.DefaultConfig(), playerLoadout, preset, ui, presenter)

		require.NoError(t, err)
		assert.Equal(t, 2, opening.InitialDistance)
		ui.AssertExpectations(t)
	})

	t.Run("wraps resolve cast errors", func(t *testing.T) {
		playerLoadout := sampleLoadout(t)
		preset := watercontexts.DefaultPresets()[0]
		resolvedContext := preset.BuildContext()
		ui := &mockOpeningUI{}
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())

		ui.On("ResolveCast", "Pesca", resolvedContext, mock.Anything).Return(encounter.CastResult{}, errors.New("cast failed")).Once()

		_, err := app.ResolveEncounterOpeningWithPreset("Pesca", encounter.DefaultConfig(), playerLoadout, preset, ui, presenter)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "resolve cast:")
		ui.AssertExpectations(t)
	})
}

type mockOpeningUI struct {
	mock.Mock
}

func (ui *mockOpeningUI) ChooseWaterContext(title string, presets []watercontexts.Preset) (watercontexts.Preset, error) {
	args := ui.Called(title, presets)
	return args.Get(0).(watercontexts.Preset), args.Error(1)
}

func (ui *mockOpeningUI) ResolveCast(title string, context encounter.WaterContext, presenter app.CastPresenter) (encounter.CastResult, error) {
	args := ui.Called(title, context, presenter)
	return args.Get(0).(encounter.CastResult), args.Error(1)
}

func (ui *mockOpeningUI) ShowEncounterOpening(title string, opening presentation.OpeningView) error {
	return ui.Called(title, opening).Error(0)
}

func sampleLoadout(t *testing.T) loadout.State {
	t.Helper()

	playerRod, err := rod.NewState(rod.DefaultConfig())
	require.NoError(t, err)
	playerLoadout, err := loadout.NewState(playerRod, nil)
	require.NoError(t, err)

	return playerLoadout
}
