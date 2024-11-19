package action_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/extensions/filesystem/ability"
	"github.com/grandper/go-screenplay/extensions/filesystem/action"
	"github.com/grandper/go-screenplay/extensions/filesystem/action/fixture"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestRemove(t *testing.T) {
	adam := screenplay.ActorNamed("Adam").WhoCan(ability.UseTheFileSystem())

	t.Run("removes a file", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			filename := fixture.CreateFilename(t)
			require.NoError(t, os.WriteFile(filename, []byte("Hello, World!"), 0o666))
			assert.NoError(t, adam.AttemptsTo(action.Remove().TheFile(filename)))
			_, err := os.Stat(filename)
			assert.True(t, os.IsNotExist(err))
		})

		t.Run("fails when the file does not exist", func(t *testing.T) {
			require.Error(t, adam.AttemptsTo(action.Remove().TheFile(fixture.CreateFilename(t))))
		})

		t.Run("fails when the file is a directory", func(t *testing.T) {
			dir := fixture.CreateDirectoryName(t)
			require.NoError(t, os.Mkdir(dir, 0o755))
			require.Error(t, adam.AttemptsTo(action.Remove().TheFile(dir)))
			require.NoError(t, os.Remove(dir))
		})

		t.Run("fails when the file name is not specified", func(t *testing.T) {
			require.ErrorIs(t,
				adam.AttemptsTo(action.Remove().TheFile("")),
				action.ErrFilenameNotProvided,
			)
		})
	})

	t.Run("removes a directory", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			dir := fixture.CreateDirectoryName(t)
			require.NoError(t, os.Mkdir(dir, 0o755))
			assert.NoError(t, adam.AttemptsTo(action.Remove().TheDirectory(dir)))
			_, err := os.Stat(dir)
			assert.True(t, os.IsNotExist(err))
		})

		t.Run("fails when the directory does not exist", func(t *testing.T) {
			require.Error(t, adam.AttemptsTo(action.Remove().TheDirectory(fixture.CreateDirectoryName(t))))
		})

		t.Run("fails when the directory is not a directory", func(t *testing.T) {
			filename := fixture.CreateFilename(t)
			require.NoError(t, os.WriteFile(filename, []byte("Hello, World!"), 0o666))
			require.Error(t, adam.AttemptsTo(action.Remove().TheDirectory(filename)))
			require.NoError(t, os.Remove(filename))
		})

		t.Run("fails when the directory name is not specified", func(t *testing.T) {
			require.ErrorIs(t,
				adam.AttemptsTo(action.Remove().TheDirectory("")),
				action.ErrDirectoryNameNotProvided,
			)
		})
	})

	t.Run("fails when the actor has not the ability to use the file system", func(t *testing.T) {
		require.ErrorIs(t,
			screenplay.ActorNamed("Bob").
				AttemptsTo(action.Remove().TheFile(fixture.CreateFilename(t))),
			screenplay.ErrActorMissingAbility,
		)
		require.ErrorIs(t,
			screenplay.ActorNamed("Bob").
				AttemptsTo(action.Remove().TheDirectory(fixture.CreateDirectoryName(t))),
			screenplay.ErrActorMissingAbility,
		)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		filename := fixture.CreateFilename(t)
		action1 := action.Remove().TheFile(filename)
		assert.Equal(t, fmt.Sprintf("remove the file '%s'", filename), action1.String())

		dir := fixture.CreateDirectoryName(t)
		action2 := action.Remove().TheDirectory(dir)
		assert.Equal(t, fmt.Sprintf("remove the directory '%s'", dir), action2.String())
	})
}
