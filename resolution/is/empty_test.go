package is_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/resolution/is"
	"github.com/grandper/go-screenplay/resolution/testdata"
)

func TestIsEmptyResolution(t *testing.T) {
	matcher := is.Empty()
	t.Run("should match if the collection is empty", func(t *testing.T) {
		testdata.AssertMatch(t, matcher, []int{})
	})

	t.Run("should match if the collection is nil", func(t *testing.T) {
		testdata.AssertMatch(t, matcher, nil)
	})

	t.Run("should match if a pointer is nil", func(t *testing.T) {
		var ptr *int
		testdata.AssertMatch(t, matcher, ptr)
	})

	t.Run("should match if a pointer points to an empty item", func(t *testing.T) {
		slice := []*int{}
		ptr := &slice
		testdata.AssertMatch(t, matcher, ptr)
	})

	t.Run("should not match if a struct is empty", func(t *testing.T) {
		testdata.AssertMatch(t, matcher, struct{}{})
	})

	t.Run("fails when the collection is not empty", func(t *testing.T) {
		testdata.AssertNoMatch(t, matcher, []int{2})
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "an empty collection", matcher.String())
	})
}
