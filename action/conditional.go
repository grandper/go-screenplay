package action

import (
	"fmt"

	"github.com/grandper/go-screenplay/screenplay"
)

// condition reports whether the guarded action should run, and describes itself.
type condition struct {
	describe string
	holds    func(theActor *screenplay.Actor) (bool, error)
}

// Conditionally starts a conditional action for the given performables.
func Conditionally(performables ...screenplay.Performable) *ConditionalAction {
	return &ConditionalAction{
		consequence: performables,
		alternative: []screenplay.Performable{},
		condition: condition{
			describe: "the condition holds",
			holds:    func(*screenplay.Actor) (bool, error) { return true, nil },
		},
	}
}

// ConditionalAction is an action that performs its performables only when a condition holds.
type ConditionalAction struct {
	consequence []screenplay.Performable
	alternative []screenplay.Performable
	condition   condition
}

// If guards the action behind a boolean expression.
func (a *ConditionalAction) If(cond bool) *ConditionalAction {
	a.condition = condition{
		describe: "the condition holds",
		holds:    func(*screenplay.Actor) (bool, error) { return cond, nil },
	}

	return a
}

// When guards the action behind a boolean expression.
func (a *ConditionalAction) When(cond bool) *ConditionalAction {
	return a.If(cond)
}

// Unless guards the action behind the negation of a boolean expression.
func (a *ConditionalAction) Unless(cond bool) *ConditionalAction {
	return a.If(!cond)
}

// IfThe guards the action behind a question and a resolution.
func (a *ConditionalAction) IfThe(
	question screenplay.Question,
	resolution screenplay.Resolution,
) *ConditionalAction {
	a.condition = questionCondition(question, resolution)

	return a
}

// questionCondition builds a condition that answers a question and matches its answer
// against a resolution. It is shared by IfThe (conditional) and WhenThe (choice).
func questionCondition(question screenplay.Question, resolution screenplay.Resolution) condition {
	return condition{
		describe: fmt.Sprintf("the %s is %s", question, resolution),
		holds: func(theActor *screenplay.Actor) (bool, error) {
			value, err := question.AnsweredBy(theActor)
			if err != nil {
				return false, err
			}

			matched, err := resolution.Resolve()(value)
			if err != nil {
				return false, fmt.Errorf(
					"an error occurred when %s evaluated the condition on %s: %w",
					theActor.Name(), question, err)
			}

			return matched, nil
		},
	}
}

// WhenThe guards the action behind a question and a resolution.
func (a *ConditionalAction) WhenThe(
	question screenplay.Question,
	resolution screenplay.Resolution,
) *ConditionalAction {
	return a.IfThe(question, resolution)
}

// UnlessThe guards the action behind the negation of a question/resolution condition.
func (a *ConditionalAction) UnlessThe(
	question screenplay.Question,
	resolution screenplay.Resolution,
) *ConditionalAction {
	a.IfThe(question, resolution)
	inner := a.condition.holds
	a.condition.describe = fmt.Sprintf("the %s is not %s", question, resolution)
	a.condition.holds = func(theActor *screenplay.Actor) (bool, error) {
		held, err := inner(theActor)

		return !held, err
	}

	return a
}

// Otherwise provides the performables to perform when the condition does not hold.
func (a *ConditionalAction) Otherwise(performables ...screenplay.Performable) *ConditionalAction {
	a.alternative = performables

	return a
}

// String describes the action.
func (a *ConditionalAction) String() string {
	description := fmt.Sprintf("%s if %s", performablesToString(a.consequence), a.condition.describe)
	if len(a.alternative) > 0 {
		description += ", otherwise " + performablesToString(a.alternative)
	}

	return description
}

// PerformAs performs the task or the action as the provided actor.
func (a *ConditionalAction) PerformAs(theActor *screenplay.Actor) error {
	holds, err := a.condition.holds(theActor)
	if err != nil {
		return err
	}

	if holds {
		return theActor.AttemptsTo(a.consequence...)
	}

	return theActor.AttemptsTo(a.alternative...)
}

// ConditionalAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*ConditionalAction)(nil)
