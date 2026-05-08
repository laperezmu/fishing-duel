package app_test

import (
	"testing"

	"pesca/internal/app"

	"github.com/stretchr/testify/assert"
)

func TestDefaultRandomizer(t *testing.T) {
	t.Run("returns a randomizer", func(t *testing.T) {
		randomizer := app.DefaultRandomizer()

		assert.NotNil(t, randomizer)
	})

	t.Run("generates different values", func(t *testing.T) {
		randomizer := app.DefaultRandomizer()

		val1 := randomizer.Intn(100)
		val2 := randomizer.Intn(100)

		assert.GreaterOrEqual(t, val1, 0)
		assert.Less(t, val1, 100)
		assert.GreaterOrEqual(t, val2, 0)
		assert.Less(t, val2, 100)
	})
}

func TestNewSeededRandom(t *testing.T) {
	t.Run("returns seeded randomizer", func(t *testing.T) {
		randomizer := app.NewSeededRandom(42)

		assert.NotNil(t, randomizer)
	})

	t.Run("produces reproducible results", func(t *testing.T) {
		randomizer1 := app.NewSeededRandom(12345)
		randomizer2 := app.NewSeededRandom(12345)

		values1 := make([]int, 10)
		values2 := make([]int, 10)

		for i := 0; i < 10; i++ {
			values1[i] = randomizer1.Intn(1000)
			values2[i] = randomizer2.Intn(1000)
		}

		assert.Equal(t, values1, values2)
	})

	t.Run("produces different results with different seeds", func(t *testing.T) {
		randomizer1 := app.NewSeededRandom(100)
		randomizer2 := app.NewSeededRandom(200)

		val1 := randomizer1.Intn(100)
		val2 := randomizer2.Intn(100)

		assert.NotEqual(t, val1, val2)
	})
}

func TestRandomizerInterface(t *testing.T) {
	t.Run("randomizer implements Float64", func(t *testing.T) {
		randomizer := app.NewSeededRandom(42)

		val := randomizer.Float64()

		assert.GreaterOrEqual(t, val, 0.0)
		assert.Less(t, val, 1.0)
	})

	t.Run("randomizer implements Shuffle", func(t *testing.T) {
		randomizer := app.NewSeededRandom(42)
		slice := []int{1, 2, 3, 4, 5}

		randomizer.Shuffle(5, func(i, j int) {
			slice[i], slice[j] = slice[j], slice[i]
		})

		assert.Len(t, slice, 5)
	})
}
