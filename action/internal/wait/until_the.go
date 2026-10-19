package wait

import (
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/v2/action/internal/retry"
	"github.com/grandper/go-screenplay/v2/internal/sentence"
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/grandper/go-screenplay/v2/timing"
)

// UntilThe makes an actor wait until the answer to a question happens to be as
// stated, for no longer than the timeout of the actor's production. Close it
// with a verb and one or more resolutions:
//
//	wait.UntilThe(logs).HappenTo(contain.Substring("ready"))
func UntilThe[T any](question screenplay.Question[T]) *UntilTheBuilder[T] {
	return &UntilTheBuilder[T]{question: question}
}

// UntilTheBuilder builds the condition an actor waits for: it holds the
// question, and waits for the verb that closes the condition.
type UntilTheBuilder[T any] struct {
	question screenplay.Question[T]
}

// HappensTo waits for the answer to satisfy every resolution. It reads "happens to ...".
func (s *UntilTheBuilder[T]) HappensTo(resolutions ...screenplay.Resolution[T]) screenplay.Performable {
	return s.waitsFor(sentence.NewCondition(sentence.Should, false, resolutions))
}

// HappensToBe waits for the answer to satisfy every resolution. It reads "happens to be ...".
func (s *UntilTheBuilder[T]) HappensToBe(resolutions ...screenplay.Resolution[T]) screenplay.Performable {
	return s.waitsFor(sentence.NewCondition(sentence.ShouldBe, false, resolutions))
}

// HappensToHave waits for the answer to satisfy every resolution. It reads "happens to have ...".
func (s *UntilTheBuilder[T]) HappensToHave(resolutions ...screenplay.Resolution[T]) screenplay.Performable {
	return s.waitsFor(sentence.NewCondition(sentence.ShouldHave, false, resolutions))
}

// DoesNotHappenTo waits for the answer to satisfy none of the resolutions. It reads "does not happen to ...".
func (s *UntilTheBuilder[T]) DoesNotHappenTo(resolutions ...screenplay.Resolution[T]) screenplay.Performable {
	return s.waitsFor(sentence.NewCondition(sentence.ShouldNot, false, resolutions))
}

// DoesNotHappenToBe waits for the answer to satisfy none of the resolutions. It reads "does not happen to be ...".
func (s *UntilTheBuilder[T]) DoesNotHappenToBe(resolutions ...screenplay.Resolution[T]) screenplay.Performable {
	return s.waitsFor(sentence.NewCondition(sentence.ShouldNotBe, false, resolutions))
}

// DoesNotHappenToHave waits for the answer to satisfy none of the resolutions. It reads "does not happen to have ...".
func (s *UntilTheBuilder[T]) DoesNotHappenToHave(resolutions ...screenplay.Resolution[T]) screenplay.Performable {
	return s.waitsFor(sentence.NewCondition(sentence.ShouldNotHave, false, resolutions))
}

// HappensToBeAnyOf waits for the answer to satisfy at least one resolution. It reads "happens to be any of: ...".
func (s *UntilTheBuilder[T]) HappensToBeAnyOf(resolutions ...screenplay.Resolution[T]) screenplay.Performable {
	return s.waitsFor(sentence.NewCondition(sentence.ShouldAnyOf, false, resolutions))
}

// HappenTo waits for the answer to satisfy every resolution. It is the plural form: it reads "happen to ...".
func (s *UntilTheBuilder[T]) HappenTo(resolutions ...screenplay.Resolution[T]) screenplay.Performable {
	return s.waitsFor(sentence.NewCondition(sentence.Should, true, resolutions))
}

// HappenToBe waits for the answer to satisfy every resolution. It is the plural form: it reads "happen to be ...".
func (s *UntilTheBuilder[T]) HappenToBe(resolutions ...screenplay.Resolution[T]) screenplay.Performable {
	return s.waitsFor(sentence.NewCondition(sentence.ShouldBe, true, resolutions))
}

