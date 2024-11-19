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

func TestFileNamedQuestion(t *testing.T) {
	t.Run("fails if the actor does not have the ability to use the file system", func(t *testing.T) {
		theActor := screenplay.ActorNamed("Adam")
		file, err := question.FileNamed("foo.bar").AnsweredBy(theActor)
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		assert.Nil(t, file)
	})

	t.Run("returns information for existing file", func(t *testing.T) {
		filename := fixture.CreateFilename(t)
		require.NoError(t, os.WriteFile(filename, []byte("foo"), 0777))
		defer func() {
			assert.NoError(t, os.Remove(filename))
		}()
		theActor := screenplay.ActorNamed("Adam").WhoCan(ability.UseTheFileSystem())
		f, err := question.FileNamed(filename).AnsweredBy(theActor)
		require.NoError(t, err)
		file, ok := f.(*question.File)
		assert.True(t, ok)
		assert.Equal(t, path.Base(filename), file.Name())
		assert.Equal(t, path.Dir(filename), file.Dir())
		assert.True(t, file.Exists())
	})

	t.Run("returns information for non-existing file", func(t *testing.T) {
		theActor := screenplay.ActorNamed("Adam").WhoCan(ability.UseTheFileSystem())
		f, err := question.FileNamed("foo.bar").AnsweredBy(theActor)
		require.NoError(t, err)
		file, ok := f.(*question.File)
		assert.True(t, ok)
		assert.Equal(t, "foo.bar", file.Name())
		assert.Equal(t, ".", file.Dir())
		assert.False(t, file.Exists())
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		theFileNamedFooBar := question.FileNamed("foo.bar")
		assert.Equal(t, "file named 'foo.bar'", theFileNamedFooBar.String())
	})
}
