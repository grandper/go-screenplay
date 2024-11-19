package question_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/extensions/filesystem/ability"
	"github.com/grandper/go-screenplay/extensions/filesystem/action/fixture"
	"github.com/grandper/go-screenplay/extensions/filesystem/question"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestContentOfTheFileNamedQuestion(t *testing.T) {
	t.Run("fails if the actor does not have the ability to use the file system", func(t *testing.T) {
		theActor := screenplay.ActorNamed("Adam")
		file, err := question.ContentOfTheFileNamed("foo.bar").AnsweredBy(theActor)
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		assert.Nil(t, file)
	})

	theActor := screenplay.ActorNamed("Adam").WhoCan(ability.UseTheFileSystem())

	t.Run("returns information for existing file", func(t *testing.T) {
		filename := fixture.CreateFilename(t)
		require.NoError(t, os.WriteFile(filename, []byte("foo"), 0777))
		defer func() {
			assert.NoError(t, os.Remove(filename))
		}()
		data, err := question.ContentOfTheFileNamed(filename).AnsweredBy(theActor)
		require.NoError(t, err)
		content, ok := data.([]byte)
		assert.True(t, ok)
		assert.Equal(t, []byte("foo"), content)
	})

	t.Run("fails for non-existing file", func(t *testing.T) {
		content, err := question.ContentOfTheFileNamed("foo.bar").AnsweredBy(theActor)
		require.Error(t, err)
		assert.Nil(t, content)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		contentOfTheFileNamedFooBar := question.ContentOfTheFileNamed("foo.bar")
		assert.Equal(t, "content of the file named 'foo.bar'", contentOfTheFileNamedFooBar.String())
	})
}
