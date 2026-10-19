package enter

import (
	"fmt"

	"github.com/grandper/go-screenplay/v2/extensions/cli/internal/verdict"

	impl "github.com/grandper/go-screenplay/v2/extensions/cli/internal/ability"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// TheText enters a text in the interactive session of the command running:
//
//	enter.TheText("Alice in Wonderland").AndPressEnter()
func TheText(format string, a ...any) *TheTextAction {
	return &TheTextAction{
		content: fmt.Sprintf(format, a...),
	}
}

// TheTextAction enters a text in an interactive session.
type TheTextAction struct {
	content    string
	pressEnter bool
}

// AndPressEnter presses enter.
func (a *TheTextAction) AndPressEnter() *TheTextAction {
	a.pressEnter = true
	return a
}

// String describes the action.
func (a *TheTextAction) String() string {
	str := fmt.Sprintf("enter the text %s", verdict.Describe(a.content))
	if a.pressEnter {
		str += " and press enter"
	}
	return str
}

// PerformAs performs the task or the action as the provided actor.
func (a *TheTextAction) PerformAs(theActor *screenplay.Actor) error {
	runCLICommands := screenplay.UseAnAbilityOf(theActor).That(impl.RunCLICommands())
	if a.pressEnter {
		return runCLICommands.ToType(a.content + "\n")
	}
	return runCLICommands.ToType(a.content)
}

// TheTextAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*TheTextAction)(nil)
