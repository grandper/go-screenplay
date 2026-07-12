package end_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/resolution/end"
	"github.com/grandper/go-screenplay/resolution/testdata"
)

func TestEndWithResolution(t *testing.T) {
	matcher := end.With("World!")

	t.Run("should match the suffix", func(t *testing.T) {
		testdata.AssertMatch(t, matcher, "Hello World!")
	})

	t.Run("fails when the suffix doesn't match", func(t *testing.T) {
		testdata.AssertNoMatch(t, matcher, "Hello")
	})

	t.Run("returns an error when the value is of the wrong type", func(t *testing.T) {
		testdata.AssertMatcherFails(t, matcher, 2)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "ending with World!", matcher.String())
	})
}
