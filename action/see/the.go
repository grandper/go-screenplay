package see

import (
	"fmt"

	"github.com/grandper/go-screenplay/screenplay"
)

// The creates an action to see if the answer to a question matches the resolution.
func The(question screenplay.Question, resolution screenplay.Resolution) *TheAction {
	return &TheAction{
		question:   question,
		resolution: resolution,
	}
}

// TheAction is an action to see if the answer to a question matches the resolution.
type TheAction struct {
	question   screenplay.Question
	resolution screenplay.Resolution
}

// String describes the action.
func (a *TheAction) String() string {
	return fmt.Sprintf("see if the %s is %s", a.question, a.resolution)
}

// PerformAs performs the task or the action as the provided actor.
func (a *TheAction) PerformAs(actor *screenplay.Actor) error {
	value, err := a.question.AnsweredBy(actor)
	if err != nil {
		return err
	}

	matched, err := a.resolution.Resolve()(value)
	if err != nil {
		return fmt.Errorf("an error occurred when %s attempted to see %s: %w", actor.Name(), a.question, err)
	}

	if !matched {
		return fmt.Errorf("%s failed to see the %s", actor.Name(), a.question)
	}
	return nil
}

// TheAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*TheAction)(nil)
