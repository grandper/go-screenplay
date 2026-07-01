package action

import (
	"fmt"
	"strings"
	"sync"

	"github.com/grandper/go-screenplay/screenplay"
)

// asynchronousErrorsKey is the key under which the errors produced by the
// asynchronous actions are stored in the actor's memory.
const asynchronousErrorsKey = "github.com/grandper/go-screenplay/action.asynchronousErrors"

// Asynchronously performs the provided tasks or actions in the background.
// Contrary to Concurrently, the actor does not wait for the actions to complete:
// PerformAs returns immediately and the actions keep running in their own
// goroutine. The errors that occur are collected and can be investigated later
// with the AsynchronousErrors question.
//
// By default, it waits for every action to complete before collecting the errors,
// which is equivalent to calling WaitingForAll. This behavior can be changed
// with CancelOnError or IgnoringErrors, and the number of actions running at the
// same time can be capped with WithLimit.
func Asynchronously(performables ...screenplay.Performable) *AsynchronouslyAction {
	return &AsynchronouslyAction{
		performables: performables,
		mode:         waitingForAll,
		limit:        0,
	}
}

// AsynchronouslyAction is an action that performs several tasks or actions in
// the background without making the actor wait for them.
type AsynchronouslyAction struct {
	performables []screenplay.Performable
	mode         concurrentlyMode
	limit        int
}

// String describes the action.
func (a *AsynchronouslyAction) String() string {
	var builder strings.Builder

	builder.WriteString("asynchronously")

	switch a.mode {
	case stoppingOnError:
		builder.WriteString(" (cancelling on error)")
	case ignoringErrors:
		builder.WriteString(" (ignoring errors)")
	}

	if a.limit > 0 {
		fmt.Fprintf(&builder, " (limited to %d at a time)", a.limit)
	}

	builder.WriteString(" ")
	builder.WriteString(performablesToString(a.performables))

	return builder.String()
}

// WaitingForAll waits for every action to complete before collecting the errors
// that occurred. This is the default behaviour.
func (a *AsynchronouslyAction) WaitingForAll() *AsynchronouslyAction {
	a.mode = waitingForAll

	return a
}

// CancelOnError cancels the remaining actions as soon as one of them fails.
// The error is collected for later investigation.
func (a *AsynchronouslyAction) CancelOnError() *AsynchronouslyAction {
	a.mode = stoppingOnError

	return a
}

// IgnoringErrors runs every action and discards the errors that occurred.
func (a *AsynchronouslyAction) IgnoringErrors() *AsynchronouslyAction {
	a.mode = ignoringErrors

	return a
}

// WithLimit caps the number of actions running at the same time.
// A limit that is zero or negative means no limit.
func (a *AsynchronouslyAction) WithLimit(limit int) *AsynchronouslyAction {
	a.limit = limit

	return a
}

// PerformAs launches the tasks or the actions in the background and returns
// immediately. The errors that occur are collected in the actor's memory and can
// be retrieved later with the AsynchronousErrors question.
func (a *AsynchronouslyAction) PerformAs(theActor *screenplay.Actor) error {
	collector := asynchronousErrorsOf(theActor)
	concurrent := a.concurrently()

	collector.track()

	go func() {
		defer collector.done()

		collector.store(concurrent.PerformAs(theActor))
	}()

	return nil
}

// concurrently builds the Concurrently action matching the configured mode and
// limit. Asynchronously delegates the actual parallel execution to it.
func (a *AsynchronouslyAction) concurrently() *ConcurrentlyAction {
	concurrent := Concurrently(a.performables...).WithLimit(a.limit)

	switch a.mode {
	case stoppingOnError:
		return concurrent.StoppingOnError()
	case ignoringErrors:
		return concurrent.IgnoringErrors()
	default:
		return concurrent.WaitingForAll()
	}
}

// AsynchronouslyAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*AsynchronouslyAction)(nil)

// AsynchronousErrors asks for all the non-nil errors produced by the actions
// launched with Asynchronously. Answering the question waits for every pending
// asynchronous action to complete before returning the collected errors.
func AsynchronousErrors() *AsynchronousErrorsQuestion {
	return &AsynchronousErrorsQuestion{}
}

// AsynchronousErrorsQuestion is a question whose answer is the slice of errors
// produced by the actions launched with Asynchronously.
type AsynchronousErrorsQuestion struct{}

// AnsweredBy waits for every pending asynchronous action to complete and returns
// the slice of non-nil errors they produced.
func (q *AsynchronousErrorsQuestion) AnsweredBy(actor *screenplay.Actor) (any, error) {
	collector, ok := actor.Recall(asynchronousErrorsKey).(*asynchronousErrors)
	if !ok {
		return []error{}, nil
	}

	return collector.collected(), nil
}

// String describes the question.
func (q *AsynchronousErrorsQuestion) String() string {
	return "the asynchronous errors"
}

// AsynchronousErrorsQuestion implements the screenplay.Question interface.
var _ screenplay.Question = (*AsynchronousErrorsQuestion)(nil)

// asynchronousErrors is a thread-safe collector of the errors produced by the
// asynchronous actions of an actor.
type asynchronousErrors struct {
	waitGroup sync.WaitGroup
	mutex     sync.Mutex
	errs      []error
}

// asynchronousErrorsOf returns the errors collector stored in the actor's
// memory, creating and storing a new one when there is none yet.
func asynchronousErrorsOf(actor *screenplay.Actor) *asynchronousErrors {
	if collector, ok := actor.Recall(asynchronousErrorsKey).(*asynchronousErrors); ok {
		return collector
	}

	collector := &asynchronousErrors{}
	actor.Remember(asynchronousErrorsKey, collector)

	return collector
}

// track registers a new pending asynchronous action.
func (c *asynchronousErrors) track() {
	c.waitGroup.Add(1)
}

// done signals that a pending asynchronous action has completed.
func (c *asynchronousErrors) done() {
	c.waitGroup.Done()
}

// store collects the non-nil errors produced by an asynchronous action.
// When the error joins several errors together (as returned by WaitingForAll),
// each of them is collected individually.
func (c *asynchronousErrors) store(err error) {
	if err == nil {
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if joined, ok := err.(interface{ Unwrap() []error }); ok {
		for _, e := range joined.Unwrap() {
			if e != nil {
				c.errs = append(c.errs, e)
			}
		}

		return
	}

	c.errs = append(c.errs, err)
}

// collected waits for every pending asynchronous action to complete and returns
// a copy of the collected errors.
func (c *asynchronousErrors) collected() []error {
	c.waitGroup.Wait()

	c.mutex.Lock()
	defer c.mutex.Unlock()

	collected := make([]error, len(c.errs))
	copy(collected, c.errs)

	return collected
}
