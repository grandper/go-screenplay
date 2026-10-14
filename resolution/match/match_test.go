package match_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/v2/resolution/fixture"
	"github.com/grandper/go-screenplay/v2/resolution/match"
)

func TestPattern(t *testing.T) {
	t.Parallel()

	resolution := match.Pattern(regexp.MustCompile(`^A-\d+$`))

	fixture.AssertMatch(t, resolution, "A-42")
	fixture.AssertNoMatch(t, resolution, "B-42")
	assert.Equal(t, "match pattern `^A-\\d+$`", resolution.String())

	t.Run("accepts the pattern as a text", func(t *testing.T) {
		t.Parallel()

		fromText := match.PatternString(`^A-\d+$`)

		fixture.AssertMatch(t, fromText, "A-42")
		fixture.AssertNoMatch(t, fromText, "B-42")
		assert.Equal(t, "match pattern `^A-\\d+$`", fromText.String())
	})

	t.Run("is invalid when the pattern is nil or does not compile", func(t *testing.T) {
		t.Parallel()

		fixture.AssertInvalid(t, match.Pattern(nil), "A-42")

		broken := match.PatternString("A-(")
		fixture.AssertInvalid(t, broken, "A-42")
		assert.Equal(t, "match pattern `A-(`", broken.String())
	})
}
