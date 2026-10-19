package overwrite_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/extensions/filesystem/ability"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/action"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/action/overwrite"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/internal/fixture"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestTheFile(t *testing.T) {
	adam := screenplay.ActorNamed("Adam").WhoCan(ability.UseTheFileSystem())

	t.Run("overwrites the content of a file", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			testOverwrite(t, func(t *testing.T, filename string) {
				require.NoError(
					t,
					adam.AttemptsTo(overwrite.TheFile(filename).WithTheContent(strings.NewReader("world"))),
				)
			})
		})

		t.Run("fails when the file does not exist", func(t *testing.T) {
			require.Error(
				t,
				adam.AttemptsTo(
					overwrite.TheFile(fixture.CreateFilename(t)).WithTheContent(strings.NewReader("some text")),
				),
			)
		})

		t.Run("fails when the file is a directory", func(t *testing.T) {
			require.Error(
				t,
				adam.AttemptsTo(overwrite.TheFile(t.TempDir()).WithTheContent(strings.NewReader("some text"))),
			)
		})

		t.Run("fails when the file is not specified", func(t *testing.T) {
			require.ErrorIs(
				t,
				adam.AttemptsTo(overwrite.TheFile("").WithTheContent(strings.NewReader("some text"))),
				action.ErrFilenameNotProvided,
			)
		})

		t.Run("fails when the content is not specified", func(t *testing.T) {
			require.ErrorIs(
				t,
				adam.AttemptsTo(overwrite.TheFile(fixture.CreateFilename(t)).WithTheContent(nil)),
				action.ErrContentNotProvided,
			)
		})

		t.Run("fails when the actor has not the ability to use the file system", func(t *testing.T) {
			require.ErrorIs(
				t,
				screenplay.ActorNamed("Bob").
					AttemptsTo(overwrite.TheFile(fixture.CreateFilename(t)).WithTheContent(strings.NewReader("some text"))),
				screenplay.ErrActorMissingAbility,
			)
		})
	})

	t.Run("overwrites the text in a file", func(t *testing.T) {
		testOverwrite(t, func(t *testing.T, filename string) {
			require.NoError(t, adam.AttemptsTo(overwrite.TheFile(filename).WithTheText("world")))
		})
	})

	t.Run("overwrites the bytes in a file", func(t *testing.T) {
		testOverwrite(t, func(t *testing.T, filename string) {
			require.NoError(t, adam.AttemptsTo(overwrite.TheFile(filename).WithTheBytes([]byte("world"))))
		})
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, `overwrite the file "notes.txt" with the text "x"`,
			overwrite.TheFile("notes.txt").WithTheText("x").String())
		assert.Equal(t, `overwrite the file "notes.txt" with the bytes "x"`,
			overwrite.TheFile("notes.txt").WithTheBytes([]byte("x")).String())
		assert.Equal(t, `overwrite the file "notes.txt" with the content "x"`,
			overwrite.TheFile("notes.txt").WithTheContent(strings.NewReader("x")).String())
	})

	t.Run("can be performed more than once", func(t *testing.T) {
		filename := filepath.Join(t.TempDir(), "notes.txt")
		require.NoError(t, os.WriteFile(filename, []byte("original"), 0o600))

		overwrite := overwrite.TheFile(filename).WithTheContent(strings.NewReader("bye"))
		require.NoError(t, adam.AttemptsTo(overwrite))
		require.NoError(t, os.WriteFile(filename, []byte("original again"), 0o600))
		require.NoError(t, adam.AttemptsTo(overwrite))

		content, err := os.ReadFile(filename)
		require.NoError(t, err)
		assert.Equal(t, "bye", string(content))
	})
}

func testOverwrite(t *testing.T, fn func(t *testing.T, filename string)) {
	t.Helper()
	filename := fixture.CreateFilename(t)
	require.NoError(t, os.WriteFile(filename, []byte("hello"), 0666))

	fn(t, filename)

	content, err := os.ReadFile(filename)
	require.NoError(t, err)
	assert.Equal(t, "world", string(content))
}
