package see

import (
	"errors"
	"strconv"
	"strings"

	"github.com/grandper/go-screenplay/screenplay"
)

// ErrInvalidAllOfArguments is returned when the arguments passed to AllOf are not valid question-resolution pairs.
var ErrInvalidAllOfArguments = errors.New("invalid arguments: you should provide question-resolution pairs to the see.AllOf function")

// AllOf creates an action to see if all the answers of different questions match their expectation.
func AllOf(tuples ...any) *AllOfAction {
	if len(tuples)%2 != 0 {
		return &AllOfAction{err: ErrInvalidAllOfArguments}
	}

	action := &AllOfAction{
		tests: []*TheAction{},
	}

	for i := 0; i < len(tuples); i++ {
		question, isAQuestion := tuples[i].(screenplay.Question)
		if !isAQuestion {
			return &AllOfAction{err: errors.New("invalid arguments: expected a Question at position " + strconv.Itoa(i))}
		}

		i++
		resolution, isAResolution := tuples[i].(screenplay.Resolution)
		if !isAResolution {
			return &AllOfAction{err: errors.New("invalid arguments: expected a Resolution at position " + strconv.Itoa(i))}
		}

		action.tests = append(action.tests, The(question, resolution))
	}

	return action
}

// AllOfAction is an action to see if all the answers to different questions match the resolution.
type AllOfAction struct {
	tests []*TheAction
	err   error
}

// String describes the action.
func (a *AllOfAction) String() string {
	strs := []string{}
	for _, test := range a.tests {
		strs = append(strs, test.String())
	}

	return strings.Join(strs, ", and ")
}

// PerformAs performs the task or the action as the provided actor.
func (a *AllOfAction) PerformAs(actor *screenplay.Actor) error {
	if a.err != nil {
		return a.err
	}

	if len(a.tests) == 0 {
		return nil
	}

	for _, test := range a.tests {
		err := actor.AttemptsTo(test)
		if err != nil {
			return err
		}
	}

	return nil
}

// AllOfAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*AllOfAction)(nil)
