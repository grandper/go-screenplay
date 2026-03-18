package action

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/grandper/go-screenplay/screenplay"
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
	return &EventuallyAction{
		performable: performable,
		polling:     screenplay.DefaultPolling,
		timeout:     screenplay.DefaultTimeout,
		caughtErr:   nil,
		uniqueErrs:  []error{},
	}
}

// EventuallyAction is an action that retries a task or action until it eventually succeed.
// This action will try until a given timeout is reached.
// If the action cannot be achieved until the timeout is reached, an error containing all
// unique failure errors encountered during retries is raised.
type EventuallyAction struct {
	performable screenplay.Performable
	polling     time.Duration
	timeout     time.Duration
	caughtErr   error
	uniqueErrs  []error
}

// String describes the action.
func (a *EventuallyAction) String() string {
	return fmt.Sprintf("eventually %s", a.performable)
}

// PerformAs performs the task or the action as the provided actor.
func (a *EventuallyAction) PerformAs(theActor *screenplay.Actor) error {
	if a.polling > a.timeout {
		return fmt.Errorf(
			"failed to eventually performed the action: %w",
			ErrPollingPeriodMustBeLessThanOrEqualToTimeout,
		)
	}

	a.uniqueErrs = []error{}

	timeoutTimer := time.NewTimer(a.timeout)
	defer timeoutTimer.Stop()

	pollingTicker := time.NewTicker(a.polling)
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
				theActor.Name(), a.performable, localCount, a.timeout.Seconds(), errors.Join(a.uniqueErrs...))
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

// For sets the time during which the actor keeps on trying.
func (a *EventuallyAction) For(amount int) *TimeFrameBuilder {
	return &TimeFrameBuilder{
		eventually: a,
		amount:     amount,
		unit:       "",
		duration:   &a.timeout,
	}
}

// TryingFor sets the time during which the actor keeps on trying.
func (a *EventuallyAction) TryingFor(amount int) *TimeFrameBuilder {
	return a.For(amount)
}

// TryingForNoLongerThan sets the time during which the actor keeps on trying.
func (a *EventuallyAction) TryingForNoLongerThan(amount int) *TimeFrameBuilder {
	return a.For(amount)
}

// WaitingFor sets the time during which the actor keeps on trying.
func (a *EventuallyAction) WaitingFor(amount int) *TimeFrameBuilder {
	return a.For(amount)
}

// Polling sets the polling frequency.
func (a *EventuallyAction) Polling(amount int) *TimeFrameBuilder {
	return &TimeFrameBuilder{
		eventually: a,
		amount:     amount,
		unit:       "",
		duration:   &a.polling,
	}
}

// PollingEvery sets the polling frequency.
func (a *EventuallyAction) PollingEvery(amount int) *TimeFrameBuilder {
	return a.Polling(amount)
}

// TryingEvery sets the polling frequency.
func (a *EventuallyAction) TryingEvery(amount int) *TimeFrameBuilder {
	return a.Polling(amount)
}

// EventuallyAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*EventuallyAction)(nil)

// TimeFrameBuilder builds a time frame combining amount and unit.
type TimeFrameBuilder struct {
	eventually *EventuallyAction
	amount     int
	unit       string
	duration   *time.Duration
}

// Milliseconds sets the timeout in milliseconds.
func (tfb *TimeFrameBuilder) Milliseconds() *EventuallyAction {
	tfb.unit = "milliseconds"
	*tfb.duration = time.Duration(tfb.amount) * time.Millisecond

	return tfb.eventually
}

// Millisecond sets the timeout in milliseconds.
func (tfb *TimeFrameBuilder) Millisecond() *EventuallyAction {
	return tfb.Milliseconds()
}

// Seconds sets the timeout in seconds.
func (tfb *TimeFrameBuilder) Seconds() *EventuallyAction {
	tfb.unit = "seconds"
	*tfb.duration = time.Duration(tfb.amount) * time.Second

	return tfb.eventually
}

// Second sets the timeout in seconds.
func (tfb *TimeFrameBuilder) Second() *EventuallyAction {
	return tfb.Seconds()
}
