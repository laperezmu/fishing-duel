package app

import (
	"pesca/internal/content/anglerprofiles"
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

type mockRunSetupUI struct{ mock.Mock }

func (ui *mockRunSetupUI) ChooseAnglerProfile(title string, profiles []anglerprofiles.Profile) (anglerprofiles.Profile, error) {
	args := ui.Called(title, profiles)
	return args.Get(0).(anglerprofiles.Profile), args.Error(1)
}
