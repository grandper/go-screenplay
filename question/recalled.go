package question

import (
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/v2/internal/verdict"

	"github.com/grandper/go-screenplay/v2/screenplay"
)

// ErrRecalledValueHasWrongType is returned when the value remembered under the
// key is not of the type the question expects.
var ErrRecalledValueHasWrongType = errors.New("the recalled value has the wrong type")

// Recalled creates a question whose answer is the value stored in the actor's
// memory under the given key, typed as T. It is the counterpart of
// remember.The: where remember.The stores the answer of a question, Recalled
// reads a stored value back as a question, so it can be checked with see.The,
// waited for with wait.UntilThe, or handed to any other component that expects a
// screenplay.Question.
//
// The answer is the stored value when it is a T, or the zero value of T (nil
// for any) when nothing is stored under the key. It fails with
// ErrRecalledValueHasWrongType when a value is stored but is not a T:
//
//	statusCode, err := actor.AsksFor(question.Recalled[int]("statusCode"))
//	anything, err := actor.AsksFor(question.Recalled[any]("statusCode"))
func Recalled[T any](key string) screenplay.Question[T] {
	return &recalledQuestion[T]{
		key: key,
	}
}

// recalledQuestion is a question answered by a value stored in the actor's memory.
type recalledQuestion[T any] struct {
	key string
}

// String describes the question.
func (q *recalledQuestion[T]) String() string {
	return fmt.Sprintf("recalled %s", verdict.Describe(q.key))
}

// AnsweredBy returns the value the actor remembers under the question's key,
// or the zero value of T when the actor remembers nothing under it.
func (q *recalledQuestion[T]) AnsweredBy(actor *screenplay.Actor) (T, error) {
	var zero T

	value := actor.Recall[any](q.key).OrElse(nil)
	if value == nil {
		return zero, nil
	}

	answer, isAT := value.(T)
	if !isAT {
		return zero, fmt.Errorf("%w: %s holds a %T, not a %T", ErrRecalledValueHasWrongType, q, value, zero)
	}

	return answer, nil
}

// recalledQuestion implements the screenplay.Question interface.
var _ screenplay.Question[any] = (*recalledQuestion[any])(nil)
