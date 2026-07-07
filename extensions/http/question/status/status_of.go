package status

import (
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/extensions/http/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

// ErrAnswerIsNotAnHTTPResponse is returned when the answer of the wrapped question is not an HTTP response.
var ErrAnswerIsNotAnHTTPResponse = errors.New("the answer of the question is not an HTTP response")

// CodeOf creates a question whose answer is the status code of the HTTP response answered by the wrapped question.
func CodeOf(question screenplay.Question) *CodeOfQuestion {
	return &CodeOfQuestion{
		question: question,
	}
}

// CodeOfQuestion is a question whose answer is the status code of the answer of the wrapped question.
type CodeOfQuestion struct {
	question screenplay.Question
}

// AnsweredBy returns the status code of the HTTP response answered by the wrapped question.
func (q *CodeOfQuestion) AnsweredBy(actor *screenplay.Actor) (any, error) {
	answer, err := q.question.AnsweredBy(actor)
	if err != nil {
		return nil, err
	}

	response, ok := answer.(*ability.HTTPResponse)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrAnswerIsNotAnHTTPResponse, q.question)
	}

	return response.StatusCode(), nil
}

// String describes the question.
func (q *CodeOfQuestion) String() string {
	return "status code of " + q.question.String()
}

// CodeOfQuestion implements the screenplay.Question interface.
var _ screenplay.Question = (*CodeOfQuestion)(nil)
