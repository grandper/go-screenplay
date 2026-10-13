package header

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/grandper/go-screenplay/v2/screenplay"
)

// ErrAnswerIsNotAStruct is returned when the wrapped question does not answer with a struct.
var ErrAnswerIsNotAStruct = errors.New("the answer of the question is not a struct")

// ErrFieldNotFound is returned when the struct answered by the wrapped question has no Header field.
var ErrFieldNotFound = errors.New("the answer of the question has no Header field")

// ErrFieldHasWrongType is returned when the Header field of the struct answered by the wrapped question
// is not of the expected type.
var ErrFieldHasWrongType = errors.New("the Header field of the answer of the question has the wrong type")

// Of creates a question whose answer is the Header field of the struct answered
// by the wrapped question. Go cannot constrain a type parameter on the fields it
// has, so the type of the field, F, is given explicitly while the answer type of
// the wrapped question is inferred:
//
//	header.Of[string](TheLastResponse)
//
// Use header.Of[any](...) to accept a field of any type.
func Of[F any, T any](question screenplay.Question[T]) *OfQuestion[F, T] {
	return &OfQuestion[F, T]{
		question: question,
	}
}

// OfQuestion is a question whose answer is the Header field of the answer of the wrapped question.
type OfQuestion[F any, T any] struct {
	question screenplay.Question[T]
}

// AnsweredBy returns the Header field of the struct answered by the wrapped question.
func (q *OfQuestion[F, T]) AnsweredBy(actor *screenplay.Actor) (F, error) {
	var zero F

	answer, err := q.question.AnsweredBy(actor)
	if err != nil {
		return zero, err
	}

	value := reflect.ValueOf(answer)
	if value.IsValid() && value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	if !value.IsValid() || value.Kind() != reflect.Struct {
		return zero, fmt.Errorf("%w: %s", ErrAnswerIsNotAStruct, q.question)
	}

	field := value.FieldByName("Header")
	if !field.IsValid() || !field.CanInterface() {
		return zero, fmt.Errorf("%w: %s", ErrFieldNotFound, q.question)
	}

	typed, isAnF := field.Interface().(F)
	if !isAnF {
		return zero, fmt.Errorf("%w: %s holds a %s, not a %T", ErrFieldHasWrongType, q.question, field.Type(), zero)
	}

	return typed, nil
}

// String describes the question.
func (q *OfQuestion[F, T]) String() string {
	return "header of " + q.question.String()
}

// OfQuestion implements the screenplay.Question interface.
var _ screenplay.Question[any] = (*OfQuestion[any, any])(nil)
