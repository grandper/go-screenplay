package last

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

// Of creates a question whose answer is the last element of the slice answered by the wrapped question.
func Of(question screenplay.Question) *OfQuestion {
	return &OfQuestion{
		question: question,
	}
}

// OfQuestion is a question whose answer is the last element of the answer of the wrapped question.
type OfQuestion struct {
	question screenplay.Question
}

// AnsweredBy returns the last element of the slice answered by the wrapped question.
func (q *OfQuestion) AnsweredBy(actor *screenplay.Actor) (any, error) {
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

	return value.Index(value.Len() - 1).Interface(), nil
}

// String describes the question.
func (q *OfQuestion) String() string {
	return "last of " + q.question.String()
}

// OfQuestion implements the screenplay.Question interface.
var _ screenplay.Question = (*OfQuestion)(nil)
