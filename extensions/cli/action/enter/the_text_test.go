package enter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action/keep"
	"github.com/grandper/go-screenplay/v2/action/see"
	"github.com/grandper/go-screenplay/v2/extensions/cli/question"
	"github.com/grandper/go-screenplay/v2/extensions/cli/question/standard"
	"github.com/grandper/go-screenplay/v2/question/last"
	"github.com/grandper/go-screenplay/v2/resolution/be"

	"github.com/grandper/go-screenplay/v2/extensions/cli/ability"

	"github.com/grandper/go-screenplay/v2/extensions/cli/action/enter"
	"github.com/grandper/go-screenplay/v2/extensions/cli/action/run"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestTheText(t *testing.T) {
	t.Run("types content", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.RunCLICommands())
			require.NoError(t,
				adam.AttemptsTo(
					run.TheCommand("sh", "-c", "read -p 'Enter your title: ' title; echo $title").
						Interactively()))

			require.NoError(t, adam.AttemptsTo(enter.TheText("Alice")))

			require.NoError(t, adam.AttemptsTo(enter.TheText(" in Wonderland").AndPressEnter()))

			require.NoError(t,
				adam.AttemptsTo(
					keep.TryingTo(see.The(standard.OutputOf(last.Of(question.Responses()))).
						Should(be.DeepEqualTo([]byte("Alice in Wonderland\n"))))))
		})

		t.Run("fails when no command is running", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.RunCLICommands())
			err := adam.AttemptsTo(enter.TheText("Hello, World!"))
			require.ErrorIs(t, err, ability.ErrNoCommandCurrentlyRunning)

			err = adam.AttemptsTo(enter.TheText("Hello, World!").AndPressEnter())
			require.ErrorIs(t, err, ability.ErrNoCommandCurrentlyRunning)
		})

		t.Run("fails when the actor lacks the ability to run CLI commands", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam")
			err := adam.AttemptsTo(enter.TheText("Hello, World!"))
			require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		})
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		action1 := enter.TheText("Hello, World!")
		assert.Equal(t, "enter the text \"Hello, World!\"", action1.String())

		action2 := enter.TheText("Hello, World!").AndPressEnter()
		assert.Equal(t, "enter the text \"Hello, World!\" and press enter", action2.String())
	})
}
