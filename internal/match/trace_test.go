package match

import (
	"testing"

	"pesca/internal/cards"
	"pesca/internal/encounter"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewResolutionTraceSnapshot(t *testing.T) {
	encounterState, err := encounter.NewState(encounter.DefaultConfig())
	require.NoError(t, err)
	before := State{Encounter: encounterState}
	after := before
	after.Encounter.Distance = 1
	trace := NewResolutionTraceSnapshot(before, after, []ResolvedEffectState{{Owner: cards.OwnerPlayer, Trigger: cards.TriggerOnDraw, Type: cards.EffectTypeAdvanceHorizontal, Priority: 50}})

	assert.Equal(t, 3, trace.Before.Track.Distance)
	assert.Equal(t, 1, trace.After.Track.Distance)
	assert.Len(t, trace.ResolvedEffects, 1)
	assert.Equal(t, cards.OwnerPlayer, trace.ResolvedEffects[0].Owner)
}
