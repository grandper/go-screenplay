package question

import (
	"fmt"

	"github.com/grandper/go-screenplay/extensions/http/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

// HeadersOfTheLastResponse asks about the headers of the last HTTP response.
func HeadersOfTheLastResponse() screenplay.Question {
	return &HeadersOfTheLastResponseQuestion{}
}

// HeadersOfTheLastResponseQuestion asks about the headers of the last HTTP response.
type HeadersOfTheLastResponseQuestion struct{}

// String describes the question.
func (q *HeadersOfTheLastResponseQuestion) String() string {
	return "headers of the last response"
}

// AnsweredBy returns the answer that an actor provided to the question.
func (q *HeadersOfTheLastResponseQuestion) AnsweredBy(theActor *screenplay.Actor) (any, error) {
	makeHTTPRequests, err := screenplay.UseAbilityTo[*ability.MakeHTTPRequestsAbility]().Of(theActor)
	if err != nil {
		return nil, err
	}
	responses := makeHTTPRequests.ToRetrieveResponses()
	if len(responses) < 1 {
		return nil, fmt.Errorf("%s has not received any HTTP responses", theActor.Name())
	}
	return responses[len(responses)-1].Headers(), nil
}

// HeadersOfTheLastResponseQuestion implements the screenplay.Question interface.
var _ screenplay.Question = &HeadersOfTheLastResponseQuestion{}
