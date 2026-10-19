package ability

import (
	impl "github.com/grandper/go-screenplay/v2/extensions/cli/internal/ability"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// RunCLICommands represents the ability of an actor to run CLI commands.
func RunCLICommands() screenplay.Ability {
	return impl.RunCLICommands()
}
