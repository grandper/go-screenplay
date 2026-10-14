package action

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/grandper/go-screenplay/v2/internal/sentence"
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
	condition screenplay.Performable // sees that the term is met; nil when waiting for an input
	term      *sentence.Term
	reader    io.Reader
}

// UntilAnInputIsProvidedBy specifies a reader to wait for.
func (a *StopAction) UntilAnInputIsProvidedBy(reader io.Reader) *StopAction {
	return &StopAction{
		reader: reader,
	}
}

// UntilThe specifies what the answer to a question happens to be for the actor
// to carry on. Close it with a verb and one or more resolutions:
//
//	Stop().UntilThe(logs).HappenTo(contain.Substring("ready"))
func (a *StopAction) UntilThe[T any](question screenplay.Question[T]) *StopBuilder[T] {
	return &StopBuilder[T]{question: question}
}

// StopBuilder builds the condition an actor stops for: it holds the question,
// and waits for the verb that closes the condition.
type StopBuilder[T any] struct {
	question screenplay.Question[T]
}

// HappensTo waits for the answer to satisfy every resolution. It reads "happens to ...".
func (s *StopBuilder[T]) HappensTo(resolutions ...screenplay.Resolution[T]) *StopAction {
	return s.waitsFor(sentence.NewCondition(sentence.Should, false, resolutions))
}

// HappensToBe waits for the answer to satisfy every resolution. It reads "happens to be ...".
func (s *StopBuilder[T]) HappensToBe(resolutions ...screenplay.Resolution[T]) *StopAction {
	return s.waitsFor(sentence.NewCondition(sentence.ShouldBe, false, resolutions))
}

// HappensToHave waits for the answer to satisfy every resolution. It reads "happens to have ...".
func (s *StopBuilder[T]) HappensToHave(resolutions ...screenplay.Resolution[T]) *StopAction {
	return s.waitsFor(sentence.NewCondition(sentence.ShouldHave, false, resolutions))
}

// DoesNotHappenTo waits for the answer to satisfy none of the resolutions. It reads "does not happen to ...".
func (s *StopBuilder[T]) DoesNotHappenTo(resolutions ...screenplay.Resolution[T]) *StopAction {
	return s.waitsFor(sentence.NewCondition(sentence.ShouldNot, false, resolutions))
}

// DoesNotHappenToBe waits for the answer to satisfy none of the resolutions. It reads "does not happen to be ...".
func (s *StopBuilder[T]) DoesNotHappenToBe(resolutions ...screenplay.Resolution[T]) *StopAction {
	return s.waitsFor(sentence.NewCondition(sentence.ShouldNotBe, false, resolutions))
}

// DoesNotHappenToHave waits for the answer to satisfy none of the resolutions. It reads "does not happen to have ...".
func (s *StopBuilder[T]) DoesNotHappenToHave(resolutions ...screenplay.Resolution[T]) *StopAction {
	return s.waitsFor(sentence.NewCondition(sentence.ShouldNotHave, false, resolutions))
}

// HappensToBeAnyOf waits for the answer to satisfy at least one resolution. It reads "happens to be any of: ...".
func (s *StopBuilder[T]) HappensToBeAnyOf(resolutions ...screenplay.Resolution[T]) *StopAction {
	return s.waitsFor(sentence.NewCondition(sentence.ShouldAnyOf, false, resolutions))
}

// HappenTo waits for the answer to satisfy every resolution. It is the plural form: it reads "happen to ...".
func (s *StopBuilder[T]) HappenTo(resolutions ...screenplay.Resolution[T]) *StopAction {
	return s.waitsFor(sentence.NewCondition(sentence.Should, true, resolutions))
}

// HappenToBe waits for the answer to satisfy every resolution. It is the plural form: it reads "happen to be ...".
func (s *StopBuilder[T]) HappenToBe(resolutions ...screenplay.Resolution[T]) *StopAction {
	return s.waitsFor(sentence.NewCondition(sentence.ShouldBe, true, resolutions))
}

// HappenToHave waits for the answer to satisfy every resolution. It is the plural form: it reads "happen to have ...".
func (s *StopBuilder[T]) HappenToHave(resolutions ...screenplay.Resolution[T]) *StopAction {
	return s.waitsFor(sentence.NewCondition(sentence.ShouldHave, true, resolutions))
}

