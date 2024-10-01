package contains_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/resolution/contains"
	"github.com/grandper/go-screenplay/resolution/testdata"
)

func TestContainsTheEntryResolution(t *testing.T) {
	matcher := contains.TheEntry("hello", "world")

	t.Run("should match if the collection contains the entry", func(t *testing.T) {
		testdata.AssertMatch(t, matcher, map[string]string{"hello": "world"})

		m := map[string]string{"hello": "world"}
		testdata.AssertMatch(t, matcher, &m)
	})

	t.Run("fails when the object is nil", func(t *testing.T) {
		testdata.AssertNoMatch(t, matcher, nil)

		var slice *[]int

		testdata.AssertNoMatch(t, matcher, slice)
	})

	t.Run("fails when the collection is empty", func(t *testing.T) {
		testdata.AssertNoMatch(t, matcher, map[string]int{})
	})

	t.Run("fails when the map doesn't contain the entry", func(t *testing.T) {
		testdata.AssertNoMatch(t, matcher, map[string]string{"foo": "bar"})
		testdata.AssertNoMatch(t, matcher, map[string]string{"hello": "bar"})
		testdata.AssertNoMatch(t, matcher, map[string]string{"foo": "world"})
	})

	t.Run("returns an error when the value is of the wrong type", func(t *testing.T) {
		testdata.AssertMatcherFails(t, matcher, 2)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "containing the entry [hello: world]", matcher.String())
	})
}
