package standard

import (
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/extensions/cli/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

// ErrAnswerIsNotAResult is returned when the wrapped question does not answer with a CLI result.
var ErrAnswerIsNotAResult = errors.New("the answer of the question is not a CLI result")

// ErrorOf creates a question whose answer is the standard error of the result answered by the wrapped question.
func ErrorOf(question screenplay.Question) *ErrorOfQuestion {
	return &ErrorOfQuestion{
		question: question,
	}
}

// ErrorOfQuestion is a question whose answer is the standard error of the answer of the wrapped question.
type ErrorOfQuestion struct {
	question screenplay.Question
}

// AnsweredBy returns the standard error of the result answered by the wrapped question.
func (q *ErrorOfQuestion) AnsweredBy(actor *screenplay.Actor) (any, error) {
	answer, err := q.question.AnsweredBy(actor)
	if err != nil {
		return nil, err
	}

	result, ok := answer.(*ability.Result)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrAnswerIsNotAResult, q.question)
	}

	return result.StdErr(), nil
}

// String describes the question.
func (q *ErrorOfQuestion) String() string {
	return "standard error of " + q.question.String()
}

// ErrorOfQuestion implements the screenplay.Question interface.
var _ screenplay.Question = (*ErrorOfQuestion)(nil)
