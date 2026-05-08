package app

import (
	"pesca/internal/content/anglerprofiles"
	"pesca/internal/content/attachmentpresets"
	"pesca/internal/content/fishprofiles"
	"pesca/internal/content/playerprofiles"
	"pesca/internal/content/rodpresets"
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

func TestResolveRunSetup(t *testing.T) {
	t.Parallel()

	ui := &mockRunSetupUI{}
	profile := anglerprofiles.DefaultUnlockedProfiles()[0]
	ui.On("ChooseAnglerProfile", "Pesca", mock.Anything).Return(profile, nil).Once()

	resolved, err := resolveRunSetup("Pesca", ui)
	require.NoError(t, err)
	assert.Equal(t, profile.ID, resolved.Profile.ID)
	assert.Equal(t, profile.StartingThread, resolved.StartingThread)
	assert.NoError(t, resolved.Loadout.Validate())
	ui.AssertExpectations(t)
}

func TestBootstrapEncounterForRun(t *testing.T) {
	t.Parallel()

	t.Run("returns error when rng is nil", func(t *testing.T) {
		t.Parallel()

		result, err := BootstrapEncounterForRun(
			"Test",
			nil,
			&mockBootstrapUI{},
			&mockSpawnUI{},
			presentation.NewPresenter(presentation.DefaultCatalog()),
			playerprofiles.DeckPreset{},
			loadout.State{},
			RunEncounterBootstrapConfig{},
		)

		require.Error(t, err)
		assert.EqualError(t, err, "randomizer is required")
		assert.Empty(t, result)
	})

	t.Run("returns error when openingUI is nil", func(t *testing.T) {
		t.Parallel()

		rng := NewSeededRandom(12345)
		result, err := BootstrapEncounterForRun(
			"Test",
			rng,
			nil,
			&mockSpawnUI{},
			presentation.NewPresenter(presentation.DefaultCatalog()),
			playerprofiles.DeckPreset{},
			loadout.State{},
			RunEncounterBootstrapConfig{},
		)

		require.Error(t, err)
		assert.EqualError(t, err, "encounter bootstrap ui is required")
		assert.Empty(t, result)
	})

	t.Run("returns error when spawnUI is nil", func(t *testing.T) {
		t.Parallel()

		rng := NewSeededRandom(12345)
		result, err := BootstrapEncounterForRun(
			"Test",
			rng,
			&mockBootstrapUI{},
			nil,
			presentation.NewPresenter(presentation.DefaultCatalog()),
			playerprofiles.DeckPreset{},
			loadout.State{},
			RunEncounterBootstrapConfig{},
		)

		require.Error(t, err)
		assert.EqualError(t, err, "spawn ui is required")
		assert.Empty(t, result)
	})
}

type mockBootstrapUI struct{ mock.Mock }

func (ui *mockBootstrapUI) ChoosePlayerDeckPreset(title string, presets []playerprofiles.DeckPreset) (playerprofiles.DeckPreset, error) {
	args := ui.Called(title, presets)
	return args.Get(0).(playerprofiles.DeckPreset), args.Error(1)
}

func (ui *mockBootstrapUI) ChooseRodPreset(title string, presets []rodpresets.Preset) (rodpresets.Preset, error) {
	args := ui.Called(title, presets)
	return args.Get(0).(rodpresets.Preset), args.Error(1)
}

func (ui *mockBootstrapUI) ChooseAttachmentPreset(title string, baseRod rod.State, presets []attachmentpresets.Preset) (attachmentpresets.Preset, error) {
	args := ui.Called(title, baseRod, presets)
	return args.Get(0).(attachmentpresets.Preset), args.Error(1)
}

func (ui *mockBootstrapUI) ResolveCast(name string, context encounter.WaterContext, presenter CastPresenter) (encounter.CastResult, error) {
	args := ui.Called(name, context, presenter)
	return args.Get(0).(encounter.CastResult), args.Error(1)
}

func (ui *mockBootstrapUI) ShowOpening(view presentation.OpeningView) error {
	return ui.Called(view).Error(0)
}

func (ui *mockBootstrapUI) ShowEncounterOpening(title string, opening presentation.OpeningView) error {
	return ui.Called(title, opening).Error(0)
}

func (ui *mockBootstrapUI) ShowFishSpawn(title string, spawn presentation.SpawnView) error {
	return ui.Called(title, spawn).Error(0)
}

func (ui *mockBootstrapUI) ChooseWaterContext(name string, contexts []watercontexts.Preset) (watercontexts.Preset, error) {
	args := ui.Called(name, contexts)
	return args.Get(0).(watercontexts.Preset), args.Error(1)
}

func (ui *mockBootstrapUI) ChooseFishDeckPreset(name string, presets []fishprofiles.Profile) (fishprofiles.Profile, error) {
	args := ui.Called(name, presets)
	return args.Get(0).(fishprofiles.Profile), args.Error(1)
}

type mockSpawnUI struct{ mock.Mock }

func (ui *mockSpawnUI) ShowFishSpawn(title string, spawn presentation.SpawnView) error {
	return ui.Called(title, spawn).Error(0)
}

func (ui *mockSpawnUI) ChooseFishDeckPreset(name string, presets []fishprofiles.Profile) (fishprofiles.Profile, error) {
	args := ui.Called(name, presets)
	return args.Get(0).(fishprofiles.Profile), args.Error(1)
}

type mockRunSetupUI struct{ mock.Mock }

func (ui *mockRunSetupUI) ChooseAnglerProfile(title string, profiles []anglerprofiles.Profile) (anglerprofiles.Profile, error) {
	args := ui.Called(title, profiles)
	return args.Get(0).(anglerprofiles.Profile), args.Error(1)
}
