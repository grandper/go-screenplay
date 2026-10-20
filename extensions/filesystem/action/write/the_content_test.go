package write_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/extensions/filesystem/ability"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/action"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/action/write"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/internal/fixture"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestTheContent(t *testing.T) {
	adam := screenplay.ActorNamed("Adam").WhoCan(ability.UseTheFileSystem())

	t.Run("appends content to a file", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			testAppend(t, func(t *testing.T, filename string) {
				require.NoError(t, adam.AttemptsTo(write.TheContent(strings.NewReader("world")).AtTheEndOf(filename)))
			})
		})

		t.Run("fails when the file does not exist", func(t *testing.T) {
			require.Error(t, adam.AttemptsTo(write.TheContent(strings.NewReader("some text")).AtTheEndOf("file.txt")))
		})

		t.Run("fails when the file is a directory", func(t *testing.T) {
			require.Error(t, adam.AttemptsTo(write.TheContent(strings.NewReader("some text")).AtTheEndOf(t.TempDir())))
		})

		t.Run("fails when the file is not specified", func(t *testing.T) {
			require.ErrorIs(
				t,
				adam.AttemptsTo(write.TheContent(strings.NewReader("some text"))),
				action.ErrFilenameNotProvided,
			)
		})

		t.Run("fails when the content is not specified", func(t *testing.T) {
			require.ErrorIs(
				t,
				adam.AttemptsTo(write.TheContent(nil).AtTheEndOf("file.txt")),
				action.ErrReaderNotProvided,
			)
		})

		t.Run("fails when the actor has not the ability to use the file system", func(t *testing.T) {
			require.ErrorIs(
				t,
				screenplay.ActorNamed("Bob").
					AttemptsTo(write.TheContent(strings.NewReader("some text")).AtTheEndOf("file.txt")),
				screenplay.ErrActorMissingAbility,
			)
		})

		t.Run("fails when the content cannot be read", func(t *testing.T) {
			errDiskGone := errors.New("the disk is gone")
			notes := filepath.Join(t.TempDir(), "notes.txt")
			require.NoError(t, os.WriteFile(notes, []byte("hello "), 0o600))

			writeTheNotes := write.TheContent(iotest.ErrReader(errDiskGone)).AtTheEndOf(notes)

			assert.Equal(t, fmt.Sprintf("write content that could not be read at the end of %q", notes),
				writeTheNotes.String())
			err := adam.AttemptsTo(writeTheNotes)
			require.ErrorIs(t, err, errDiskGone)
			require.EqualError(t, err, "failed to read the content: the disk is gone")
			content, err := os.ReadFile(notes)
			require.NoError(t, err)
			assert.Equal(t, "hello ", string(content))
		})
	})

	t.Run("appends text to a file", func(t *testing.T) {
		testAppend(t, func(t *testing.T, filename string) {
			require.NoError(t, adam.AttemptsTo(write.TheText("world").AtTheEndOf(filename)))
		})
	})

	t.Run("appends bytes to a file", func(t *testing.T) {
		testAppend(t, func(t *testing.T, filename string) {
			require.NoError(t, adam.AttemptsTo(write.TheBytes([]byte("world")).AtTheEndOf(filename)))
		})
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, `write the text "some text" at the end of "file.txt"`,
			write.TheText("some text").AtTheEndOf("file.txt").String())
		assert.Equal(t, `write the bytes "some bytes" at the end of "file.txt"`,
			write.TheBytes([]byte("some bytes")).AtTheEndOf("file.txt").String())
		assert.Equal(t, `write the content "some content" at the end of "file.txt"`,
			write.TheContent(strings.NewReader("some content")).AtTheEndOf("file.txt").String())
	})

	t.Run("can be performed more than once", func(t *testing.T) {
		filename := filepath.Join(t.TempDir(), "notes.txt")
		require.NoError(t, os.WriteFile(filename, nil, 0o600))

		appendHello := write.TheContent(strings.NewReader("hello ")).AtTheEndOf(filename)
		require.NoError(t, adam.AttemptsTo(appendHello, appendHello))

		content, err := os.ReadFile(filename)
		require.NoError(t, err)
		assert.Equal(t, "hello hello ", string(content))
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
