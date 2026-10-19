package run

import impl "github.com/grandper/go-screenplay/v2/extensions/cli/action/internal/run"

// TheCommand runs a given command with arguments:
//
//	run.TheCommand("echo", "Hello World")
func TheCommand(name string, args ...string) *impl.TheCommandAction {
	return impl.TheCommand(name, args...)
}
