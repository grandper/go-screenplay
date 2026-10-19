package term

import (
	"github.com/grandper/go-screenplay/v2/internal/sentence"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// That starts a term about the answer to a question. Close it with a verb and
// one or more resolutions:
//
//	term.That(statusCode).Should(be.EqualTo(200))
//	term.That(logs).ShouldNot(contain.Substring("panic"))
//
// The resolutions share the answer type of the question, so a mismatch does
// not compile.
func That[T any](question screenplay.Question[T]) *TheSubject[T] {
	return &TheSubject[T]{
		question: question,
	}
}

// TheSubject is the question a term is about, waiting for its verb.
type TheSubject[T any] struct {
	question screenplay.Question[T]
}

// Should states that the answer satisfies every resolution. It reads "should ...".
func (s *TheSubject[T]) Should(resolutions ...screenplay.Resolution[T]) *sentence.Term {
	return s.states(sentence.Should, resolutions)
}

// ShouldBe states that the answer satisfies every resolution. It reads "should be ...".
func (s *TheSubject[T]) ShouldBe(resolutions ...screenplay.Resolution[T]) *sentence.Term {
	return s.states(sentence.ShouldBe, resolutions)
}

// ShouldHave states that the answer satisfies every resolution. It reads "should have ...".
func (s *TheSubject[T]) ShouldHave(resolutions ...screenplay.Resolution[T]) *sentence.Term {
	return s.states(sentence.ShouldHave, resolutions)
}

// ShouldNot states that the answer satisfies none of the resolutions. It reads "should not ...".
func (s *TheSubject[T]) ShouldNot(resolutions ...screenplay.Resolution[T]) *sentence.Term {
	return s.states(sentence.ShouldNot, resolutions)
}

// ShouldNotBe states that the answer satisfies none of the resolutions. It reads "should not be ...".
func (s *TheSubject[T]) ShouldNotBe(resolutions ...screenplay.Resolution[T]) *sentence.Term {
	return s.states(sentence.ShouldNotBe, resolutions)
}

// ShouldNotHave states that the answer satisfies none of the resolutions. It reads "should not have ...".
func (s *TheSubject[T]) ShouldNotHave(resolutions ...screenplay.Resolution[T]) *sentence.Term {
	return s.states(sentence.ShouldNotHave, resolutions)
}

// ShouldAnyOf states that the answer satisfies at least one resolution. It reads "should be any of: ...".
func (s *TheSubject[T]) ShouldAnyOf(resolutions ...screenplay.Resolution[T]) *sentence.Term {
	return s.states(sentence.ShouldAnyOf, resolutions)
}

// states creates the term stating that the answer satisfies the resolutions.
func (s *TheSubject[T]) states(verb sentence.Verb, resolutions []screenplay.Resolution[T]) *sentence.Term {
	return sentence.NewTerm(s.question, sentence.NewClause(verb, resolutions))
}
