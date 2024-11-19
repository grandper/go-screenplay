package resolution_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/extensions/filesystem/question"
	"github.com/grandper/go-screenplay/extensions/filesystem/resolution"
	"github.com/grandper/go-screenplay/resolution/testdata"
)

func TestExistsResolution(t *testing.T) {
	matcher := resolution.Exists()
	existingFile := question.CreateFakeExistingFile("existing_file.txt")
	missingFile := question.CreateFakeMissingFile("missing_file.txt")
	existingDirectory := question.CreateFakeExistingDirectory("existing_directory")
	missingDirectory := question.CreateFakeMissingDirectory("missing_directory")

	t.Run("should match the suffix", func(t *testing.T) {
		testdata.AssertMatch(t, matcher, existingFile)
		testdata.AssertMatch(t, matcher, existingDirectory)
	})

	t.Run("fails when the suffix doesn't match", func(t *testing.T) {
		testdata.AssertNoMatch(t, matcher, missingFile)
		testdata.AssertNoMatch(t, matcher, missingDirectory)
		testdata.AssertNoMatch(t, matcher, nil)
	})

	t.Run("returns an error when the value is of the wrong type", func(t *testing.T) {
		testdata.AssertMatcherFails(t, matcher, 2)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "exists", matcher.String())
	})
}
