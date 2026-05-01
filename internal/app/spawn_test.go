package app_test

import (
	"errors"
	"pesca/internal/app"
	"pesca/internal/content/fishprofiles"
	"pesca/internal/content/habitats"
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
	assert.Equal(t, "surface-control", spawn.Profile.ID)
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
