package action

import (
	"errors"
	"fmt"
	"time"

	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/grandper/go-screenplay/v2/timing"
)

// ErrReasonNotProvided is returned when a pause is performed without a reason.
var ErrReasonNotProvided = errors.New("cannot pause without a reason: you must call the .Because() method")

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
		return fmt.Errorf("failed to pause: %w", ErrReasonNotProvided)
	}

	time.Sleep(a.duration)

	return nil
}

// PauseAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*PauseAction)(nil)
