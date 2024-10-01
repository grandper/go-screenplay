package contains_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/resolution/contains"
	"github.com/grandper/go-screenplay/resolution/testdata"
)

func TestContainsTheTextResolution(t *testing.T) {
	matcher := contains.TheText("lo Wo")

	t.Run("should contain the text", func(t *testing.T) {
		testdata.AssertMatch(t, matcher, "Hello World!")
	})

	t.Run("fails when the text is not in the string", func(t *testing.T) {
		testdata.AssertNoMatch(t, matcher, "World!")
	})

	t.Run("returns an error when the value is of the wrong type", func(t *testing.T) {
		testdata.AssertMatcherFails(t, matcher, 2)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "containing the text lo Wo", matcher.String())
	})
}
