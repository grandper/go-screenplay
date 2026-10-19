package remember

import (
	impl "github.com/grandper/go-screenplay/v2/action/internal/remember"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// ErrAnswerIsNil is returned when the question to remember returns a nil answer.
var ErrAnswerIsNil = impl.ErrAnswerIsNil

// The creates an action that asks a question and stores its answer in the
// actor's memory under a key:
//
//	remember.The(bearerToken).As("session token")
func The[T any](question screenplay.Question[T]) *impl.TheBuilder[T] {
	return impl.The(question)
}
