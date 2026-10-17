package wait

import (
	"fmt"
	"io"
	"os"

	"github.com/grandper/go-screenplay/v2/screenplay"
)

// UntilTheEnterKeyIsPressed makes an actor wait until the 'enter' key is
// pressed, so that a test can be paused while debugging it.
func UntilTheEnterKeyIsPressed() *UntilTheEnterKeyIsPressedAction {
	return &UntilTheEnterKeyIsPressedAction{}
}

// UntilTheEnterKeyIsPressedAction is an action that makes an actor wait until
// the 'enter' key is pressed.
type UntilTheEnterKeyIsPressedAction struct{}

// String describes the action.
func (a *UntilTheEnterKeyIsPressedAction) String() string {
	return "wait until the enter key is pressed"
}

// PerformAs blocks until a line is read from the standard input, or until the
// actor's context ends.
func (a *UntilTheEnterKeyIsPressedAction) PerformAs(actor *screenplay.Actor) error {
	return waitForAnInput(actor, os.Stdin)
}

// UntilTheEnterKeyIsPressedAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*UntilTheEnterKeyIsPressedAction)(nil)

// UntilAnInputIsProvidedBy makes an actor wait until a line is read from the
// reader.
func UntilAnInputIsProvidedBy(reader io.Reader) *UntilAnInputIsProvidedByAction {
	return &UntilAnInputIsProvidedByAction{
		reader: reader,
	}
}

// UntilAnInputIsProvidedByAction is an action that makes an actor wait until a
// line is read from a reader.
type UntilAnInputIsProvidedByAction struct {
	reader io.Reader
}

// String describes the action. A reader has no name to narrate.
func (a *UntilAnInputIsProvidedByAction) String() string {
	return "wait until an input is provided"
}

// PerformAs blocks until a line is read from the reader, or until the actor's
// context ends.
func (a *UntilAnInputIsProvidedByAction) PerformAs(actor *screenplay.Actor) error {
	return waitForAnInput(actor, a.reader)
}

// UntilAnInputIsProvidedByAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*UntilAnInputIsProvidedByAction)(nil)

// waitForAnInput blocks until a line is read from the reader, or until the
// actor's context ends. The prompt is written to the standard error of the
// process, where the person debugging the test reads it; the narrator already
// states that the actor waits.
func waitForAnInput(actor *screenplay.Actor, reader io.Reader) error {
	fmt.Fprintln(os.Stderr, "the actor stops, waiting for your go (press enter to continue)")

	read := make(chan error, 1)
	go func() {
		_, err := fmt.Fscanln(reader)
		read <- err
	}()

	ctx := actor.Context()
	select {
	case err := <-read:
		return err
	case <-ctx.Done():
		return fmt.Errorf("%s stopped waiting for an input: %w", actor.Name(), ctx.Err())
	}
}
