package choose

import (
	"fmt"
	"strings"

	"github.com/grandper/go-screenplay/v2/action/internal/condition"
	"github.com/grandper/go-screenplay/v2/action/internal/performables"
	"github.com/grandper/go-screenplay/v2/internal/sentence"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// ordinal renders 1 -> "first", 2 -> "second", ... for the boolean-branch
// description; from the sixth on, the rank is written with digits: "6th", "21st".
func ordinal(position int) string {
	names := []string{"first", "second", "third", "fourth", "fifth"}
	if position >= 1 && position <= len(names) {
		return names[position-1]
	}

	return sentence.Rank(position)
}

// joinChoices assembles "choose to A when x, to B when y, or to Z otherwise".
// Each clause is already in the infinitive form "to <action> when <condition>"; the
// fallback becomes the last alternative "to <action> <closing>", where closing is the
// word written ("otherwise" or "default"), and the list is closed with "or" so it reads
// as picking exactly one branch (not a to-do list).
func joinChoices(clauses []string, fallback []screenplay.Performable, closing string) string {
	alternatives := make([]string, len(clauses), len(clauses)+1)
	copy(alternatives, clauses)
	if len(fallback) > 0 {
		alternatives = append(alternatives, "to "+performables.Describe(fallback)+" "+closing)
	}

	if len(alternatives) <= 1 {
		return "choose " + strings.Join(alternatives, "")
	}

	last := len(alternatives) - 1

	return "choose " + strings.Join(alternatives[:last], ", ") + ", or " + alternatives[last]
}

// choiceBranch pairs a condition with the performables run when it is the first to hold.
// condition is the same type proceed.To uses, so When reuses its boolean logic and
// WhenThe reuses its question/resolution logic.
type choiceBranch struct {
	condition    condition.Condition
	performables []screenplay.Performable
}

// To starts a multi-way choice with the performables of its first branch; close
// the branch with When or WhenThe, and name the next one with To:
//
//	choose.To(ReadTheResponse()).WhenThe(statusCode).HappensTo(be.EqualTo(200)).
//		To(AbortTheRequest()).Otherwise()
func To(actions ...screenplay.Performable) *ToAction {
	return &ToAction{
		branches: []choiceBranch{},
		pending:  actions,
		fallback: []screenplay.Performable{},
		closing:  "otherwise",
	}
}

// ToAction performs the performables of the first branch whose condition holds.
type ToAction struct {
	branches []choiceBranch
	pending  []screenplay.Performable
	fallback []screenplay.Performable
	closing  string // the word that closed the fallback: "otherwise" or "default"
}

// To names the performables of the next branch; close it with When or WhenThe.
func (a *ToAction) To(actions ...screenplay.Performable) *ToAction {
	a.pending = actions

	return a
}

// When closes the pending branch with a boolean condition.
func (a *ToAction) When(cond bool) *ToAction {
	branchCondition := condition.Boolean(cond)
	branchCondition.Description = fmt.Sprintf("the %s condition holds", ordinal(len(a.branches)+1))

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
//	choose.To(ReadTheResponse()).WhenThe(statusCode).HappensTo(be.EqualTo(200))
func (a *ToAction) WhenThe[T any](question screenplay.Question[T]) *ToBuilder[T] {
	return &ToBuilder[T]{action: a, question: question}
}

// ToBuilder builds the condition of a branch of a choice: it holds the
// question the branch depends on, and waits for the verb that closes it.
type ToBuilder[T any] struct {
	action   *ToAction
	question screenplay.Question[T]
}

// HappensTo closes the branch with the answer satisfying every resolution. It reads "happens to ...".
func (s *ToBuilder[T]) HappensTo(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.closes(sentence.NewCondition(sentence.Should, false, resolutions))
}

// HappensToBe closes the branch with the answer satisfying every resolution. It reads "happens to be ...".
func (s *ToBuilder[T]) HappensToBe(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.closes(sentence.NewCondition(sentence.ShouldBe, false, resolutions))
}

// HappensToHave closes the branch with the answer satisfying every resolution. It reads "happens to have ...".
func (s *ToBuilder[T]) HappensToHave(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.closes(sentence.NewCondition(sentence.ShouldHave, false, resolutions))
}

// DoesNotHappenTo closes the branch with the answer satisfying none of the resolutions. It reads "does not happen to ...".
func (s *ToBuilder[T]) DoesNotHappenTo(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.closes(sentence.NewCondition(sentence.ShouldNot, false, resolutions))
}

// DoesNotHappenToBe closes the branch with the answer satisfying none of the resolutions. It reads "does not happen to be ...".
func (s *ToBuilder[T]) DoesNotHappenToBe(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.closes(sentence.NewCondition(sentence.ShouldNotBe, false, resolutions))
}

// DoesNotHappenToHave closes the branch with the answer satisfying none of the resolutions. It reads "does not happen to have ...".
func (s *ToBuilder[T]) DoesNotHappenToHave(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.closes(sentence.NewCondition(sentence.ShouldNotHave, false, resolutions))
}

// HappensToBeAnyOf closes the branch with the answer satisfying at least one resolution. It reads "happens to be any of: ...".
func (s *ToBuilder[T]) HappensToBeAnyOf(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.closes(sentence.NewCondition(sentence.ShouldAnyOf, false, resolutions))
}

// HappenTo closes the branch with the answer satisfying every resolution. It is the plural form: it reads "happen to ...".
func (s *ToBuilder[T]) HappenTo(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.closes(sentence.NewCondition(sentence.Should, true, resolutions))
}

// HappenToBe closes the branch with the answer satisfying every resolution. It is the plural form: it reads "happen to be ...".
func (s *ToBuilder[T]) HappenToBe(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.closes(sentence.NewCondition(sentence.ShouldBe, true, resolutions))
}

// HappenToHave closes the branch with the answer satisfying every resolution. It is the plural form: it reads "happen to have ...".
func (s *ToBuilder[T]) HappenToHave(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.closes(sentence.NewCondition(sentence.ShouldHave, true, resolutions))
}

// DoNotHappenTo closes the branch with the answer satisfying none of the resolutions. It is the plural form: it reads "do not happen to ...".
func (s *ToBuilder[T]) DoNotHappenTo(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.closes(sentence.NewCondition(sentence.ShouldNot, true, resolutions))
}

// DoNotHappenToBe closes the branch with the answer satisfying none of the resolutions. It is the plural form: it reads "do not happen to be ...".
func (s *ToBuilder[T]) DoNotHappenToBe(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.closes(sentence.NewCondition(sentence.ShouldNotBe, true, resolutions))
}

// DoNotHappenToHave closes the branch with the answer satisfying none of the resolutions. It is the plural form: it reads "do not happen to have ...".
func (s *ToBuilder[T]) DoNotHappenToHave(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.closes(sentence.NewCondition(sentence.ShouldNotHave, true, resolutions))
}

// HappenToBeAnyOf closes the branch with the answer satisfying at least one resolution. It is the plural form: it reads "happen to be any of: ...".
func (s *ToBuilder[T]) HappenToBeAnyOf(resolutions ...screenplay.Resolution[T]) *ToAction {
	return s.closes(sentence.NewCondition(sentence.ShouldAnyOf, true, resolutions))
}

// closes closes the pending branch: the answer satisfies the clause.
func (s *ToBuilder[T]) closes(clause sentence.Clause[T]) *ToAction {
	s.action.branches = append(s.action.branches, choiceBranch{
		condition:    condition.Term(sentence.NewTerm(s.question, clause)),
		performables: s.action.pending,
	})
	s.action.pending = []screenplay.Performable{}

	return s.action
}

// Otherwise closes the pending To branch as the fallback, run when no branch holds.
func (a *ToAction) Otherwise() *ToAction {
	return a.closedBy("otherwise")
}

// Default closes the pending To branch as the fallback, run when no branch holds.
// It is narrated "default".
func (a *ToAction) Default() *ToAction {
	return a.closedBy("default")
}

// closedBy closes the pending To branch as the fallback, with the word written.
func (a *ToAction) closedBy(closing string) *ToAction {
	a.fallback = a.pending
	a.pending = []screenplay.Performable{}
	a.closing = closing

	return a
}

// String describes the action.
func (a *ToAction) String() string {
	clauses := make([]string, 0, len(a.branches))
	for _, b := range a.branches {
		clauses = append(clauses, fmt.Sprintf(
			"to %s when %s", performables.Describe(b.performables), b.condition.Description))
	}

	return joinChoices(clauses, a.fallback, a.closing)
}

// PerformAs performs the task or the action as the provided actor.
func (a *ToAction) PerformAs(theActor *screenplay.Actor) error {
	for _, b := range a.branches {
		holds, err := b.condition.HoldsFor(theActor)
		if err != nil {
			return err
		}

		if holds {
			return theActor.AttemptsTo(b.performables...)
		}
	}

	return theActor.AttemptsTo(a.fallback...)
}

// ToAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*ToAction)(nil)
