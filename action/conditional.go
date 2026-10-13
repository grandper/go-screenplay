package action

import (
	"fmt"

	"github.com/grandper/go-screenplay/v2/internal/sentence"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// condition reports whether the guarded action should run, and describes itself.
// keyword introduces the description: "if" or "unless".
type condition struct {
	keyword  string
	describe string
	holds    func(theActor *screenplay.Actor) (bool, error)
}

// always is the condition of a conditional action that has not been guarded yet.
func always() condition {
	return booleanCondition(true)
}

// booleanCondition builds a condition from a boolean expression.
func booleanCondition(cond bool) condition {
	return condition{
		keyword:  "if",
		describe: "the condition holds",
		holds:    func(*screenplay.Actor) (bool, error) { return cond, nil },
	}
}

// termCondition builds a condition that holds when a term is met. It is shared
// by IfThe (conditional) and WhenThe (choice).
func termCondition(t *sentence.Term) condition {
	return condition{
		keyword:  "if",
		describe: t.String(),
		holds: func(actor *screenplay.Actor) (bool, error) {
			verdict, err := t.JudgedBy(actor)
			if err != nil {
				return false, err
			}

			return verdict.Met, nil
		},
	}
}

// negated returns the condition that holds when the receiver does not.
func (c condition) negated() condition {
	return condition{
		keyword:  "unless",
		describe: c.describe,
		holds: func(theActor *screenplay.Actor) (bool, error) {
			held, err := c.holds(theActor)

			return !held, err
		},
	}
}

// Conditionally starts a conditional action for the given performables.
func Conditionally(performables ...screenplay.Performable) *ConditionalAction {
	return &ConditionalAction{
		consequence: withoutNil(performables),
		alternative: []screenplay.Performable{},
		condition:   always(),
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
	a.condition = booleanCondition(cond)

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

// IfThe guards the action behind what the answer to a question happens to be.
// Close it with a verb and one or more resolutions:
//
//	Conditionally(task).IfThe(statusCode).HappensTo(be.EqualTo(200))
func (a *ConditionalAction) IfThe[T any](question screenplay.Question[T]) *ConditionalBuilder[T] {
	return &ConditionalBuilder[T]{action: a, question: question, unless: false}
}

// WhenThe guards the action behind what the answer to a question happens to be.
// It is an alias of IfThe.
func (a *ConditionalAction) WhenThe[T any](question screenplay.Question[T]) *ConditionalBuilder[T] {
	return a.IfThe(question)
}

// UnlessThe guards the action behind the negation of what the answer to a
// question happens to be: the action is performed unless the sentence holds.
//
//	Conditionally(task).UnlessThe(statusCode).HappensTo(be.EqualTo(200))
func (a *ConditionalAction) UnlessThe[T any](question screenplay.Question[T]) *ConditionalBuilder[T] {
	return &ConditionalBuilder[T]{action: a, question: question, unless: true}
}

// ConditionalBuilder builds the condition of a conditional action: it holds
// the question the action is guarded by, and waits for the verb that closes
// the condition.
type ConditionalBuilder[T any] struct {
	action   *ConditionalAction
	question screenplay.Question[T]
	unless   bool
}

// HappensTo guards the action behind the answer satisfying every resolution. It reads "happens to ...".
func (s *ConditionalBuilder[T]) HappensTo(resolutions ...screenplay.Resolution[T]) *ConditionalAction {
	return s.guards(sentence.NewCondition(sentence.Should, false, resolutions))
}

// HappensToBe guards the action behind the answer satisfying every resolution. It reads "happens to be ...".
func (s *ConditionalBuilder[T]) HappensToBe(resolutions ...screenplay.Resolution[T]) *ConditionalAction {
	return s.guards(sentence.NewCondition(sentence.ShouldBe, false, resolutions))
}

// HappensToHave guards the action behind the answer satisfying every resolution. It reads "happens to have ...".
func (s *ConditionalBuilder[T]) HappensToHave(resolutions ...screenplay.Resolution[T]) *ConditionalAction {
	return s.guards(sentence.NewCondition(sentence.ShouldHave, false, resolutions))
}

// DoesNotHappenTo guards the action behind the answer satisfying none of the resolutions. It reads "does not happen to ...".
func (s *ConditionalBuilder[T]) DoesNotHappenTo(resolutions ...screenplay.Resolution[T]) *ConditionalAction {
	return s.guards(sentence.NewCondition(sentence.ShouldNot, false, resolutions))
}

// DoesNotHappenToBe guards the action behind the answer satisfying none of the resolutions. It reads "does not happen to be ...".
func (s *ConditionalBuilder[T]) DoesNotHappenToBe(resolutions ...screenplay.Resolution[T]) *ConditionalAction {
	return s.guards(sentence.NewCondition(sentence.ShouldNotBe, false, resolutions))
}

// DoesNotHappenToHave guards the action behind the answer satisfying none of the resolutions. It reads "does not happen to have ...".
func (s *ConditionalBuilder[T]) DoesNotHappenToHave(resolutions ...screenplay.Resolution[T]) *ConditionalAction {
	return s.guards(sentence.NewCondition(sentence.ShouldNotHave, false, resolutions))
}

// HappensToBeAnyOf guards the action behind the answer satisfying at least one resolution. It reads "happens to be any of: ...".
func (s *ConditionalBuilder[T]) HappensToBeAnyOf(resolutions ...screenplay.Resolution[T]) *ConditionalAction {
	return s.guards(sentence.NewCondition(sentence.ShouldAnyOf, false, resolutions))
}

// HappenTo guards the action behind the answer satisfying every resolution. It is the plural form: it reads "happen to ...".
func (s *ConditionalBuilder[T]) HappenTo(resolutions ...screenplay.Resolution[T]) *ConditionalAction {
	return s.guards(sentence.NewCondition(sentence.Should, true, resolutions))
}

// HappenToBe guards the action behind the answer satisfying every resolution. It is the plural form: it reads "happen to be ...".
func (s *ConditionalBuilder[T]) HappenToBe(resolutions ...screenplay.Resolution[T]) *ConditionalAction {
	return s.guards(sentence.NewCondition(sentence.ShouldBe, true, resolutions))
}

// HappenToHave guards the action behind the answer satisfying every resolution. It is the plural form: it reads "happen to have ...".
func (s *ConditionalBuilder[T]) HappenToHave(resolutions ...screenplay.Resolution[T]) *ConditionalAction {
	return s.guards(sentence.NewCondition(sentence.ShouldHave, true, resolutions))
}

// DoNotHappenTo guards the action behind the answer satisfying none of the resolutions. It is the plural form: it reads "do not happen to ...".
func (s *ConditionalBuilder[T]) DoNotHappenTo(resolutions ...screenplay.Resolution[T]) *ConditionalAction {
	return s.guards(sentence.NewCondition(sentence.ShouldNot, true, resolutions))
}

// DoNotHappenToBe guards the action behind the answer satisfying none of the resolutions. It is the plural form: it reads "do not happen to be ...".
func (s *ConditionalBuilder[T]) DoNotHappenToBe(resolutions ...screenplay.Resolution[T]) *ConditionalAction {
	return s.guards(sentence.NewCondition(sentence.ShouldNotBe, true, resolutions))
}

// DoNotHappenToHave guards the action behind the answer satisfying none of the resolutions. It is the plural form: it reads "do not happen to have ...".
func (s *ConditionalBuilder[T]) DoNotHappenToHave(resolutions ...screenplay.Resolution[T]) *ConditionalAction {
	return s.guards(sentence.NewCondition(sentence.ShouldNotHave, true, resolutions))
}

// HappenToBeAnyOf guards the action behind the answer satisfying at least one resolution. It is the plural form: it reads "happen to be any of: ...".
func (s *ConditionalBuilder[T]) HappenToBeAnyOf(resolutions ...screenplay.Resolution[T]) *ConditionalAction {
	return s.guards(sentence.NewCondition(sentence.ShouldAnyOf, true, resolutions))
}

// guards sets the condition of the action: the answer satisfies the clause.
func (s *ConditionalBuilder[T]) guards(clause sentence.Clause[T]) *ConditionalAction {
	guard := termCondition(sentence.NewTerm(s.question, clause))
	if s.unless {
		guard = guard.negated()
	}

	s.action.condition = guard

	return s.action
}

// Otherwise provides the performables to perform when the condition does not hold.
func (a *ConditionalAction) Otherwise(performables ...screenplay.Performable) *ConditionalAction {
	a.alternative = withoutNil(performables)

	return a
}

// String describes the action.
func (a *ConditionalAction) String() string {
	description := fmt.Sprintf(
		"conditionally %s %s %s",
		performablesToString(a.consequence),
		a.condition.keyword,
		a.condition.describe,
	)
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
