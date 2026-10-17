package wait

import (
	"errors"
	"fmt"
	"time"

	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/grandper/go-screenplay/v2/timing"
)

// ErrReasonNotProvided is returned when a wait is performed without a reason.
var ErrReasonNotProvided = errors.New("cannot wait without a reason: you must call the .Because() method")

// For creates an action that waits for a given time. The returned
// DurationBuilder lets the caller pick a unit, for example
// wait.For(20).Milliseconds().
// Hard wait should be avoided when possible.
func For(number int) *timing.DurationBuilder[ForAction] {
	action := &ForAction{}
	action.durationBuilder = timing.NewDurationBuilder(action, &action.duration)

	return action.durationBuilder.For(number)
}

// ForAction is an action that waits for a given duration.
// Hard wait should be avoided when possible.
type ForAction struct {
	durationBuilder *timing.DurationBuilder[ForAction]
	duration        time.Duration
	reason          string
}

// Because specifies the reason for waiting.
func (a *ForAction) Because(reason string) *ForAction {
	a.reason = reason

	return a
}

// String describes the action.
func (a *ForAction) String() string {
	return fmt.Sprintf("wait for %s because %s", a.durationBuilder, a.reason)
}

// PerformAs performs the task or the action as the provided actor.
func (a *ForAction) PerformAs(_ *screenplay.Actor) error {
	if a.reason == "" {
		return fmt.Errorf("failed to wait: %w", ErrReasonNotProvided)
	}

	time.Sleep(a.duration)

	return nil
}

// ForAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*ForAction)(nil)
