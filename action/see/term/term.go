package term

import (
	impl "github.com/grandper/go-screenplay/v2/action/internal/see/term"
	"github.com/grandper/go-screenplay/v2/internal/sentence"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// ErrInvalidTerm is wrapped by the error an action returns when a sentence was
// built wrongly: a verb without a resolution, a nil resolution, a nil
// question, or a nil term. Nothing panics when the sentence is written; the
// mistake is reported when the actor performs the action.
var ErrInvalidTerm = sentence.ErrInvalidTerm

// Term is a statement about the answer to a question: "the status code should
// be equal to 200". It is built with That, and handed to see.AllOfTheseMet,
// see.AtLeastOneOfTheseMet and see.NoneOfTheseMet.
type Term = sentence.Term

// That starts a term about the answer to a question. Close it with a verb and
// one or more resolutions:
//
//	term.That(statusCode).Should(be.EqualTo(200))
//	term.That(logs).ShouldNot(contain.Substring("panic"))
//
// The resolutions share the answer type of the question, so a mismatch does
// not compile.
func That[T any](question screenplay.Question[T]) *impl.TheSubject[T] {
	return impl.That(question)
}
