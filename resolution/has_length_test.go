package resolution_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/resolution"
	"github.com/grandper/go-screenplay/resolution/testdata"
)

func TestHasLengthResolution(t *testing.T) {
	matcher0 := resolution.HasLength(0)
	matcher1 := resolution.HasLength(1)
	matcher2 := resolution.HasLength(2)

	t.Run("should match if the collection has the right length", func(t *testing.T) {
		testdata.AssertMatch(t, matcher0, nil)
		testdata.AssertMatch(t, matcher1, 1)

		var a *int
		b := new(int)
		*b = 1
		testdata.AssertMatch(t, matcher0, a)
		testdata.AssertMatch(t, matcher1, b)

		testdata.AssertMatch(t, matcher0, []int{})
		testdata.AssertMatch(t, matcher1, []int{1})
		testdata.AssertMatch(t, matcher2, []int{1, 2})

		testdata.AssertMatch(t, matcher0, map[int]int{})
		testdata.AssertMatch(t, matcher1, map[int]int{1: 1})
		testdata.AssertMatch(t, matcher2, map[int]int{1: 1, 2: 2})
	})

	t.Run("fails when the collection has not the right size", func(t *testing.T) {
		testdata.AssertNoMatch(t, matcher1, []int{})
		testdata.AssertNoMatch(t, matcher1, []int{1, 2})
		testdata.AssertNoMatch(t, matcher2, []int{})
		testdata.AssertNoMatch(t, matcher2, []int{1})

		testdata.AssertNoMatch(t, matcher1, map[int]int{})
		testdata.AssertNoMatch(t, matcher1, map[int]int{1: 1, 2: 2})
		testdata.AssertNoMatch(t, matcher2, map[int]int{})
		testdata.AssertNoMatch(t, matcher2, map[int]int{1: 1})
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "0 item long", matcher0.String())
		assert.Equal(t, "1 item long", matcher1.String())
		assert.Equal(t, "2 items long", matcher2.String())
	})
}
