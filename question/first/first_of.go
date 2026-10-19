package first

import (
	"fmt"

	"github.com/grandper/go-screenplay/v2/question"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// ErrAnswerIsEmpty is returned when the wrapped question answers with an empty
// slice. It is question.ErrAnswerIsEmpty, shared by first.Of and last.Of.
var ErrAnswerIsEmpty = question.ErrAnswerIsEmpty

// Of creates a question whose answer is the first element of the slice answered
// by the wrapped question. The wrapped question must answer with a slice; its
// element type is the answer type of the question.
func Of[S ~[]E, E any](question screenplay.Question[S]) screenplay.Question[E] {
	return &ofQuestion[S, E]{
		question: question,
	}
}

// ofQuestion is a question whose answer is the first element of the answer of the wrapped question.
type ofQuestion[S ~[]E, E any] struct {
	question screenplay.Question[S]
}

// AnsweredBy returns the first element of the slice answered by the wrapped question.
func (q *ofQuestion[S, E]) AnsweredBy(actor *screenplay.Actor) (E, error) {
	var zero E

	answer, err := q.question.AnsweredBy(actor)
	if err != nil {
		return zero, err
	}

	if len(answer) == 0 {
		return zero, fmt.Errorf("%w: %s", ErrAnswerIsEmpty, q.question)
	}

	return answer[0], nil
}

// String describes the question.
func (q *ofQuestion[S, E]) String() string {
	return "first of " + q.question.String()
}

// ofQuestion implements the screenplay.Question interface.
var _ screenplay.Question[any] = (*ofQuestion[[]any, any])(nil)
