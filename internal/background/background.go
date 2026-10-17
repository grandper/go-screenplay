// Package background holds the collector of the errors produced by the actions
// an actor performs in the background, shared by proceed.AsynchronouslyTo and
// the question asynchronous.Errors.
package background

import (
	"sync"

	"github.com/grandper/go-screenplay/v2/screenplay"
)

// errorsKey is the key under which the collector is stored in the actor's
// memory.
const errorsKey = "github.com/grandper/go-screenplay/v2/internal/background.errors"

// Errors is a thread-safe collector of the errors produced by the background
// actions of an actor. The pending actions are counted under the mutex and
// awaited through the condition, so that an action may start while another
// goroutine waits: a sync.WaitGroup forbids Add concurrent with Wait.
type Errors struct {
	mutex   sync.Mutex
	settled *sync.Cond // signaled when the pending count drops to zero
	pending int
	errs    []error
}

func newErrors() *Errors {
	collector := &Errors{}
	collector.settled = sync.NewCond(&collector.mutex)

	return collector
}

// creation makes the creation of an actor's collector atomic: two background
// actions launched at the same time on one actor share one collector.
var creation sync.Mutex //nolint:gochecknoglobals // shared by every actor so that the creation stays atomic

// ErrorsOf returns the collector stored in the actor's memory, creating and
// storing a new one when there is none yet.
func ErrorsOf(actor *screenplay.Actor) *Errors {
	creation.Lock()
	defer creation.Unlock()

	if collector := actor.Recall[*Errors](errorsKey).OrElse(nil); collector != nil {
		return collector
	}

	collector := newErrors()
	actor.Remember(errorsKey, collector)

	return collector
}

// CollectedBy waits for every pending background action of the actor to
// complete and returns the errors they produced.
func CollectedBy(actor *screenplay.Actor) []error {
	collector := actor.Recall[*Errors](errorsKey).OrElse(nil)
	if collector == nil {
		return []error{}
	}

	return collector.collected()
}

// Track registers a new pending background action.
func (c *Errors) Track() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.pending++
}

// Done signals that a pending background action has completed.
func (c *Errors) Done() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.pending--
	if c.pending == 0 {
		c.settled.Broadcast()
	}
}

// Store collects the non-nil errors produced by a background action. When the
// error joins several errors together, each of them is collected individually.
func (c *Errors) Store(err error) {
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

// collected waits for every pending background action to complete and returns
// a copy of the collected errors.
func (c *Errors) collected() []error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for c.pending > 0 {
		c.settled.Wait()
	}

	collected := make([]error, len(c.errs))
	copy(collected, c.errs)

	return collected
}
