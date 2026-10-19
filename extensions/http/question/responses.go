package question

import (
	"github.com/grandper/go-screenplay/v2/extensions/http/ability"
	impl "github.com/grandper/go-screenplay/v2/extensions/http/internal/ability"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// Responses asks about all the HTTP responses received by the actor.
func Responses() screenplay.Question[[]*ability.HTTPResponse] {
	return &responsesQuestion{}
}

// responsesQuestion asks about all the HTTP responses received by the actor.
type responsesQuestion struct{}

// String describes the question.
func (q *responsesQuestion) String() string {
	return "HTTP responses"
}

// AnsweredBy returns the answer that an actor provided to the question.
func (q *responsesQuestion) AnsweredBy(theActor *screenplay.Actor) ([]*ability.HTTPResponse, error) {
	return impl.TheAbilityOf(theActor).ToRetrieveResponses()
}

// responsesQuestion implements the screenplay.Question interface.
var _ screenplay.Question[[]*ability.HTTPResponse] = (*responsesQuestion)(nil)
