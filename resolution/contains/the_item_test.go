package contains_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/resolution/contains"
	"github.com/grandper/go-screenplay/resolution/testdata"
)

func TestContainsTheItemResolution(t *testing.T) {
	matcher := contains.TheItem(3)

	createChannel := func(numResults int) <-chan int {
		out := make(chan int)
		go func() {
			for i := range numResults {
				out <- i
			}
			close(out)
		}()

		return out
	}

	t.Run("should match if we pass a matching value", func(t *testing.T) {
		testdata.AssertMatch(t, matcher, 3)
	})

	t.Run("should match if the iterable contains the value", func(t *testing.T) {
		testdata.AssertMatch(t, matcher, createChannel(4))
		testdata.AssertMatch(t, matcher, []int{3})

		slice := []int{3}
		testdata.AssertMatch(t, matcher, &slice)
		testdata.AssertMatch(t, matcher, []int{1, 3, 5})
	})

	t.Run("fails when the object is nil", func(t *testing.T) {
		testdata.AssertNoMatch(t, matcher, nil)

		var slice *[]int

		testdata.AssertNoMatch(t, matcher, slice)
	})

	t.Run("fails when the iterable is empty", func(t *testing.T) {
		testdata.AssertNoMatch(t, matcher, []int{})
		testdata.AssertNoMatch(t, matcher, map[string]int{})
	})

	t.Run("fails when the iterable doesn't contains the item", func(t *testing.T) {
		testdata.AssertNoMatch(t, matcher, createChannel(3))
		testdata.AssertNoMatch(t, matcher, []int{2})

		slice := []int{2}
		testdata.AssertNoMatch(t, matcher, &slice)
		testdata.AssertNoMatch(t, matcher, []int{1, 5})
		testdata.AssertNoMatch(t, matcher, map[string]int{"a": 1, "c": 5})
	})

	t.Run("fails when the map are used", func(t *testing.T) {
		testdata.AssertNoMatch(t, matcher, map[string]int{"a": 1, "b": 2, "c": 3})
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "containing the item 3", matcher.String())
	})
}
