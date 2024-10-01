package is_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/resolution/is"
	"github.com/grandper/go-screenplay/resolution/testdata"
)

func TestIsEqualToResolution(t *testing.T) {
	matcher := is.EqualTo(2.5)

	t.Run("should match equal value", func(t *testing.T) {
		testdata.AssertMatch(t, matcher, 2.5)
	})

	t.Run("fails when the values don't match", func(t *testing.T) {
		testdata.AssertNoMatch(t, matcher, 3.0)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "equal to 2.5", matcher.String())
	})
}
