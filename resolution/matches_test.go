package resolution_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/resolution"
	"github.com/grandper/go-screenplay/resolution/testdata"
)

func TestMatchesResolution(t *testing.T) {
	const regexStr = `H\w{4}`
	regex := regexp.MustCompile(regexStr)
	matcher1 := resolution.Matches(regex)
	matcher2 := resolution.MatchesRegexString(regexStr)

	t.Run("should match the text", func(t *testing.T) {
		testdata.AssertMatch(t, matcher1, "Hello World!")
		testdata.AssertMatch(t, matcher2, "Hello World!")
	})

	t.Run("fails when the text don't match", func(t *testing.T) {
		testdata.AssertNoMatch(t, matcher1, "Hi!")
		testdata.AssertNoMatch(t, matcher2, "Hi!")
	})

	t.Run("returns an error when the value is of the wrong type", func(t *testing.T) {
		testdata.AssertMatcherFails(t, matcher1, 2)
		testdata.AssertMatcherFails(t, matcher2, 2)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, `text matching the pattern H\w{4}`, matcher1.String())
		assert.Equal(t, `text matching the pattern H\w{4}`, matcher2.String())
	})
}
