package contains_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/resolution/contains"
	"github.com/grandper/go-screenplay/resolution/testdata"
)

func TestContainsTheKeyResolution(t *testing.T) {
	matcher := contains.TheKey("hello")

	t.Run("should match if the collection contains the value", func(t *testing.T) {
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

	t.Run("fails when the map doesn't contain the key", func(t *testing.T) {
		testdata.AssertNoMatch(t, matcher, map[string]string{"foo": "bar"})
	})

	t.Run("returns an error when the value is of the wrong type", func(t *testing.T) {
		testdata.AssertMatcherFails(t, matcher, 2)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "containing the key hello", matcher.String())
	})
}
