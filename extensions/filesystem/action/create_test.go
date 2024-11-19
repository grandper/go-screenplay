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

func TestCreate(t *testing.T) {
	adam := screenplay.ActorNamed("Adam").WhoCan(ability.UseTheFileSystem())

	t.Run("can create a file", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			filename := fixture.CreateFilename(t)
			defer func() {
				assert.NoError(t, os.Remove(filename))
			}()
			require.NoError(t, adam.AttemptsTo(action.Create().TheFile(filename)))
			file, err := os.Open(filename)
			require.NoError(t, err)
			require.NoError(t, file.Close())
		})

		t.Run("fail if the name is empty", func(t *testing.T) {
			require.ErrorIs(t, adam.AttemptsTo(action.Create().TheFile("")), action.ErrFilenameNotProvided)
		})

		t.Run("fails when the actor has not the ability to use the file system", func(t *testing.T) {
			require.ErrorIs(
				t,
				screenplay.ActorNamed("Bob").AttemptsTo(action.Create().TheFile(fixture.CreateFilename(t))),
				screenplay.ErrActorMissingAbility,
			)
		})

		t.Run("fails when the file already exists", func(t *testing.T) {
			filename := fixture.CreateFilename(t)
			defer func() {
				assert.NoError(t, os.Remove(filename))
			}()
			require.NoError(t, adam.AttemptsTo(action.Create().TheFile(filename)))
			require.Error(t, adam.AttemptsTo(action.Create().TheFile(filename)))
		})

		t.Run("with content", func(t *testing.T) {
			filename := fixture.CreateFilename(t)
			defer func() {
				assert.NoError(t, os.Remove(filename))
			}()
			require.NoError(
				t,
				adam.AttemptsTo(action.Create().TheFile(filename).Containing(strings.NewReader("Hello World"))),
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
				adam.AttemptsTo(action.Create().TheFile(filename).ContainingBytes([]byte("Hello World"))),
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
			require.NoError(t, adam.AttemptsTo(action.Create().TheFile(filename).ContainingTheText("Hello World")))
			content, err := os.ReadFile(filename)
			require.NoError(t, err)
			assert.Equal(t, "Hello World", string(content))
		})

		t.Run("implements the stringer interface", func(t *testing.T) {
			filename := fixture.CreateFilename(t)
			action1 := action.Create().TheFile(filename)
			assert.Equal(t, "create the file "+filename, action1.String())

			action2 := action.Create().TheTemporaryFile(filename)
			assert.Equal(t, "create the temporary file "+filename, action2.String())
		})
	})

	t.Run("can create a temporary file", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			var filename string
			require.NoError(t, adam.AttemptsTo(action.Create().TheTemporaryFile("test-*.txt").AndSaveNameTo(&filename)))
			assert.Contains(t, filename, "test-")
		})

		t.Run("successfully with content", func(t *testing.T) {
			var filename string
			require.NoError(
				t,
				adam.AttemptsTo(
					action.Create().
						TheTemporaryFile("test-*.txt").
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
			require.NoError(t, adam.AttemptsTo(action.Create().TheDirectory(dir)))
			require.DirExists(t, dir)
		})

		t.Run("fail if the name is empty", func(t *testing.T) {
			require.ErrorIs(t, adam.AttemptsTo(action.Create().TheDirectory("")), action.ErrDirectoryNameNotProvided)
		})

		t.Run("fails when the actor has not the ability to use the file system", func(t *testing.T) {
			require.ErrorIs(
				t,
				screenplay.ActorNamed("Bob").AttemptsTo(action.Create().TheDirectory(fixture.CreateDirectoryName(t))),
				screenplay.ErrActorMissingAbility,
			)
		})

		t.Run("fails when the directory already exists", func(t *testing.T) {
			dir := fixture.CreateDirectoryName(t)
			defer func() {
				assert.NoError(t, os.Remove(dir))
			}()
			require.NoError(t, adam.AttemptsTo(action.Create().TheDirectory(dir)))
			require.Error(t, adam.AttemptsTo(action.Create().TheDirectory(dir)))
		})

		t.Run("implements the stringer interface", func(t *testing.T) {
			dir := fixture.CreateDirectoryName(t)
			action1 := action.Create().TheDirectory(dir)
			assert.Equal(t, "create the directory "+dir, action1.String())

			action2 := action.Create().TheTemporaryDirectory(dir)
			assert.Equal(t, "create the temporary directory "+dir, action2.String())
		})
	})

	t.Run("can create a temporary directory", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			var dir string
			require.NoError(t, adam.AttemptsTo(action.Create().TheTemporaryDirectory("test-*").AndSaveNameTo(&dir)))
			assert.Contains(t, dir, "test-")
		})
	})
}
