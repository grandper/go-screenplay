package action

import (
	"fmt"
	"strings"

	"github.com/grandper/go-screenplay/v2/internal/sentence"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// ordinal renders 1 -> "first", 2 -> "second", ... for the boolean-branch description.
func ordinal(position int) string {
	names := []string{"first", "second", "third", "fourth", "fifth"}
	if position >= 1 && position <= len(names) {
		return names[position-1]
	}

	return fmt.Sprintf("%dth", position)
}

// joinChoices assembles "choose to A when x, to B when y, or to Z otherwise".
// Each clause is already in the infinitive form "to <action> when <condition>"; the
// fallback becomes the last alternative "to <action> otherwise", and the list is closed
// with "or" so it reads as picking exactly one branch (not a to-do list).
func joinChoices(clauses []string, fallback []screenplay.Performable) string {
	alternatives := make([]string, len(clauses), len(clauses)+1)
	copy(alternatives, clauses)
	if len(fallback) > 0 {
		alternatives = append(alternatives, "to "+performablesToString(fallback)+" otherwise")
	}

	if len(alternatives) <= 1 {
		return "choose " + strings.Join(alternatives, "")
	}

	last := len(alternatives) - 1

	return "choose " + strings.Join(alternatives[:last], ", ") + ", or " + alternatives[last]
}

// choiceBranch pairs a condition with the performables run when it is the first to hold.
// condition is the same type ConditionalAction uses, so When reuses its boolean logic and
// WhenThe reuses its question/resolution logic.
type choiceBranch struct {
	condition    condition
	performables []screenplay.Performable
}

// Choose starts a multi-way choice; each branch's condition is a boolean or a question/resolution.
func Choose() *ChoiceAction {
	return &ChoiceAction{
		branches: []choiceBranch{},
		pending:  []screenplay.Performable{},
		fallback: []screenplay.Performable{},
	}
}

// ChoiceAction performs the performables of the first branch whose condition holds.
type ChoiceAction struct {
	branches []choiceBranch
	pending  []screenplay.Performable
	fallback []screenplay.Performable
}

// To names the performables of the next branch; close it with When or WhenThe.
func (a *ChoiceAction) To(performables ...screenplay.Performable) *ChoiceAction {
	a.pending = performables

	return a
}

// When closes the pending branch with a boolean condition.
func (a *ChoiceAction) When(cond bool) *ChoiceAction {
	branchCondition := booleanCondition(cond)
	branchCondition.describe = fmt.Sprintf("the %s condition holds", ordinal(len(a.branches)+1))

	a.branches = append(a.branches, choiceBranch{
		condition:    branchCondition,
		performables: a.pending,
	})
	a.pending = []screenplay.Performable{}

	return a
}

// WhenThe closes the pending branch with what the answer to a question happens
// to be. Close it with a verb and one or more resolutions:
//
//	Choose().To(celebrate).WhenThe(statusCode).HappensTo(be.EqualTo(200))
func (a *ChoiceAction) WhenThe[T any](question screenplay.Question[T]) *ChoiceBuilder[T] {
	return &ChoiceBuilder[T]{action: a, question: question}
}

// ChoiceBuilder builds the condition of a branch of a choice: it holds the
// question the branch depends on, and waits for the verb that closes it.
type ChoiceBuilder[T any] struct {
	action   *ChoiceAction
	question screenplay.Question[T]
}

// HappensTo closes the branch with the answer satisfying every resolution. It reads "happens to ...".
func (s *ChoiceBuilder[T]) HappensTo(resolutions ...screenplay.Resolution[T]) *ChoiceAction {
	return s.closes(sentence.NewCondition(sentence.Should, false, resolutions))
}

// HappensToBe closes the branch with the answer satisfying every resolution. It reads "happens to be ...".
func (s *ChoiceBuilder[T]) HappensToBe(resolutions ...screenplay.Resolution[T]) *ChoiceAction {
	return s.closes(sentence.NewCondition(sentence.ShouldBe, false, resolutions))
}

// HappensToHave closes the branch with the answer satisfying every resolution. It reads "happens to have ...".
func (s *ChoiceBuilder[T]) HappensToHave(resolutions ...screenplay.Resolution[T]) *ChoiceAction {
	return s.closes(sentence.NewCondition(sentence.ShouldHave, false, resolutions))
}

// DoesNotHappenTo closes the branch with the answer satisfying none of the resolutions. It reads "does not happen to ...".
func (s *ChoiceBuilder[T]) DoesNotHappenTo(resolutions ...screenplay.Resolution[T]) *ChoiceAction {
	return s.closes(sentence.NewCondition(sentence.ShouldNot, false, resolutions))
}

// DoesNotHappenToBe closes the branch with the answer satisfying none of the resolutions. It reads "does not happen to be ...".
func (s *ChoiceBuilder[T]) DoesNotHappenToBe(resolutions ...screenplay.Resolution[T]) *ChoiceAction {
	return s.closes(sentence.NewCondition(sentence.ShouldNotBe, false, resolutions))
}

// DoesNotHappenToHave closes the branch with the answer satisfying none of the resolutions. It reads "does not happen to have ...".
func (s *ChoiceBuilder[T]) DoesNotHappenToHave(resolutions ...screenplay.Resolution[T]) *ChoiceAction {
	return s.closes(sentence.NewCondition(sentence.ShouldNotHave, false, resolutions))
}

// HappensToBeAnyOf closes the branch with the answer satisfying at least one resolution. It reads "happens to be any of: ...".
func (s *ChoiceBuilder[T]) HappensToBeAnyOf(resolutions ...screenplay.Resolution[T]) *ChoiceAction {
	return s.closes(sentence.NewCondition(sentence.ShouldAnyOf, false, resolutions))
}

// HappenTo closes the branch with the answer satisfying every resolution. It is the plural form: it reads "happen to ...".
func (s *ChoiceBuilder[T]) HappenTo(resolutions ...screenplay.Resolution[T]) *ChoiceAction {
	return s.closes(sentence.NewCondition(sentence.Should, true, resolutions))
}

// HappenToBe closes the branch with the answer satisfying every resolution. It is the plural form: it reads "happen to be ...".
func (s *ChoiceBuilder[T]) HappenToBe(resolutions ...screenplay.Resolution[T]) *ChoiceAction {
	return s.closes(sentence.NewCondition(sentence.ShouldBe, true, resolutions))
}

// HappenToHave closes the branch with the answer satisfying every resolution. It is the plural form: it reads "happen to have ...".
func (s *ChoiceBuilder[T]) HappenToHave(resolutions ...screenplay.Resolution[T]) *ChoiceAction {
	return s.closes(sentence.NewCondition(sentence.ShouldHave, true, resolutions))
}

// DoNotHappenTo closes the branch with the answer satisfying none of the resolutions. It is the plural form: it reads "do not happen to ...".
func (s *ChoiceBuilder[T]) DoNotHappenTo(resolutions ...screenplay.Resolution[T]) *ChoiceAction {
	return s.closes(sentence.NewCondition(sentence.ShouldNot, true, resolutions))
}

// DoNotHappenToBe closes the branch with the answer satisfying none of the resolutions. It is the plural form: it reads "do not happen to be ...".
func (s *ChoiceBuilder[T]) DoNotHappenToBe(resolutions ...screenplay.Resolution[T]) *ChoiceAction {
	return s.closes(sentence.NewCondition(sentence.ShouldNotBe, true, resolutions))
}

// DoNotHappenToHave closes the branch with the answer satisfying none of the resolutions. It is the plural form: it reads "do not happen to have ...".
func (s *ChoiceBuilder[T]) DoNotHappenToHave(resolutions ...screenplay.Resolution[T]) *ChoiceAction {
	return s.closes(sentence.NewCondition(sentence.ShouldNotHave, true, resolutions))
}

// HappenToBeAnyOf closes the branch with the answer satisfying at least one resolution. It is the plural form: it reads "happen to be any of: ...".
func (s *ChoiceBuilder[T]) HappenToBeAnyOf(resolutions ...screenplay.Resolution[T]) *ChoiceAction {
	return s.closes(sentence.NewCondition(sentence.ShouldAnyOf, true, resolutions))
}

// closes closes the pending branch: the answer satisfies the clause.
func (s *ChoiceBuilder[T]) closes(clause sentence.Clause[T]) *ChoiceAction {
	s.action.branches = append(s.action.branches, choiceBranch{
		condition:    termCondition(sentence.NewTerm(s.question, clause)),
		performables: s.action.pending,
	})
	s.action.pending = []screenplay.Performable{}

	return s.action
}

// Otherwise closes the pending To branch as the fallback, run when no branch holds.
func (a *ChoiceAction) Otherwise() *ChoiceAction {
	a.fallback = a.pending
	a.pending = []screenplay.Performable{}

	return a
}

// Default closes the pending To branch as the fallback, run when no branch holds.
func (a *ChoiceAction) Default() *ChoiceAction {
	return a.Otherwise()
}

// String describes the action.
func (a *ChoiceAction) String() string {
	clauses := make([]string, 0, len(a.branches))
	for _, b := range a.branches {
		clauses = append(clauses, fmt.Sprintf(
			"to %s when %s", performablesToString(b.performables), b.condition.describe))
	}

	return joinChoices(clauses, a.fallback)
}

// PerformAs performs the task or the action as the provided actor.
func (a *ChoiceAction) PerformAs(theActor *screenplay.Actor) error {
	for _, b := range a.branches {
		holds, err := b.condition.holds(theActor)
		if err != nil {
			return err
		}

		if holds {
			return theActor.AttemptsTo(b.performables...)
		}
	}

	return theActor.AttemptsTo(a.fallback...)
}

// ChoiceAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*ChoiceAction)(nil)

// clauseBranch pairs a clause with the performables run when the answer is the first to satisfy it.
type clauseBranch[T any] struct {
	clause       sentence.Clause[T]
	performables []screenplay.Performable
}

// ChooseBasedOnThe starts a multi-way choice on the answer of a question; each branch states what
// the answer should be.
func ChooseBasedOnThe[T any](question screenplay.Question[T]) *ChoiceOnAction[T] {
	return &ChoiceOnAction[T]{
		question: question,
		branches: []clauseBranch[T]{},
		pending:  []screenplay.Performable{},
		fallback: []screenplay.Performable{},
	}
}

// ChoiceOnAction performs the performables of the first branch whose clause the answer satisfies.
type ChoiceOnAction[T any] struct {
	question screenplay.Question[T]
	branches []clauseBranch[T]
	pending  []screenplay.Performable
	fallback []screenplay.Performable
}

// To names the performables of the next branch; close it with WhenIt.
func (a *ChoiceOnAction[T]) To(performables ...screenplay.Performable) *ChoiceOnAction[T] {
	a.pending = performables

	return a
}

// WhenIt closes the pending branch with what the answer happens to be. Close
// it with a verb and one or more resolutions:
//
//	ChooseBasedOnThe(statusCode).To(celebrate).WhenIt().HappensTo(be.EqualTo(200))
func (a *ChoiceOnAction[T]) WhenIt() *ChoiceOnBuilder[T] {
	return &ChoiceOnBuilder[T]{action: a}
}

// ChoiceOnBuilder builds the condition of a branch of a choice on the answer
// to a question: it waits for the verb that closes it.
type ChoiceOnBuilder[T any] struct {
	action *ChoiceOnAction[T]
}

// HappensTo closes the branch with the answer satisfying every resolution. It reads "happens to ...".
func (s *ChoiceOnBuilder[T]) HappensTo(resolutions ...screenplay.Resolution[T]) *ChoiceOnAction[T] {
	return s.closes(sentence.NewCondition(sentence.Should, false, resolutions))
}

// HappensToBe closes the branch with the answer satisfying every resolution. It reads "happens to be ...".
func (s *ChoiceOnBuilder[T]) HappensToBe(resolutions ...screenplay.Resolution[T]) *ChoiceOnAction[T] {
	return s.closes(sentence.NewCondition(sentence.ShouldBe, false, resolutions))
}

// HappensToHave closes the branch with the answer satisfying every resolution. It reads "happens to have ...".
func (s *ChoiceOnBuilder[T]) HappensToHave(resolutions ...screenplay.Resolution[T]) *ChoiceOnAction[T] {
	return s.closes(sentence.NewCondition(sentence.ShouldHave, false, resolutions))
}

// DoesNotHappenTo closes the branch with the answer satisfying none of the resolutions. It reads "does not happen to ...".
func (s *ChoiceOnBuilder[T]) DoesNotHappenTo(resolutions ...screenplay.Resolution[T]) *ChoiceOnAction[T] {
	return s.closes(sentence.NewCondition(sentence.ShouldNot, false, resolutions))
}

// DoesNotHappenToBe closes the branch with the answer satisfying none of the resolutions. It reads "does not happen to be ...".
func (s *ChoiceOnBuilder[T]) DoesNotHappenToBe(resolutions ...screenplay.Resolution[T]) *ChoiceOnAction[T] {
	return s.closes(sentence.NewCondition(sentence.ShouldNotBe, false, resolutions))
}

// DoesNotHappenToHave closes the branch with the answer satisfying none of the resolutions. It reads "does not happen to have ...".
func (s *ChoiceOnBuilder[T]) DoesNotHappenToHave(resolutions ...screenplay.Resolution[T]) *ChoiceOnAction[T] {
	return s.closes(sentence.NewCondition(sentence.ShouldNotHave, false, resolutions))
}

// HappensToBeAnyOf closes the branch with the answer satisfying at least one resolution. It reads "happens to be any of: ...".
func (s *ChoiceOnBuilder[T]) HappensToBeAnyOf(resolutions ...screenplay.Resolution[T]) *ChoiceOnAction[T] {
	return s.closes(sentence.NewCondition(sentence.ShouldAnyOf, false, resolutions))
}

// closes closes the pending branch: the answer satisfies the clause.
func (s *ChoiceOnBuilder[T]) closes(clause sentence.Clause[T]) *ChoiceOnAction[T] {
	s.action.branches = append(s.action.branches, clauseBranch[T]{
		clause:       clause,
		performables: s.action.pending,
	})
	s.action.pending = []screenplay.Performable{}

	return s.action
}

// Otherwise closes the pending To branch as the fallback, run when no branch matches.
func (a *ChoiceOnAction[T]) Otherwise() *ChoiceOnAction[T] {
	a.fallback = a.pending
	a.pending = []screenplay.Performable{}

	return a
}

// Default closes the pending To branch as the fallback, run when no branch matches.
func (a *ChoiceOnAction[T]) Default() *ChoiceOnAction[T] {
	return a.Otherwise()
}

// String describes the action.
func (a *ChoiceOnAction[T]) String() string {
	clauses := make([]string, 0, len(a.branches))
	for _, b := range a.branches {
		clauses = append(clauses, fmt.Sprintf(
			"to %s when it %s", performablesToString(b.performables), b.clause))
	}

	choices := strings.TrimPrefix(joinChoices(clauses, a.fallback), "choose ")

	return fmt.Sprintf("choose based on the %s %s", a.question, choices)
}

// PerformAs performs the task or the action as the provided actor.
func (a *ChoiceOnAction[T]) PerformAs(theActor *screenplay.Actor) error {
	value, err := a.question.AnsweredBy(theActor)
	if err != nil {
		return err
	}

	for _, b := range a.branches {
		verdict, judgeErr := b.clause.Judge("it", value)
		if judgeErr != nil {
			return judgeErr
		}

		if verdict.Met {
			return theActor.AttemptsTo(b.performables...)
		}
	}

	return theActor.AttemptsTo(a.fallback...)
}

// ChoiceOnAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*ChoiceOnAction[any])(nil)
