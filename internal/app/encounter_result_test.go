package app

import (
	"pesca/internal/content/fishprofiles"
	"pesca/internal/encounter"
	"pesca/internal/match"
	"pesca/internal/player/loadout"
	"pesca/internal/player/rod"
	"pesca/internal/run"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveEncounterResultCaptured(t *testing.T) {
	t.Parallel()

	result, err := ResolveEncounterResult(match.State{
		Encounter: encounter.State{Status: encounter.StatusCaptured, EndReason: encounter.EndReasonTrackCapture},
		Player:    match.PlayerState{Loadout: mustBuildLoadout(t, 5)},
		Lifecycle: match.LifecycleState{Finished: true},
	}, fishprofiles.Spawn{Profile: fishprofiles.Profile{ID: "tuna", Name: "Atun"}})

	require.NoError(t, err)
	assert.Equal(t, run.EncounterOutcomeCaptured, result.Outcome)
	assert.True(t, result.NodeResolved)
	assert.NotNil(t, result.Capture)
	assert.Equal(t, "tuna", result.Capture.FishID)
	assert.Equal(t, "Atun", result.Capture.FishName)
	assert.Equal(t, 0, result.ThreadDamage)
}

func TestResolveEncounterResultTrackEscapeAlwaysAddsOneThreadDamage(t *testing.T) {
	t.Parallel()

	result, err := ResolveEncounterResult(match.State{
		Encounter: encounter.State{Status: encounter.StatusEscaped, EndReason: encounter.EndReasonTrackEscape, Distance: 7},
		Player:    match.PlayerState{Loadout: mustBuildLoadout(t, 5)},
		Lifecycle: match.LifecycleState{Finished: true},
	}, fishprofiles.Spawn{Profile: fishprofiles.Profile{ID: "tuna", Name: "Atun"}})

	require.NoError(t, err)
	assert.Equal(t, run.EncounterOutcomeEscaped, result.Outcome)
	assert.True(t, result.NodeResolved)
	assert.Equal(t, 1, result.ThreadDamage)
	assert.Nil(t, result.Capture)
}

func TestResolveEncounterResultDepthEscapeAlwaysAddsOneThreadDamage(t *testing.T) {
	t.Parallel()

	result, err := ResolveEncounterResult(match.State{
		Encounter: encounter.State{Status: encounter.StatusEscaped, EndReason: encounter.EndReasonDepthEscape, Depth: 4},
		Player:    match.PlayerState{Loadout: mustBuildLoadout(t, 5)},
		Lifecycle: match.LifecycleState{Finished: true},
	}, fishprofiles.Spawn{Profile: fishprofiles.Profile{ID: "tuna", Name: "Atun"}})

	require.NoError(t, err)
	assert.Equal(t, run.EncounterOutcomeEscaped, result.Outcome)
	assert.True(t, result.NodeResolved)
	assert.Equal(t, 1, result.ThreadDamage)
	assert.Nil(t, result.Capture)
}

func TestResolveEncounterResultSplashEscapeDoesNotAddThreadDamage(t *testing.T) {
	t.Parallel()

	result, err := ResolveEncounterResult(match.State{
		Encounter: encounter.State{Status: encounter.StatusEscaped, EndReason: encounter.EndReasonSplashEscape},
		Player:    match.PlayerState{Loadout: mustBuildLoadout(t, 5)},
		Lifecycle: match.LifecycleState{Finished: true},
	}, fishprofiles.Spawn{Profile: fishprofiles.Profile{ID: "tuna", Name: "Atun"}})

	require.NoError(t, err)
	assert.Equal(t, run.EncounterOutcomeEscaped, result.Outcome)
	assert.True(t, result.NodeResolved)
	assert.Equal(t, 0, result.ThreadDamage)
	assert.Nil(t, result.Capture)
}

func TestResolveEncounterResultRejectsUnfinishedMatch(t *testing.T) {
	t.Parallel()

	_, err := ResolveEncounterResult(match.State{}, fishprofiles.Spawn{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "finished match")
}

func mustBuildLoadout(t *testing.T, trackMaxDistance int) loadout.State {
	t.Helper()

	playerRod, err := rod.NewState(rod.Config{
		OpeningMaxDistance: trackMaxDistance,
		OpeningMaxDepth:    1,
		TrackMaxDistance:   trackMaxDistance,
		TrackMaxDepth:      2,
	})
	require.NoError(t, err)
	playerLoadout, err := loadout.NewState(playerRod, nil)
	require.NoError(t, err)

	return playerLoadout
}
