package action

import (
	"fmt"

	"github.com/grandper/go-screenplay/extensions/cli/ability"
	"github.com/grandper/go-screenplay/screenplay"
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
	str := fmt.Sprintf("type the content '%s'", a.content)
	if a.pressEnter {
		str += " and press enter"
	}
	return str
}

// PerformAs performs the task or the action as the provided actor.
func (a *TypeAction) PerformAs(theActor *screenplay.Actor) error {
	runCLICommands, err := screenplay.UseAbilityTo[*ability.RunCLICommandsAbility]().Of(theActor)
	if err != nil {
		return err
	}
	if a.pressEnter {
		return runCLICommands.Type(a.content + "\n")
	}
	return runCLICommands.Type(a.content)
}

// TypeAction implements the screenplay.Action interface.
var _ screenplay.Performable = (*TypeAction)(nil)
