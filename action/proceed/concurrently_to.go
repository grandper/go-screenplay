package proceed

import (
	impl "github.com/grandper/go-screenplay/v2/action/internal/proceed"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// ConcurrentlyTo performs the provided tasks or actions in parallel.
// By default it waits for every action to complete and joins their errors,
// which is equivalent to calling WaitingForAll. This behavior can be changed
// with StoppingOnError or IgnoringErrors, and the number of actions running at
// the same time can be capped with WithLimit.
func ConcurrentlyTo(actions ...screenplay.Performable) *impl.ConcurrentlyToAction {
	return impl.ConcurrentlyTo(actions...)
}

// SimultaneouslyTo performs the provided tasks or actions in parallel. It
// behaves like ConcurrentlyTo and is narrated with "simultaneously".
func SimultaneouslyTo(actions ...screenplay.Performable) *impl.ConcurrentlyToAction {
	return impl.SimultaneouslyTo(actions...)
}
