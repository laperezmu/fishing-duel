package waterpools

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultCatalog(t *testing.T) {
	catalog := DefaultCatalog()

	require.Len(t, catalog, 3)
	assert.Equal(t, Shoreline, catalog[0].ID)
	assert.Equal(t, Offshore, catalog[1].ID)
	assert.Equal(t, MixedCurrent, catalog[2].ID)
}

func TestIDValidate(t *testing.T) {
	require.NoError(t, Shoreline.Validate())
	assert.EqualError(t, ID("bogus").Validate(), "unknown water pool \"bogus\"")
}
