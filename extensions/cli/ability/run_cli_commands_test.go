package ability_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/extensions/cli/ability"
	"github.com/grandper/go-screenplay/v2/extensions/cli/action/run"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestRunCLICommands(t *testing.T) {
	t.Run("gives an actor the ability to run CLI commands", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.RunCLICommands())
		require.NoError(t, adam.AttemptsTo(run.TheCommand("echo", "Hello World")))
	})
}
