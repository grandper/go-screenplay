package question_test

import (
	"testing"

	"github.com/grandper/go-screenplay/extensions/cli/ability"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/extensions/cli/question"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestEnvironmentVariable(t *testing.T) {
	t.Run("returns a nil interface if the environment variable does not exist", func(t *testing.T) {
		theActor := screenplay.ActorNamed("Adam").WhoCan(ability.RunCLICommands())
		value, err := question.EnvironmentVariableNamed("HELLO").AnsweredBy(theActor)
		require.ErrorIs(t, err, question.ErrEnvironmentVariableNotFound)
		assert.Nil(t, value)
	})

	t.Run("returns the value of the environment variable if it exists", func(t *testing.T) {
		t.Setenv("HELLO", "world")
		theActor := screenplay.ActorNamed("Adam").WhoCan(ability.RunCLICommands())
		value, err := question.EnvironmentVariableNamed("HELLO").AnsweredBy(theActor)
		require.NoError(t, err)
		assert.Equal(t, "world", value)
	})

	t.Run("fails when fails to get the error code when the actor does not have the ability MakeHttpRequest",
		func(t *testing.T) {
			theActor := screenplay.ActorNamed("Adam")
			value, err := question.EnvironmentVariableNamed("HELLO").AnsweredBy(theActor)
			require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
			assert.Nil(t, value)
		})

	t.Run("implements the stringer interface", func(t *testing.T) {
		q := question.EnvironmentVariableNamed("HELLO")
		assert.Equal(t, "environment variable named 'HELLO'", q.String())
	})
}
