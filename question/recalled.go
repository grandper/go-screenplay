package question

import (
	"fmt"

	"github.com/grandper/go-screenplay/screenplay"
)

// Recalled creates a question whose answer is the value stored in the actor's
// memory under the given key. It is the counterpart of action.Remember: where
// Remember stores the answer of a question, Recalled reads a stored value back
// as a question, so it can be checked with see.The, wrapped in Eventually, or
// handed to any other component that expects a screenplay.Question.
//
// The answer is exactly what actor.Recall returns: the stored value, or nil
// when nothing is stored under the key.
func Recalled(key string) *RecalledQuestion {
	return &RecalledQuestion{
		key: key,
	}
}

// RecalledQuestion is a question answered by a value stored in the actor's memory.
type RecalledQuestion struct {
	key string
}

// String describes the question.
func (q *RecalledQuestion) String() string {
	return fmt.Sprintf("recalled '%s'", q.key)
}

// AnsweredBy returns the value the actor remembers under the question's key,
// or nil when the actor remembers nothing under it.
func (q *RecalledQuestion) AnsweredBy(actor *screenplay.Actor) (any, error) {
	return actor.Recall(q.key), nil
}

// RecalledQuestion implements the screenplay.Question interface.
var _ screenplay.Question = (*RecalledQuestion)(nil)
