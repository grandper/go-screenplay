package log

import "github.com/grandper/go-screenplay/v2/screenplay"

// The answers a question and narrates its value as an aside, which is handy
// for debugging:
//
//	log.The(statusCode)
//
// It performs screenplay.Log, which holds the narration.
func The[T any](question screenplay.Question[T]) screenplay.Performable {
	return &theAction[T]{
		log: screenplay.Log(question),
	}
}

// theAction is an action that narrates the answer to a question.
type theAction[T any] struct {
	log screenplay.Performable
}

// String describes the action.
func (a *theAction[T]) String() string {
	return a.log.String()
}

// PerformAs answers the question and, when the actor is narrating, whispers its
// value as an aside. It returns any error the question produced.
func (a *theAction[T]) PerformAs(actor *screenplay.Actor) error {
	return a.log.PerformAs(actor)
}

// theAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*theAction[any])(nil)
