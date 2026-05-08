package encounter_test

import (
	"testing"

	"pesca/internal/encounter"

	"github.com/stretchr/testify/assert"
)

func TestNewCastController(t *testing.T) {
	t.Run("creates controller with slot width of 1 when slot width is not positive", func(t *testing.T) {
		controller := encounter.NewCastController(0, nil)

		assert.Equal(t, 5, controller.TotalSlots())
	})

	t.Run("calculates total slots as bands times slot width", func(t *testing.T) {
		controller := encounter.NewCastController(3, nil)

		assert.Equal(t, 15, controller.TotalSlots())
	})

	t.Run("uses provided positions", func(t *testing.T) {
		positions := []int{2, 4, 6, 8}
		controller := encounter.NewCastController(2, positions)

		assert.Equal(t, 10, controller.TotalSlots())
	})
}

func TestCastControllerTotalSlots(t *testing.T) {
	t.Run("returns total slots as bands times slot width", func(t *testing.T) {
		controller := encounter.NewCastController(3, nil)

		assert.Equal(t, 15, controller.TotalSlots())
	})
}

func TestCastControllerCurrentPosition(t *testing.T) {
	t.Run("returns current position from positions array", func(t *testing.T) {
		positions := []int{1, 3, 5, 7, 5, 3, 1}
		controller := encounter.NewCastController(2, positions)

		position := controller.CurrentPosition()

		assert.Equal(t, 1, position)
	})

	t.Run("returns 0 when positions array is empty", func(t *testing.T) {
		controller := encounter.NewCastController(1, []int{})

		position := controller.CurrentPosition()

		assert.Equal(t, 0, position)
	})
}

func TestCastControllerResolveBand(t *testing.T) {
	t.Run("returns first band for negative position", func(t *testing.T) {
		controller := encounter.NewCastController(2, nil)

		band := controller.ResolveBand(-1)

		assert.NotEmpty(t, band)
	})

	t.Run("resolves band based on position and slot width", func(t *testing.T) {
		controller := encounter.NewCastController(3, nil)

		band1 := controller.ResolveBand(0)
		band2 := controller.ResolveBand(3)
		band3 := controller.ResolveBand(6)

		assert.NotEqual(t, band1, band2)
		assert.NotEqual(t, band2, band3)
	})

	t.Run("clamps to last band when position exceeds total", func(t *testing.T) {
		controller := encounter.NewCastController(2, nil)

		band := controller.ResolveBand(100)

		assert.NotEmpty(t, band)
	})
}

func TestCastControllerAdvance(t *testing.T) {
	t.Run("advances to next position", func(t *testing.T) {
		positions := []int{1, 3, 5}
		controller := encounter.NewCastController(1, positions)

		assert.Equal(t, 1, controller.CurrentPosition())

		next := controller.Advance()

		assert.Equal(t, 3, next)
	})

	t.Run("advances to next position with default positions", func(t *testing.T) {
		controller := encounter.NewCastController(1, nil)

		assert.Equal(t, 0, controller.CurrentPosition())

		next := controller.Advance()

		assert.Equal(t, 1, next)
	})

	t.Run("wraps around to start", func(t *testing.T) {
		positions := []int{1, 3}
		controller := encounter.NewCastController(1, positions)

		controller.Advance()
		next := controller.Advance()

		assert.Equal(t, 1, next)
	})
}
