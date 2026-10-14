package action_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action"
	"github.com/grandper/go-screenplay/v2/action/see"
	"github.com/grandper/go-screenplay/v2/extensions/cli/question"
	"github.com/grandper/go-screenplay/v2/extensions/cli/question/standard"
	"github.com/grandper/go-screenplay/v2/question/last"
	"github.com/grandper/go-screenplay/v2/resolution/be"

	"github.com/grandper/go-screenplay/v2/extensions/cli/ability"

	cliaction "github.com/grandper/go-screenplay/v2/extensions/cli/action"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestType(t *testing.T) {
	t.Run("types content", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.RunCLICommands())
			require.NoError(t,
				adam.AttemptsTo(
					cliaction.RunTheCommand("sh", "-c", "read -p 'Enter your title: ' title; echo $title").
						Interactively()))

			require.NoError(t, adam.AttemptsTo(cliaction.Type("Alice")))

			require.NoError(t, adam.AttemptsTo(cliaction.Type(" in Wonderland").AndPressEnter()))

			require.NoError(t,
				adam.AttemptsTo(
					action.Eventually(see.The(standard.OutputOf(last.Of(question.Responses()))).
						Should(be.DeepEqualTo([]byte("Alice in Wonderland\n"))))))
		})

		t.Run("fails when no command is running", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.RunCLICommands())
			err := adam.AttemptsTo(cliaction.Type("Hello, World!"))
			require.ErrorIs(t, err, ability.ErrNoCommandCurrentlyRunning)

			err = adam.AttemptsTo(cliaction.Type("Hello, World!").AndPressEnter())
			require.ErrorIs(t, err, ability.ErrNoCommandCurrentlyRunning)
		})

		t.Run("fails when the actor lacks the ability to run CLI commands", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam")
			err := adam.AttemptsTo(cliaction.Type("Hello, World!"))
			require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		})
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		action1 := cliaction.Type("Hello, World!")
		assert.Equal(t, "type the content \"Hello, World!\"", action1.String())

		action2 := cliaction.Type("Hello, World!").AndPressEnter()
		assert.Equal(t, "type the content \"Hello, World!\" and press enter", action2.String())
	})
}
