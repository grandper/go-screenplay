package keep

import (
	impl "github.com/grandper/go-screenplay/v2/action/internal/keep"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

var (
	// ErrPollingPeriodMustBeLessThanOrEqualToTimeout is returned when the polling period is greater than the timeout.
	ErrPollingPeriodMustBeLessThanOrEqualToTimeout = impl.ErrPollingPeriodMustBeLessThanOrEqualToTimeout
	// ErrPollingPeriodAndTimeoutMustBePositive is returned when the polling period or the timeout is zero or negative.
	ErrPollingPeriodAndTimeoutMustBePositive = impl.ErrPollingPeriodAndTimeoutMustBePositive
	// ErrNoAttemptCompleted is returned when the timeout is reached while the only attempt is still running.
	ErrNoAttemptCompleted = impl.ErrNoAttemptCompleted
)

// TryingTo retries a task or action until it succeeds, for the time and at the
// pace set by its window, or by the actor's production when the window is not
// set:
//
//	keep.TryingTo(CancelTheOrder()).For(5).Seconds().PollingEvery(500).Milliseconds()
//
// If the action does not succeed before the end of the window, an error
// containing all the distinct errors of the tries is returned.
func TryingTo(performable screenplay.Performable) *impl.TryingToAction {
	return impl.TryingTo(performable)
}
