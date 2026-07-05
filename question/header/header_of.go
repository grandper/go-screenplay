package header

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/grandper/go-screenplay/screenplay"
)

// ErrAnswerIsNotAStruct is returned when the wrapped question does not answer with a struct.
var ErrAnswerIsNotAStruct = errors.New("the answer of the question is not a struct")

// ErrFieldNotFound is returned when the struct answered by the wrapped question has no Header field.
var ErrFieldNotFound = errors.New("the answer of the question has no Header field")

// Of creates a question whose answer is the Header field of the struct answered by the wrapped question.
func Of(question screenplay.Question) *OfQuestion {
	return &OfQuestion{
		question: question,
	}
}

// OfQuestion is a question whose answer is the Header field of the answer of the wrapped question.
type OfQuestion struct {
	question screenplay.Question
}

// AnsweredBy returns the Header field of the struct answered by the wrapped question.
func (q *OfQuestion) AnsweredBy(actor *screenplay.Actor) (any, error) {
	answer, err := q.question.AnsweredBy(actor)
	if err != nil {
		return nil, err
	}

	value := reflect.ValueOf(answer)
	if value.IsValid() && value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	if !value.IsValid() || value.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%w: %s", ErrAnswerIsNotAStruct, q.question)
	}

	field := value.FieldByName("Header")
	if !field.IsValid() || !field.CanInterface() {
		return nil, fmt.Errorf("%w: %s", ErrFieldNotFound, q.question)
	}

	return field.Interface(), nil
}

// String describes the question.
func (q *OfQuestion) String() string {
	return "header of " + q.question.String()
}

// OfQuestion implements the screenplay.Question interface.
var _ screenplay.Question = (*OfQuestion)(nil)
