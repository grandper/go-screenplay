package greater_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/resolution/greater"
	"github.com/grandper/go-screenplay/resolution/testdata"
)

func TestGreaterThanResolution(t *testing.T) {
	matcher := greater.Than(2.5)

	t.Run("should match value greater than 2.5", func(t *testing.T) {
		testdata.AssertMatch(t, matcher, 3.0)
	})

	t.Run("fails when the values don't match", func(t *testing.T) {
		testdata.AssertNoMatch(t, matcher, 2.0)
		testdata.AssertNoMatch(t, matcher, 2.5)
	})

	t.Run("returns an error when the types don't match", func(t *testing.T) {
		testdata.AssertMatcherFails(t, matcher, false)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "greater than 2.5", matcher.String())
	})
}
