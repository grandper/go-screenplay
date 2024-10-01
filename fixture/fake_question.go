package fixture

import (
	"sync"

	"github.com/grandper/go-screenplay/screenplay"
)

// FakeQuestion represents a fake question that an actor can ask.
type FakeQuestion struct {
	mutex       sync.RWMutex
	answer      any
	err         error
	description string
}

// NewFakeQuestion creates a new fake question.
func NewFakeQuestion(description string, answer any) *FakeQuestion {
	return &FakeQuestion{
		mutex:       sync.RWMutex{},
		answer:      answer,
		err:         nil,
		description: description,
	}
}

// NewFailingFakeQuestion creates a new fake question that fails.
func NewFailingFakeQuestion(description string, err error) *FakeQuestion {
	return &FakeQuestion{
		mutex:       sync.RWMutex{},
		answer:      nil,
		err:         err,
		description: description,
	}
}

// AnsweredBy returns the answer that an actor provided to the question.
func (q *FakeQuestion) AnsweredBy(_ *screenplay.Actor) (any, error) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	if q.err != nil {
		return nil, q.err
	}

	return q.answer, nil
}

// AnswerWith sets the answer that an actor provided to the question.
func (q *FakeQuestion) AnswerWith(answer string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.answer = answer
}

// String describes the question.
func (q *FakeQuestion) String() string {
	return q.description
}

// FakeQuestion implements the screenplay.Question interface.
var _ screenplay.Question = (*FakeQuestion)(nil)
