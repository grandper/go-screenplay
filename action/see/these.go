package see

import (
	impl "github.com/grandper/go-screenplay/v2/action/internal/see"
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/grandper/go-screenplay/v2/term"
)

// AllOfTheseMet creates an action to see that every term is met. The terms may
// be about questions of different answer types:
//
//	see.AllOfTheseMet(
//		term.That(statusCode).Should(be.EqualTo(200)),
//		term.That(logs).Should(contain.Substring("order placed")),
//	)
//
// The actor checks the terms in order and stops at the first one that is not
// met.
func AllOfTheseMet(terms ...*term.Term) screenplay.Performable {
	return impl.AllOfTheseMet(terms...)
}

// AtLeastOneOfTheseMet creates an action to see that at least one term is met.
// The actor checks the terms in order and stops at the first one that is met;
// a term whose question cannot be answered counts as not met.
func AtLeastOneOfTheseMet(terms ...*term.Term) screenplay.Performable {
	return impl.AtLeastOneOfTheseMet(terms...)
}

// NoneOfTheseMet creates an action to see that no term is met. The action
// fails when a question cannot be answered: a term that could not be checked
// is not known to be unmet.
func NoneOfTheseMet(terms ...*term.Term) screenplay.Performable {
	return impl.NoneOfTheseMet(terms...)
}

// ThatAllOfThe starts an assertion that the answers to all the questions
// satisfy the same resolutions. Every question must share the same answer
// type; use AllOfTheseMet when they differ.
//
//	see.ThatAllOfThe(profileWidget, loginForm).Should(contain.Substring("Adam"))
func ThatAllOfThe[T any](questions ...screenplay.Question[T]) *impl.TheSubjects[T] {
	return impl.ThatAllOfThe(questions...)
}

// ThatAnyOfThe starts an assertion that the answer to at least one of the
// questions satisfies the resolutions. Every question must share the same
// answer type; use AtLeastOneOfTheseMet when they differ.
//
//	see.ThatAnyOfThe(profileWidget, loginForm).Should(contain.Substring("Adam"))
func ThatAnyOfThe[T any](questions ...screenplay.Question[T]) *impl.TheSubjects[T] {
	return impl.ThatAnyOfThe(questions...)
}

// ThatNoneOfThe starts an assertion that the answer to none of the questions
// satisfies the resolutions. Every question must share the same answer type;
// use NoneOfTheseMet when they differ.
//
//	see.ThatNoneOfThe(profileWidget, loginForm).Should(contain.Substring("error"))
func ThatNoneOfThe[T any](questions ...screenplay.Question[T]) *impl.TheSubjects[T] {
	return impl.ThatNoneOfThe(questions...)
}
