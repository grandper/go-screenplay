package action

import (
	"errors"
	"fmt"
	"strings"

	"github.com/grandper/go-screenplay/screenplay"
)

// Either executes one of the two provided actions or tasks.
func Either(performables ...screenplay.Performable) *EitherAction {
	return &EitherAction{
		performables:          performables,
		alternatePerformables: []screenplay.Performable{},
	}
}

// EitherAction is an action that makes a note about the answer of a question.
type EitherAction struct {
	performables          []screenplay.Performable
	alternatePerformables []screenplay.Performable
}

// String describes the action.
func (a *EitherAction) String() string {
	return fmt.Sprintf("either %s or %s",
		performablesToString(a.performables),
		performablesToString(a.alternatePerformables))
}

func performablesToString(performables []screenplay.Performable) string {
	strs := performablesToStringSlice(performables)
	switch len(strs) {
	case 1:
		return strs[0]
	default:
		return strings.Join(strs, ", ")
	}
}

func performablesToStringSlice(performables []screenplay.Performable) []string {
	strs := make([]string, 0, len(performables))
	for _, performable := range performables {
		strs = append(strs, performable.String())
	}

	return strs
}

// PerformAs performs the task or the action as the provided actor.
func (a *EitherAction) PerformAs(theActor *screenplay.Actor) error {
	if len(a.performables) == 0 {
		return errors.New("you must provide at least one performable")
	}

	if len(a.alternatePerformables) == 0 {
		return errors.New("you must provide at least one alternate performable")
	}

	if err := theActor.AttemptsTo(a.performables...); err == nil {
		return nil
	}

	return theActor.AttemptsTo(a.alternatePerformables...)
}

// Or provides an alternative task or action to perform.
func (a *EitherAction) Or(alternatePerformables ...screenplay.Performable) *EitherAction {
	a.alternatePerformables = alternatePerformables

	return a
}

// Except provides an alternative task or action to perform.
func (a *EitherAction) Except(alternatePerformables ...screenplay.Performable) *EitherAction {
	return a.Or(alternatePerformables...)
}

// Else provides an alternative task or action to perform.
func (a *EitherAction) Else(alternatePerformables ...screenplay.Performable) *EitherAction {
	return a.Or(alternatePerformables...)
}

// Otherwise provides an alternative task or action to perform.
func (a *EitherAction) Otherwise(alternatePerformables ...screenplay.Performable) *EitherAction {
	return a.Or(alternatePerformables...)
}

// Alternatively provides an alternative task or action to perform.
func (a *EitherAction) Alternatively(alternatePerformables ...screenplay.Performable) *EitherAction {
	return a.Or(alternatePerformables...)
}

// FailingThat provides an alternative task or action to perform.
func (a *EitherAction) FailingThat(alternatePerformables ...screenplay.Performable) *EitherAction {
	return a.Or(alternatePerformables...)
}

// EitherAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*EitherAction)(nil)
