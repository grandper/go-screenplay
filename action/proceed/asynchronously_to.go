package proceed

import (
	impl "github.com/grandper/go-screenplay/v2/action/internal/proceed"
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
func AsynchronouslyTo(actions ...screenplay.Performable) *impl.AsynchronouslyToAction {
	return impl.AsynchronouslyTo(actions...)
}
