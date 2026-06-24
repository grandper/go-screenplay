package first

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/grandper/go-screenplay/screenplay"
)

// ErrAnswerIsNotASlice is returned when the wrapped question does not answer with a slice.
var ErrAnswerIsNotASlice = errors.New("the answer of the question is not a slice")

// ErrAnswerIsEmpty is returned when the wrapped question answers with an empty slice.
var ErrAnswerIsEmpty = errors.New("the answer of the question is an empty slice")

// Of creates a question whose answer is the first element of the slice answered by the wrapped question.
func Of(question screenplay.Question) *FirstOfQuestion {
	return &FirstOfQuestion{
		question: question,
	}
}

// FirstOfQuestion is a question whose answer is the first element of the answer of the wrapped question.
type FirstOfQuestion struct {
	question screenplay.Question
}

// AnsweredBy returns the first element of the slice answered by the wrapped question.
func (q *FirstOfQuestion) AnsweredBy(actor *screenplay.Actor) (any, error) {
	answer, err := q.question.AnsweredBy(actor)
	if err != nil {
		return nil, err
	}

	value := reflect.ValueOf(answer)
	if !value.IsValid() || (value.Kind() != reflect.Slice && value.Kind() != reflect.Array) {
		return nil, fmt.Errorf("%w: %s", ErrAnswerIsNotASlice, q.question)
	}

	if value.Len() == 0 {
		return nil, fmt.Errorf("%w: %s", ErrAnswerIsEmpty, q.question)
	}

	return value.Index(0).Interface(), nil
}

// String describes the question.
func (q *FirstOfQuestion) String() string {
	return "first of " + q.question.String()
}

// FirstOfQuestion implements the screenplay.Question interface.
var _ screenplay.Question = (*FirstOfQuestion)(nil)
