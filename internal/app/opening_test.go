package app_test

import (
	"errors"
	"pesca/internal/app"
	"pesca/internal/content/watercontexts"
	"pesca/internal/encounter"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestResolveEncounterOpening(t *testing.T) {
	baseConfig := encounter.DefaultConfig()
	preset := watercontexts.DefaultPresets()[0]
	resolvedContext := preset.BuildContext()
	castResult := encounter.CastResult{Band: encounter.CastBandShort}
	ui := &mockOpeningUI{}

	ui.On("ChooseWaterContext", "Pesca", []watercontexts.Preset{preset}).Return(preset, nil).Once()
	ui.On("ResolveCast", "Pesca", resolvedContext).Return(castResult, nil).Once()
	ui.On("ShowEncounterOpening", "Pesca", mock.MatchedBy(func(opening encounter.Opening) bool {
		return opening.InitialDistance == 1 && opening.InitialDepth == preset.InitialDepth && opening.CastResult.Band == encounter.CastBandShort
	})).Return(nil).Once()

	opening, err := app.ResolveEncounterOpening("Pesca", baseConfig, []watercontexts.Preset{preset}, ui)

	require.NoError(t, err)
	assert.Equal(t, 1, opening.InitialDistance)
	assert.Equal(t, preset.InitialDepth, opening.InitialDepth)
	ui.AssertExpectations(t)
}

func TestResolveEncounterOpeningWrapsErrors(t *testing.T) {
	t.Run("returns an error when ui is missing", func(t *testing.T) {
		_, err := app.ResolveEncounterOpening("Pesca", encounter.DefaultConfig(), watercontexts.DefaultPresets(), nil)

		require.Error(t, err)
		assert.EqualError(t, err, "opening ui is required")
	})

	t.Run("wraps selection errors", func(t *testing.T) {
		preset := watercontexts.DefaultPresets()[0]
		ui := &mockOpeningUI{}
		ui.On("ChooseWaterContext", "Pesca", []watercontexts.Preset{preset}).Return(watercontexts.Preset{}, errors.New("selection failed")).Once()

		_, err := app.ResolveEncounterOpening("Pesca", encounter.DefaultConfig(), []watercontexts.Preset{preset}, ui)

		require.Error(t, err)
		assert.EqualError(t, err, "choose water context: selection failed")
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

func (ui *mockOpeningUI) ResolveCast(title string, context encounter.WaterContext) (encounter.CastResult, error) {
	args := ui.Called(title, context)
	return args.Get(0).(encounter.CastResult), args.Error(1)
}

func (ui *mockOpeningUI) ShowEncounterOpening(title string, opening encounter.Opening) error {
	return ui.Called(title, opening).Error(0)
}
