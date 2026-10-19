package standard

import (
	"fmt"

	"github.com/grandper/go-screenplay/v2/extensions/cli/ability"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// ErrorOf creates a question whose answer is the standard error of the result answered by the wrapped question.
func ErrorOf(question screenplay.Question[*ability.Result]) screenplay.Question[[]byte] {
	return &errorOfQuestion{
		question: question,
	}
}

// errorOfQuestion is a question whose answer is the standard error of the answer of the wrapped question.
type errorOfQuestion struct {
	question screenplay.Question[*ability.Result]
}

// AnsweredBy returns the standard error of the result answered by the wrapped question.
func (q *errorOfQuestion) AnsweredBy(actor *screenplay.Actor) ([]byte, error) {
	result, err := q.question.AnsweredBy(actor)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("%w: %s", ErrResultIsNil, q.question)
	}

	return result.StdErr(), nil
}

// String describes the question.
func (q *errorOfQuestion) String() string {
	return "standard error of " + q.question.String()
}

// errorOfQuestion implements the screenplay.Question interface.
var _ screenplay.Question[[]byte] = (*errorOfQuestion)(nil)
