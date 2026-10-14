package number

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/grandper/go-screenplay/v2/screenplay"
)

// ErrAnswerIsNotACollection is returned when the wrapped question does not answer with a slice,
// an array or a map.
var ErrAnswerIsNotACollection = errors.New("the answer of the question is not a slice, an array or a map")

// Of creates a question whose answer is the number of elements in the slice,
// array or map answered by the wrapped question. A single type parameter cannot
// range over slices, arrays and maps, so the collection type is left
// unconstrained and checked when the question is answered.
func Of[C any](question screenplay.Question[C]) *OfQuestion[C] {
	return &OfQuestion[C]{
		question: question,
	}
}

// OfQuestion is a question whose answer is the number of elements in the answer of the wrapped question.
type OfQuestion[C any] struct {
	question screenplay.Question[C]
}

// AnsweredBy returns the number of elements in the slice, array or map answered by the wrapped question.
func (q *OfQuestion[C]) AnsweredBy(actor *screenplay.Actor) (int, error) {
	answer, err := q.question.AnsweredBy(actor)
	if err != nil {
		return 0, err
	}

	value := reflect.ValueOf(answer)
	if !value.IsValid() ||
		(value.Kind() != reflect.Slice && value.Kind() != reflect.Array && value.Kind() != reflect.Map) {
		return 0, fmt.Errorf("%w: %s", ErrAnswerIsNotACollection, q.question)
	}

	return value.Len(), nil
}

// String describes the question.
func (q *OfQuestion[C]) String() string {
	return "number of " + q.question.String()
}

// OfQuestion implements the screenplay.Question interface.
var _ screenplay.Question[int] = (*OfQuestion[any])(nil)
