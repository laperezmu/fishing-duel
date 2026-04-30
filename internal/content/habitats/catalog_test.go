package habitats

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultCatalog(t *testing.T) {
	catalog := DefaultCatalog()

	require.Len(t, catalog, 6)
	assert.Equal(t, Bottom, catalog[0].Tag)
	assert.Equal(t, Surface, catalog[3].Tag)
	assert.Equal(t, Rock, catalog[5].Tag)
}

func TestTagValidate(t *testing.T) {
	require.NoError(t, Bottom.Validate())
	assert.EqualError(t, Tag("bogus").Validate(), "unknown habitat tag \"bogus\"")
}

func TestStrings(t *testing.T) {
	assert.Equal(t, []string{"bottom", "surface"}, Strings([]Tag{Bottom, Surface}))
}

func TestNames(t *testing.T) {
	assert.Equal(t, []string{"Fondo", "Superficie"}, Names([]Tag{Bottom, Surface}))
	assert.Equal(t, "bogus", Name(Tag("bogus")))
}
