package see

import (
	"strings"

	"github.com/grandper/go-screenplay/screenplay"
)

// AllOf creates an action to see if all the answers of different questions match their expectation.
func AllOf(tuples ...any) *AllOfAction {
	if len(tuples)%2 != 0 {
		panic("you should provide question-resolution pairs to the see.AllOf function")
	}

	action := &AllOfAction{
		tests: []*TheAction{},
	}

	for i := 0; i < len(tuples); i++ {
		question, isAQuestion := tuples[i].(screenplay.Question)
		if !isAQuestion {
			panic("the tuple must contain a question")
		}

		i++
		resolution, isAResolution := tuples[i].(screenplay.Resolution)

		if !isAResolution {
			panic("the tuple must contain a resolution")
		}

		action.tests = append(action.tests, The(question, resolution))
	}

	return action
}

// AllOfAction is an action to see if all the answers to different questions match the resolution.
type AllOfAction struct {
	tests []*TheAction
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
