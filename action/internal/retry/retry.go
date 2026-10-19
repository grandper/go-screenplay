// Package retry holds the loop that keeps trying an attempt until it
// succeeds, shared by keep.TryingTo and wait.UntilThe.
package retry

import (
	"errors"
	"sync"
	"time"

	"github.com/grandper/go-screenplay/v2/action/timing"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

var (
	// ErrPollingPeriodMustBeLessThanOrEqualToTimeout is returned when the polling period is greater than the timeout.
	ErrPollingPeriodMustBeLessThanOrEqualToTimeout = errors.New("polling period must be less than or equal to timeout")
	// ErrPollingPeriodAndTimeoutMustBePositive is returned when the polling period or the timeout is zero or negative.
	ErrPollingPeriodAndTimeoutMustBePositive = errors.New("polling period and timeout must be positive")
	// ErrNoAttemptCompleted is returned when the timeout is reached while the only attempt is still running.
	ErrNoAttemptCompleted = errors.New("no attempt completed within the timeout")
)

// WindowFor returns a copy of the window with any side left unset (zero)
// filled from the actor's production timeout and polling interval.
func WindowFor(theActor *screenplay.Actor, window timing.Window) timing.Window {
	if window.Total == 0 {
		window.Total = theActor.Timeout()
	}

	if window.Interval == 0 {
		window.Interval = theActor.Polling()
	}

	return window
}

// Check returns why a window cannot be used, or nil when it can.
func Check(window timing.Window) error {
	if window.Interval <= 0 || window.Total <= 0 {
		return ErrPollingPeriodAndTimeoutMustBePositive
	}

	if !window.Valid() {
		return ErrPollingPeriodMustBeLessThanOrEqualToTimeout
	}

	return nil
}

// Until keeps trying the attempt within the window until it succeeds. It
// returns the number of tries and nil on success. When the actor's context
// ends, it returns the error of the context. When the window ends, it returns
// the distinct errors of the failed tries joined together, or
// ErrNoAttemptCompleted when no try completed. The window must pass Check.
func Until(theActor *screenplay.Actor, window timing.Window, attempt func() error) (int, error) {
	// The state of one run is local to it, so the same action can be performed
	// by several actors at the same time.
	var uniqueErrs []error

	timeoutTimer := time.NewTimer(window.Total)
	defer timeoutTimer.Stop()

	pollingTicker := time.NewTicker(window.Interval)
	defer pollingTicker.Stop()

	errCh := make(chan error, 1)
	var mutex sync.RWMutex
	count := 0

	try := func() {
		mutex.Lock()
		count++
		mutex.Unlock()

		errCh <- attempt()
	}

	tries := func() int {
		mutex.RLock()
		defer mutex.RUnlock()

		return count
	}

	// The first attempt races the timer like the retries do: a slow first
	// attempt never holds the action past its window.
	go try()

	var tick <-chan time.Time

	ctx := theActor.Context()

	for {
		select {
		case <-ctx.Done():
			return tries(), ctx.Err()
		case <-timeoutTimer.C:
			if len(uniqueErrs) == 0 {
				return tries(), ErrNoAttemptCompleted
			}

			return tries(), errors.Join(uniqueErrs...)
		case <-tick:
			tick = nil

			go try()
		case caughtErr := <-errCh:
			if caughtErr == nil {
				return tries(), nil
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

// containsErr reports whether err is already represented in errs by message.
func containsErr(errs []error, err error) bool {
	for _, e := range errs {
		if e.Error() == err.Error() {
			return true
		}
	}

	return false
}
