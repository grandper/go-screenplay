package action

import (
	"fmt"

	"github.com/grandper/go-screenplay/screenplay"
)

// Log asks a question so its answer is narrated, which is handy for debugging.
// The line is rendered by whatever adapters are attached to the actor's
// narrator; with no adapter attached it is silent.
func Log(question screenplay.Question) *LogAction {
	return &LogAction{
		question: question,
	}
}

// LogAction is an action to log the answer to a question.
type LogAction struct {
	question screenplay.Question
}

// String describes the action.
func (a *LogAction) String() string {
	return fmt.Sprintf("log the %s", a.question.String())
}

// PerformAs asks the question through the actor so the narrator reveals its
// answer, and returns any error the question produced.
func (a *LogAction) PerformAs(actor *screenplay.Actor) error {
	_, err := actor.AsksFor(a.question)

	return err
}

// LogAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*LogAction)(nil)
