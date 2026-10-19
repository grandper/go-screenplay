package proceed

import (
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/v2/action/internal/performables"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

var (
	// ErrNoPerformableProvided is returned when EitherTo is performed without performables.
	ErrNoPerformableProvided = errors.New("you must provide at least one performable")

	// ErrNoAlternateProvided is returned when EitherTo is performed without alternate performables.
	ErrNoAlternateProvided = errors.New("you must provide at least one alternate performable")
)

// EitherTo performs the provided actions or tasks, or the alternate ones given
// to Or when the first ones fail:
//
//	proceed.EitherTo(RestoreTheSession()).Or(LogInAs("wanda"))
func EitherTo(actions ...screenplay.Performable) *EitherToAction {
	return &EitherToAction{
		performables:          performables.WithoutNil(actions),
		alternatePerformables: []screenplay.Performable{},
	}
}

// EitherToAction is an action that performs its performables, or an alternate
// set of performables when the first ones fail.
type EitherToAction struct {
	performables          []screenplay.Performable
	alternatePerformables []screenplay.Performable
}

// Or provides the alternative tasks or actions to perform.
func (a *EitherToAction) Or(alternatePerformables ...screenplay.Performable) *EitherToAction {
	a.alternatePerformables = performables.WithoutNil(alternatePerformables)

	return a
}

// String describes the action.
func (a *EitherToAction) String() string {
	return fmt.Sprintf("proceed either to %s or %s",
		performables.Describe(a.performables),
		performables.Describe(a.alternatePerformables))
}

// PerformAs performs the task or the action as the provided actor.
func (a *EitherToAction) PerformAs(theActor *screenplay.Actor) error {
	if len(a.performables) == 0 {
		return ErrNoPerformableProvided
	}

	if len(a.alternatePerformables) == 0 {
		return ErrNoAlternateProvided
	}

	if err := theActor.AttemptsTo(a.performables...); err == nil {
		return nil
	}

	return theActor.AttemptsTo(a.alternatePerformables...)
}

// EitherToAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*EitherToAction)(nil)
