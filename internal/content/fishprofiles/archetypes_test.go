package fishprofiles

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultArchetypes(t *testing.T) {
	archetypes := DefaultArchetypes()

	require.Len(t, archetypes, 7)
	assert.Equal(t, ArchetypeBaselineCycle, archetypes[0].ID)
	assert.Equal(t, ArchetypeHybridPressure, archetypes[6].ID)
}

func TestArchetypeIDValidate(t *testing.T) {
	require.NoError(t, ArchetypeBaselineCycle.Validate())
	assert.EqualError(t, ArchetypeID("bogus").Validate(), "unknown fish archetype \"bogus\"")
}

func TestArchetypeName(t *testing.T) {
	assert.Equal(t, "Tempo de apertura", Name(ArchetypeDrawTempo))
	assert.Equal(t, "bogus", Name(ArchetypeID("bogus")))
}
