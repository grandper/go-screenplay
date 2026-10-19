package proceed

import (
	"fmt"

	"github.com/grandper/go-screenplay/v2/action/internal/condition"
	"github.com/grandper/go-screenplay/v2/action/internal/performables"
	"github.com/grandper/go-screenplay/v2/internal/sentence"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// To starts an action that performs the given performables. Guard it with a
// condition to perform them only when the condition holds:
//
//	proceed.To(RetryTheRequest()).IfThe(statusCode).HappensTo(be.AtLeast(500))
func To(actions ...screenplay.Performable) *ToAction {
	return &ToAction{
		consequence: performables.WithoutNil(actions),
		alternative: []screenplay.Performable{},
	}
}

// ToAction is an action that performs its performables, only when its
// condition holds when it has one.
type ToAction struct {
	consequence []screenplay.Performable
	alternative []screenplay.Performable
	condition   *condition.Condition
}

// If guards the action behind a boolean expression.
func (a *ToAction) If(cond bool) *ToAction {
	return a.guardedBy(condition.Boolean(cond))
}

// When guards the action behind a boolean expression. It is narrated "when
// the condition holds".
func (a *ToAction) When(cond bool) *ToAction {
	return a.guardedBy(condition.Boolean(cond).IntroducedBy("when"))
}

// Unless guards the action behind the negation of a boolean expression. It is
// narrated "unless the condition holds".
func (a *ToAction) Unless(cond bool) *ToAction {
	return a.guardedBy(condition.Boolean(cond).Negated())
}

// IfThe guards the action behind what the answer to a question happens to be.
// Close it with a verb and one or more resolutions:
//
//	proceed.To(RetryTheRequest()).IfThe(statusCode).HappensTo(be.AtLeast(500))
func (a *ToAction) IfThe[T any](question screenplay.Question[T]) *ToBuilder[T] {
	return &ToBuilder[T]{action: a, question: question, keyword: "if"}
}

// WhenThe guards the action behind what the answer to a question happens to be.
// It behaves like IfThe and is narrated with "when":
//
//	proceed.To(RetryTheRequest()).WhenThe(statusCode).HappensTo(be.AtLeast(500))
func (a *ToAction) WhenThe[T any](question screenplay.Question[T]) *ToBuilder[T] {
	return &ToBuilder[T]{action: a, question: question, keyword: "when"}
}

// UnlessThe guards the action behind the negation of what the answer to a
// question happens to be: the action is performed unless the sentence holds.
//
//	proceed.To(LogTheResponse()).UnlessThe(statusCode).HappensTo(be.EqualTo(200))
func (a *ToAction) UnlessThe[T any](question screenplay.Question[T]) *ToBuilder[T] {
	return &ToBuilder[T]{action: a, question: question, keyword: "unless"}
}

// ToBuilder builds the condition of a proceed.To action: it holds the question
// the action is guarded by, and waits for the verb that closes the condition.
type ToBuilder[T any] struct {
	action   *ToAction
	question screenplay.Question[T]
	keyword  string
}

// HappensTo guards the action behind the answer satisfying every resolution. It reads "happens to ...".
func (s *ToBuilder[T]) HappensTo(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.guards(sentence.NewCondition(sentence.Should, false, resolutions))
}

// HappensToBe guards the action behind the answer satisfying every resolution. It reads "happens to be ...".
func (s *ToBuilder[T]) HappensToBe(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.guards(sentence.NewCondition(sentence.ShouldBe, false, resolutions))
}

// HappensToHave guards the action behind the answer satisfying every resolution. It reads "happens to have ...".
func (s *ToBuilder[T]) HappensToHave(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.guards(sentence.NewCondition(sentence.ShouldHave, false, resolutions))
}

// DoesNotHappenTo guards the action behind the answer satisfying none of the resolutions. It reads "does not happen to ...".
func (s *ToBuilder[T]) DoesNotHappenTo(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.guards(sentence.NewCondition(sentence.ShouldNot, false, resolutions))
}

// DoesNotHappenToBe guards the action behind the answer satisfying none of the resolutions. It reads "does not happen to be ...".
func (s *ToBuilder[T]) DoesNotHappenToBe(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.guards(sentence.NewCondition(sentence.ShouldNotBe, false, resolutions))
}

// DoesNotHappenToHave guards the action behind the answer satisfying none of the resolutions. It reads "does not happen to have ...".
func (s *ToBuilder[T]) DoesNotHappenToHave(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.guards(sentence.NewCondition(sentence.ShouldNotHave, false, resolutions))
}

// HappensToBeAnyOf guards the action behind the answer satisfying at least one resolution. It reads "happens to be any of: ...".
func (s *ToBuilder[T]) HappensToBeAnyOf(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.guards(sentence.NewCondition(sentence.ShouldAnyOf, false, resolutions))
}

// HappenTo guards the action behind the answer satisfying every resolution. It is the plural form: it reads "happen to ...".
func (s *ToBuilder[T]) HappenTo(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.guards(sentence.NewCondition(sentence.Should, true, resolutions))
}

// HappenToBe guards the action behind the answer satisfying every resolution. It is the plural form: it reads "happen to be ...".
func (s *ToBuilder[T]) HappenToBe(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.guards(sentence.NewCondition(sentence.ShouldBe, true, resolutions))
}

// HappenToHave guards the action behind the answer satisfying every resolution. It is the plural form: it reads "happen to have ...".
func (s *ToBuilder[T]) HappenToHave(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.guards(sentence.NewCondition(sentence.ShouldHave, true, resolutions))
}

// DoNotHappenTo guards the action behind the answer satisfying none of the resolutions. It is the plural form: it reads "do not happen to ...".
func (s *ToBuilder[T]) DoNotHappenTo(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.guards(sentence.NewCondition(sentence.ShouldNot, true, resolutions))
}

// DoNotHappenToBe guards the action behind the answer satisfying none of the resolutions. It is the plural form: it reads "do not happen to be ...".
func (s *ToBuilder[T]) DoNotHappenToBe(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.guards(sentence.NewCondition(sentence.ShouldNotBe, true, resolutions))
}

// DoNotHappenToHave guards the action behind the answer satisfying none of the resolutions. It is the plural form: it reads "do not happen to have ...".
func (s *ToBuilder[T]) DoNotHappenToHave(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.guards(sentence.NewCondition(sentence.ShouldNotHave, true, resolutions))
}

// HappenToBeAnyOf guards the action behind the answer satisfying at least one resolution. It is the plural form: it reads "happen to be any of: ...".
func (s *ToBuilder[T]) HappenToBeAnyOf(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.guards(sentence.NewCondition(sentence.ShouldAnyOf, true, resolutions))
}

// guards sets the condition of the action: the answer satisfies the clause.
func (s *ToBuilder[T]) guards(clause sentence.Clause[T]) *ToAction {
	guard := condition.Term(sentence.NewTerm(s.question, clause))
	if s.keyword == "unless" {
		guard = guard.Negated()
	} else {
		guard = guard.IntroducedBy(s.keyword)
	}

	return s.action.guardedBy(guard)
}

// guardedBy sets the condition of the action.
func (a *ToAction) guardedBy(guard condition.Condition) *ToAction {
	a.condition = &guard

	return a
}

// Otherwise provides the performables to perform when the condition does not hold.
func (a *ToAction) Otherwise(actions ...screenplay.Performable) *ToAction {
	a.alternative = performables.WithoutNil(actions)

	return a
}

// String describes the action.
func (a *ToAction) String() string {
	description := "proceed to " + performables.Describe(a.consequence)
	if a.condition != nil {
		description += fmt.Sprintf(" %s %s", a.condition.Keyword, a.condition.Description)
	}

	if len(a.alternative) > 0 {
		description += ", otherwise " + performables.Describe(a.alternative)
	}

	return description
}

// PerformAs performs the consequence when the condition holds, or when there
// is no condition, and the alternative otherwise.
func (a *ToAction) PerformAs(theActor *screenplay.Actor) error {
	if a.condition == nil {
		return theActor.AttemptsTo(a.consequence...)
	}

	holds, err := a.condition.HoldsFor(theActor)
	if err != nil {
		return err
	}

	if holds {
		return theActor.AttemptsTo(a.consequence...)
	}

	return theActor.AttemptsTo(a.alternative...)
}

// ToAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*ToAction)(nil)
