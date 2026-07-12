package timing_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/timing"
)

// retryParent is a minimal parent builder that embeds a WindowBuilder, the
// way EventuallyAction does, so we can exercise the fluent API end to end.
type retryParent struct {
	*timing.WindowBuilder[retryParent]

	window timing.Window
}

func newRetryParent() *retryParent {
	parent := &retryParent{}
	parent.WindowBuilder = timing.NewWindowBuilder(parent, &parent.window)

	return parent
}

func TestWindowBuilder(t *testing.T) {
	t.Run("writes the total time and the polling interval", func(t *testing.T) {
		parent := newRetryParent()

		parent.For(100).Milliseconds().Polling(10).Milliseconds()

		assert.Equal(t, 100*time.Millisecond, parent.window.Total)
		assert.Equal(t, 10*time.Millisecond, parent.window.Interval)
	})

	t.Run("returns the parent so the chain keeps flowing", func(t *testing.T) {
		parent := newRetryParent()

		assert.Same(t, parent, parent.For(1).Second())
	})

	t.Run("timeout wordings are all aliases", func(t *testing.T) {
		wordings := map[string]func(*retryParent) *timing.DurationBuilder[retryParent]{
			"For":                   func(p *retryParent) *timing.DurationBuilder[retryParent] { return p.For(5) },
			"TryingFor":             func(p *retryParent) *timing.DurationBuilder[retryParent] { return p.TryingFor(5) },
			"TryingForNoLongerThan": func(p *retryParent) *timing.DurationBuilder[retryParent] { return p.TryingForNoLongerThan(5) },
			"WaitingFor":            func(p *retryParent) *timing.DurationBuilder[retryParent] { return p.WaitingFor(5) },
		}

		for name, configure := range wordings {
			t.Run(name, func(t *testing.T) {
				parent := newRetryParent()

				configure(parent).Seconds()

				assert.Equal(t, 5*time.Second, parent.window.Total)
				assert.Zero(t, parent.window.Interval)
			})
		}
	})

	t.Run("polling wordings are all aliases", func(t *testing.T) {
		wordings := map[string]func(*retryParent) *timing.DurationBuilder[retryParent]{
			"Polling":      func(p *retryParent) *timing.DurationBuilder[retryParent] { return p.Polling(5) },
			"PollingEvery": func(p *retryParent) *timing.DurationBuilder[retryParent] { return p.PollingEvery(5) },
			"TryingEvery":  func(p *retryParent) *timing.DurationBuilder[retryParent] { return p.TryingEvery(5) },
		}

		for name, configure := range wordings {
			t.Run(name, func(t *testing.T) {
				parent := newRetryParent()

				configure(parent).Seconds()

				assert.Equal(t, 5*time.Second, parent.window.Interval)
				assert.Zero(t, parent.window.Total)
			})
		}
	})
}
