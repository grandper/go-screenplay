package question

import (
	"github.com/grandper/go-screenplay/extensions/cli/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

// Responses asks about all the responses recorded so far.
func Responses() screenplay.Question {
	return &ResponsesQuestion{}
}

// ResponsesQuestion asks about all the responses recorded so far.
type ResponsesQuestion struct{}

// String describes the question.
func (q *ResponsesQuestion) String() string {
	return "responses"
}

// AnsweredBy returns the answer that an actor provided to the question.
func (q *ResponsesQuestion) AnsweredBy(theActor *screenplay.Actor) (any, error) {
	theActorRunsCLICommands, err := screenplay.UseAbilityTo[*ability.RunCLICommandsAbility]().Of(theActor)
	if err != nil {
		return nil, err
	}
	return theActorRunsCLICommands.Responses(), nil
}

// ResponsesQuestion implements the screenplay.Question interface.
var _ screenplay.Question = (*ResponsesQuestion)(nil)
