package contains_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/resolution/contains"
	"github.com/grandper/go-screenplay/resolution/testdata"
)

func TestContainsTheBytesResolution(t *testing.T) {
	matcher := contains.TheBytes([]byte("lo Wo"))

	t.Run("should contain the bytes", func(t *testing.T) {
		testdata.AssertMatch(t, matcher, []byte("Hello World!"))
	})

	t.Run("fails when the bytes are not in the slice", func(t *testing.T) {
		testdata.AssertNoMatch(t, matcher, []byte("World!"))
	})

	t.Run("returns an error when the value is of the wrong type", func(t *testing.T) {
		testdata.AssertMatcherFails(t, matcher, 2)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "containing the bytes lo Wo", matcher.String())
	})
}
