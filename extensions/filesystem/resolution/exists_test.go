package resolution_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/v2/extensions/filesystem/question"
	filesystemfixture "github.com/grandper/go-screenplay/v2/extensions/filesystem/question/fixture"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/resolution"
	"github.com/grandper/go-screenplay/v2/resolution/fixture"
)

func TestExistResolution(t *testing.T) {
	fileMatcher := resolution.Exist[*question.File]()
	directoryMatcher := resolution.Exist[*question.Directory]()
	existingFile := filesystemfixture.CreateFakeExistingFile("existing_file.txt")
	missingFile := filesystemfixture.CreateFakeMissingFile("missing_file.txt")
	existingDirectory := filesystemfixture.CreateFakeExistingDirectory("existing_directory")
	missingDirectory := filesystemfixture.CreateFakeMissingDirectory("missing_directory")

	t.Run("matches an existing file or directory", func(t *testing.T) {
		fixture.AssertMatch(t, fileMatcher, existingFile)
		fixture.AssertMatch(t, directoryMatcher, existingDirectory)
	})

	t.Run("does not match a missing file or directory", func(t *testing.T) {
		fixture.AssertNoMatch(t, fileMatcher, missingFile)
		fixture.AssertNoMatch(t, directoryMatcher, missingDirectory)
	})

	t.Run("does not match a nil file or directory", func(t *testing.T) {
		fixture.AssertNoMatch(t, fileMatcher, nil)
		fixture.AssertNoMatch(t, directoryMatcher, nil)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "exist", fileMatcher.String())
	})
}
