package question_test

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/extensions/filesystem/action/fixture"

	"github.com/grandper/go-screenplay/extensions/filesystem/ability"
	"github.com/grandper/go-screenplay/extensions/filesystem/question"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestDirectoryNamedQuestion(t *testing.T) {
	t.Run("fails if the actor does not have the ability to use the file system", func(t *testing.T) {
		theActor := screenplay.ActorNamed("Adam")
		file, err := question.DirectoryNamed("foobar").AnsweredBy(theActor)
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		assert.Nil(t, file)
	})

	t.Run("returns information for existing directory", func(t *testing.T) {
		directoryName := fixture.CreateDirectoryName(t)
		require.NoError(t, os.Mkdir(directoryName, 0777))
		defer func() {
			assert.NoError(t, os.Remove(directoryName))
		}()
		theActor := screenplay.ActorNamed("Adam").WhoCan(ability.UseTheFileSystem())
		d, err := question.DirectoryNamed(directoryName).AnsweredBy(theActor)
		require.NoError(t, err)
		directory, ok := d.(*question.Directory)
		assert.True(t, ok)
		assert.Equal(t, path.Base(directoryName), directory.Name())
		assert.Equal(t, path.Dir(directoryName), directory.Dir())
		assert.True(t, directory.Exists())
	})

	t.Run("returns information for non-existing file", func(t *testing.T) {
		theActor := screenplay.ActorNamed("Adam").WhoCan(ability.UseTheFileSystem())
		d, err := question.DirectoryNamed("foobar").AnsweredBy(theActor)
		require.NoError(t, err)
		directory, ok := d.(*question.Directory)
		assert.True(t, ok)
		assert.Equal(t, "foobar", directory.Name())
		assert.Equal(t, ".", directory.Dir())
		assert.False(t, directory.Exists())
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		directoryNamedFoobar := question.DirectoryNamed("foobar")
		assert.Equal(t, "directory named 'foobar'", directoryNamedFoobar.String())
	})
}
