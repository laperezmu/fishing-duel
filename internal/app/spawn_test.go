package app_test

import (
	"errors"
	"pesca/internal/app"
	"pesca/internal/content/attachmentpresets"
	"pesca/internal/content/fishprofiles"
	"pesca/internal/content/habitats"
	"pesca/internal/content/playerprofiles"
	"pesca/internal/content/rodpresets"
	"pesca/internal/content/watercontexts"
	"pesca/internal/content/waterpools"
	"pesca/internal/encounter"
	"pesca/internal/player/loadout"
	"pesca/internal/player/rod"
	"pesca/internal/presentation"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestResolveFishSpawn(t *testing.T) {
	playerLoadout := sampleSurfaceLoadout(t)
	opening := encounter.Opening{
		WaterContext:    encounter.WaterContext{PoolTag: waterpools.Shoreline},
		InitialDistance: 1,
		InitialDepth:    1,
	}
	ui := &mockSpawnUI{}
	presenter := presentation.NewPresenter(presentation.DefaultCatalog())

	ui.On("ShowFishSpawn", "Pesca", mock.MatchedBy(func(spawn presentation.SpawnView) bool {
		return spawn.ProfileLabel == "Control de superficie" && spawn.CandidateCount >= 2
	})).Return(nil).Once()

	spawn, err := app.ResolveFishSpawn("Pesca", opening, playerLoadout, fishprofiles.DefaultProfiles(), ui, presenter)

	require.NoError(t, err)
	assert.Equal(t, fishprofiles.ProfileID("surface-control"), spawn.Profile.ID)
	ui.AssertExpectations(t)
}

func TestBootstrapEncounterWithConfigUsesClosedFishPool(t *testing.T) {
	ui := &bootstrapSpawnUI{}
	rng := app.NewSeededRandom(7)
	catalog := fishprofiles.DefaultCatalog()

	deckPreset := samplePlayerDeckPreset()
	rodPreset := sampleRodPreset()
	attachmentPreset := sampleAttachmentPreset()
	waterPreset := sampleWaterContextPreset()

	ui.On("ChoosePlayerDeckPreset", "Pesca", mock.Anything).Return(deckPreset, nil).Once()
	ui.On("ChooseRodPreset", "Pesca", mock.Anything).Return(rodPreset, nil).Once()
	ui.On("ChooseAttachmentPreset", "Pesca", mock.Anything, mock.Anything).Return(attachmentPreset, nil).Once()
	ui.On("ChooseWaterContext", "Pesca", mock.Anything).Return(waterPreset, nil).Once()
	ui.On("ResolveCast", "Pesca", mock.Anything, mock.Anything).Return(encounter.CastResult{Band: encounter.CastBandShort}, nil).Once()
	ui.On("ShowEncounterOpening", "Pesca", mock.Anything).Return(nil).Once()
	ui.On("ShowFishSpawn", "Pesca", mock.MatchedBy(func(spawn presentation.SpawnView) bool {
		return spawn.ProfileLabel != "Presion horizontal" && spawn.ProfileLabel != "Corriente mixta"
	})).Return(nil).Once()

	engine, err := app.BootstrapEncounterWithConfig("Pesca", rng, ui, app.EncounterBootstrapConfig{
		FishCatalog: catalog,
		FishPoolID:  fishprofiles.PoolID("shoreline-basics"),
	})

	require.NoError(t, err)
	require.NotNil(t, engine)
	ui.AssertExpectations(t)
}

func TestBootstrapEncounterWithConfigUsesDirectFishPreset(t *testing.T) {
	ui := &bootstrapSpawnUI{}
	rng := app.NewSeededRandom(7)

	deckPreset := samplePlayerDeckPreset()
	rodPreset := sampleRodPreset()
	attachmentPreset := sampleAttachmentPreset()
	waterPreset := sampleWaterContextPreset()

	ui.On("ChoosePlayerDeckPreset", "Pesca", mock.Anything).Return(deckPreset, nil).Once()
	ui.On("ChooseRodPreset", "Pesca", mock.Anything).Return(rodPreset, nil).Once()
	ui.On("ChooseAttachmentPreset", "Pesca", mock.Anything, mock.Anything).Return(attachmentPreset, nil).Once()
	ui.On("ChooseWaterContext", "Pesca", mock.Anything).Return(waterPreset, nil).Once()
	ui.On("ResolveCast", "Pesca", mock.Anything, mock.Anything).Return(encounter.CastResult{Band: encounter.CastBandShort}, nil).Once()
	ui.On("ShowEncounterOpening", "Pesca", mock.Anything).Return(nil).Once()
	ui.On("ShowFishSpawn", "Pesca", mock.MatchedBy(func(spawn presentation.SpawnView) bool {
		return spawn.ProfileLabel == "Clasico"
	})).Return(nil).Once()

	engine, err := app.BootstrapEncounterWithConfig("Pesca", rng, ui, app.EncounterBootstrapConfig{FishPresetID: fishprofiles.ProfileID("classic")})

	require.NoError(t, err)
	require.NotNil(t, engine)
	ui.AssertExpectations(t)
}

func TestBootstrapEncounterWithConfigUsesPresetIDsWithoutSelectionPrompts(t *testing.T) {
	ui := &bootstrapSpawnUI{}
	rng := app.NewSeededRandom(7)

	ui.On("ResolveCast", "Pesca", mock.Anything, mock.Anything).Return(encounter.CastResult{Band: encounter.CastBandShort}, nil).Once()
	ui.On("ShowEncounterOpening", "Pesca", mock.Anything).Return(nil).Once()
	ui.On("ShowFishSpawn", "Pesca", mock.Anything).Return(nil).Once()

	engine, err := app.BootstrapEncounterWithConfig("Pesca", rng, ui, app.EncounterBootstrapConfig{
		PlayerDeckPresetID: samplePlayerDeckPreset().ID,
		RodPresetID:        sampleRodPreset().ID,
		AttachmentPresetID: sampleAttachmentPreset().ID,
		WaterContextID:     sampleWaterContextPreset().ID,
		FishPresetID:       fishprofiles.ProfileID("classic"),
	})

	require.NoError(t, err)
	require.NotNil(t, engine)
	ui.AssertExpectations(t)
}

func TestBootstrapEncounterWithConfigAppliesStateOverrides(t *testing.T) {
	ui := &bootstrapSpawnUI{}
	rng := app.NewSeededRandom(7)

	initialDistance := 2
	initialDepth := 0
	captureDistance := 1
	recycleCount := 3

	ui.On("ResolveCast", "Pesca", mock.Anything, mock.Anything).Return(encounter.CastResult{Band: encounter.CastBandShort}, nil).Once()
	ui.On("ShowEncounterOpening", "Pesca", mock.Anything).Return(nil).Once()
	ui.On("ShowFishSpawn", "Pesca", mock.Anything).Return(nil).Once()

	engine, err := app.BootstrapEncounterWithConfig("Pesca", rng, ui, app.EncounterBootstrapConfig{
		PlayerDeckPresetID: samplePlayerDeckPreset().ID,
		RodPresetID:        sampleRodPreset().ID,
		AttachmentPresetID: sampleAttachmentPreset().ID,
		WaterContextID:     sampleWaterContextPreset().ID,
		FishPresetID:       fishprofiles.ProfileID("classic"),
		StateOverrides: app.SandboxStateOverrides{
			InitialDistance: &initialDistance,
			InitialDepth:    &initialDepth,
			CaptureDistance: &captureDistance,
			RecycleCount:    &recycleCount,
		},
	})

	require.NoError(t, err)
	state := engine.State()
	assert.Equal(t, 2, state.Encounter.Distance)
	assert.Equal(t, 0, state.Encounter.Depth)
	assert.Equal(t, 1, state.Encounter.Config.CaptureDistance)
	assert.Equal(t, 3, state.Deck.RecycleCount)
	ui.AssertExpectations(t)
}

func TestBootstrapEncounterWithConfigReturnsPoolErrors(t *testing.T) {
	ui := &bootstrapSpawnUI{}
	rng := app.NewSeededRandom(7)

	deckPreset := samplePlayerDeckPreset()
	rodPreset := sampleRodPreset()
	attachmentPreset := sampleAttachmentPreset()
	waterPreset := sampleWaterContextPreset()

	ui.On("ChoosePlayerDeckPreset", "Pesca", mock.Anything).Return(deckPreset, nil).Once()
	ui.On("ChooseRodPreset", "Pesca", mock.Anything).Return(rodPreset, nil).Once()
	ui.On("ChooseAttachmentPreset", "Pesca", mock.Anything, mock.Anything).Return(attachmentPreset, nil).Once()
	ui.On("ChooseWaterContext", "Pesca", mock.Anything).Return(waterPreset, nil).Once()
	ui.On("ResolveCast", "Pesca", mock.Anything, mock.Anything).Return(encounter.CastResult{Band: encounter.CastBandShort}, nil).Once()
	ui.On("ShowEncounterOpening", "Pesca", mock.Anything).Return(nil).Once()

	_, err := app.BootstrapEncounterWithConfig("Pesca", rng, ui, app.EncounterBootstrapConfig{FishPoolID: fishprofiles.PoolID("missing-pool")})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "resolve encounter fish pool: unknown fish pool id missing-pool")
	ui.AssertExpectations(t)
}

