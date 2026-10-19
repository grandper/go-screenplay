package header

import (
	"fmt"
	"reflect"

	"github.com/grandper/go-screenplay/v2/question"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

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
		return zero, fmt.Errorf("%w: %s", question.ErrAnswerIsNotAStruct, q.question)
	}

	field := value.FieldByName("Header")
	if !field.IsValid() || !field.CanInterface() {
		return zero, fmt.Errorf("%w: Header: %s", question.ErrFieldNotFound, q.question)
	}

	typed, isAnF := field.Interface().(F)
	if !isAnF {
		return zero, fmt.Errorf(
			"%w: Header: %s holds a %s, not a %T",
			question.ErrFieldHasWrongType,
			q.question,
			field.Type(),
			zero,
		)
	}

	return typed, nil
}

// String describes the question.
func (q *OfQuestion[F, T]) String() string {
	return "header of " + q.question.String()
}

// OfQuestion implements the screenplay.Question interface.
var _ screenplay.Question[any] = (*OfQuestion[any, any])(nil)
