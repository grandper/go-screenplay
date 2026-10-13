package standard

import (
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/v2/extensions/cli/ability"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// ErrAnswerIsNil is returned when the wrapped question answers with a nil CLI result.
var ErrAnswerIsNil = errors.New("the answer of the question is a nil CLI result")

// ErrorOf creates a question whose answer is the standard error of the result answered by the wrapped question.
func ErrorOf(question screenplay.Question[*ability.Result]) *ErrorOfQuestion {
	return &ErrorOfQuestion{
		question: question,
	}
}

// ErrorOfQuestion is a question whose answer is the standard error of the answer of the wrapped question.
type ErrorOfQuestion struct {
	question screenplay.Question[*ability.Result]
}

// AnsweredBy returns the standard error of the result answered by the wrapped question.
func (q *ErrorOfQuestion) AnsweredBy(actor *screenplay.Actor) ([]byte, error) {
	result, err := q.question.AnsweredBy(actor)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("%w: %s", ErrAnswerIsNil, q.question)
	}

	return result.StdErr(), nil
}

// String describes the question.
func (q *ErrorOfQuestion) String() string {
	return "standard error of " + q.question.String()
}

// ErrorOfQuestion implements the screenplay.Question interface.
var _ screenplay.Question[[]byte] = (*ErrorOfQuestion)(nil)
