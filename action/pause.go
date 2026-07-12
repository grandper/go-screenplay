package action

import (
	"errors"
	"fmt"
	"time"

	"github.com/grandper/go-screenplay/screenplay"
	"github.com/grandper/go-screenplay/timing"
)

// PauseFor creates an action that pauses for a given time. The returned
// DurationBuilder lets the caller pick a unit, for example
// PauseFor(20).Milliseconds().
// Hard wait should be avoided when possible.
func PauseFor(number int) *timing.DurationBuilder[PauseAction] {
	action := &PauseAction{}
	action.durationBuilder = timing.NewDurationBuilder(action, &action.duration)

	return action.durationBuilder.For(number)
}

// PauseAction is an action that pauses for a given duration.
// Hard wait should be avoided when possible.
type PauseAction struct {
	durationBuilder *timing.DurationBuilder[PauseAction]
	duration        time.Duration
	reason          string
}

// Because specifies the reason for pausing.
func (a *PauseAction) Because(reason string) *PauseAction {
	a.reason = reason

	return a
}

// String describes the action.
func (a *PauseAction) String() string {
	return fmt.Sprintf("pause for %s because %s", a.durationBuilder, a.reason)
}

// PerformAs performs the task or the action as the provided actor.
func (a *PauseAction) PerformAs(_ *screenplay.Actor) error {
	if a.reason == "" {
		return errors.New(
			"failed to pause: cannot pause without a reason: you must call the .Because() method",
		)
	}

	time.Sleep(a.duration)

	return nil
}

// PauseAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*PauseAction)(nil)
