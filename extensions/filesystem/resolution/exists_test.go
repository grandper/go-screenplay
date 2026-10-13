package resolution_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/v2/extensions/filesystem/question"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/resolution"
	"github.com/grandper/go-screenplay/v2/resolution/testdata"
)

func TestExistResolution(t *testing.T) {
	fileMatcher := resolution.Exist[*question.File]()
	directoryMatcher := resolution.Exist[*question.Directory]()
	existingFile := question.CreateFakeExistingFile("existing_file.txt")
	missingFile := question.CreateFakeMissingFile("missing_file.txt")
	existingDirectory := question.CreateFakeExistingDirectory("existing_directory")
	missingDirectory := question.CreateFakeMissingDirectory("missing_directory")

	t.Run("matches an existing file or directory", func(t *testing.T) {
		testdata.AssertMatch(t, fileMatcher, existingFile)
		testdata.AssertMatch(t, directoryMatcher, existingDirectory)
	})

	t.Run("does not match a missing file or directory", func(t *testing.T) {
		testdata.AssertNoMatch(t, fileMatcher, missingFile)
		testdata.AssertNoMatch(t, directoryMatcher, missingDirectory)
	})

	t.Run("does not match a nil file or directory", func(t *testing.T) {
		testdata.AssertNoMatch(t, fileMatcher, nil)
		testdata.AssertNoMatch(t, directoryMatcher, nil)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "exist", fileMatcher.String())
	})
}
