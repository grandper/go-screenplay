package timing_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/timing"
)

// fakeParent is a stand-in parent builder used to verify that the fluent chain
// returns the parent and that the computed duration is written into its field.
type fakeParent struct {
	duration time.Duration
}

func TestDurationBuilder(t *testing.T) {
	t.Run("returns the parent so the chain can continue", func(t *testing.T) {
		parent := &fakeParent{}
		builder := timing.NewDurationBuilder(parent, &parent.duration).For(5)
		assert.Same(t, parent, builder.Seconds())
	})

	t.Run("writes the computed duration for each unit", func(t *testing.T) {
		cases := []struct {
			name     string
			build    func(p *fakeParent) *fakeParent
			expected time.Duration
		}{
			{"milliseconds", func(p *fakeParent) *fakeParent {
				return timing.NewDurationBuilder(p, &p.duration).For(100).Milliseconds()
			}, 100 * time.Millisecond},
			{"seconds", func(p *fakeParent) *fakeParent {
				return timing.NewDurationBuilder(p, &p.duration).For(3).Seconds()
			}, 3 * time.Second},
			{"minutes", func(p *fakeParent) *fakeParent {
				return timing.NewDurationBuilder(p, &p.duration).For(2).Minutes()
			}, 2 * time.Minute},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				p := &fakeParent{}
				tc.build(p)
				assert.Equal(t, tc.expected, p.duration)
			})
		}
	})

	t.Run("singular and plural units are equivalent", func(t *testing.T) {
		p1, p2 := &fakeParent{}, &fakeParent{}
		timing.NewDurationBuilder(p1, &p1.duration).For(1).Second()
		timing.NewDurationBuilder(p2, &p2.duration).For(1).Seconds()
		assert.Equal(t, p1.duration, p2.duration)

		p3, p4 := &fakeParent{}, &fakeParent{}
		timing.NewDurationBuilder(p3, &p3.duration).For(1).Millisecond()
		timing.NewDurationBuilder(p4, &p4.duration).For(1).Milliseconds()
		assert.Equal(t, p3.duration, p4.duration)

		p5, p6 := &fakeParent{}, &fakeParent{}
		timing.NewDurationBuilder(p5, &p5.duration).For(1).Minute()
		timing.NewDurationBuilder(p6, &p6.duration).For(1).Minutes()
		assert.Equal(t, p5.duration, p6.duration)
	})

	t.Run("every wording alias sets the amount identically", func(t *testing.T) {
		aliases := map[string]func(*timing.DurationBuilder[fakeParent], int) *timing.DurationBuilder[fakeParent]{
			"For":                   (*timing.DurationBuilder[fakeParent]).For,
			"During":                (*timing.DurationBuilder[fakeParent]).During,
			"TryingFor":             (*timing.DurationBuilder[fakeParent]).TryingFor,
			"TryingForNoLongerThan": (*timing.DurationBuilder[fakeParent]).TryingForNoLongerThan,
			"WaitingFor":            (*timing.DurationBuilder[fakeParent]).WaitingFor,
			"Every":                 (*timing.DurationBuilder[fakeParent]).Every,
			"PollingEvery":          (*timing.DurationBuilder[fakeParent]).PollingEvery,
			"TryingEvery":           (*timing.DurationBuilder[fakeParent]).TryingEvery,
		}
		for name, setAmount := range aliases {
			t.Run(name, func(t *testing.T) {
				p := &fakeParent{}
				parent := setAmount(timing.NewDurationBuilder(p, &p.duration), 3).Seconds()
				assert.Same(t, p, parent)
				assert.Equal(t, 3*time.Second, p.duration)
			})
		}
	})

	t.Run("describes the time frame using the amount and the chosen unit", func(t *testing.T) {
		parent := &fakeParent{}
		builder := timing.NewDurationBuilder(parent, &parent.duration).For(30)
		builder.Seconds()
		assert.Equal(t, "30 seconds", builder.String())
	})

	t.Run("uses the singular unit for a single time unit and the plural otherwise", func(t *testing.T) {
		p1 := &fakeParent{}
		singular := timing.NewDurationBuilder(p1, &p1.duration).For(1)
		singular.Seconds()
		assert.Equal(t, "1 second", singular.String())

		p2 := &fakeParent{}
		plural := timing.NewDurationBuilder(p2, &p2.duration).For(20)
		plural.Milliseconds()
		assert.Equal(t, "20 milliseconds", plural.String())
	})

	t.Run("a zero amount yields a zero duration", func(t *testing.T) {
		p := &fakeParent{}
		timing.NewDurationBuilder(p, &p.duration).Seconds()
		assert.Equal(t, time.Duration(0), p.duration)
	})

	t.Run("does not panic when the duration pointer is nil", func(t *testing.T) {
		parent := &fakeParent{}
		assert.NotPanics(t, func() {
			timing.NewDurationBuilder(parent, nil).For(1).Seconds()
		})
	})
}
