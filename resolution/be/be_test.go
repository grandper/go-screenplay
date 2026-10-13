package be_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/resolution/be"
	"github.com/grandper/go-screenplay/v2/resolution/testdata"
)

func TestBe(t *testing.T) {

	t.Run("equal to", func(t *testing.T) {
		t.Parallel()

		t.Run("is satisfied by an equal value", func(t *testing.T) {
			t.Parallel()
			testdata.AssertMatch(t, be.EqualTo(2.5), 2.5)
			testdata.AssertMatch(t, be.EqualTo("hello"), "hello")
		})

		t.Run("is not satisfied by another value", func(t *testing.T) {
			t.Parallel()
			testdata.AssertNoMatch(t, be.EqualTo(2.5), 3.0)
			require.EqualError(t, be.EqualTo("hello").Inspect("bye"), `got "bye"`)
		})

		t.Run("describes itself", func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, "be equal to 2.5", be.EqualTo(2.5).String())
			assert.Equal(t, `be equal to "hello"`, be.EqualTo("hello").String())
		})
	})

	t.Run("deep equal to", func(t *testing.T) {
		t.Parallel()

		resolution := be.DeepEqualTo([]string{"Adam", "Eve"})

		t.Run("is satisfied by a deeply equal value", func(t *testing.T) {
			t.Parallel()
			testdata.AssertMatch(t, resolution, []string{"Adam", "Eve"})
			testdata.AssertMatch(t, be.DeepEqualTo(map[string]int{"a": 1}), map[string]int{"a": 1})
		})

		t.Run("is not satisfied by another value", func(t *testing.T) {
			t.Parallel()
			testdata.AssertNoMatch(t, resolution, []string{"Adam"})
			testdata.AssertNoMatch(t, resolution, nil)
		})

		t.Run("describes itself", func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, "be deeply equal to [Adam Eve]", resolution.String())
		})
	})

	t.Run("comparisons", func(t *testing.T) {
		t.Parallel()

		t.Run("GreaterThan", func(t *testing.T) {
			t.Parallel()
			testdata.AssertMatch(t, be.GreaterThan(5), 6)
			testdata.AssertNoMatch(t, be.GreaterThan(5), 5)
			testdata.AssertMatch(t, be.GreaterThan("a"), "b")
			assert.Equal(t, "be greater than 5", be.GreaterThan(5).String())
		})

		t.Run("AtLeast", func(t *testing.T) {
			t.Parallel()
			testdata.AssertMatch(t, be.AtLeast(5), 5)
			testdata.AssertNoMatch(t, be.AtLeast(5), 4)
			assert.Equal(t, "be at least 5", be.AtLeast(5).String())
		})

		t.Run("LessThan", func(t *testing.T) {
			t.Parallel()
			testdata.AssertMatch(t, be.LessThan(5), 4)
			testdata.AssertNoMatch(t, be.LessThan(5), 5)
			assert.Equal(t, "be less than 5", be.LessThan(5).String())
		})

		t.Run("AtMost", func(t *testing.T) {
			t.Parallel()
			testdata.AssertMatch(t, be.AtMost(5), 5)
			testdata.AssertNoMatch(t, be.AtMost(5), 6)
			assert.Equal(t, "be at most 5", be.AtMost(5).String())
		})
	})

	t.Run("between", func(t *testing.T) {
		t.Parallel()

		resolution := be.Between(1, 10)

		t.Run("is satisfied by a value within the bounds, both included", func(t *testing.T) {
			t.Parallel()
			testdata.AssertMatch(t, resolution, 1)
			testdata.AssertMatch(t, resolution, 5)
			testdata.AssertMatch(t, resolution, 10)
		})

		t.Run("is not satisfied by a value outside the bounds", func(t *testing.T) {
			t.Parallel()
			testdata.AssertNoMatch(t, resolution, 0)
			testdata.AssertNoMatch(t, resolution, 11)
		})

		t.Run("describes itself", func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, "be between 1 and 10", resolution.String())
		})

		t.Run("is invalid when the lower bound is greater than the upper bound", func(t *testing.T) {
			t.Parallel()

			broken := be.Between(10, 1)

			testdata.AssertInvalid(t, broken, 5)
			assert.Equal(t, "be between 10 and 1", broken.String())
			require.EqualError(t, broken.Inspect(5),
				`invalid resolution "be between 10 and 1": the lower bound is greater than the upper bound`)
		})
	})

	t.Run("close to", func(t *testing.T) {
		t.Parallel()

		resolution := be.CloseTo(3.0, 0.5)

		t.Run("is satisfied by a number within the tolerance, the bounds included", func(t *testing.T) {
			t.Parallel()
			testdata.AssertMatch(t, resolution, 3.0)
			testdata.AssertMatch(t, resolution, 2.5)
			testdata.AssertMatch(t, resolution, 3.5)
			testdata.AssertMatch(t, be.CloseTo(10, 2), 12)
		})

		t.Run("is not satisfied by a number outside the tolerance", func(t *testing.T) {
			t.Parallel()
			testdata.AssertNoMatch(t, resolution, 3.6)
			testdata.AssertNoMatch(t, resolution, math.NaN())
			testdata.AssertNoMatch(t, be.CloseTo(10, 2), 13)
		})

		t.Run("describes itself", func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, "be close to 3 within 0.5", resolution.String())
		})

		t.Run("is invalid when the tolerance is not a positive number", func(t *testing.T) {
			t.Parallel()
			testdata.AssertInvalid(t, be.CloseTo(3.0, -0.1), 3.0)
			testdata.AssertInvalid(t, be.CloseTo(3.0, math.NaN()), 3.0)
		})
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		t.Run("is satisfied by an empty value", func(t *testing.T) {
			t.Parallel()
			testdata.AssertMatch(t, be.Empty[[]string](), []string{})
			testdata.AssertMatch(t, be.Empty[[]string](), nil)
			testdata.AssertMatch(t, be.Empty[map[string]int](), map[string]int{})
			testdata.AssertMatch(t, be.Empty[string](), "")
			testdata.AssertMatch(t, be.Empty[*[]int](), nil)
			testdata.AssertMatch[any](t, be.Empty[any](), []int{})
		})

		t.Run("is not satisfied by a value holding something", func(t *testing.T) {
			t.Parallel()
			testdata.AssertNoMatch(t, be.Empty[[]string](), []string{"Adam"})
			testdata.AssertNoMatch(t, be.Empty[map[string]int](), map[string]int{"a": 1})
			testdata.AssertNoMatch(t, be.Empty[string](), "hello")
			testdata.AssertNoMatch(t, be.Empty[*[]int](), &[]int{1})
		})

		t.Run("is not satisfied by a value without a length", func(t *testing.T) {
			t.Parallel()
			testdata.AssertNoMatch[any](t, be.Empty[any](), 3)
		})

		t.Run("describes itself", func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, "be empty", be.Empty[string]().String())
		})

		t.Run("is invalid when the values of the type have no length", func(t *testing.T) {
			t.Parallel()
			testdata.AssertInvalid(t, be.Empty[int](), 0)
		})
	})

	t.Run("zero", func(t *testing.T) {
		t.Parallel()

		t.Run("is satisfied by the zero value", func(t *testing.T) {
			t.Parallel()
			testdata.AssertMatch(t, be.Zero[int](), 0)
			testdata.AssertMatch(t, be.Zero[string](), "")
			testdata.AssertMatch(t, be.Zero[*int](), nil)
		})

		t.Run("is not satisfied by another value", func(t *testing.T) {
			t.Parallel()
			testdata.AssertNoMatch(t, be.Zero[int](), 3)
			testdata.AssertNoMatch(t, be.Zero[string](), "hello")
		})

		t.Run("describes itself", func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, "be zero", be.Zero[int]().String())
		})
	})

	t.Run("true and false", func(t *testing.T) {
		t.Parallel()

		testdata.AssertMatch(t, be.True(), true)
		testdata.AssertNoMatch(t, be.True(), false)
		assert.Equal(t, "be true", be.True().String())

		testdata.AssertMatch(t, be.False(), false)
		testdata.AssertNoMatch(t, be.False(), true)
		assert.Equal(t, "be false", be.False().String())
	})

	t.Run("equal to, with an interface type parameter", func(t *testing.T) {
		testdata.AssertMatch[any](t, be.EqualTo[any](3), 3)
		testdata.AssertNoMatch[any](t, be.EqualTo[any](3), 4)
		// A slice held in an interface cannot be compared: the resolution is invalid, no panic.
		testdata.AssertInvalid[any](t, be.EqualTo[any]([]int{1}), []int{1})
	})
}
