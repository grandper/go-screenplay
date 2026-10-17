package log

import "github.com/grandper/go-screenplay/v2/screenplay"

// The answers a question and narrates its value as an aside, which is handy
// for debugging:
//
//	log.The(statusCode)
//
// It performs screenplay.Log, which holds the narration.
func The[T any](question screenplay.Question[T]) *TheAction[T] {
	return &TheAction[T]{
		log: screenplay.Log(question),
	}
}

// TheAction is an action that narrates the answer to a question.
type TheAction[T any] struct {
	log *screenplay.LogAction[T]
}

// String describes the action.
func (a *TheAction[T]) String() string {
	return a.log.String()
}

// PerformAs answers the question and, when the actor is narrating, whispers its
// value as an aside. It returns any error the question produced.
func (a *TheAction[T]) PerformAs(actor *screenplay.Actor) error {
	return a.log.PerformAs(actor)
}

// TheAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*TheAction[any])(nil)
