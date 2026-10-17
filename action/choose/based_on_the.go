package choose

import (
	"fmt"
	"strings"

	"github.com/grandper/go-screenplay/v2/internal/performables"
	"github.com/grandper/go-screenplay/v2/internal/sentence"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// clauseBranch pairs a clause with the performables run when the answer is the first to satisfy it.
type clauseBranch[T any] struct {
	clause       sentence.Clause[T]
	performables []screenplay.Performable
}

// BasedOnThe starts a multi-way choice on the answer of a question; each branch states what
// the answer happens to be:
//
//	choose.BasedOnThe(statusCode).
//		To(ReadTheResponse()).WhenIt().HappensTo(be.EqualTo(200)).
//		To(AbortTheRequest()).Otherwise()
func BasedOnThe[T any](question screenplay.Question[T]) *BasedOnTheAction[T] {
	return &BasedOnTheAction[T]{
		question: question,
		branches: []clauseBranch[T]{},
		pending:  []screenplay.Performable{},
		fallback: []screenplay.Performable{},
		closing:  "otherwise",
	}
}

// BasedOnTheAction performs the performables of the first branch whose clause the answer satisfies.
type BasedOnTheAction[T any] struct {
	question screenplay.Question[T]
	branches []clauseBranch[T]
	pending  []screenplay.Performable
	fallback []screenplay.Performable
	closing  string // the word that closed the fallback: "otherwise" or "default"
}

// To names the performables of the next branch; close it with WhenIt.
func (a *BasedOnTheAction[T]) To(actions ...screenplay.Performable) *BasedOnTheAction[T] {
	a.pending = actions

	return a
}

// WhenIt closes the pending branch with what the answer happens to be. Close
// it with a verb and one or more resolutions:
//
//	choose.BasedOnThe(statusCode).To(ReadTheResponse()).WhenIt().HappensTo(be.EqualTo(200))
func (a *BasedOnTheAction[T]) WhenIt() *BasedOnTheBuilder[T] {
	return &BasedOnTheBuilder[T]{action: a}
}

// BasedOnTheBuilder builds the condition of a branch of a choice on the answer
// to a question: it waits for the verb that closes it.
type BasedOnTheBuilder[T any] struct {
	action *BasedOnTheAction[T]
}

// HappensTo closes the branch with the answer satisfying every resolution. It reads "happens to ...".
func (s *BasedOnTheBuilder[T]) HappensTo(resolutions ...screenplay.Resolution[T]) *BasedOnTheAction[T] {
	return s.closes(sentence.NewCondition(sentence.Should, false, resolutions))
}

// HappensToBe closes the branch with the answer satisfying every resolution. It reads "happens to be ...".
func (s *BasedOnTheBuilder[T]) HappensToBe(resolutions ...screenplay.Resolution[T]) *BasedOnTheAction[T] {
	return s.closes(sentence.NewCondition(sentence.ShouldBe, false, resolutions))
}

// HappensToHave closes the branch with the answer satisfying every resolution. It reads "happens to have ...".
func (s *BasedOnTheBuilder[T]) HappensToHave(resolutions ...screenplay.Resolution[T]) *BasedOnTheAction[T] {
	return s.closes(sentence.NewCondition(sentence.ShouldHave, false, resolutions))
}

// DoesNotHappenTo closes the branch with the answer satisfying none of the resolutions. It reads "does not happen to ...".
func (s *BasedOnTheBuilder[T]) DoesNotHappenTo(resolutions ...screenplay.Resolution[T]) *BasedOnTheAction[T] {
	return s.closes(sentence.NewCondition(sentence.ShouldNot, false, resolutions))
}

// DoesNotHappenToBe closes the branch with the answer satisfying none of the resolutions. It reads "does not happen to be ...".
func (s *BasedOnTheBuilder[T]) DoesNotHappenToBe(resolutions ...screenplay.Resolution[T]) *BasedOnTheAction[T] {
	return s.closes(sentence.NewCondition(sentence.ShouldNotBe, false, resolutions))
}

// DoesNotHappenToHave closes the branch with the answer satisfying none of the resolutions. It reads "does not happen to have ...".
func (s *BasedOnTheBuilder[T]) DoesNotHappenToHave(resolutions ...screenplay.Resolution[T]) *BasedOnTheAction[T] {
	return s.closes(sentence.NewCondition(sentence.ShouldNotHave, false, resolutions))
}

// HappensToBeAnyOf closes the branch with the answer satisfying at least one resolution. It reads "happens to be any of: ...".
func (s *BasedOnTheBuilder[T]) HappensToBeAnyOf(resolutions ...screenplay.Resolution[T]) *BasedOnTheAction[T] {
	return s.closes(sentence.NewCondition(sentence.ShouldAnyOf, false, resolutions))
}

// closes closes the pending branch: the answer satisfies the clause.
func (s *BasedOnTheBuilder[T]) closes(clause sentence.Clause[T]) *BasedOnTheAction[T] {
	s.action.branches = append(s.action.branches, clauseBranch[T]{
		clause:       clause,
		performables: s.action.pending,
	})
	s.action.pending = []screenplay.Performable{}

	return s.action
}

// Otherwise closes the pending To branch as the fallback, run when no branch matches.
func (a *BasedOnTheAction[T]) Otherwise() *BasedOnTheAction[T] {
	return a.closedBy("otherwise")
}

// Default closes the pending To branch as the fallback, run when no branch matches.
// It is narrated "default".
func (a *BasedOnTheAction[T]) Default() *BasedOnTheAction[T] {
	return a.closedBy("default")
}

// closedBy closes the pending To branch as the fallback, with the word written.
func (a *BasedOnTheAction[T]) closedBy(closing string) *BasedOnTheAction[T] {
	a.fallback = a.pending
	a.pending = []screenplay.Performable{}
	a.closing = closing

	return a
}

// String describes the action.
func (a *BasedOnTheAction[T]) String() string {
	clauses := make([]string, 0, len(a.branches))
	for _, b := range a.branches {
		clauses = append(clauses, fmt.Sprintf(
			"to %s when it %s", performables.Describe(b.performables), b.clause))
	}

	choices := strings.TrimPrefix(joinChoices(clauses, a.fallback, a.closing), "choose ")

	return fmt.Sprintf("choose based on the %s %s", a.question, choices)
}

// PerformAs performs the task or the action as the provided actor.
func (a *BasedOnTheAction[T]) PerformAs(theActor *screenplay.Actor) error {
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

// BasedOnTheAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*BasedOnTheAction[any])(nil)
