package create_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/extensions/filesystem/ability"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/action"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/action/create"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/action/fixture"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestTheFile(t *testing.T) {
	adam := screenplay.ActorNamed("Adam").WhoCan(ability.UseTheFileSystem())

	t.Run("can create a file", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			filename := fixture.CreateFilename(t)
			defer func() {
				assert.NoError(t, os.Remove(filename))
			}()
			require.NoError(t, adam.AttemptsTo(create.TheFile(filename)))
			file, err := os.Open(filename)
			require.NoError(t, err)
			require.NoError(t, file.Close())
		})

		t.Run("fail if the name is empty", func(t *testing.T) {
			require.ErrorIs(t, adam.AttemptsTo(create.TheFile("")), action.ErrFilenameNotProvided)
		})

		t.Run("fails when the actor has not the ability to use the file system", func(t *testing.T) {
			require.ErrorIs(
				t,
				screenplay.ActorNamed("Bob").AttemptsTo(create.TheFile(fixture.CreateFilename(t))),
				screenplay.ErrActorMissingAbility,
			)
		})

		t.Run("fails when the file already exists", func(t *testing.T) {
			filename := fixture.CreateFilename(t)
			defer func() {
				assert.NoError(t, os.Remove(filename))
			}()
			require.NoError(t, adam.AttemptsTo(create.TheFile(filename)))
			require.Error(t, adam.AttemptsTo(create.TheFile(filename)))
		})

		t.Run("with content", func(t *testing.T) {
			filename := fixture.CreateFilename(t)
			defer func() {
				assert.NoError(t, os.Remove(filename))
			}()
			require.NoError(
				t,
				adam.AttemptsTo(create.TheFile(filename).Containing(strings.NewReader("Hello World"))),
			)
			content, err := os.ReadFile(filename)
			require.NoError(t, err)
			assert.Equal(t, "Hello World", string(content))
		})

		t.Run("with bytes", func(t *testing.T) {
			filename := fixture.CreateFilename(t)
			defer func() {
				assert.NoError(t, os.Remove(filename))
			}()
			require.NoError(
				t,
				adam.AttemptsTo(create.TheFile(filename).ContainingBytes([]byte("Hello World"))),
			)
			content, err := os.ReadFile(filename)
			require.NoError(t, err)
			assert.Equal(t, "Hello World", string(content))
		})

		t.Run("with the text", func(t *testing.T) {
			filename := fixture.CreateFilename(t)
			defer func() {
				assert.NoError(t, os.Remove(filename))
			}()
			require.NoError(t, adam.AttemptsTo(create.TheFile(filename).ContainingTheText("Hello World")))
			content, err := os.ReadFile(filename)
			require.NoError(t, err)
			assert.Equal(t, "Hello World", string(content))
		})

		t.Run("implements the stringer interface", func(t *testing.T) {
			filename := fixture.CreateFilename(t)
			action1 := create.TheFile(filename)
			assert.Equal(t, fmt.Sprintf("create the file %q", filename), action1.String())

			action2 := create.TheTemporaryFile(filename)
			assert.Equal(t, fmt.Sprintf("create the temporary file %q", filename), action2.String())

			assert.Equal(t, `create the file "notes.txt" containing the text "x"`,
				create.TheFile("notes.txt").ContainingTheText("x").String())
			assert.Equal(t, `create the file "notes.txt" containing the bytes "x"`,
				create.TheFile("notes.txt").ContainingBytes([]byte("x")).String())
			assert.Equal(t, `create the file "notes.txt" containing the content "x"`,
				create.TheFile("notes.txt").Containing(strings.NewReader("x")).String())
		})
	})

	t.Run("can create a temporary file", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			var filename string
			require.NoError(t, adam.AttemptsTo(create.TheTemporaryFile("test-*.txt").AndSaveNameTo(&filename)))
			assert.Contains(t, filename, "test-")
		})

		t.Run("successfully with content", func(t *testing.T) {
			var filename string
			require.NoError(
				t,
				adam.AttemptsTo(
					create.TheTemporaryFile("test-*.txt").
						ContainingTheText("Hello World").
						AndSaveNameTo(&filename),
				),
			)
			assert.Contains(t, filename, "test-")
			content, err := os.ReadFile(filename)
			require.NoError(t, err)
			assert.Equal(t, "Hello World", string(content))
		})
	})

	t.Run("can create a directory", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			dir := fixture.CreateDirectoryName(t)
			defer func() {
				assert.NoError(t, os.Remove(dir))
			}()
			require.NoError(t, adam.AttemptsTo(create.TheDirectory(dir)))
			require.DirExists(t, dir)
		})

		t.Run("fail if the name is empty", func(t *testing.T) {
			require.ErrorIs(t, adam.AttemptsTo(create.TheDirectory("")), action.ErrDirectoryNameNotProvided)
		})

		t.Run("fails when the actor has not the ability to use the file system", func(t *testing.T) {
			require.ErrorIs(
				t,
				screenplay.ActorNamed("Bob").AttemptsTo(create.TheDirectory(fixture.CreateDirectoryName(t))),
				screenplay.ErrActorMissingAbility,
			)
		})

		t.Run("fails when the directory already exists", func(t *testing.T) {
			dir := fixture.CreateDirectoryName(t)
			defer func() {
				assert.NoError(t, os.Remove(dir))
			}()
			require.NoError(t, adam.AttemptsTo(create.TheDirectory(dir)))
			require.Error(t, adam.AttemptsTo(create.TheDirectory(dir)))
		})

		t.Run("implements the stringer interface", func(t *testing.T) {
			dir := fixture.CreateDirectoryName(t)
			action1 := create.TheDirectory(dir)
			assert.Equal(t, fmt.Sprintf("create the directory %q", dir), action1.String())

			action2 := create.TheTemporaryDirectory(dir)
			assert.Equal(t, fmt.Sprintf("create the temporary directory %q", dir), action2.String())
		})
	})

	t.Run("can create a temporary directory", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			var dir string
			require.NoError(t, adam.AttemptsTo(create.TheTemporaryDirectory("test-*").AndSaveNameTo(&dir)))
			assert.Contains(t, dir, "test-")
		})
	})
}