func TestResolveFishSpawnWrapsErrors(t *testing.T) {
	t.Run("returns an error when ui is missing", func(t *testing.T) {
		_, err := app.ResolveFishSpawn("Pesca", encounter.Opening{}, sampleSurfaceLoadout(t), fishprofiles.DefaultProfiles(), nil, presentation.NewPresenter(presentation.DefaultCatalog()))

		require.Error(t, err)
		assert.EqualError(t, err, "spawn ui is required")
	})

	t.Run("returns an error when no profile matches", func(t *testing.T) {
		ui := &mockSpawnUI{}
		opening := encounter.Opening{
			WaterContext:    encounter.WaterContext{PoolTag: waterpools.Offshore},
			InitialDistance: 5,
			InitialDepth:    4,
		}

		_, err := app.ResolveFishSpawn("Pesca", opening, sampleLoadout(t), fishprofiles.DefaultProfiles(), ui, presentation.NewPresenter(presentation.DefaultCatalog()))

		require.Error(t, err)
		assert.Contains(t, err.Error(), "resolve fish spawn: no fish profile matches water offshore")
	})

	t.Run("wraps ui errors", func(t *testing.T) {
		ui := &mockSpawnUI{}
		opening := encounter.Opening{
			WaterContext:    encounter.WaterContext{PoolTag: waterpools.Shoreline},
			InitialDistance: 0,
			InitialDepth:    1,
		}
		ui.On("ShowFishSpawn", "Pesca", mock.Anything).Return(errors.New("ui failed")).Once()

		_, err := app.ResolveFishSpawn("Pesca", opening, sampleLoadout(t), fishprofiles.DefaultProfiles(), ui, presentation.NewPresenter(presentation.DefaultCatalog()))

		require.Error(t, err)
		assert.EqualError(t, err, "show fish spawn: ui failed")
		ui.AssertExpectations(t)
	})
}

