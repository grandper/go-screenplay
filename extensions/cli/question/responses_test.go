package question_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/extensions/cli/ability"
	"github.com/grandper/go-screenplay/v2/extensions/cli/action/run"
	"github.com/grandper/go-screenplay/v2/extensions/cli/question"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestResponses(t *testing.T) {
	t.Run("returns all the responses", func(t *testing.T) {
		theActor := screenplay.ActorNamed("Adam").WhoCan(ability.RunCLICommands())
		require.NoError(t, theActor.AttemptsTo(run.TheCommand("echo", "Hello World")))
		require.NoError(t, theActor.AttemptsTo(run.TheCommand("echo", "Goodbye World")))
		responses, err := question.Responses().AnsweredBy(theActor)
		require.NoError(t, err)
		require.Len(t, responses, 2)
		assert.Equal(t, []byte("Hello World\n"), responses[0].StdOut())
		assert.Equal(t, []byte("Goodbye World\n"), responses[1].StdOut())
	})

	t.Run("returns an empty slice when no command has been run", func(t *testing.T) {
		theActor := screenplay.ActorNamed("Adam").WhoCan(ability.RunCLICommands())
		value, err := question.Responses().AnsweredBy(theActor)
		require.NoError(t, err)
		assert.Empty(t, value)
	})

	t.Run("fails when the actor does not have the ability RunCLICommands",
		func(t *testing.T) {
			theActor := screenplay.ActorNamed("Adam")
			value, err := question.Responses().AnsweredBy(theActor)
			require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
			assert.Nil(t, value)
		})

	t.Run("implements the stringer interface", func(t *testing.T) {
		q := question.Responses()
		assert.Equal(t, "command results", q.String())
	})

	t.Run("implements the question interface", func(t *testing.T) {
		assert.Implements(t, (*screenplay.Question[[]*ability.Result])(nil), question.Responses())
	})
}
