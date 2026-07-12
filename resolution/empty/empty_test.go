package empty_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/resolution/empty"
	"github.com/grandper/go-screenplay/resolution/testdata"
)

func TestEmptyCollectionResolution(t *testing.T) {
	matcher := empty.Collection()

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

	t.Run("should match if a pointer points to an empty collection", func(t *testing.T) {
		slice := []*int{}
		ptr := &slice
		testdata.AssertMatch(t, matcher, ptr)
	})

	t.Run("fails when the collection is not empty", func(t *testing.T) {
		testdata.AssertNoMatch(t, matcher, []int{2})
	})

	t.Run("returns an error when the object is not a collection", func(t *testing.T) {
		testdata.AssertMatcherFails(t, matcher, 42)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "an empty collection", matcher.String())
	})
}

func TestEmptyValueResolution(t *testing.T) {
	matcher := empty.Value()

	t.Run("should match the zero value of a number", func(t *testing.T) {
		testdata.AssertMatch(t, matcher, 0)
	})

	t.Run("should match the false boolean", func(t *testing.T) {
		testdata.AssertMatch(t, matcher, false)
	})

	t.Run("should match if the object is nil", func(t *testing.T) {
		testdata.AssertMatch(t, matcher, nil)
	})

	t.Run("should match if a pointer is nil", func(t *testing.T) {
		var ptr *int
		testdata.AssertMatch(t, matcher, ptr)
	})

	t.Run("should match if a struct is empty", func(t *testing.T) {
		testdata.AssertMatch(t, matcher, struct{}{})
	})

	t.Run("fails when the value is not empty", func(t *testing.T) {
		testdata.AssertNoMatch(t, matcher, 2)
		testdata.AssertNoMatch(t, matcher, true)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "an empty value", matcher.String())
	})
}
