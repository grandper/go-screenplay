package question

import (
	"github.com/grandper/go-screenplay/extensions/http/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

// Responses asks about all the HTTP responses received by the actor.
func Responses() screenplay.Question {
	return &ResponsesQuestion{}
}

// ResponsesQuestion asks about all the HTTP responses received by the actor.
type ResponsesQuestion struct{}

// String describes the question.
func (q *ResponsesQuestion) String() string {
	return "the HTTP responses"
}

// AnsweredBy returns the answer that an actor provided to the question.
func (q *ResponsesQuestion) AnsweredBy(theActor *screenplay.Actor) (any, error) {
	makeHTTPRequests, err := screenplay.UseAbilityTo[*ability.MakeHTTPRequestsAbility]().Of(theActor)
	if err != nil {
		return nil, err
	}
	return makeHTTPRequests.ToRetrieveResponses(), nil
}

// ResponsesQuestion implements the screenplay.Question interface.
var _ screenplay.Question = &ResponsesQuestion{}
