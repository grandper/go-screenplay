package keep

import (
	"fmt"

	"github.com/grandper/go-screenplay/v2/action/internal/retry"
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/grandper/go-screenplay/v2/timing"
)

var (
	// ErrPollingPeriodMustBeLessThanOrEqualToTimeout is returned when the polling period is greater than the timeout.
	ErrPollingPeriodMustBeLessThanOrEqualToTimeout = retry.ErrPollingPeriodMustBeLessThanOrEqualToTimeout
	// ErrPollingPeriodAndTimeoutMustBePositive is returned when the polling period or the timeout is zero or negative.
	ErrPollingPeriodAndTimeoutMustBePositive = retry.ErrPollingPeriodAndTimeoutMustBePositive
	// ErrNoAttemptCompleted is returned when the timeout is reached while the only attempt is still running.
	ErrNoAttemptCompleted = retry.ErrNoAttemptCompleted
)

// TryingTo retries a task or action until it succeeds, for the time and at the
// pace set by its window, or by the actor's production when the window is not
// set:
//
//	keep.TryingTo(CancelTheOrder()).For(5).Seconds().PollingEvery(500).Milliseconds()
//
// If the action does not succeed before the end of the window, an error
// containing all the distinct errors of the tries is returned.
func TryingTo(performable screenplay.Performable) *TryingToAction {
	// The window starts unset (zero durations): each side is filled from the
	// actor's production at PerformAs unless the caller sets it through the
	// WindowBuilder (For, PollingEvery, ...).
	action := &TryingToAction{
		performable: performable,
		window:      timing.NewWindow(0, 0),
	}
	action.WindowBuilder = timing.NewWindowBuilder(action, &action.window)

	return action
}

// TryingToAction is an action that retries a task or action until it succeeds.
type TryingToAction struct {
	*timing.WindowBuilder[TryingToAction]

	performable screenplay.Performable
	window      timing.Window
}

// String describes the action, followed by its window as it was written.
func (a *TryingToAction) String() string {
	description := fmt.Sprintf("keep trying to %s", a.performable)
	if window := a.WindowBuilder.String(); window != "" {
		description += " " + window
	}

	return description
}

// PerformAs performs the task or the action as the provided actor until it
// succeeds or the window ends.
func (a *TryingToAction) PerformAs(theActor *screenplay.Actor) error {
	window := retry.WindowFor(theActor, a.window)
	if err := retry.Check(window); err != nil {
		return fmt.Errorf("failed to keep trying to perform the action: %w", err)
	}

	tries, err := retry.Until(theActor, window, func() error {
		return theActor.AttemptsTo(a.performable)
	})
	if err != nil {
		return fmt.Errorf("%s kept trying to %s for %s (%s) without success: %w",
			theActor.Name(), a.performable, window.Total, triesOf(tries), err)
	}

	return nil
}

// triesOf counts the tries: "1 try", "3 tries".
func triesOf(count int) string {
	if count == 1 {
		return "1 try"
	}

	return fmt.Sprintf("%d tries", count)
}

// TryingToAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*TryingToAction)(nil)
