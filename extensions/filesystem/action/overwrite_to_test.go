package action_test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/extensions/filesystem/ability"
	"github.com/grandper/go-screenplay/extensions/filesystem/action"
	"github.com/grandper/go-screenplay/extensions/filesystem/action/fixture"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestOverwriteTo(t *testing.T) {
	adam := screenplay.ActorNamed("Adam").WhoCan(ability.UseTheFileSystem())

	t.Run("overwrites the content of a file", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			testOverwrite(t, func(t *testing.T, filename string) {
				require.NoError(
					t,
					adam.AttemptsTo(action.OverwriteTo(filename).WithTheContent(strings.NewReader("world"))),
				)
			})
		})

		t.Run("fails when the file does not exist", func(t *testing.T) {
			require.Error(
				t,
				adam.AttemptsTo(
					action.OverwriteTo(fixture.CreateFilename(t)).WithTheContent(strings.NewReader("some text")),
				),
			)
		})

		t.Run("fails when the file is a directory", func(t *testing.T) {
			require.Error(
				t,
				adam.AttemptsTo(action.OverwriteTo(t.TempDir()).WithTheContent(strings.NewReader("some text"))),
			)
		})

		t.Run("fails when the file is not specified", func(t *testing.T) {
			require.ErrorIs(
				t,
				adam.AttemptsTo(action.OverwriteTo("").WithTheContent(strings.NewReader("some text"))),
				action.ErrFilenameNotProvided,
			)
		})

		t.Run("fails when the content is not specified", func(t *testing.T) {
			require.ErrorIs(
				t,
				adam.AttemptsTo(action.OverwriteTo(fixture.CreateFilename(t)).WithTheContent(nil)),
				action.ErrContentNotProvided,
			)
		})

		t.Run("fails when the actor has not the ability to use the file system", func(t *testing.T) {
			require.ErrorIs(
				t,
				screenplay.ActorNamed("Bob").
					AttemptsTo(action.OverwriteTo(fixture.CreateFilename(t)).WithTheContent(strings.NewReader("some text"))),
				screenplay.ErrActorMissingAbility,
			)
		})
	})

	t.Run("overwrites the text in a file", func(t *testing.T) {
		testOverwrite(t, func(t *testing.T, filename string) {
			require.NoError(t, adam.AttemptsTo(action.OverwriteTo(filename).WithTheText("world")))
		})
	})

	t.Run("overwrites the bytes in a file", func(t *testing.T) {
		testOverwrite(t, func(t *testing.T, filename string) {
			require.NoError(t, adam.AttemptsTo(action.OverwriteTo(filename).WithTheBytes([]byte("world"))))
		})
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		filename := fixture.CreateFilename(t)
		action1 := action.OverwriteTo(filename).WithTheText("some text")
		assert.Equal(t, "overwrite the file "+filename, action1.String())
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
