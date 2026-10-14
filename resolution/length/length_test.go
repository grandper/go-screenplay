package length_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/resolution/fixture"
	"github.com/grandper/go-screenplay/v2/resolution/length"
)

func TestLength(t *testing.T) {
	t.Parallel()

	t.Run("exactly", func(t *testing.T) {
		t.Parallel()

		resolution := length.Exactly[[]string](2)

		fixture.AssertMatch(t, resolution, []string{"Adam", "Eve"})
		fixture.AssertNoMatch(t, resolution, []string{"Adam"})
		fixture.AssertNoMatch(t, resolution, nil)
		require.EqualError(t, resolution.Inspect([]string{"Adam"}), "got [Adam]")
		assert.Equal(t, "length exactly 2", resolution.String())

		t.Run("measures every kind of value with a length", func(t *testing.T) {
			t.Parallel()

			values := make(chan int, 3)
			values <- 1

			fixture.AssertMatch(t, length.Exactly[string](5), "hello")
			fixture.AssertMatch(t, length.Exactly[map[string]int](1), map[string]int{"a": 1})
			fixture.AssertMatch(t, length.Exactly[[2]int](2), [2]int{1, 2})
			fixture.AssertMatch(t, length.Exactly[chan int](1), values)
			fixture.AssertMatch(t, length.Exactly[*[]int](2), &[]int{1, 2})
			fixture.AssertMatch(t, length.Exactly[*[]int](0), nil)
			fixture.AssertMatch[any](t, length.Exactly[any](2), []int{1, 2})
		})

		t.Run("is not satisfied by a value without a length", func(t *testing.T) {
			t.Parallel()
			fixture.AssertNoMatch[any](t, length.Exactly[any](1), 3)
		})
	})

	t.Run("at least", func(t *testing.T) {
		t.Parallel()

		resolution := length.AtLeast[string](2)

		fixture.AssertMatch(t, resolution, "ab")
		fixture.AssertMatch(t, resolution, "abc")
		fixture.AssertNoMatch(t, resolution, "a")
		fixture.AssertNoMatch[any](t, length.AtLeast[any](0), 3)
		assert.Equal(t, "length at least 2", resolution.String())
	})

	t.Run("at most", func(t *testing.T) {
		t.Parallel()

		resolution := length.AtMost[string](2)

		fixture.AssertMatch(t, resolution, "")
		fixture.AssertMatch(t, resolution, "ab")
		fixture.AssertNoMatch(t, resolution, "abc")
		assert.Equal(t, "length at most 2", resolution.String())
	})

	t.Run("between", func(t *testing.T) {
		t.Parallel()

		resolution := length.Between[string](2, 3)

		fixture.AssertMatch(t, resolution, "ab")
		fixture.AssertMatch(t, resolution, "abc")
		fixture.AssertNoMatch(t, resolution, "a")
		fixture.AssertNoMatch(t, resolution, "abcd")
		assert.Equal(t, "length between 2 and 3", resolution.String())
	})

	t.Run("a length that can never hold is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("when the values of the type have no length", func(t *testing.T) {
			t.Parallel()
			fixture.AssertInvalid(t, length.Exactly[int](1), 1)
			fixture.AssertInvalid(t, length.AtLeast[bool](1), true)
			fixture.AssertInvalid(t, length.AtMost[float64](1), 1.0)
			fixture.AssertInvalid(t, length.Between[struct{}](1, 2), struct{}{})
			require.EqualError(t, length.Exactly[int](1).Inspect(1),
				`invalid resolution "length exactly 1": values of type int have no length`)
		})

		t.Run("when a length is negative", func(t *testing.T) {
			t.Parallel()
			fixture.AssertInvalid(t, length.Exactly[string](-1), "")
			fixture.AssertInvalid(t, length.AtLeast[string](-1), "")
			fixture.AssertInvalid(t, length.AtMost[string](-1), "")
			fixture.AssertInvalid(t, length.Between[string](-1, 2), "")
		})

		t.Run("when the shortest length is greater than the longest one", func(t *testing.T) {
			t.Parallel()
			fixture.AssertInvalid(t, length.Between[string](3, 2), "ab")
		})
	})
}
