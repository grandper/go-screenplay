package status

import (
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/v2/extensions/http/ability"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// ErrResponseIsNil is returned when the wrapped question answers with a nil HTTP response.
var ErrResponseIsNil = errors.New("the answer of the question is a nil HTTP response")

// CodeOf creates a question whose answer is the status code of the HTTP response answered by the wrapped question.
func CodeOf(question screenplay.Question[*ability.HTTPResponse]) *CodeOfQuestion {
	return &CodeOfQuestion{
		question: question,
	}
}

// CodeOfQuestion is a question whose answer is the status code of the answer of the wrapped question.
type CodeOfQuestion struct {
	question screenplay.Question[*ability.HTTPResponse]
}

// AnsweredBy returns the status code of the HTTP response answered by the wrapped question.
func (q *CodeOfQuestion) AnsweredBy(actor *screenplay.Actor) (int, error) {
	response, err := q.question.AnsweredBy(actor)
	if err != nil {
		return 0, err
	}

	if response == nil {
		return 0, fmt.Errorf("%w: %s", ErrResponseIsNil, q.question)
	}

	return response.StatusCode(), nil
}

// String describes the question.
func (q *CodeOfQuestion) String() string {
	return "status code of " + q.question.String()
}

// CodeOfQuestion implements the screenplay.Question interface.
var _ screenplay.Question[int] = (*CodeOfQuestion)(nil)
