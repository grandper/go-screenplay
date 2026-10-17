package proceed

import (
	"fmt"

	"github.com/grandper/go-screenplay/v2/internal/background"
	"github.com/grandper/go-screenplay/v2/internal/performables"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// AsynchronouslyTo performs the provided tasks or actions in the background.
// Contrary to ConcurrentlyTo, the actor does not wait for the actions to complete:
// PerformAs returns immediately and the actions keep running in their own
// goroutine. The errors that occur are collected and can be investigated later
// with the question asynchronous.Errors.
//
// By default, it waits for every action to complete before collecting the errors,
// which is equivalent to calling WaitingForAll. This behavior can be changed
// with StoppingOnError or IgnoringErrors, and the number of actions running at the
// same time can be capped with WithLimit.
func AsynchronouslyTo(actions ...screenplay.Performable) *AsynchronouslyToAction {
	return &AsynchronouslyToAction{
		performables: performables.WithoutNil(actions),
		mode:         waitingForAll,
		limit:        0,
		options:      []string{},
	}
}

// AsynchronouslyToAction is an action that performs several tasks or actions in
// the background without making the actor wait for them.
type AsynchronouslyToAction struct {
	performables []screenplay.Performable
	mode         concurrentlyMode
	limit        int
	options      []string // the options as they were written, in call order
}

// String describes the action, followed by its options as they were written.
func (a *AsynchronouslyToAction) String() string {
	return describe("asynchronously", a.performables, a.options)
}

// WaitingForAll waits for every action to complete before collecting the errors
// that occurred. This is the default behavior.
func (a *AsynchronouslyToAction) WaitingForAll() *AsynchronouslyToAction {
	a.mode = waitingForAll
	a.options = append(a.options, "waiting for all")

	return a
}

// StoppingOnError stops the remaining actions as soon as one of them fails.
// The error is collected for later investigation. An action still running is
// canceled through the actor's context: it observes the cancellation with
// actor.Context(). An action not started yet is skipped.
func (a *AsynchronouslyToAction) StoppingOnError() *AsynchronouslyToAction {
	a.mode = stoppingOnError
	a.options = append(a.options, "stopping on error")

	return a
}

// IgnoringErrors runs every action and discards the errors that occurred.
func (a *AsynchronouslyToAction) IgnoringErrors() *AsynchronouslyToAction {
	a.mode = ignoringErrors
	a.options = append(a.options, "ignoring errors")

	return a
}

// WithLimit caps the number of actions running at the same time.
// A limit that is zero or negative means no limit.
func (a *AsynchronouslyToAction) WithLimit(limit int) *AsynchronouslyToAction {
	a.limit = limit
	a.options = append(a.options, fmt.Sprintf("with limit %d", limit))

	return a
}

// PerformAs launches the tasks or the actions in the background and returns
// immediately. The errors that occur are collected in the actor's memory and can
// be retrieved later with the question asynchronous.Errors.
func (a *AsynchronouslyToAction) PerformAs(theActor *screenplay.Actor) error {
	collector := background.ErrorsOf(theActor)
	concurrent := a.concurrently()

	collector.Track()

	go func() {
		defer collector.Done()

		collector.Store(concurrent.PerformAs(theActor))
	}()

	return nil
}

// concurrently builds the ConcurrentlyTo action matching the configured mode
// and limit. AsynchronouslyTo delegates the actual parallel execution to it.
func (a *AsynchronouslyToAction) concurrently() *ConcurrentlyToAction {
	concurrent := ConcurrentlyTo(a.performables...)
	concurrent.mode = a.mode
	concurrent.limit = a.limit

	return concurrent
}

// AsynchronouslyToAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*AsynchronouslyToAction)(nil)
