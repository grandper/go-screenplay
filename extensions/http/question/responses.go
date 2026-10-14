package question

import (
	"github.com/grandper/go-screenplay/v2/extensions/http/ability"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// Responses asks about all the HTTP responses received by the actor.
func Responses() *ResponsesQuestion {
	return &ResponsesQuestion{}
}

// ResponsesQuestion asks about all the HTTP responses received by the actor.
type ResponsesQuestion struct{}

// String describes the question.
func (q *ResponsesQuestion) String() string {
	return "HTTP responses"
}

// AnsweredBy returns the answer that an actor provided to the question.
func (q *ResponsesQuestion) AnsweredBy(theActor *screenplay.Actor) ([]*ability.HTTPResponse, error) {
	return screenplay.UseAnAbilityOf(theActor).That(ability.MakeHTTPRequests()).ToRetrieveResponses()
}

// ResponsesQuestion implements the screenplay.Question interface.
var _ screenplay.Question[[]*ability.HTTPResponse] = (*ResponsesQuestion)(nil)
