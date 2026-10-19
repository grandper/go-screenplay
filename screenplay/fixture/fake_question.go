package fixture

import (
	"sync"

	"github.com/grandper/go-screenplay/v2/screenplay"
)

// FakeQuestion represents a fake question that an actor can ask. T is the type
// of its answer.
type FakeQuestion[T any] struct {
	mutex       sync.RWMutex
	answer      T
	err         error
	description string
}

// NewFakeQuestion creates a new fake question answering with the given value.
// The answer type is inferred from the value; spell it out when the value is
// untyped, for example NewFakeQuestion[any]("nil answer", nil).
func NewFakeQuestion[T any](description string, answer T) *FakeQuestion[T] {
	return &FakeQuestion[T]{
		mutex:       sync.RWMutex{},
		answer:      answer,
		err:         nil,
		description: description,
	}
}

// NewFailingFakeQuestion creates a new fake question that fails. The answer
// type cannot be inferred and must be provided, for example
// NewFailingFakeQuestion[string]("the token", err).
func NewFailingFakeQuestion[T any](description string, err error) *FakeQuestion[T] {
	return &FakeQuestion[T]{
		mutex:       sync.RWMutex{},
		err:         err,
		description: description,
	}
}

// AnsweredBy returns the answer that an actor provided to the question.
func (q *FakeQuestion[T]) AnsweredBy(_ *screenplay.Actor) (T, error) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	if q.err != nil {
		var zero T

		return zero, q.err
	}

	return q.answer, nil
}

// AnswerWith sets the answer that an actor provided to the question.
func (q *FakeQuestion[T]) AnswerWith(answer T) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.answer = answer
}

// String describes the question.
func (q *FakeQuestion[T]) String() string {
	return q.description
}

// FakeQuestion implements the screenplay.Question interface.
var _ screenplay.Question[any] = (*FakeQuestion[any])(nil)
