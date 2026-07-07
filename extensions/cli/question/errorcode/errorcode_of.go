package errorcode

import (
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/extensions/cli/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

// ErrAnswerIsNotAResult is returned when the wrapped question does not answer with a CLI result.
var ErrAnswerIsNotAResult = errors.New("the answer of the question is not a CLI result")

// Of creates a question whose answer is the error code of the result answered by the wrapped question.
func Of(question screenplay.Question) *OfQuestion {
	return &OfQuestion{
		question: question,
	}
}

// OfQuestion is a question whose answer is the error code of the answer of the wrapped question.
type OfQuestion struct {
	question screenplay.Question
}

// AnsweredBy returns the error code of the result answered by the wrapped question.
func (q *OfQuestion) AnsweredBy(actor *screenplay.Actor) (any, error) {
	answer, err := q.question.AnsweredBy(actor)
	if err != nil {
		return nil, err
	}

	result, ok := answer.(*ability.Result)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrAnswerIsNotAResult, q.question)
	}

	return result.ExitCode(), nil
}

// String describes the question.
func (q *OfQuestion) String() string {
	return "error code of " + q.question.String()
}

// OfQuestion implements the screenplay.Question interface.
var _ screenplay.Question = (*OfQuestion)(nil)
