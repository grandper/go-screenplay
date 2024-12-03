package question

import (
	"github.com/grandper/go-screenplay/extensions/cli/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

// StandardErrorOfTheLastResponse asks about the standard error of the last response.
func StandardErrorOfTheLastResponse() screenplay.Question {
	return &StandardErrorOfTheLastResponseQuestion{}
}

// StandardErrorOfTheLastResponseQuestion asks about the standard error of the last HTTP response.
type StandardErrorOfTheLastResponseQuestion struct{}

// String describes the question.
func (q *StandardErrorOfTheLastResponseQuestion) String() string {
	return "standard error of the last response"
}

// AnsweredBy returns the answer that an actor provided to the question.
func (q *StandardErrorOfTheLastResponseQuestion) AnsweredBy(theActor *screenplay.Actor) (any, error) {
	theActorMakesHTTPRequests, err := screenplay.UseAbilityTo[*ability.RunCLICommandsAbility]().Of(theActor)
	if err != nil {
		return nil, err
	}
	responses := theActorMakesHTTPRequests.Responses()
	if len(responses) < 1 {
		return nil, ErrNoResponses
	}
	return responses[len(responses)-1].StdErr(), nil
}

// StandardErrorOfTheLastResponseQuestion implements the screenplay.Question interface.
var _ screenplay.Question = &StandardErrorOfTheLastResponseQuestion{}
