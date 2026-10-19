package question

import (
	"github.com/grandper/go-screenplay/v2/extensions/cli/ability"
	impl "github.com/grandper/go-screenplay/v2/extensions/cli/internal/ability"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// Responses asks about all the responses recorded so far.
func Responses() screenplay.Question[[]*ability.Result] {
	return &responsesQuestion{}
}

// responsesQuestion asks about all the responses recorded so far.
type responsesQuestion struct{}

// String describes the question.
func (q *responsesQuestion) String() string {
	return "command results"
}

// AnsweredBy returns the answer that an actor provided to the question.
func (q *responsesQuestion) AnsweredBy(theActor *screenplay.Actor) ([]*ability.Result, error) {
	return screenplay.UseAnAbilityOf(theActor).That(impl.RunCLICommands()).ToRetrieveResponses()
}

// responsesQuestion implements the screenplay.Question interface.
var _ screenplay.Question[[]*ability.Result] = (*responsesQuestion)(nil)