type mockSpawnUI struct {
	mock.Mock
}

func (ui *mockSpawnUI) ShowFishSpawn(title string, spawn presentation.SpawnView) error {
	return ui.Called(title, spawn).Error(0)
}

type bootstrapSpawnUI struct {
	mock.Mock
}

func (ui *bootstrapSpawnUI) ChoosePlayerDeckPreset(title string, presets []playerprofiles.DeckPreset) (playerprofiles.DeckPreset, error) {
	args := ui.Called(title, presets)
	return args.Get(0).(playerprofiles.DeckPreset), args.Error(1)
}

func (ui *bootstrapSpawnUI) ChooseRodPreset(title string, presets []rodpresets.Preset) (rodpresets.Preset, error) {
	args := ui.Called(title, presets)
	return args.Get(0).(rodpresets.Preset), args.Error(1)
}

func (ui *bootstrapSpawnUI) ChooseAttachmentPreset(title string, baseRod rod.State, presets []attachmentpresets.Preset) (attachmentpresets.Preset, error) {
	args := ui.Called(title, baseRod, presets)
	return args.Get(0).(attachmentpresets.Preset), args.Error(1)
}

func (ui *bootstrapSpawnUI) ChooseWaterContext(title string, presets []watercontexts.Preset) (watercontexts.Preset, error) {
	args := ui.Called(title, presets)
	return args.Get(0).(watercontexts.Preset), args.Error(1)
}

func (ui *bootstrapSpawnUI) ResolveCast(title string, context encounter.WaterContext, presenter app.CastPresenter) (encounter.CastResult, error) {
	args := ui.Called(title, context, presenter)
	return args.Get(0).(encounter.CastResult), args.Error(1)
}

func (ui *bootstrapSpawnUI) ShowEncounterOpening(title string, opening presentation.OpeningView) error {
	return ui.Called(title, opening).Error(0)
}

func (ui *bootstrapSpawnUI) ShowFishSpawn(title string, spawn presentation.SpawnView) error {
	return ui.Called(title, spawn).Error(0)
}

func samplePlayerDeckPreset() playerprofiles.DeckPreset {
	return playerprofiles.DefaultPresets()[0]
}

func sampleRodPreset() rodpresets.Preset {
	return rodpresets.DefaultPresets()[1]
}

func sampleAttachmentPreset() attachmentpresets.Preset {
	return attachmentpresets.DefaultPresets()[0]
}

func sampleWaterContextPreset() watercontexts.Preset {
	return watercontexts.DefaultPresets()[0]
}

func sampleSurfaceLoadout(t *testing.T) loadout.State {
	t.Helper()

	playerRod, err := rod.NewState(rod.DefaultConfig())
	require.NoError(t, err)
	playerLoadout, err := loadout.NewState(playerRod, []loadout.Attachment{{
		ID:          "floating-line",
		Name:        "Linea flotante",
		HabitatTags: []habitats.Tag{habitats.Surface},
	}})
	require.NoError(t, err)

	return playerLoadout
}
