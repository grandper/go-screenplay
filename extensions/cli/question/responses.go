package question

import (
	"github.com/grandper/go-screenplay/v2/extensions/cli/ability"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// Responses asks about all the responses recorded so far.
func Responses() *ResponsesQuestion {
	return &ResponsesQuestion{}
}

// ResponsesQuestion asks about all the responses recorded so far.
type ResponsesQuestion struct{}

// String describes the question.
func (q *ResponsesQuestion) String() string {
	return "command results"
}

// AnsweredBy returns the answer that an actor provided to the question.
func (q *ResponsesQuestion) AnsweredBy(theActor *screenplay.Actor) ([]*ability.Result, error) {
	return screenplay.UseAnAbilityOf(theActor).That(ability.RunCLICommands()).ToRetrieveResponses()
}

// ResponsesQuestion implements the screenplay.Question interface.
var _ screenplay.Question[[]*ability.Result] = (*ResponsesQuestion)(nil)
