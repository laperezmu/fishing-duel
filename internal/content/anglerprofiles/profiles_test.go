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

func TestProfileValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		profile Profile
		wantErr bool
	}{
		{
			name: "valid profile",
			profile: Profile{
				ID:                 "test-profile",
				Name:               "Test Profile",
				DeckPresetID:       "classic",
				RodPresetID:        "coastal-control",
				AttachmentPresetID: "none",
				StartingThread:     5,
				UnlockedByDefault:  true,
			},
			wantErr: false,
		},
		{
			name: "missing id",
			profile: Profile{
				Name:              "Test Profile",
				StartingThread:    5,
				UnlockedByDefault: true,
			},
			wantErr: true,
		},
		{
			name: "missing name",
			profile: Profile{
				ID:                 "test-profile",
				DeckPresetID:       "classic",
				RodPresetID:        "coastal-control",
				AttachmentPresetID: "none",
				StartingThread:     5,
				UnlockedByDefault:  true,
			},
			wantErr: true,
		},
		{
			name: "negative thread",
			profile: Profile{
				ID:                 "test-profile",
				Name:               "Test Profile",
				DeckPresetID:       "classic",
				RodPresetID:        "coastal-control",
				AttachmentPresetID: "none",
				StartingThread:     -1,
				UnlockedByDefault:  true,
			},
			wantErr: true,
		},
		{
			name: "zero thread",
			profile: Profile{
				ID:                 "test-profile",
				Name:               "Test Profile",
				DeckPresetID:       "classic",
				RodPresetID:        "coastal-control",
				AttachmentPresetID: "none",
				StartingThread:     0,
				UnlockedByDefault:  true,
			},
			wantErr: true,
		},
		{
			name: "missing deck preset",
			profile: Profile{
				ID:                 "test-profile",
				Name:               "Test Profile",
				RodPresetID:        "coastal-control",
				AttachmentPresetID: "none",
				StartingThread:     5,
				UnlockedByDefault:  true,
			},
			wantErr: true,
		},
		{
			name: "missing rod preset",
			profile: Profile{
				ID:                 "test-profile",
				Name:               "Test Profile",
				DeckPresetID:       "classic",
				AttachmentPresetID: "none",
				StartingThread:     5,
				UnlockedByDefault:  true,
			},
			wantErr: true,
		},
		{
			name: "missing attachment preset",
			profile: Profile{
				ID:                "test-profile",
				Name:              "Test Profile",
				DeckPresetID:      "classic",
				RodPresetID:       "coastal-control",
				StartingThread:    5,
				UnlockedByDefault: true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.profile.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestResolveStartWithDifferentProfiles(t *testing.T) {
	t.Parallel()

	profiles := DefaultUnlockedProfiles()
	require.NotEmpty(t, profiles)

	for _, profile := range profiles {
		t.Run(profile.ID, func(t *testing.T) {
			t.Parallel()

			resolved, err := ResolveStart(profile)
			require.NoError(t, err)
			assert.Equal(t, profile.ID, resolved.Profile.ID)
			assert.NoError(t, resolved.Loadout.Validate())
		})
	}
}
