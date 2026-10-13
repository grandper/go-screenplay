package see

import (
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/v2/internal/sentence"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// The starts an assertion about the answer to a question. Close it with a verb
// and one or more resolutions:
//
//	see.The(statusCode).Should(be.EqualTo(200))
//	see.The(logs).ShouldNot(contain.Substring("panic"))
//	see.The(customers).ShouldHave(length.Exactly[[]string](2))
//
// The resolutions share the answer type of the question, so a mismatch does
// not compile.
func The[T any](question screenplay.Question[T]) *TheSubject[T] {
	return &TheSubject[T]{
		question: question,
	}
}

// TheSubject is the question see.The is about, waiting for its verb.
type TheSubject[T any] struct {
	question screenplay.Question[T]
}

// Should sees that the answer satisfies every resolution. It reads "should ...".
func (s *TheSubject[T]) Should(resolutions ...screenplay.Resolution[T]) *TheAction {
	return s.sees(sentence.Should, resolutions)
}

// ShouldBe sees that the answer satisfies every resolution. It reads "should be ...".
func (s *TheSubject[T]) ShouldBe(resolutions ...screenplay.Resolution[T]) *TheAction {
	return s.sees(sentence.ShouldBe, resolutions)
}

// ShouldHave sees that the answer satisfies every resolution. It reads "should have ...".
func (s *TheSubject[T]) ShouldHave(resolutions ...screenplay.Resolution[T]) *TheAction {
	return s.sees(sentence.ShouldHave, resolutions)
}

// ShouldNot sees that the answer satisfies none of the resolutions. It reads "should not ...".
func (s *TheSubject[T]) ShouldNot(resolutions ...screenplay.Resolution[T]) *TheAction {
	return s.sees(sentence.ShouldNot, resolutions)
}

// ShouldNotBe sees that the answer satisfies none of the resolutions. It reads "should not be ...".
func (s *TheSubject[T]) ShouldNotBe(resolutions ...screenplay.Resolution[T]) *TheAction {
	return s.sees(sentence.ShouldNotBe, resolutions)
}

// ShouldNotHave sees that the answer satisfies none of the resolutions. It reads "should not have ...".
func (s *TheSubject[T]) ShouldNotHave(resolutions ...screenplay.Resolution[T]) *TheAction {
	return s.sees(sentence.ShouldNotHave, resolutions)
}

// ShouldAnyOf sees that the answer satisfies at least one resolution. It reads "should be any of: ...".
func (s *TheSubject[T]) ShouldAnyOf(resolutions ...screenplay.Resolution[T]) *TheAction {
	return s.sees(sentence.ShouldAnyOf, resolutions)
}

// sees creates the action seeing that the answer satisfies the resolutions.
func (s *TheSubject[T]) sees(verb sentence.Verb, resolutions []screenplay.Resolution[T]) *TheAction {
	return newTheAction(sentence.NewTerm(s.question, sentence.NewClause(verb, resolutions)))
}

// newTheAction creates the action seeing that a single term is met.
func newTheAction(t *sentence.Term) *TheAction {
	return &TheAction{
		term: t,
	}
}

// TheAction sees whether a term about the answer to a question is met.
type TheAction struct {
	term *sentence.Term
}

// String describes the action.
func (a *TheAction) String() string {
	return "see that " + a.term.String()
}

// PerformAs performs the task or the action as the provided actor.
func (a *TheAction) PerformAs(actor *screenplay.Actor) error {
	err := a.term.MetBy(actor)

	var unmet *sentence.UnmetError
	if errors.As(err, &unmet) {
		return fmt.Errorf("%s failed to see that %w", actor.Name(), err)
	}

	return err
}

// TheAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*TheAction)(nil)
