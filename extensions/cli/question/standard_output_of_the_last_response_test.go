package question_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/extensions/cli/ability"
	"github.com/grandper/go-screenplay/extensions/cli/action"

	"github.com/grandper/go-screenplay/extensions/cli/question"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestStandardOutputOfTheLastResponse(t *testing.T) {
	t.Run("returns the standard output of the last response", func(t *testing.T) {
		theActor := screenplay.ActorNamed("Adam").WhoCan(ability.RunCLICommands())
		require.NoError(t, theActor.AttemptsTo(action.RunTheCommand("echo", "Hello World")))
		value, err := question.StandardOutputOfTheLastResponse().AnsweredBy(theActor)
		require.NoError(t, err)
		assert.Equal(t, []byte("Hello World\n"), value)

		require.Error(t, theActor.AttemptsTo(action.RunTheCommand("ls", "/foobar")))
		value, err = question.StandardOutputOfTheLastResponse().AnsweredBy(theActor)
		require.NoError(t, err)
		assert.Equal(t, []byte(""), value)
	})

	t.Run("fails when no command has been run", func(t *testing.T) {
		theActor := screenplay.ActorNamed("Adam").WhoCan(ability.RunCLICommands())
		value, err := question.StandardOutputOfTheLastResponse().AnsweredBy(theActor)
		require.ErrorIs(t, err, question.ErrNoResponses)
		assert.Nil(t, value)
	})

	t.Run("fails when the actor does not have the ability RunCLICommands",
		func(t *testing.T) {
			theActor := screenplay.ActorNamed("Adam")
			value, err := question.StandardOutputOfTheLastResponse().AnsweredBy(theActor)
			require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
			assert.Nil(t, value)
		})

	t.Run("implements the stringer interface", func(t *testing.T) {
		q := question.StandardOutputOfTheLastResponse()
		assert.Equal(t, "standard output of the last response", q.String())
	})
}
