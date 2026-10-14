package action

import (
	"fmt"

	"github.com/grandper/go-screenplay/v2/internal/verdict"

	"github.com/grandper/go-screenplay/v2/extensions/cli/ability"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// Type types input in an interactive session.
func Type(format string, a ...any) *TypeAction {
	return &TypeAction{
		content: fmt.Sprintf(format, a...),
	}
}

// TypeAction types input in an interactive session.
type TypeAction struct {
	content    string
	pressEnter bool
}

// AndPressEnter presses enter.
func (a *TypeAction) AndPressEnter() *TypeAction {
	a.pressEnter = true
	return a
}

// String describes the action.
func (a *TypeAction) String() string {
	str := fmt.Sprintf("type the content %s", verdict.Describe(a.content))
	if a.pressEnter {
		str += " and press enter"
	}
	return str
}

// PerformAs performs the task or the action as the provided actor.
func (a *TypeAction) PerformAs(theActor *screenplay.Actor) error {
	runCLICommands := screenplay.UseAnAbilityOf(theActor).That(ability.RunCLICommands())
	if a.pressEnter {
		return runCLICommands.ToType(a.content + "\n")
	}
	return runCLICommands.ToType(a.content)
}

// TypeAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*TypeAction)(nil)
