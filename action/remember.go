package action

import (
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/screenplay"
)

// ErrAnswerIsNil is returned when the question to remember returns a nil answer.
var ErrAnswerIsNil = errors.New("the answer of the question is nil")

// Remember creates an action that asks a question and stores its answer in the actor's memory.
func Remember(question screenplay.Question) *RememberBuilder {
	return &RememberBuilder{
		question: question,
	}
}

// RememberBuilder builds a Remember action.
type RememberBuilder struct {
	question screenplay.Question
}

// As sets the key under which the answer is stored.
func (b *RememberBuilder) As(key string) *RememberAction {
	return &RememberAction{
		question: b.question,
		key:      key,
		allowNil: false,
	}
}

// RememberAction is an action that asks a question and remembers its answer.
type RememberAction struct {
	question screenplay.Question
	key      string
	allowNil bool
}

// AllowingNil permits storing a nil answer without raising an error.
func (a *RememberAction) AllowingNil() *RememberAction {
	a.allowNil = true
	return a
}

// String describes the action.
func (a *RememberAction) String() string {
	return fmt.Sprintf("remember the %s as '%s'", a.question, a.key)
}

// PerformAs performs the task or the action as the provided actor.
func (a *RememberAction) PerformAs(actor *screenplay.Actor) error {
	answer, err := a.question.AnsweredBy(actor)
	if err != nil {
		return err
	}

	if answer == nil && !a.allowNil {
		return fmt.Errorf("%w: %s", ErrAnswerIsNil, a.question)
	}

	actor.Remember(a.key, answer)
	return nil
}

// RememberAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*RememberAction)(nil)
