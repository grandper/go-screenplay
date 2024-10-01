package is_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/resolution/is"
	"github.com/grandper/go-screenplay/resolution/testdata"
)

func TestIsInRangeResolution(t *testing.T) {
	const lowerBound = 10.0
	const upperBound = 21.2
	matcher := is.InRange(lowerBound, upperBound)

	t.Run("should match value less than 2.5", func(t *testing.T) {
		testdata.AssertMatch(t, matcher, 10.0)
		testdata.AssertMatch(t, matcher, 21.2)
		testdata.AssertMatch(t, matcher, 15.0)
	})

	t.Run("fails when the values don't match", func(t *testing.T) {
		testdata.AssertNoMatch(t, matcher, 4.0)
		testdata.AssertNoMatch(t, matcher, 9.99)
		testdata.AssertNoMatch(t, matcher, 21.3)
		testdata.AssertNoMatch(t, matcher, 100.0)
	})

	t.Run("returns an error when the types don't match", func(t *testing.T) {
		testdata.AssertMatcherFails(t, matcher, false)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "in the range [10, 21.2]", matcher.String())
	})
}
