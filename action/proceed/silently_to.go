package proceed

import (
	"errors"

	"github.com/grandper/go-screenplay/v2/screenplay"
)

// SilentlyTo wraps a performable so that the narration produced while it runs
// is muted: the nested beats and asides it would emit never reach the
// narrator's adapters, collapsing a whole sub-tree of actions into a single,
// opaque line that reads "proceed silently to perform an action". The wrapped
// performable's own description is never narrated, so a value it carries never
// reaches an adapter.
//
// It is meant to keep secrets out of the logs — an action that enters a password
// or sets an authorization header can be silenced so its value never reaches an
// adapter — and to trim the noise of a verbose task down to a single line.
//
// A production configured with screenplay.WithForceAllNarration neutralizes
// this decorator, so a scenario can be debugged with every step narrated again.
func SilentlyTo(performable screenplay.Performable) screenplay.Performable {
	return &silentlyToAction{
		performable: performable,
	}
}

// silentlyToAction is a performable that mutes the narration produced while the
// wrapped performable runs.
type silentlyToAction struct {
	performable screenplay.Performable
}

// String describes the action without describing the wrapped performable, so
// that a value the wrapped performable carries is never narrated.
func (a *silentlyToAction) String() string {
	return "proceed silently to perform an action"
}

// PerformAs performs the wrapped performable through a muted view of the actor,
// so the narration it would produce is dropped. The view leaves the actor
// itself untouched, so a SilentlyTo running in one concurrent branch never
// mutes another. When the actor's narrator cannot be muted (see
// screenplay.WithForceAllNarration), the steps narrate as usual.
func (a *silentlyToAction) PerformAs(theActor *screenplay.Actor) error {
	if a.performable == nil {
		return errors.New("the performable given to SilentlyTo is nil")
	}

	return a.performable.PerformAs(theActor.Muted())
}

// silentlyToAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*silentlyToAction)(nil)
