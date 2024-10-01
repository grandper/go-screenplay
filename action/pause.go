package action

import (
	"errors"
	"fmt"
	"time"

	"github.com/grandper/go-screenplay/screenplay"
)

// PauseFor creates an action that pauses for a given time.
// Hard wait should be avoided when possible.
func PauseFor(number int) *PauseActionBuilder {
	return &PauseActionBuilder{
		number: number,
	}
}

// PauseActionBuilder is a builder for creating a PauseAction.
type PauseActionBuilder struct {
	number int
}

// Seconds specifies that the count is in seconds.
func (p *PauseActionBuilder) Seconds() *PauseAction {
	return &PauseAction{
		number:   p.number,
		unit:     singularOrPlural(p.number, "second", "seconds"),
		duration: time.Duration(p.number) * time.Second,
		reason:   "",
	}
}

// Second specifies that the count is in seconds.
func (p *PauseActionBuilder) Second() *PauseAction {
	return p.Seconds()
}

// Milliseconds specifies that the count is in milliseconds.
func (p *PauseActionBuilder) Milliseconds() *PauseAction {
	return &PauseAction{
		number:   p.number,
		unit:     singularOrPlural(p.number, "millisecond", "milliseconds"),
		duration: time.Duration(p.number) * time.Millisecond,
		reason:   "",
	}
}

// Millisecond specifies that the count is in milliseconds.
func (p *PauseActionBuilder) Millisecond() *PauseAction {
	return p.Milliseconds()
}

func singularOrPlural(count int, singular, plural string) string {
	if count > 1 {
		return plural
	}
	return singular
}

// PauseAction is an action to see if the answer to a question matches the resolution.
// Hard wait should be avoided when possible.
type PauseAction struct {
	number   int
	unit     string
	duration time.Duration
	reason   string
}

// Because specifies the reason for making a PauseActionBuilder.
func (a *PauseAction) Because(reason string) *PauseAction {
	a.reason = reason
	return a
}

// String describes the action.
func (a *PauseAction) String() string {
	return fmt.Sprintf("PauseActionBuilder for %d %s because %s", a.number, a.unit, a.reason)
}

// PerformAs performs the task or the action as the provided actor.
func (a *PauseAction) PerformAs(_ *screenplay.Actor) error {
	if a.reason == "" {
		return errors.New(
			"failed to PauseActionBuilder: cannot PauseActionBuilder with a reason: you must call the .Because() method",
		)
	}
	time.Sleep(a.duration)
	return nil
}

// PauseAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*PauseAction)(nil)
