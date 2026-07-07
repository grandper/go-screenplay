package standard

import (
	"fmt"

	"github.com/grandper/go-screenplay/extensions/cli/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

// OutputOf creates a question whose answer is the standard output of the result answered by the wrapped question.
func OutputOf(question screenplay.Question) *OutputOfQuestion {
	return &OutputOfQuestion{
		question: question,
	}
}

// OutputOfQuestion is a question whose answer is the standard output of the answer of the wrapped question.
type OutputOfQuestion struct {
	question screenplay.Question
}

// AnsweredBy returns the standard output of the result answered by the wrapped question.
func (q *OutputOfQuestion) AnsweredBy(actor *screenplay.Actor) (any, error) {
	answer, err := q.question.AnsweredBy(actor)
	if err != nil {
		return nil, err
	}

	result, ok := answer.(*ability.Result)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrAnswerIsNotAResult, q.question)
	}

	return result.StdOut(), nil
}

// String describes the question.
func (q *OutputOfQuestion) String() string {
	return "standard output of " + q.question.String()
}

// OutputOfQuestion implements the screenplay.Question interface.
var _ screenplay.Question = (*OutputOfQuestion)(nil)
