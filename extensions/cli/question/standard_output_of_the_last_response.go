package question

import (
	"github.com/grandper/go-screenplay/extensions/cli/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

// StandardOutputOfTheLastResponse asks about the standard output of the last response.
func StandardOutputOfTheLastResponse() screenplay.Question {
	return &StandardOutputOfTheLastResponseQuestion{}
}

// StandardOutputOfTheLastResponseQuestion asks about the standard output of the last HTTP response.
type StandardOutputOfTheLastResponseQuestion struct{}

// String describes the question.
func (q *StandardOutputOfTheLastResponseQuestion) String() string {
	return "standard output of the last response"
}

// AnsweredBy returns the answer that an actor provided to the question.
func (q *StandardOutputOfTheLastResponseQuestion) AnsweredBy(theActor *screenplay.Actor) (any, error) {
	theActorMakesHTTPRequests, err := screenplay.UseAbilityTo[*ability.RunCLICommandsAbility]().Of(theActor)
	if err != nil {
		return nil, err
	}
	responses := theActorMakesHTTPRequests.Responses()
	if len(responses) < 1 {
		return nil, ErrNoResponses
	}
	return responses[len(responses)-1].StdOut(), nil
}

// StandardOutputOfTheLastResponseQuestion implements the screenplay.Question interface.
var _ screenplay.Question = &StandardOutputOfTheLastResponseQuestion{}
