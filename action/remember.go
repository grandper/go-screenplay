package action

import (
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/v2/internal/verdict"

	"github.com/grandper/go-screenplay/v2/screenplay"
)

// ErrAnswerIsNil is returned when the question to remember returns a nil answer.
var ErrAnswerIsNil = errors.New("the answer of the question is nil")

// Remember creates an action that asks a question and stores its answer in the actor's memory.
func Remember[T any](question screenplay.Question[T]) *RememberBuilder[T] {
	return &RememberBuilder[T]{
		question: question,
	}
}

// RememberBuilder builds a Remember action.
type RememberBuilder[T any] struct {
	question screenplay.Question[T]
}

// As sets the key under which the answer is stored.
func (b *RememberBuilder[T]) As(key string) *RememberAction[T] {
	return &RememberAction[T]{
		question: b.question,
		key:      key,
		allowNil: false,
	}
}

// RememberAction is an action that asks a question and remembers its answer.
type RememberAction[T any] struct {
	question screenplay.Question[T]
	key      string
	allowNil bool
}

// AllowingNil permits storing a nil answer without raising an error.
func (a *RememberAction[T]) AllowingNil() *RememberAction[T] {
	a.allowNil = true
	return a
}

// String describes the action.
func (a *RememberAction[T]) String() string {
	return fmt.Sprintf("remember the %s as %s", a.question, verdict.Describe(a.key))
}

// PerformAs performs the task or the action as the provided actor.
func (a *RememberAction[T]) PerformAs(actor *screenplay.Actor) error {
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

// RememberAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*RememberAction[any])(nil)
