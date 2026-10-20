package create_test

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
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/action/create"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/internal/fixture"
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

		t.Run("fails when the content cannot be read", func(t *testing.T) {
			errDiskGone := errors.New("the disk is gone")
			notes := filepath.Join(t.TempDir(), "notes.txt")

			createTheNotes := create.TheFile(notes).Containing(iotest.ErrReader(errDiskGone))

			assert.Equal(t, fmt.Sprintf("create the file %q containing content that could not be read", notes),
				createTheNotes.String())
			err := adam.AttemptsTo(createTheNotes)
			require.ErrorIs(t, err, errDiskGone)
			require.EqualError(t, err, "failed to read the content: the disk is gone")
			assert.NoFileExists(t, notes)
		})

		t.Run("ignores a nil content", func(t *testing.T) {
			notes := filepath.Join(t.TempDir(), "notes.txt")

			createTheNotes := create.TheFile(notes).Containing(nil)

			assert.Equal(t, fmt.Sprintf("create the file %q", notes), createTheNotes.String())
			require.NoError(t, adam.AttemptsTo(createTheNotes))
			content, err := os.ReadFile(notes)
			require.NoError(t, err)
			assert.Empty(t, content)
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
			t.Cleanup(func() { assert.NoError(t, os.Remove(filename)) })
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
			t.Cleanup(func() { assert.NoError(t, os.Remove(filename)) })
			assert.Contains(t, filename, "test-")
			content, err := os.ReadFile(filename)
			require.NoError(t, err)
			assert.Equal(t, "Hello World", string(content))
		})

		t.Run("fails when the temporary file pattern holds a path separator", func(t *testing.T) {
			require.EqualError(t,
				adam.AttemptsTo(create.TheTemporaryFile("reports/notes-*.txt")),
				"createtemp reports/notes-*.txt: pattern contains path separator")
			require.EqualError(t,
				adam.AttemptsTo(create.TheTemporaryFile("reports/notes-*.txt").ContainingTheText("hello")),
				"createtemp reports/notes-*.txt: pattern contains path separator")
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
			t.Cleanup(func() { assert.NoError(t, os.Remove(dir)) })
			assert.Contains(t, dir, "test-")
		})

		t.Run("fails when the temporary directory pattern holds a path separator", func(t *testing.T) {
			require.EqualError(t,
				adam.AttemptsTo(create.TheTemporaryDirectory("reports/2026-*")),
				"mkdirtemp reports/2026-*: pattern contains path separator")
		})
	})
}
