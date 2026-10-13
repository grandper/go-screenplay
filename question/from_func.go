package question

import (
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// AnsweredFn represents a function that answers a question as an actor.
type AnsweredFn[T any] func(theActor *screenplay.Actor) (T, error)

// FromFunc creates a new question from a function. The answer type is inferred
// from the function's return type.
func FromFunc[T any](description string, fn AnsweredFn[T]) screenplay.Question[T] {
	return &funcQuestion[T]{
		answeredBy:  fn,
		description: description,
	}
}

type funcQuestion[T any] struct {
	answeredBy  AnsweredFn[T]
	description string
}

// String describes the question.
func (q *funcQuestion[T]) String() string {
	return q.description
}

// AnsweredBy returns the answer that the provided actor gives to the question.
func (q *funcQuestion[T]) AnsweredBy(actor *screenplay.Actor) (T, error) {
	return q.answeredBy(actor)
}

// funcQuestion implements the screenplay.Question interface.
var _ screenplay.Question[any] = (*funcQuestion[any])(nil)
