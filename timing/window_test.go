package timing_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/timing"
)

func TestWindow(t *testing.T) {
	t.Run("captures the total time and the interval between tries", func(t *testing.T) {
		window := timing.NewWindow(2*time.Second, 500*time.Millisecond)
		assert.Equal(t, 2*time.Second, window.Total)
		assert.Equal(t, 500*time.Millisecond, window.Interval)
	})

	t.Run("is valid when the interval is smaller than the total", func(t *testing.T) {
		window := timing.NewWindow(2*time.Second, 500*time.Millisecond)
		assert.True(t, window.Valid())
	})

	t.Run("is valid when the interval equals the total", func(t *testing.T) {
		window := timing.NewWindow(time.Second, time.Second)
		assert.True(t, window.Valid())
	})

	t.Run("is invalid when the interval is larger than the total", func(t *testing.T) {
		window := timing.NewWindow(500*time.Millisecond, 2*time.Second)
		assert.False(t, window.Valid())
	})
}