// DoNotHappenTo waits for the answer to satisfy none of the resolutions. It is the plural form: it reads "do not happen to ...".
func (s *StopBuilder[T]) DoNotHappenTo(resolutions ...screenplay.Resolution[T]) *StopAction {
	return s.waitsFor(sentence.NewCondition(sentence.ShouldNot, true, resolutions))
}

// DoNotHappenToBe waits for the answer to satisfy none of the resolutions. It is the plural form: it reads "do not happen to be ...".
func (s *StopBuilder[T]) DoNotHappenToBe(resolutions ...screenplay.Resolution[T]) *StopAction {
	return s.waitsFor(sentence.NewCondition(sentence.ShouldNotBe, true, resolutions))
}

// DoNotHappenToHave waits for the answer to satisfy none of the resolutions. It is the plural form: it reads "do not happen to have ...".
func (s *StopBuilder[T]) DoNotHappenToHave(resolutions ...screenplay.Resolution[T]) *StopAction {
	return s.waitsFor(sentence.NewCondition(sentence.ShouldNotHave, true, resolutions))
}

// HappenToBeAnyOf waits for the answer to satisfy at least one resolution. It is the plural form: it reads "happen to be any of: ...".
func (s *StopBuilder[T]) HappenToBeAnyOf(resolutions ...screenplay.Resolution[T]) *StopAction {
	return s.waitsFor(sentence.NewCondition(sentence.ShouldAnyOf, true, resolutions))
}

// waitsFor creates the action stopping until the answer satisfies the clause.
func (s *StopBuilder[T]) waitsFor(clause sentence.Clause[T]) *StopAction {
	condition := sentence.NewTerm(s.question, clause)

	return &StopAction{
		condition: &awaitedCondition{term: condition},
		term:      condition,
	}
}

// String describes the action.
func (a *StopAction) String() string {
	if a.condition == nil {
		return `stop until the "enter" key is pressed`
	}
	return fmt.Sprintf("stop until %s", a.term)
}

// PerformAs performs the task or the action as the provided actor.
func (a *StopAction) PerformAs(actor *screenplay.Actor) error {
	if a.condition == nil {
		return a.waitForAnInput(actor)
	}

	// A condition that was built wrongly can never be seen: it is reported at
	// once rather than after waiting for the whole timeout.
	err := a.term.MetBy(actor)
	if err == nil {
		return nil
	}

	if errors.Is(err, sentence.ErrInvalidTerm) || errors.Is(err, screenplay.ErrInvalidResolution) {
		return err
	}

	if err = actor.AttemptsTo(Eventually(a.condition)); err != nil {
		return a.neverSeen(actor)
	}
	return nil
}

// waitForAnInput blocks until a line is read from the reader, or until the
// actor's context ends. The prompt is written to the standard error of the
// process, where the person debugging the test reads it; the narrator already
// states that the actor stops.
func (a *StopAction) waitForAnInput(actor *screenplay.Actor) error {
	fmt.Fprintln(os.Stderr, "the actor stops, waiting for your go (press enter to continue)")

	read := make(chan error, 1)
	go func() {
		_, err := fmt.Fscanln(a.reader)
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

// awaitedCondition is what the actor keeps trying to see while it is stopped.
type awaitedCondition struct {
	term *sentence.Term
}

// String describes the condition.
func (c *awaitedCondition) String() string {
	return "see that " + c.term.String()
}

// PerformAs sees whether the condition is met.
func (c *awaitedCondition) PerformAs(actor *screenplay.Actor) error {
	return c.term.MetBy(actor)
}

// neverSeen reports that the actor stopped for the whole timeout without
// seeing the condition, with the reason of the last try.
func (a *StopAction) neverSeen(actor *screenplay.Actor) error {
	err := a.term.MetBy(actor)
	if err == nil {
		// Met now, but only after the timeout: the actor still failed to see it in time.
		return fmt.Errorf("%s stopped for %s without seeing that %s in time", actor.Name(), actor.Timeout(), a.term)
	}

	var unmet *sentence.UnmetError
	if errors.As(err, &unmet) {
		return fmt.Errorf("%s stopped for %s without seeing that %w", actor.Name(), actor.Timeout(), err)
	}

	return fmt.Errorf("%s stopped for %s without seeing that %s: %w", actor.Name(), actor.Timeout(), a.term, err)
}

// StopAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*StopAction)(nil)
