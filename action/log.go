package action

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/grandper/go-screenplay/screenplay"
)

// Log logs the answer to a question.
// This can be used for debugging purpose.
func Log(question screenplay.Question) *LogAction {
	return &LogAction{
		question: question,
		logFn: func(format string, args ...any) {
			slog.Default().InfoContext(context.Background(), format, args...)
		},
	}
}

// LogAction is an action to log the answer to a question.
type LogAction struct {
	question screenplay.Question
	logFn    func(string, ...any)
}

// String describes the action.
func (a *LogAction) String() string {
	return fmt.Sprintf("log the %s", a.question.String())
}

// PerformAs performs the task or the action as the provided actor.
func (a *LogAction) PerformAs(actor *screenplay.Actor) error {
	answer, err := a.question.AnsweredBy(actor)
	if err != nil {
		return err
	}

	a.logFn("the value of the '%s' is '%s'", a.question.String(), answer)
	return nil
}

// LogAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*LogAction)(nil)