// HappenToHave waits for the answer to satisfy every resolution. It is the plural form: it reads "happen to have ...".
func (s *UntilTheBuilder[T]) HappenToHave(resolutions ...screenplay.Resolution[T]) screenplay.Performable {
	return s.waitsFor(sentence.NewCondition(sentence.ShouldHave, true, resolutions))
}

// DoNotHappenTo waits for the answer to satisfy none of the resolutions. It is the plural form: it reads "do not happen to ...".
func (s *UntilTheBuilder[T]) DoNotHappenTo(resolutions ...screenplay.Resolution[T]) screenplay.Performable {
	return s.waitsFor(sentence.NewCondition(sentence.ShouldNot, true, resolutions))
}

// DoNotHappenToBe waits for the answer to satisfy none of the resolutions. It is the plural form: it reads "do not happen to be ...".
func (s *UntilTheBuilder[T]) DoNotHappenToBe(resolutions ...screenplay.Resolution[T]) screenplay.Performable {
	return s.waitsFor(sentence.NewCondition(sentence.ShouldNotBe, true, resolutions))
}

// DoNotHappenToHave waits for the answer to satisfy none of the resolutions. It is the plural form: it reads "do not happen to have ...".
func (s *UntilTheBuilder[T]) DoNotHappenToHave(resolutions ...screenplay.Resolution[T]) screenplay.Performable {
	return s.waitsFor(sentence.NewCondition(sentence.ShouldNotHave, true, resolutions))
}

// HappenToBeAnyOf waits for the answer to satisfy at least one resolution. It is the plural form: it reads "happen to be any of: ...".
func (s *UntilTheBuilder[T]) HappenToBeAnyOf(resolutions ...screenplay.Resolution[T]) screenplay.Performable {
	return s.waitsFor(sentence.NewCondition(sentence.ShouldAnyOf, true, resolutions))
}

// waitsFor creates the action waiting until the answer satisfies the clause.
func (s *UntilTheBuilder[T]) waitsFor(clause sentence.Clause[T]) *untilTheAction {
	return &untilTheAction{
		term: sentence.NewTerm(s.question, clause),
	}
}

// untilTheAction is an action that makes an actor wait until a condition is met.
type untilTheAction struct {
	term *sentence.Term
}

// String describes the action.
func (a *untilTheAction) String() string {
	return fmt.Sprintf("wait until %s", a.term)
}

// PerformAs keeps checking the condition until it is met, or until the
// timeout of the actor's production.
func (a *untilTheAction) PerformAs(actor *screenplay.Actor) error {
	// A condition that was built wrongly can never be seen: it is reported at
	// once rather than after waiting for the whole timeout.
	err := a.term.MetBy(actor)
	if err == nil {
		return nil
	}

	if errors.Is(err, sentence.ErrInvalidTerm) || errors.Is(err, screenplay.ErrInvalidResolution) {
		return err
	}

	window := retry.WindowFor(actor, timing.NewWindow(0, 0))
	if retry.Check(window) != nil {
		return a.neverSeen(actor)
	}

	if _, err = retry.Until(actor, window, func() error { return a.term.MetBy(actor) }); err != nil {
		return a.neverSeen(actor)
	}

	return nil
}

// neverSeen reports that the actor waited for the whole timeout without
// seeing the condition, with the reason of the last try.
func (a *untilTheAction) neverSeen(actor *screenplay.Actor) error {
	err := a.term.MetBy(actor)
	if err == nil {
		// Met now, but only after the timeout: the actor still failed to see it in time.
		return fmt.Errorf("%s waited for %s without seeing that %s in time", actor.Name(), actor.Timeout(), a.term)
	}

	var unmet *sentence.UnmetError
	if errors.As(err, &unmet) {
		return fmt.Errorf("%s waited for %s without seeing that %w", actor.Name(), actor.Timeout(), err)
	}

	return fmt.Errorf("%s waited for %s without seeing that %s: %w", actor.Name(), actor.Timeout(), a.term, err)
}

// untilTheAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*untilTheAction)(nil)
