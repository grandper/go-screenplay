package action

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/grandper/go-screenplay/v2/action/see"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// Stop makes an actor stop until a condition is met.
// If stop is used directly, the actor will stop until you press 'enter'
// to continue the test. Otherwise, the actor will stop until
// the condition is met or raise a timeout error.
func Stop() *StopAction {
	return &StopAction{
		reader: os.Stdin,
	}
}

// StopAction is an action to make an actor stop until a condition is met.
type StopAction struct {
	condition  screenplay.Performable // sees that the answer matches the resolution; nil when waiting for an input
	question   fmt.Stringer
	resolution screenplay.Resolution
	reader     io.Reader
}

// UntilAnInputIsProvidedBy specifies a reader to wait for.
func (a *StopAction) UntilAnInputIsProvidedBy(reader io.Reader) *StopAction {
	return &StopAction{
		reader: reader,
	}
}

// UntilThe specifies a condition to wait for.
func (a *StopAction) UntilThe[T any](question screenplay.Question[T], resolution screenplay.Resolution) *StopAction {
	return &StopAction{
		condition:  see.The(question).Is(resolution),
		question:   question,
		resolution: resolution,
	}
}

// String describes the action.
func (a *StopAction) String() string {
	if a.condition == nil {
		return "stop until the 'enter' key is pressed"
	}
	return fmt.Sprintf("stop until the %s is %s", a.question, a.resolution)
}

// PerformAs performs the task or the action as the provided actor.
func (a *StopAction) PerformAs(actor *screenplay.Actor) error {
	if a.condition == nil {
		slog.Default().
			InfoContext(context.Background(), "the actor stops, waiting for your go\n(press enter to continue)\n")

		_, err := fmt.Fscanln(a.reader)
		if err != nil {
			return err
		}
		return nil
	}
	if err := actor.AttemptsTo(Eventually(a.condition)); err != nil {
		return fmt.Errorf("%s stopped for %s, but %s was never %s",
			actor.Name(), actor.Timeout(), a.question.String(), a.resolution.String())
	}
	return nil
}

// StopAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*StopAction)(nil)
