package errorcode

import (
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/v2/extensions/cli/ability"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// ErrAnswerIsNil is returned when the wrapped question answers with a nil CLI result.
var ErrAnswerIsNil = errors.New("the answer of the question is a nil CLI result")

// Of creates a question whose answer is the error code of the result answered by the wrapped question.
func Of(question screenplay.Question[*ability.Result]) *OfQuestion {
	return &OfQuestion{
		question: question,
	}
}

// OfQuestion is a question whose answer is the error code of the answer of the wrapped question.
type OfQuestion struct {
	question screenplay.Question[*ability.Result]
}

// AnsweredBy returns the error code of the result answered by the wrapped question.
func (q *OfQuestion) AnsweredBy(actor *screenplay.Actor) (int, error) {
	result, err := q.question.AnsweredBy(actor)
	if err != nil {
		return 0, err
	}

	if result == nil {
		return 0, fmt.Errorf("%w: %s", ErrAnswerIsNil, q.question)
	}

	return result.ExitCode(), nil
}

// String describes the question.
func (q *OfQuestion) String() string {
	return "error code of " + q.question.String()
}

// OfQuestion implements the screenplay.Question interface.
var _ screenplay.Question[int] = (*OfQuestion)(nil)
