// Package sentence holds what the sentences of the library share: the verbs,
// the clause a verb builds with its resolutions, and the term stating what the
// answer to a question should be.
package sentence

import (
	"fmt"

	"github.com/grandper/go-screenplay/v2/screenplay"
)

// Term is a statement about the answer to a question: "the status code should
// be equal to 200". It is met when the answer satisfies its clause.
type Term struct {
	description string
	// mistake is why the term can never be judged, when its sentence was
	// built wrongly.
	mistake error
	// judge answers the question and judges its answer.
	judge func(actor *screenplay.Actor) (Verdict, error)
}

// NewTerm creates the term stating that the answer to a question satisfies a
// clause. A term about a nil question can never be judged: it fails with
// ErrInvalidTerm each time it is judged.
func NewTerm[T any](question screenplay.Question[T], clause Clause[T]) *Term {
	if question == nil {
		mistake := fmt.Errorf("%w: the question is nil", ErrInvalidTerm)

		return &Term{
			description: fmt.Sprintf("the <nil> %s", clause),
			mistake:     mistake,
			judge:       func(*screenplay.Actor) (Verdict, error) { return Verdict{}, mistake },
		}
	}

	subject := fmt.Sprintf("the %s", question)

	return &Term{
		description: fmt.Sprintf("%s %s", subject, clause),
		mistake:     clause.err,
		judge: func(actor *screenplay.Actor) (Verdict, error) {
			if clause.err != nil {
				return Verdict{}, clause.err
			}

			answer, err := question.AnsweredBy(actor)
			if err != nil {
				return Verdict{}, err
			}

			return clause.Judge(subject, answer)
		},
	}
}

// String describes the term: "the status code should be equal to 200".
func (t *Term) String() string {
	return t.description
}

// Mistake returns why the term can never be judged, when its sentence was
// built wrongly: a verb without a resolution, a nil resolution, or a nil
// question. It returns nil for a term that was built correctly.
func (t *Term) Mistake() error {
	return t.mistake
}

// JudgedBy answers the question as an actor and judges the answer. It fails
// when the question cannot be answered, or when the sentence was built
// wrongly.
func (t *Term) JudgedBy(actor *screenplay.Actor) (Verdict, error) {
	return t.judge(actor)
}

// MetBy returns nil when the term is met for an actor. When it is not, it
// returns an *UnmetError explaining why; when the question cannot be answered
// or the sentence was built wrongly, it returns that error.
func (t *Term) MetBy(actor *screenplay.Actor) error {
	verdict, err := t.judge(actor)
	if err != nil {
		return err
	}

	if !verdict.Met {
		return &UnmetError{
			Failure: verdict.Failure,
			Reason:  verdict.Reason,
		}
	}

	return nil
}

// UnmetError is returned when a term is not met.
type UnmetError struct {
	// Failure describes the part of the sentence that is not satisfied and
	// why: "the status code should be less than 300: got 404".
	Failure string
	// Reason is the error of the resolution that is not satisfied, when there
	// is one.
	Reason error
}

// Error describes the part of the sentence that is not satisfied and why.
func (e *UnmetError) Error() string {
	return e.Failure
}

// Unwrap returns the error of the resolution that is not satisfied.
func (e *UnmetError) Unwrap() error {
	return e.Reason
}
