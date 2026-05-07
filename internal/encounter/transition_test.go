package encounter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyDelta(t *testing.T) {
	tests := []struct {
		title         string
		initialState  State
		delta         Delta
		wantDistance  int
		wantDepth     int
		wantEvent     Event
		wantStatus    Status
		wantEndReason EndReason
	}{
		{
			title:        "moves the fish horizontally and vertically when no surface event is triggered",
			initialState: newEncounterState(t),
			delta: Delta{
				DistanceShift: 1,
				DepthShift:    1,
			},
			wantDistance:  4,
			wantDepth:     2,
			wantEvent:     Event{},
			wantStatus:    StatusOngoing,
			wantEndReason: EndReasonNone,
		},
		{
			title:        "triggers a splash event without escape when the fish tries to rise above the surface",
			initialState: newEncounterState(t),
			delta: Delta{
				DepthShift: -2,
			},
			wantDistance:  3,
			wantDepth:     0,
			wantEvent:     Event{Kind: EventKindSplash, Escaped: false},
			wantStatus:    StatusOngoing,
			wantEndReason: EndReasonNone,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			state := test.initialState

			ApplyDelta(&state, test.delta)

			assert.Equal(t, test.wantDistance, state.Distance)
			assert.Equal(t, test.wantDepth, state.Depth)
			assert.Equal(t, test.wantEvent, state.LastEvent)
			assert.Equal(t, test.wantStatus, state.Status)
			assert.Equal(t, test.wantEndReason, state.EndReason)
		})
	}
}

func TestApplySplashResolution(t *testing.T) {
	state := newEncounterState(t)
	ApplyDelta(&state, Delta{DepthShift: -2})
	require.NotNil(t, state.Splash)

	ApplySplashResolution(&state, SplashResolution{SuccessfulJumps: 1})

	assert.Nil(t, state.Splash)
	assert.Equal(t, Event{Kind: EventKindSplash, Escaped: false}, state.LastEvent)
	assert.Equal(t, StatusOngoing, state.Status)

	state = newEncounterState(t)
	ApplyDelta(&state, Delta{DepthShift: -2})
	ApplySplashResolution(&state, SplashResolution{Escaped: true})
	assert.Nil(t, state.Splash)
	assert.Equal(t, Event{Kind: EventKindSplash, Escaped: true}, state.LastEvent)
	assert.Equal(t, StatusEscaped, state.Status)
	assert.Equal(t, EndReasonSplashEscape, state.EndReason)
}

func newEncounterState(t *testing.T) State {
	t.Helper()

	state, err := NewState(DefaultConfig())
	require.NoError(t, err)

	return state
}
