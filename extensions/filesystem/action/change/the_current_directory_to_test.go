package change_test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/ability"
	filesystemaction "github.com/grandper/go-screenplay/v2/extensions/filesystem/action"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/action/change"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/action/create"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/action/fixture"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestTheCurrentDirectoryTo(t *testing.T) {
	t.Run("can change directory", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			processDir, err := os.Getwd()
			require.NoError(t, err)

			dir := t.TempDir()
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.UseTheFileSystem())
			require.NoError(t, adam.AttemptsTo(change.TheCurrentDirectoryTo(dir)))

			var useTheFileSystem *ability.UseTheFileSystemAbility
			err = adam.AttemptsTo(action.FromFunc("use the ability", func(theActor *screenplay.Actor) error {
				useTheFileSystem = screenplay.UseAnAbilityOf(theActor).That(ability.UseTheFileSystem())

				return nil
			}))
			require.NoError(t, err)
			currentDir, err := useTheFileSystem.ToRetrieveCurrentPath()
			require.NoError(t, err)
			assert.Equal(t, dir, currentDir)

			// The working directory of the process is never changed.
			stillProcessDir, err := os.Getwd()
			require.NoError(t, err)
			assert.Equal(t, processDir, stillProcessDir)
		})

		t.Run("so that later relative paths resolve against the new directory", func(t *testing.T) {
			dir := t.TempDir()
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.UseTheFileSystem())
			require.NoError(t, adam.AttemptsTo(
				change.TheCurrentDirectoryTo(dir),
				create.TheFile("notes.txt").ContainingTheText("hello"),
			))

			content, err := os.ReadFile(filepath.Join(dir, "notes.txt"))
			require.NoError(t, err)
			assert.Equal(t, "hello", string(content))
		})

		t.Run("for this actor only", func(t *testing.T) {
			// Run with -race: each actor has its own current directory, so two
			// actors changing directory at the same time never disturb each other.
			dirA, dirB := t.TempDir(), t.TempDir()
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.UseTheFileSystem())
			wanda := screenplay.ActorNamed("Wanda").WhoCan(ability.UseTheFileSystem())

			var group sync.WaitGroup
			errs := make([]error, 2)
			group.Add(2)
			go func() {
				defer group.Done()
				errs[0] = adam.AttemptsTo(change.TheCurrentDirectoryTo(dirA), create.TheFile("notes.txt"))
			}()
			go func() {
				defer group.Done()
				errs[1] = wanda.AttemptsTo(change.TheCurrentDirectoryTo(dirB), create.TheFile("notes.txt"))
			}()
			group.Wait()
			require.NoError(t, errs[0])
			require.NoError(t, errs[1])

			assert.FileExists(t, filepath.Join(dirA, "notes.txt"))
			assert.FileExists(t, filepath.Join(dirB, "notes.txt"))
		})

		t.Run("fail if the directory does not exist", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.UseTheFileSystem())
			require.ErrorIs(t, adam.AttemptsTo(change.TheCurrentDirectoryTo("")), filesystemaction.ErrDirectoryNameNotProvided)
		})

		t.Run("fails when the actor has not the ability to use the file system", func(t *testing.T) {
			require.ErrorIs(t, screenplay.ActorNamed("Bob").AttemptsTo(change.TheCurrentDirectoryTo(t.TempDir())),
				screenplay.ErrActorMissingAbility)
		})

		t.Run("fails when the directory does not exist", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.UseTheFileSystem())
			require.Error(t, adam.AttemptsTo(change.TheCurrentDirectoryTo(fixture.CreateDirectoryName(t))))
		})
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		action1 := change.TheCurrentDirectoryTo("new_dir")
		assert.Equal(t, `change the current directory to "new_dir"`, action1.String())
	})
}
