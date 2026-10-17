package proceed

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/grandper/go-screenplay/v2/internal/performables"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// concurrentlyMode describes how a ConcurrentlyToAction handles errors.
type concurrentlyMode int

const (
	// waitingForAll waits for every action to complete and joins their errors.
	waitingForAll concurrentlyMode = iota
	// stoppingOnError cancels the remaining actions as soon as one fails.
	stoppingOnError
	// ignoringErrors runs every action and discards their errors.
	ignoringErrors
)

// ConcurrentlyTo performs the provided tasks or actions in parallel.
// By default it waits for every action to complete and joins their errors,
// which is equivalent to calling WaitingForAll. This behavior can be changed
// with StoppingOnError or IgnoringErrors, and the number of actions running at
// the same time can be capped with WithLimit.
func ConcurrentlyTo(actions ...screenplay.Performable) *ConcurrentlyToAction {
	return &ConcurrentlyToAction{
		adverb:       "concurrently",
		performables: performables.WithoutNil(actions),
		mode:         waitingForAll,
		limit:        0,
		options:      []string{},
	}
}

// SimultaneouslyTo performs the provided tasks or actions in parallel. It
// behaves like ConcurrentlyTo and is narrated with "simultaneously".
func SimultaneouslyTo(actions ...screenplay.Performable) *ConcurrentlyToAction {
	action := ConcurrentlyTo(actions...)
	action.adverb = "simultaneously"

	return action
}

// ConcurrentlyToAction is an action that performs several tasks or actions in parallel.
type ConcurrentlyToAction struct {
	adverb       string
	performables []screenplay.Performable
	mode         concurrentlyMode
	limit        int
	options      []string // the options as they were written, in call order
}

// String describes the action, followed by its options as they were written.
func (a *ConcurrentlyToAction) String() string {
	return describe(a.adverb, a.performables, a.options)
}

// describe narrates "proceed <adverb> to A, B, <option>, <option>".
func describe(adverb string, actions []screenplay.Performable, options []string) string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "proceed %s to %s", adverb, performables.Describe(actions))

	for _, option := range options {
		builder.WriteString(", ")
		builder.WriteString(option)
	}

	return builder.String()
}

// WaitingForAll waits for every action to complete and joins the errors that
// occurred using errors.Join. This is the default behavior.
func (a *ConcurrentlyToAction) WaitingForAll() *ConcurrentlyToAction {
	a.mode = waitingForAll
	a.options = append(a.options, "waiting for all")

	return a
}

// StoppingOnError cancels the remaining actions as soon as one of them fails
// and returns that error. An action still running is canceled through the
// actor's context: it observes the cancellation with actor.Context(), as the
// README's "Working with contexts" section describes. An action not started yet
// is skipped.
func (a *ConcurrentlyToAction) StoppingOnError() *ConcurrentlyToAction {
	a.mode = stoppingOnError
	a.options = append(a.options, "stopping on error")

	return a
}

// StoppingOnFirstError cancels the remaining actions as soon as one of them
// fails. It behaves like StoppingOnError and is narrated "stopping on first error".
func (a *ConcurrentlyToAction) StoppingOnFirstError() *ConcurrentlyToAction {
	a.mode = stoppingOnError
	a.options = append(a.options, "stopping on first error")

	return a
}

// IgnoringErrors runs every action and discards the errors that occurred.
func (a *ConcurrentlyToAction) IgnoringErrors() *ConcurrentlyToAction {
	a.mode = ignoringErrors
	a.options = append(a.options, "ignoring errors")

	return a
}

// WithLimit caps the number of actions running at the same time.
// A limit that is zero or negative means no limit.
func (a *ConcurrentlyToAction) WithLimit(limit int) *ConcurrentlyToAction {
	a.limit = limit
	a.options = append(a.options, fmt.Sprintf("with limit %d", limit))

	return a
}

// PerformAs performs the tasks or the actions in parallel as the provided actor.
func (a *ConcurrentlyToAction) PerformAs(theActor *screenplay.Actor) error {
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
func (a *ConcurrentlyToAction) performWaitingForAll(theActor *screenplay.Actor) error {
	errs := make([]error, len(a.performables))

	a.run(func(index int, performable screenplay.Performable) {
		if err := theActor.AttemptsTo(performable); err != nil {
			errs[index] = fmt.Errorf("%s: %w", performable, err)
		}
	})

	return errors.Join(errs...)
}

// performIgnoringErrors performs every action, waits for all of them and
// discards the errors that occurred.
func (a *ConcurrentlyToAction) performIgnoringErrors(theActor *screenplay.Actor) error {
	a.run(func(_ int, performable screenplay.Performable) {
		_ = theActor.AttemptsTo(performable)
	})

	return nil
}

// performStoppingOnError performs every action and cancels the remaining ones
// as soon as one fails. The actions are performed through a view of the actor
// whose context is the errgroup's, so they observe the cancellation with
// actor.Context() while the actor itself, and every other goroutine holding
// it, keep their own context.
func (a *ConcurrentlyToAction) performStoppingOnError(theActor *screenplay.Actor) error {
	group, ctx := errgroup.WithContext(theActor.Context())
	actorInTheRun := theActor.WithContext(ctx)

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

			if err := actorInTheRun.AttemptsTo(performable); err != nil {
				return err
			}

			return nil
		})
	}

	return group.Wait()
}

// run performs every action in its own goroutine, waits for all of them and
// caps the number of actions running at the same time when a limit is set.
// Each action is performed through AttemptsTo, so it is narrated as a beat.
func (a *ConcurrentlyToAction) run(perform func(index int, performable screenplay.Performable)) {
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

// ConcurrentlyToAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*ConcurrentlyToAction)(nil)
