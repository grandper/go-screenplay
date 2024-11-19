package action_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/extensions/filesystem/action/fixture"

	"github.com/grandper/go-screenplay/extensions/filesystem/ability"
	"github.com/grandper/go-screenplay/extensions/filesystem/action"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestChangeDirectory(t *testing.T) {
	adam := screenplay.ActorNamed("Adam").WhoCan(ability.UseTheFileSystem())

	t.Run("can change directory", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			dir := t.TempDir()
			require.NoError(t, adam.AttemptsTo(action.ChangeDirectoryTo(dir)))
			currentDir, err := os.Getwd()
			require.NoError(t, err)
			assert.Contains(t, currentDir, dir)
		})

		t.Run("fail if the directory does not exist", func(t *testing.T) {
			require.ErrorIs(t, adam.AttemptsTo(action.ChangeDirectoryTo("")), action.ErrDirectoryNotProvided)
		})

		t.Run("fails when the actor has not the ability to use the file system", func(t *testing.T) {
			require.ErrorIs(t, screenplay.ActorNamed("Bob").AttemptsTo(action.ChangeDirectoryTo(t.TempDir())),
				screenplay.ErrActorMissingAbility)
		})

		t.Run("fails when the directory does not exist", func(t *testing.T) {
			require.Error(t, adam.AttemptsTo(action.ChangeDirectoryTo(fixture.CreateDirectoryName(t))))
		})
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		action1 := action.ChangeDirectoryTo("new_dir")
		assert.Equal(t, "change the directory to new_dir", action1.String())
	})
}
