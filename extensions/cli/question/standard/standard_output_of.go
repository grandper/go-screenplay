package standard

import (
	"fmt"

	"github.com/grandper/go-screenplay/v2/extensions/cli/ability"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// OutputOf creates a question whose answer is the standard output of the result answered by the wrapped question.
func OutputOf(question screenplay.Question[*ability.Result]) *OutputOfQuestion {
	return &OutputOfQuestion{
		question: question,
	}
}

// OutputOfQuestion is a question whose answer is the standard output of the answer of the wrapped question.
type OutputOfQuestion struct {
	question screenplay.Question[*ability.Result]
}

// AnsweredBy returns the standard output of the result answered by the wrapped question.
func (q *OutputOfQuestion) AnsweredBy(actor *screenplay.Actor) ([]byte, error) {
	result, err := q.question.AnsweredBy(actor)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("%w: %s", ErrAnswerIsNil, q.question)
	}

	return result.StdOut(), nil
}

// String describes the question.
func (q *OutputOfQuestion) String() string {
	return "standard output of " + q.question.String()
}

// OutputOfQuestion implements the screenplay.Question interface.
var _ screenplay.Question[[]byte] = (*OutputOfQuestion)(nil)
