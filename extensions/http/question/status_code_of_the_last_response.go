package question

import (
	"fmt"

	"github.com/grandper/go-screenplay/extensions/http/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

// StatusCodeOfTheLastResponse asks about the status code of the last HTTP response.
func StatusCodeOfTheLastResponse() screenplay.Question {
	return &StatusCodeOfTheLastResponseQuestion{}
}

// StatusCodeOfTheLastResponseQuestion asks about the status code of the last HTTP response.
type StatusCodeOfTheLastResponseQuestion struct{}

// String describes the question.
func (q *StatusCodeOfTheLastResponseQuestion) String() string {
	return "HTTP status code of the last response"
}

// AnsweredBy returns the answer that an actor provided to the question.
func (q *StatusCodeOfTheLastResponseQuestion) AnsweredBy(theActor *screenplay.Actor) (any, error) {
	makeHTTPRequests, err := screenplay.UseAbilityTo[*ability.MakeHTTPRequestsAbility]().Of(theActor)
	if err != nil {
		return nil, err
	}
	responses := makeHTTPRequests.ToRetrieveResponses()
	if len(responses) < 1 {
		return nil, fmt.Errorf("%s has not received any HTTP responses", theActor.Name())
	}
	return responses[len(responses)-1].StatusCode(), nil
}

// StatusCodeTheLastResponseQuestion implements the screenplay.Question interface.
var _ screenplay.Question = &StatusCodeOfTheLastResponseQuestion{}
