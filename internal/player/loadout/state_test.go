package loadout

import (
	"pesca/internal/content/habitats"
	"pesca/internal/encounter"
	"pesca/internal/player/rod"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewState(t *testing.T) {
	tests := []struct {
		title       string
		playerRod   rod.State
		attachments []Attachment
		wantState   State
		wantErrText string
	}{
		{
			title:     "returns a loadout state when the rod is valid",
			playerRod: sampleRodState(),
			attachments: []Attachment{{
				ID:          "float-basic",
				Name:        "Flotador basico",
				Description: "Placeholder para futuros aditamentos.",
				HabitatTags: []habitats.Tag{habitats.Surface},
			}},
			wantState: State{
				Rod: sampleRodState(),
				Attachments: []Attachment{{
					ID:          "float-basic",
					Name:        "Flotador basico",
					Description: "Placeholder para futuros aditamentos.",
					HabitatTags: []habitats.Tag{habitats.Surface},
				}},
			},
		},
		{
			title:       "returns an error when the rod is invalid",
			playerRod:   rod.State{OpeningMaxDistance: 7, OpeningMaxDepth: 3, TrackMaxDistance: 5, TrackMaxDepth: 3},
			wantErrText: "rod: opening max distance must be less than or equal to track max distance",
		},
		{
			title:     "returns an error when an attachment has no id",
			playerRod: sampleRodState(),
			attachments: []Attachment{{
				Name: "Sin id",
			}},
			wantErrText: "attachment id is required",
		},
		{
			title:     "returns an error when an attachment has no name",
			playerRod: sampleRodState(),
			attachments: []Attachment{{
				ID: "sinker",
			}},
			wantErrText: "attachment name is required",
		},
		{
			title:     "returns an error when attachment ids repeat",
			playerRod: sampleRodState(),
			attachments: []Attachment{{
				ID:   "dup",
				Name: "Uno",
			}, {
				ID:   "dup",
				Name: "Dos",
			}},
			wantErrText: "attachment ids must be unique",
		},
		{
			title:     "returns an error when an attachment creates invalid effective rod limits",
			playerRod: sampleRodState(),
			attachments: []Attachment{{
				ID:                      "bad-distance",
				Name:                    "Mala distancia",
				OpeningDistanceModifier: -5,
			}},
			wantErrText: "rod: opening max distance must be greater than or equal to 0",
		},
		{
			title:     "returns an error when habitat tags include empty values",
			playerRod: sampleRodState(),
			attachments: []Attachment{{
				ID:          "bad-tag",
				Name:        "Tag vacio",
				HabitatTags: []habitats.Tag{""},
			}},
			wantErrText: "unknown habitat tag \"\"",
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			state, err := NewState(test.playerRod, test.attachments)

			if test.wantErrText != "" {
				require.Error(t, err)
				assert.EqualError(t, err, test.wantErrText)
				assert.Equal(t, State{}, state)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.wantState, state)
		})
	}
}

func TestOpeningLimits(t *testing.T) {
	state, err := NewState(sampleRodState(), []Attachment{{
		ID:                      "weighted",
		Name:                    "Plomada",
		OpeningDistanceModifier: -1,
		OpeningDepthModifier:    1,
	}})
	require.NoError(t, err)

	assert.Equal(t, encounter.OpeningLimits{MaxInitialDistance: 3, MaxInitialDepth: 3}, state.OpeningLimits())
}

func TestEffectiveTrackAndHabitats(t *testing.T) {
	state, err := NewState(sampleRodState(), []Attachment{{
		ID:                    "line-reinforced",
		Name:                  "Linea reforzada",
		TrackDistanceModifier: 1,
		TrackDepthModifier:    2,
		SplashBonusDistance:   1,
		HabitatTags:           []habitats.Tag{habitats.Rock, habitats.Weed, habitats.Rock},
	}})
	require.NoError(t, err)

	assert.Equal(t, 6, state.TrackMaxDistance())
	assert.Equal(t, 5, state.TrackMaxDepth())
	assert.Equal(t, 2, state.SplashSuccessDistanceBonus())
	assert.Equal(t, []habitats.Tag{habitats.Rock, habitats.Weed}, state.HabitatTags())
}

func TestNewStateClonesAttachmentHabitatTags(t *testing.T) {
	attachments := []Attachment{{
		ID:          "line-reinforced",
		Name:        "Linea reforzada",
		HabitatTags: []habitats.Tag{habitats.Rock},
	}}

	state, err := NewState(sampleRodState(), attachments)
	require.NoError(t, err)

	attachments[0].HabitatTags[0] = "changed"

	assert.Equal(t, []habitats.Tag{habitats.Rock}, state.Attachments[0].HabitatTags)
}

func sampleRodState() rod.State {
	return rod.State{
		OpeningMaxDistance:  4,
		OpeningMaxDepth:     2,
		TrackMaxDistance:    5,
		TrackMaxDepth:       3,
		SplashBonusDistance: 1,
	}
}
