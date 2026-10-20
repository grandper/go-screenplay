package remove_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/extensions/filesystem/ability"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/action"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/action/change"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/action/remove"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/internal/fixture"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestTheFile(t *testing.T) {
	adam := screenplay.ActorNamed("Adam").WhoCan(ability.UseTheFileSystem())

	t.Run("removes a file", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			filename := fixture.CreateFilename(t)
			require.NoError(t, os.WriteFile(filename, []byte("Hello, World!"), 0o666))
			assert.NoError(t, adam.AttemptsTo(remove.TheFile(filename)))
			_, err := os.Stat(filename)
			assert.True(t, os.IsNotExist(err))
		})

		t.Run("fails when the file does not exist", func(t *testing.T) {
			filename := fixture.CreateFilename(t)
			err := adam.AttemptsTo(remove.TheFile(filename))
			require.ErrorIs(t, err, action.ErrFileDoesNotExist)
			assert.EqualError(t, err, fmt.Sprintf("file %q does not exist", filename))
		})

		t.Run("fails when the file is a directory", func(t *testing.T) {
			dir := fixture.CreateDirectoryName(t)
			require.NoError(t, os.Mkdir(dir, 0o755))
			require.Error(t, adam.AttemptsTo(remove.TheFile(dir)))
			require.NoError(t, os.Remove(dir))
		})

		t.Run("fails when the file name is not specified", func(t *testing.T) {
			require.ErrorIs(t,
				adam.AttemptsTo(remove.TheFile("")),
				action.ErrFilenameNotProvided,
			)
		})
	})

	t.Run("removes a directory", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			dir := fixture.CreateDirectoryName(t)
			require.NoError(t, os.Mkdir(dir, 0o755))
			assert.NoError(t, adam.AttemptsTo(remove.TheDirectory(dir)))
			_, err := os.Stat(dir)
			assert.True(t, os.IsNotExist(err))
		})

		t.Run("fails when the directory does not exist", func(t *testing.T) {
			dir := fixture.CreateDirectoryName(t)
			err := adam.AttemptsTo(remove.TheDirectory(dir))
			require.ErrorIs(t, err, action.ErrDirectoryDoesNotExist)
			assert.EqualError(t, err, fmt.Sprintf("directory %q does not exist", dir))
		})

		t.Run("fails when the directory is not a directory", func(t *testing.T) {
			filename := fixture.CreateFilename(t)
			require.NoError(t, os.WriteFile(filename, []byte("Hello, World!"), 0o666))
			require.Error(t, adam.AttemptsTo(remove.TheDirectory(filename)))
			require.NoError(t, os.Remove(filename))
		})

		t.Run("fails when the directory name is not specified", func(t *testing.T) {
			require.ErrorIs(t,
				adam.AttemptsTo(remove.TheDirectory("")),
				action.ErrDirectoryNameNotProvided,
			)
		})

		t.Run("refuses to remove the parent of the current directory", func(t *testing.T) {
			err := adam.AttemptsTo(remove.TheDirectory(".."))
			require.ErrorIs(t, err, ability.ErrProtectedDirectory)
			require.EqualError(
				t,
				err,
				`refused to remove a protected directory: ".." contains the current working directory`,
			)

			_, err = os.Stat("..")
			require.NoError(t, err)
		})

		t.Run("refuses to remove the current directory", func(t *testing.T) {
			err := adam.AttemptsTo(remove.TheDirectory("."))
			require.ErrorIs(t, err, ability.ErrProtectedDirectory)
			assert.EqualError(t, err, `refused to remove a protected directory: "." is the current working directory`)
		})

		t.Run("refuses to remove the working directory of the process, however the actor moved", func(t *testing.T) {
			project := filepath.Join(t.TempDir(), "project")
			reports := filepath.Join(project, "reports")
			require.NoError(t, os.MkdirAll(reports, 0o755))
			t.Chdir(reports)

			bob := screenplay.ActorNamed("Bob").WhoCan(ability.UseTheFileSystem())
			require.NoError(t, bob.AttemptsTo(change.TheCurrentDirectoryTo(t.TempDir())))

			require.EqualError(
				t,
				bob.AttemptsTo(remove.TheDirectory(reports)),
				fmt.Sprintf(
					"refused to remove a protected directory: %q is the working directory of the process",
					reports,
				),
			)
			require.EqualError(
				t,
				bob.AttemptsTo(remove.TheDirectory(project)),
				fmt.Sprintf(
					"refused to remove a protected directory: %q contains the working directory of the process",
					project,
				),
			)
			assert.DirExists(t, reports)
		})
	})

	t.Run("fails when the actor has not the ability to use the file system", func(t *testing.T) {
		require.ErrorIs(t,
			screenplay.ActorNamed("Bob").
				AttemptsTo(remove.TheFile(fixture.CreateFilename(t))),
			screenplay.ErrActorMissingAbility,
		)
		require.ErrorIs(t,
			screenplay.ActorNamed("Bob").
				AttemptsTo(remove.TheDirectory(fixture.CreateDirectoryName(t))),
			screenplay.ErrActorMissingAbility,
		)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		filename := fixture.CreateFilename(t)
		action1 := remove.TheFile(filename)
		assert.Equal(t, fmt.Sprintf("remove the file %q", filename), action1.String())

		dir := fixture.CreateDirectoryName(t)
		action2 := remove.TheDirectory(dir)
		assert.Equal(t, fmt.Sprintf("remove the directory %q", dir), action2.String())
	})
}
