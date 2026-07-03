package action

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/grandper/go-screenplay/screenplay"
	"github.com/grandper/go-screenplay/utils"
)

var (
	// ErrPollingPeriodMustBeLessThanOrEqualToTimeout is returned when the polling period is greater than the timeout.
	ErrPollingPeriodMustBeLessThanOrEqualToTimeout = errors.New("polling period must be less than or equal to timeout")
)

// Eventually retries a task or action until it eventually succeed.
// This action will try until a given timeout is reached.
// If the action cannot be achieved until the timeout is reached, an error containing all
// unique failure errors encountered during retries is raised.
func Eventually(performable screenplay.Performable) *EventuallyAction {
	action := &EventuallyAction{
		performable: performable,
		window:      utils.NewRetryWindow(screenplay.DefaultTimeout, screenplay.DefaultPolling),
		caughtErr:   nil,
		uniqueErrs:  []error{},
	}
	action.RetryWindowBuilder = utils.NewRetryWindowBuilder(action, &action.window)

	return action
}

// EventuallyAction is an action that retries a task or action until it eventually succeed.
// This action will try until a given timeout is reached.
// If the action cannot be achieved until the timeout is reached, an error containing all
// unique failure errors encountered during retries is raised.
type EventuallyAction struct {
	*utils.RetryWindowBuilder[EventuallyAction]
	performable screenplay.Performable
	window      utils.RetryWindow
	caughtErr   error
	uniqueErrs  []error
}

// String describes the action.
func (a *EventuallyAction) String() string {
	return fmt.Sprintf("eventually %s", a.performable)
}

// PerformAs performs the task or the action as the provided actor.
func (a *EventuallyAction) PerformAs(theActor *screenplay.Actor) error {
	if !a.window.Valid() {
		return fmt.Errorf(
			"failed to eventually performed the action: %w",
			ErrPollingPeriodMustBeLessThanOrEqualToTimeout,
		)
	}

	a.uniqueErrs = []error{}

	timeoutTimer := time.NewTimer(a.window.Total)
	defer timeoutTimer.Stop()

	pollingTicker := time.NewTicker(a.window.Interval)
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

	tryToPerformTheAction()

	var tick <-chan time.Time

	for {
		select {
		case <-timeoutTimer.C:
			mutex.RLock()
			localCount := count
			mutex.RUnlock()
			return fmt.Errorf("an error occurred when %s tried to eventually %s %d times over %f seconds: %w",
				theActor.Name(), a.performable, localCount, a.window.Total.Seconds(), errors.Join(a.uniqueErrs...))
		case <-tick:
			tick = nil

			go tryToPerformTheAction()
		case a.caughtErr = <-errCh:
			if a.caughtErr == nil {
				return nil
			}

			if !containsErr(a.uniqueErrs, a.caughtErr) {
				a.uniqueErrs = append(a.uniqueErrs, a.caughtErr)
			}

			tick = pollingTicker.C
		}
	}
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
