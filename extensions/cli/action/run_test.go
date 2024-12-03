package action_test

import (
	"testing"
	"time"

	"github.com/grandper/go-screenplay/extensions/cli/ability"
	"github.com/grandper/go-screenplay/extensions/cli/action"
	"github.com/grandper/go-screenplay/extensions/cli/question"
	"github.com/grandper/go-screenplay/screenplay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunTheCommand(t *testing.T) {
	t.Run("runs a command", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.RunCLICommands())
			require.NoError(t, adam.AttemptsTo(action.RunTheCommand("echo", "Hello, World!")))

			value, err := question.StandardOutputOfTheLastResponse().AnsweredBy(adam)
			require.NoError(t, err)
			assert.Equal(t, []byte("Hello, World!\n"), value)
		})

		t.Run("in a workdir", func(t *testing.T) {
			dir := t.TempDir()
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.RunCLICommands())
			require.NoError(t, adam.AttemptsTo(action.RunTheCommand("pwd").InTheWorkingDirectory(dir)))

			value, err := question.StandardOutputOfTheLastResponse().AnsweredBy(adam)
			require.NoError(t, err)
			assert.Contains(t, string(value.([]byte)), dir+"\n")
		})

		t.Run("interactively", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.RunCLICommands())
			require.NoError(t, adam.AttemptsTo(action.RunTheCommand("echo", "Hello, World!").Interactively()))

			time.Sleep(time.Second)
			value, err := question.StandardOutputOfTheLastResponse().AnsweredBy(adam)
			require.NoError(t, err)
			assert.Equal(t, []byte("Hello, World!\n"), value)
		})

		t.Run("with environment variables", func(t *testing.T) {
			env := map[string]string{
				"GREETING": "Hello, World!",
			}
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.RunCLICommands())
			require.NoError(t, adam.AttemptsTo(action.RunTheCommand("sh", "-c", "echo $GREETING").WithEnv(env)))

			value, err := question.StandardOutputOfTheLastResponse().AnsweredBy(adam)
			require.NoError(t, err)
			assert.Equal(t, []byte("Hello, World!\n"), value)
		})

		t.Run("with a single environment variable", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.RunCLICommands())
			require.NoError(t, adam.AttemptsTo(action.RunTheCommand("sh", "-c", "echo $GREETING").
				WithEnvVar("GREETING", "Hello, World!")))

			value, err := question.StandardOutputOfTheLastResponse().AnsweredBy(adam)
			require.NoError(t, err)
			assert.Equal(t, []byte("Hello, World!\n"), value)
		})

		t.Run("fails when the command is empty", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.RunCLICommands())
			err := adam.AttemptsTo(action.RunTheCommand(""))
			require.ErrorIs(t, err, action.ErrEmptyCommandName)
		})

		t.Run("fails when the command is invalid", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.RunCLICommands())
			err := adam.AttemptsTo(action.RunTheCommand("nonexistent_command"))
			require.Error(t, err)
		})

		t.Run("fails when the actor lacks the ability to run CLI commands", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam")
			err := adam.AttemptsTo(action.RunTheCommand("echo", "Hello, World!"))
			require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		})
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		action1 := action.RunTheCommand("echo", "Hello, World!")
		assert.Equal(t, "run the command 'echo Hello, World!'", action1.String())

		action2 := action.RunTheCommand("echo", "Hello, World!").Interactively()
		assert.Equal(t, "run the command 'echo Hello, World!' interactively", action2.String())

		action3 := action.RunTheCommand("echo", "Hello, World!").InTheWorkingDirectory("/tmp")
		assert.Equal(t, "run the command 'echo Hello, World!' in the working directory /tmp", action3.String())

		action4 := action.RunTheCommand("echo", "Hello, World!").Interactively().InTheWorkingDirectory("/tmp")
		assert.Equal(
			t,
			"run the command 'echo Hello, World!' interactively in the working directory /tmp",
			action4.String(),
		)
	})
}
