package anglerprofiles

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultUnlockedProfilesOnlyReturnUnlockedByDefault(t *testing.T) {
	t.Parallel()

	profiles := DefaultUnlockedProfiles()
	require.NotEmpty(t, profiles)
	for _, profile := range profiles {
		assert.True(t, profile.UnlockedByDefault)
	}
}

func TestResolveStartBuildsLoadoutAndThread(t *testing.T) {
	t.Parallel()

	resolved, err := ResolveStart(DefaultUnlockedProfiles()[0])
	require.NoError(t, err)
	assert.NotEmpty(t, resolved.Profile.ID)
	assert.NotEmpty(t, resolved.DeckPreset.ID)
	assert.NotEmpty(t, resolved.RodPreset.ID)
	assert.Equal(t, resolved.Profile.StartingThread, resolved.StartingThread)
	assert.NoError(t, resolved.Loadout.Validate())
}
