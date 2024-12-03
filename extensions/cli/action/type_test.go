package action_test

import (
	"testing"

	action2 "github.com/grandper/go-screenplay/action"
	"github.com/grandper/go-screenplay/action/see"
	"github.com/grandper/go-screenplay/extensions/cli/question"
	"github.com/grandper/go-screenplay/resolution/is"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/extensions/cli/ability"

	"github.com/grandper/go-screenplay/extensions/cli/action"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestType(t *testing.T) {
	t.Run("types content", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.RunCLICommands())
			require.NoError(t,
				adam.AttemptsTo(
					action.RunTheCommand("sh", "-c", "read -p 'Enter your title: ' title; echo $title").
						Interactively()))

			require.NoError(t, adam.AttemptsTo(action.Type("Alice")))

			require.NoError(t, adam.AttemptsTo(action.Type(" in Wonderland").AndPressEnter()))

			require.NoError(t,
				adam.AttemptsTo(
					action2.Eventually(see.The(question.StandardOutputOfTheLastResponse(),
						is.EqualTo([]byte("Alice in Wonderland\n"))))))
		})

		t.Run("fails when no command is running", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.RunCLICommands())
			err := adam.AttemptsTo(action.Type("Hello, World!"))
			require.ErrorIs(t, err, ability.ErrNoCommandCurrentlyRunning)

			err = adam.AttemptsTo(action.Type("Hello, World!").AndPressEnter())
			require.ErrorIs(t, err, ability.ErrNoCommandCurrentlyRunning)
		})

		t.Run("fails when the actor lacks the ability to run CLI commands", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam")
			err := adam.AttemptsTo(action.Type("Hello, World!"))
			require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		})
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		action1 := action.Type("Hello, World!")
		assert.Equal(t, "type the content 'Hello, World!'", action1.String())

		action2 := action.Type("Hello, World!").AndPressEnter()
		assert.Equal(t, "type the content 'Hello, World!' and press enter", action2.String())
	})
}
