package is_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/resolution/is"
	"github.com/grandper/go-screenplay/resolution/testdata"
)

func TestIsGreaterThanOrEqualToResolution(t *testing.T) {
	matcher := is.GreaterThanOrEqualTo(2.5)

	t.Run("should match value greater than or equal to 2.5", func(t *testing.T) {
		testdata.AssertMatch(t, matcher, 2.5)
		testdata.AssertMatch(t, matcher, 3.0)
	})

	t.Run("fails when the values don't match", func(t *testing.T) {
		testdata.AssertNoMatch(t, matcher, 2.0)
	})

	t.Run("returns an error when the types don't match", func(t *testing.T) {
		testdata.AssertMatcherFails(t, matcher, false)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "greater than or equal to 2.5", matcher.String())
	})
}
