package action

import "github.com/grandper/go-screenplay/screenplay"

// PerformableFn represents is function that makes an actor performs an action or a task.
type PerformableFn func(theActor *screenplay.Actor) error

// FromFunc creates a new performable from a function.
func FromFunc(description string, fn PerformableFn) screenplay.Performable {
	return &funcPerformable{
		performAS:   fn,
		description: description,
	}
}

type funcPerformable struct {
	performAS   PerformableFn
	description string
}

// String describes the action.
func (p *funcPerformable) String() string {
	return p.description
}

// PerformAs performs the task or the action as the provided actor.
func (p *funcPerformable) PerformAs(actor *screenplay.Actor) error {
	return p.performAS(actor)
}

// PauseAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*funcPerformable)(nil)
