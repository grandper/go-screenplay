package question

import (
	"github.com/grandper/go-screenplay/screenplay"
)

// AnsweredFn represents a function that answers a question as an actor.
type AnsweredFn func(theActor *screenplay.Actor) (any, error)

// FromFunc creates a new question from a function.
func FromFunc(description string, fn AnsweredFn) screenplay.Question {
	return &funcQuestion{
		answeredBy:  fn,
		description: description,
	}
}

type funcQuestion struct {
	answeredBy  AnsweredFn
	description string
}

// String describes the question.
func (q *funcQuestion) String() string {
	return q.description
}

// AnsweredBy returns the answer that the provided actor gives to the question.
func (q *funcQuestion) AnsweredBy(actor *screenplay.Actor) (any, error) {
	return q.answeredBy(actor)
}

// funcQuestion implements the screenplay.Question interface.
var _ screenplay.Question = (*funcQuestion)(nil)
