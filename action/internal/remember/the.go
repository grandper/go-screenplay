package remember

import (
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/v2/internal/verdict"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// ErrAnswerIsNil is returned when the question to remember returns a nil answer.
var ErrAnswerIsNil = errors.New("the answer of the question is nil")

// The creates an action that asks a question and stores its answer in the
// actor's memory under a key:
//
//	remember.The(bearerToken).As("session token")
func The[T any](question screenplay.Question[T]) *TheBuilder[T] {
	return &TheBuilder[T]{
		question: question,
	}
}

// TheBuilder waits for the key under which the answer is stored.
type TheBuilder[T any] struct {
	question screenplay.Question[T]
}

// As sets the key under which the answer is stored.
func (b *TheBuilder[T]) As(key string) *TheAction[T] {
	return &TheAction[T]{
		question: b.question,
		key:      key,
		allowNil: false,
	}
}

// TheAction is an action that asks a question and remembers its answer.
type TheAction[T any] struct {
	question screenplay.Question[T]
	key      string
	allowNil bool
}

// AllowingNil permits storing a nil answer without raising an error.
func (a *TheAction[T]) AllowingNil() *TheAction[T] {
	a.allowNil = true
	return a
}

// String describes the action.
func (a *TheAction[T]) String() string {
	description := fmt.Sprintf("remember the %s as %s", a.question, verdict.Describe(a.key))
	if a.allowNil {
		description += " allowing nil"
	}

	return description
}

// PerformAs performs the task or the action as the provided actor.
func (a *TheAction[T]) PerformAs(actor *screenplay.Actor) error {
	answer, err := a.question.AnsweredBy(actor)
	if err != nil {
		return err
	}

	if any(answer) == nil && !a.allowNil {
		return fmt.Errorf("%w: %s", ErrAnswerIsNil, a.question)
	}

	actor.Remember(a.key, answer)
	return nil
}

// TheAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*TheAction[any])(nil)
