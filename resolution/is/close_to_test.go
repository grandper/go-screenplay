package is_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/resolution/is"
	"github.com/grandper/go-screenplay/resolution/testdata"
)

func TestIsCloseToResolution(t *testing.T) {
	const delta = 3.5
	matcher := is.CloseTo(100.0, delta)

	t.Run("should match value less than 2.5", func(t *testing.T) {
		testdata.AssertMatch(t, matcher, 102.0)
		testdata.AssertMatch(t, matcher, 103.5)
	})

	t.Run("fails when the values don't match", func(t *testing.T) {
		testdata.AssertNoMatch(t, matcher, 103.51)
		testdata.AssertNoMatch(t, matcher, 110.0)
	})

	t.Run("returns an error when the types don't match", func(t *testing.T) {
		testdata.AssertMatcherFails(t, matcher, false)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "at most 3.5 away from 100", matcher.String())
	})
}
