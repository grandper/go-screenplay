package enter

import impl "github.com/grandper/go-screenplay/v2/extensions/cli/action/internal/enter"

// TheText enters a text in the interactive session of the command running:
//
//	enter.TheText("Alice in Wonderland").AndPressEnter()
func TheText(format string, a ...any) *impl.TheTextAction {
	return impl.TheText(format, a...)
}
