package see

import (
	"fmt"
	"strings"

	"github.com/grandper/go-screenplay/screenplay"
)

// AnyOf creates an action to see if any of the answers to different questions match its expectation.
func AnyOf(tuples ...any) *AnyOfAction {
	if len(tuples)%2 != 0 {
		panic("you should provide question-resolution pairs to the see.AnyOf function")
	}

	action := &AnyOfAction{
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

// AnyOfAction is an action to see if any of the answers to different questions matches the resolution.
type AnyOfAction struct {
	tests []*TheAction
}

// String describes the action.
func (a *AnyOfAction) String() string {
	strs := make([]string, 0, len(a.tests))
	for _, test := range a.tests {
		strs = append(strs, test.String())
	}

	return strings.Join(strs, ", or ")
}

// PerformAs performs the task or the action as the provided actor.
func (a *AnyOfAction) PerformAs(actor *screenplay.Actor) error {
	if len(a.tests) == 0 {
		return nil
	}

	for _, test := range a.tests {
		err := actor.AttemptsTo(test)
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("%s did not find any expected answers", actor.Name())
}

// AnyOfAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*AnyOfAction)(nil)
