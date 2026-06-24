package number

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/grandper/go-screenplay/screenplay"
)

// ErrAnswerIsNotACollection is returned when the wrapped question does not answer with a slice or a map.
var ErrAnswerIsNotACollection = errors.New("the answer of the question is not a slice or a map")

// Of creates a question whose answer is the number of elements in the slice or map answered by the wrapped question.
func Of(question screenplay.Question) *OfQuestion {
	return &OfQuestion{
		question: question,
	}
}

// OfQuestion is a question whose answer is the number of elements in the answer of the wrapped question.
type OfQuestion struct {
	question screenplay.Question
}

// AnsweredBy returns the number of elements in the slice or map answered by the wrapped question.
func (q *OfQuestion) AnsweredBy(actor *screenplay.Actor) (any, error) {
	answer, err := q.question.AnsweredBy(actor)
	if err != nil {
		return nil, err
	}

	value := reflect.ValueOf(answer)
	if !value.IsValid() ||
		(value.Kind() != reflect.Slice && value.Kind() != reflect.Array && value.Kind() != reflect.Map) {
		return nil, fmt.Errorf("%w: %s", ErrAnswerIsNotACollection, q.question)
	}

	return value.Len(), nil
}

// String describes the question.
func (q *OfQuestion) String() string {
	return "number of " + q.question.String()
}

// OfQuestion implements the screenplay.Question interface.
var _ screenplay.Question = (*OfQuestion)(nil)
