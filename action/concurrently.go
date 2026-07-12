package action

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/grandper/go-screenplay/screenplay"
)

// concurrentlyMode describes how a ConcurrentlyAction handles errors.
type concurrentlyMode int

const (
	// waitingForAll waits for every action to complete and joins their errors.
	waitingForAll concurrentlyMode = iota
	// stoppingOnError cancels the remaining actions as soon as one fails.
	stoppingOnError
	// ignoringErrors runs every action and discards their errors.
	ignoringErrors
)

// Concurrently performs the provided tasks or actions in parallel.
// By default it waits for every action to complete and joins their errors,
// which is equivalent to calling WaitingForAll. This behavior can be changed
// with StoppingOnError or IgnoringErrors, and the number of actions running at
// the same time can be capped with WithLimit.
func Concurrently(performables ...screenplay.Performable) *ConcurrentlyAction {
	return &ConcurrentlyAction{
		performables: performables,
		mode:         waitingForAll,
		limit:        0,
	}
}

// Simultaneously performs the provided tasks or actions in parallel.
// It is an alias for Concurrently.
func Simultaneously(performables ...screenplay.Performable) *ConcurrentlyAction {
	return Concurrently(performables...)
}

// ConcurrentlyAction is an action that performs several tasks or actions in parallel.
type ConcurrentlyAction struct {
	performables []screenplay.Performable
	mode         concurrentlyMode
	limit        int
}

// String describes the action.
func (a *ConcurrentlyAction) String() string {
	var builder strings.Builder

	builder.WriteString("concurrently")

	switch a.mode {
	case waitingForAll:
	case stoppingOnError:
		builder.WriteString(" (stopping on error)")
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

// WaitingForAll waits for every action to complete and joins the errors that
// occurred using errors.Join. This is the default behavior.
func (a *ConcurrentlyAction) WaitingForAll() *ConcurrentlyAction {
	a.mode = waitingForAll

	return a
}

// StoppingOnError cancels the remaining actions as soon as one of them fails
// and returns that error.
func (a *ConcurrentlyAction) StoppingOnError() *ConcurrentlyAction {
	a.mode = stoppingOnError

	return a
}

// StoppingOnFirstError cancels the remaining actions as soon as one of them
// fails. It is an alias for StoppingOnError.
func (a *ConcurrentlyAction) StoppingOnFirstError() *ConcurrentlyAction {
	return a.StoppingOnError()
}

// IgnoringErrors runs every action and discards the errors that occurred.
func (a *ConcurrentlyAction) IgnoringErrors() *ConcurrentlyAction {
	a.mode = ignoringErrors

	return a
}

// WithLimit caps the number of actions running at the same time.
// A limit that is zero or negative means no limit.
func (a *ConcurrentlyAction) WithLimit(limit int) *ConcurrentlyAction {
	a.limit = limit

	return a
}

// PerformAs performs the tasks or the actions in parallel as the provided actor.
func (a *ConcurrentlyAction) PerformAs(theActor *screenplay.Actor) error {
	switch a.mode {
	case stoppingOnError:
		return a.performStoppingOnError(theActor)
	case ignoringErrors:
		return a.performIgnoringErrors(theActor)
	case waitingForAll:
		return a.performWaitingForAll(theActor)
	}

	return a.performWaitingForAll(theActor)
}

// performWaitingForAll performs every action, waits for all of them and joins
// the errors that occurred.
func (a *ConcurrentlyAction) performWaitingForAll(theActor *screenplay.Actor) error {
	errs := make([]error, len(a.performables))

	a.run(func(index int, performable screenplay.Performable) {
		if err := performable.PerformAs(theActor); err != nil {
			errs[index] = fmt.Errorf("%s: %w", performable, err)
		}
	})

	return errors.Join(errs...)
}

// performIgnoringErrors performs every action, waits for all of them and
// discards the errors that occurred.
func (a *ConcurrentlyAction) performIgnoringErrors(theActor *screenplay.Actor) error {
	a.run(func(_ int, performable screenplay.Performable) {
		_ = performable.PerformAs(theActor)
	})

	return nil
}

// performStoppingOnError performs every action and cancels the remaining ones
// as soon as one fails, relying on the errgroup context for cancellation.
func (a *ConcurrentlyAction) performStoppingOnError(theActor *screenplay.Actor) error {
	group, ctx := errgroup.WithContext(theActor.Context())

	if a.limit > 0 {
		group.SetLimit(a.limit)
	}

	for _, performable := range a.performables {
		group.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
			}

			if err := performable.PerformAs(theActor); err != nil {
				return err
			}

			return nil
		})
	}

	return group.Wait()
}

// run performs every action in its own goroutine, waits for all of them and
// caps the number of actions running at the same time when a limit is set.
func (a *ConcurrentlyAction) run(perform func(index int, performable screenplay.Performable)) {
	var (
		waitGroup sync.WaitGroup
		semaphore chan struct{}
	)

	if a.limit > 0 {
		semaphore = make(chan struct{}, a.limit)
	}

	for index, performable := range a.performables {
		waitGroup.Add(1)

		go func(index int, performable screenplay.Performable) {
			defer waitGroup.Done()

			if semaphore != nil {
				semaphore <- struct{}{}
				defer func() { <-semaphore }()
			}

			perform(index, performable)
		}(index, performable)
	}

	waitGroup.Wait()
}

// ConcurrentlyAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*ConcurrentlyAction)(nil)
