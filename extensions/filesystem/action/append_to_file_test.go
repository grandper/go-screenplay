package action_test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/extensions/filesystem/action/fixture"

	"github.com/grandper/go-screenplay/extensions/filesystem/ability"
	"github.com/grandper/go-screenplay/extensions/filesystem/action"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestAppendTo(t *testing.T) {
	adam := screenplay.ActorNamed("Adam").WhoCan(ability.UseTheFileSystem())

	t.Run("appends content to a file", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			testAppend(t, func(t *testing.T, filename string) {
				require.NoError(t, adam.AttemptsTo(action.AppendTheContent(strings.NewReader("world")).To(filename)))
			})
		})

		t.Run("fails when the file does not exist", func(t *testing.T) {
			require.Error(t, adam.AttemptsTo(action.AppendTheContent(strings.NewReader("some text")).To("file.txt")))
		})

		t.Run("fails when the file is a directory", func(t *testing.T) {
			require.Error(t, adam.AttemptsTo(action.AppendTheContent(strings.NewReader("some text")).To(t.TempDir())))
		})

		t.Run("fails when the file is not specified", func(t *testing.T) {
			require.ErrorIs(
				t,
				adam.AttemptsTo(action.AppendTheContent(strings.NewReader("some text"))),
				action.ErrFilenameNotProvided,
			)
		})

		t.Run("fails when the content is not specified", func(t *testing.T) {
			require.ErrorIs(
				t,
				adam.AttemptsTo(action.AppendTheContent(nil).To("file.txt")),
				action.ErrReaderNotProvided,
			)
		})

		t.Run("fails when the actor has not the ability to use the file system", func(t *testing.T) {
			require.ErrorIs(
				t,
				screenplay.ActorNamed("Bob").
					AttemptsTo(action.AppendTheContent(strings.NewReader("some text")).To("file.txt")),
				screenplay.ErrActorMissingAbility,
			)
		})
	})

	t.Run("appends text to a file", func(t *testing.T) {
		testAppend(t, func(t *testing.T, filename string) {
			require.NoError(t, adam.AttemptsTo(action.AppendTheText("world").To(filename)))
		})
	})

	t.Run("appends bytes to a file", func(t *testing.T) {
		testAppend(t, func(t *testing.T, filename string) {
			require.NoError(t, adam.AttemptsTo(action.AppendTheBytes([]byte("world")).To(filename)))
		})
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		action1 := action.AppendTheText("some text").To("file.txt")
		assert.Equal(t, "append content to the file file.txt", action1.String())
	})
}

func testAppend(t *testing.T, fn func(t *testing.T, filename string)) {
	t.Helper()
	filename := fixture.CreateFilename(t)
	require.NoError(t, os.WriteFile(filename, []byte("hello "), 0644))

	fn(t, filename)

	content, err := os.ReadFile(filename)
	require.NoError(t, err)
	assert.Equal(t, "hello world", string(content))
}
