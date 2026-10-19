package asynchronous

import (
	"github.com/grandper/go-screenplay/v2/internal/background"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// Errors asks for all the non-nil errors produced by the actions launched with
// proceed.AsynchronouslyTo. Answering the question waits for every pending
// asynchronous action to complete before returning the collected errors.
func Errors() screenplay.Question[[]error] {
	return &errorsQuestion{}
}

// errorsQuestion is a question whose answer is the slice of errors produced by
// the actions launched with proceed.AsynchronouslyTo.
type errorsQuestion struct{}

// AnsweredBy waits for every pending asynchronous action to complete and returns
// the slice of non-nil errors they produced.
func (q *errorsQuestion) AnsweredBy(actor *screenplay.Actor) ([]error, error) {
	return background.CollectedBy(actor), nil
}

// String describes the question.
func (q *errorsQuestion) String() string {
	return "asynchronous errors"
}

// errorsQuestion implements the screenplay.Question interface.
var _ screenplay.Question[[]error] = (*errorsQuestion)(nil)
