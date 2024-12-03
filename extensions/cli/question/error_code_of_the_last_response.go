package question

import (
	"fmt"

	"github.com/grandper/go-screenplay/extensions/cli/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

var (
	// ErrNoResponses is returned when no responses have been recorded.
	ErrNoResponses = fmt.Errorf("no responses have been recorded")
)

// ErrorCodeOfTheLastResponse asks about the error code of the last response.
func ErrorCodeOfTheLastResponse() screenplay.Question {
	return &ErrorCodeOfTheLastResponseQuestion{}
}

// ErrorCodeOfTheLastResponseQuestion asks about the error code of the last HTTP response.
type ErrorCodeOfTheLastResponseQuestion struct{}

// String describes the question.
func (q *ErrorCodeOfTheLastResponseQuestion) String() string {
	return "error code of the last response"
}

// AnsweredBy returns the answer that an actor provided to the question.
func (q *ErrorCodeOfTheLastResponseQuestion) AnsweredBy(theActor *screenplay.Actor) (any, error) {
	theActorMakesHTTPRequests, err := screenplay.UseAbilityTo[*ability.RunCLICommandsAbility]().Of(theActor)
	if err != nil {
		return nil, err
	}
	responses := theActorMakesHTTPRequests.Responses()
	if len(responses) < 1 {
		return nil, ErrNoResponses
	}
	return responses[len(responses)-1].ExitCode(), nil
}

// ErrorCodeOfTheLastResponseQuestion implements the screenplay.Question interface.
var _ screenplay.Question = &ErrorCodeOfTheLastResponseQuestion{}
