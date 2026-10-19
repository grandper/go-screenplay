package wait

import (
	impl "github.com/grandper/go-screenplay/v2/action/internal/wait"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// UntilThe makes an actor wait until the answer to a question happens to be as
// stated, for no longer than the timeout of the actor's production. Close it
// with a verb and one or more resolutions:
//
//	wait.UntilThe(logs).HappenTo(contain.Substring("ready"))
func UntilThe[T any](question screenplay.Question[T]) *impl.UntilTheBuilder[T] {
	return impl.UntilThe(question)
}
