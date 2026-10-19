package wait

import (
	impl "github.com/grandper/go-screenplay/v2/action/internal/wait"
	"github.com/grandper/go-screenplay/v2/timing"
)

// ErrReasonNotProvided is returned when a wait is performed without a reason.
var ErrReasonNotProvided = impl.ErrReasonNotProvided

// For creates an action that waits for a given time. The returned
// DurationBuilder lets the caller pick a unit, for example
// wait.For(20).Milliseconds().
// Hard wait should be avoided when possible.
func For(number int) *timing.DurationBuilder[impl.ForAction] {
	return impl.For(number)
}
