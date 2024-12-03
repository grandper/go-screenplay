package question

import (
	"fmt"

	"github.com/grandper/go-screenplay/extensions/http/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

// BodyOfTheLastResponse asks about the body of the last HTTP response.
func BodyOfTheLastResponse() screenplay.Question {
	return &BodyOfTheLastResponseQuestion{}
}

// BodyOfTheLastResponseQuestion asks about the body of the last HTTP response.
type BodyOfTheLastResponseQuestion struct{}

// String describes the question.
func (q *BodyOfTheLastResponseQuestion) String() string {
	return "body of the last response"
}

// AnsweredBy returns the answer that an actor provided to the question.
func (q *BodyOfTheLastResponseQuestion) AnsweredBy(theActor *screenplay.Actor) (any, error) {
	theActorMakesHTTPRequests, err := screenplay.UseAbilityTo[*ability.MakeHTTPRequestsAbility]().Of(theActor)
	if err != nil {
		return nil, err
	}
	responses := theActorMakesHTTPRequests.ToRetrieveResponses()
	if len(responses) < 1 {
		return nil, fmt.Errorf("%s has not received any HTTP responses", theActor.Name())
	}
	return responses[len(responses)-1].Body(), nil
}

// BodyOfTheLastResponseQuestion implements the screenplay.Question interface.
var _ screenplay.Question = &BodyOfTheLastResponseQuestion{}
