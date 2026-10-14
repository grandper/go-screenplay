package action

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/grandper/go-screenplay/v2/timing"
)

var (
	// ErrPollingPeriodMustBeLessThanOrEqualToTimeout is returned when the polling period is greater than the timeout.
	ErrPollingPeriodMustBeLessThanOrEqualToTimeout = errors.New("polling period must be less than or equal to timeout")
	// ErrPollingPeriodAndTimeoutMustBePositive is returned when the polling period or the timeout is zero or negative.
	ErrPollingPeriodAndTimeoutMustBePositive = errors.New("polling period and timeout must be positive")
	// ErrNoAttemptCompleted is returned when the timeout is reached while the only attempt is still running.
	ErrNoAttemptCompleted = errors.New("no attempt completed within the timeout")
)

// Eventually retries a task or action until it eventually succeed.
// This action will try until a given timeout is reached.
// If the action cannot be achieved until the timeout is reached, an error containing all
// unique failure errors encountered during retries is raised.
func Eventually(performable screenplay.Performable) *EventuallyAction {
	// The window starts unset (zero durations): each side is filled from the
	// actor's production at PerformAs unless the caller overrides it through the
	// WindowBuilder (For, Polling, ...).
	action := &EventuallyAction{
		performable: performable,
		window:      timing.NewWindow(0, 0),
	}
	action.WindowBuilder = timing.NewWindowBuilder(action, &action.window)

	return action
}

// EventuallyAction is an action that retries a task or action until it eventually succeed.
// This action will try until a given timeout is reached.
// If the action cannot be achieved until the timeout is reached, an error containing all
// unique failure errors encountered during retries is raised.
type EventuallyAction struct {
	*timing.WindowBuilder[EventuallyAction]

	performable screenplay.Performable
	window      timing.Window
}

// String describes the action.
func (a *EventuallyAction) String() string {
	return fmt.Sprintf("eventually %s", a.performable)
}

// PerformAs performs the task or the action as the provided actor.
func (a *EventuallyAction) PerformAs(theActor *screenplay.Actor) error {
	window := a.windowFor(theActor)

	if window.Interval <= 0 || window.Total <= 0 {
		return fmt.Errorf(
			"failed to eventually perform the action: %w",
			ErrPollingPeriodAndTimeoutMustBePositive,
		)
	}

	if !window.Valid() {
		return fmt.Errorf(
			"failed to eventually perform the action: %w",
			ErrPollingPeriodMustBeLessThanOrEqualToTimeout,
		)
	}

	// The state of one run is local to it, so the same instance can be
	// performed by several actors at the same time.
	var uniqueErrs []error

	timeoutTimer := time.NewTimer(window.Total)
	defer timeoutTimer.Stop()

	pollingTicker := time.NewTicker(window.Interval)
	defer pollingTicker.Stop()

	errCh := make(chan error, 1)
	var mutex sync.RWMutex
	count := 0

	tryToPerformTheAction := func() {
		mutex.Lock()
		count++
		mutex.Unlock()

		errCh <- theActor.AttemptsTo(a.performable)
	}

	// The first attempt races the timer like the retries do: a slow first
	// attempt never holds the action past its window.
	go tryToPerformTheAction()

	var tick <-chan time.Time

	ctx := theActor.Context()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("an error occurred when %s tried to eventually %s: %w",
				theActor.Name(), a.performable, ctx.Err())
		case <-timeoutTimer.C:
			mutex.RLock()
			localCount := count
			mutex.RUnlock()
			if len(uniqueErrs) == 0 {
				return fmt.Errorf("an error occurred when %s tried to eventually %s %d times over %s: %w",
					theActor.Name(), a.performable, localCount, window.Total, ErrNoAttemptCompleted)
			}
			return fmt.Errorf("an error occurred when %s tried to eventually %s %d times over %s: %w",
				theActor.Name(), a.performable, localCount, window.Total, errors.Join(uniqueErrs...))
		case <-tick:
			tick = nil

			go tryToPerformTheAction()
		case caughtErr := <-errCh:
			if caughtErr == nil {
				return nil
			}

			if !containsErr(uniqueErrs, caughtErr) {
				uniqueErrs = append(uniqueErrs, caughtErr)
			}

			// The polling interval counts from the end of the failed attempt, so an
			// attempt longer than the interval never triggers an immediate retry.
			pollingTicker.Reset(window.Interval)
			tick = pollingTicker.C
		}
	}
}

// windowFor resolves the retry window used for a given actor. It returns a copy
// of the configured window with any side left unset (zero) filled from the
// actor's production timeout and polling interval, so the action reads timeout
// and polling from the production while still honoring explicit overrides made
// through the WindowBuilder.
func (a *EventuallyAction) windowFor(theActor *screenplay.Actor) timing.Window {
	window := a.window

	if window.Total == 0 {
		window.Total = theActor.Timeout()
	}

	if window.Interval == 0 {
		window.Interval = theActor.Polling()
	}

	return window
}

// containsErr reports whether err is already represented in errs by message.
func containsErr(errs []error, err error) bool {
	for _, e := range errs {
		if e.Error() == err.Error() {
			return true
		}
	}
	return false
}

// EventuallyAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*EventuallyAction)(nil)
