package waterpools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	t.Run("returns name for known pool ID", func(t *testing.T) {
		name := Name(Shoreline)

		assert.Equal(t, "Costa cercana", name)
	})

	t.Run("returns ID string for unknown pool ID", func(t *testing.T) {
		name := Name(ID("unknown"))

		assert.Equal(t, "unknown", name)
	})

	t.Run("returns name for offshore", func(t *testing.T) {
		name := Name(Offshore)

		assert.Equal(t, "Mar abierto", name)
	})

	t.Run("returns name for mixed current", func(t *testing.T) {
		name := Name(MixedCurrent)

		assert.Equal(t, "Corriente mixta", name)
	})
}
