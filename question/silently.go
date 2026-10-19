package question

import (
	"errors"

	"github.com/grandper/go-screenplay/v2/screenplay"
)

// Silently wraps a question so that the narration produced while it is answered
// is muted: any beats or asides emitted while computing the answer never reach
// the narrator's adapters. The question is narrated as "the secret": neither the
// wrapped question's description nor its answer is narrated.
//
// It is meant to keep secrets out of the logs — a question that reads a token or
// a password can be answered without its answer or its intermediate questions
// reaching an adapter — and to trim the noise of a chatty question down to
// nothing.
//
// A production configured with screenplay.WithForceAllNarration neutralizes
// this decorator, so a scenario can be debugged with every step narrated again.
func Silently[T any](q screenplay.Question[T]) screenplay.Question[T] {
	return &silentlyQuestion[T]{
		question: q,
	}
}

// silentlyQuestion is a question that mutes the narration produced while the
// wrapped question is answered.
type silentlyQuestion[T any] struct {
	question screenplay.Question[T]
}

// String describes the question without describing the wrapped question, so
// that a value the wrapped question carries is never narrated.
func (q *silentlyQuestion[T]) String() string {
	return "secret"
}

// KeepsItsAnswerSecret marks the question as one whose answer must not be
// narrated. The actor checks for this method before revealing an answer.
func (q *silentlyQuestion[T]) KeepsItsAnswerSecret() {}

// AnsweredBy answers the wrapped question through a muted view of the actor, so
// the narration it would produce is dropped. The view leaves the actor itself
// untouched, so a Silently running in one concurrent branch never mutes
// another. When the actor's narrator cannot be muted (see
// screenplay.WithForceAllNarration), the steps narrate as usual.
func (q *silentlyQuestion[T]) AnsweredBy(actor *screenplay.Actor) (T, error) {
	if q.question == nil {
		var zero T

		return zero, errors.New("the question given to Silently is nil")
	}

	return q.question.AnsweredBy(actor.Muted())
}

// silentlyQuestion implements the screenplay.Question interface.
var _ screenplay.Question[any] = (*silentlyQuestion[any])(nil)
