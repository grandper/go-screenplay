package see

import (
	"errors"
	"fmt"
	"strings"

	"github.com/grandper/go-screenplay/v2/internal/sentence"
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/grandper/go-screenplay/v2/term"
)

// quantifier tells how many of the terms of a TheseAction should be met.
type quantifier int

const (
	all quantifier = iota
	atLeastOne
	none
)

// String describes the quantifier as it reads in a sentence.
func (q quantifier) String() string {
	switch q {
	case all:
		return "all of these terms are met"
	case atLeastOne:
		return "at least one of these terms is met"
	case none:
		return "none of these terms is met"
	default:
		return "all of these terms are met"
	}
}

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
func AllOfTheseMet(terms ...*term.Term) *TheseAction {
	return these(all, "see.AllOfTheseMet", terms)
}

// AtLeastOneOfTheseMet creates an action to see that at least one term is met.
// The actor checks the terms in order and stops at the first one that is met;
// a term whose question cannot be answered counts as not met.
func AtLeastOneOfTheseMet(terms ...*term.Term) *TheseAction {
	return these(atLeastOne, "see.AtLeastOneOfTheseMet", terms)
}

// NoneOfTheseMet creates an action to see that no term is met. The action
// fails when a question cannot be answered: a term that could not be checked
// is not known to be unmet.
func NoneOfTheseMet(terms ...*term.Term) *TheseAction {
	return these(none, "see.NoneOfTheseMet", terms)
}

// these creates the action judging several terms. A nil term, or a term whose
// sentence was built wrongly, can never be judged: the action remembers it and
// fails with term.ErrInvalidTerm when it is performed, before judging any
// term.
func these(q quantifier, name string, terms []*term.Term) *TheseAction {
	action := &TheseAction{
		quantifier: q,
		terms:      make([]*term.Term, 0, len(terms)),
	}

	for position, t := range terms {
		if t == nil {
			action.err = fmt.Errorf("%w: the term at position %d given to %s is nil",
				term.ErrInvalidTerm, position, name)

			continue
		}

		if mistake := t.Mistake(); mistake != nil && action.err == nil {
			action.err = mistake
		}

		action.terms = append(action.terms, t)
	}

	return action
}

// TheseAction sees whether all, at least one, or none of several terms are met.
type TheseAction struct {
	quantifier quantifier
	terms      []*term.Term
	err        error // why the action can never be performed, when it was built wrongly
}

// String describes the action. The terms are not repeated here: each of them
// is narrated on its own line when the actor checks it.
func (a *TheseAction) String() string {
	return fmt.Sprintf("see that %s", a.quantifier)
}

// PerformAs performs the task or the action as the provided actor.
func (a *TheseAction) PerformAs(actor *screenplay.Actor) error {
	if a.err != nil {
		return a.err
	}

	switch a.quantifier {
	case atLeastOne:
		return a.seeAtLeastOne(actor)
	case none:
		return a.seeNone(actor)
	case all:
		return a.seeAll(actor)
	default:
		return a.seeAll(actor)
	}
}

func (a *TheseAction) seeAll(actor *screenplay.Actor) error {
	for _, t := range a.terms {
		if err := actor.AttemptsTo(newTheAction(t)); err != nil {
			return err
		}
	}

	return nil
}

func (a *TheseAction) seeAtLeastOne(actor *screenplay.Actor) error {
	failures := make([]error, 0, len(a.terms))
	reasons := make([]string, 0, len(a.terms))

	for _, t := range a.terms {
		err := actor.AttemptsTo(newTheAction(t))
		if err == nil {
			return nil
		}

		if isAMistake(err) {
			return err
		}

		failures = append(failures, err)
		reasons = append(reasons, "\n- "+whyNotMet(t, err))
	}

	return &notSeenError{
		message: fmt.Sprintf(
			"%s failed to see that at least one of these terms was met:%s",
			actor.Name(),
			strings.Join(reasons, ""),
		),
		failures: failures,
	}
}

// whyNotMet describes why a term is not met: the part of its sentence that is
// not satisfied, or the error of the question that could not be answered.
func whyNotMet(t *term.Term, err error) string {
	var unmet *sentence.UnmetError
	if errors.As(err, &unmet) {
		return unmet.Error()
	}

	return fmt.Sprintf("%s: %s", t, err)
}

// notSeenError is returned when at least one of several terms should be met
// and none is. It keeps the failure of every term.
type notSeenError struct {
	message  string
	failures []error
}

// Error lists the terms and why each of them is not met.
func (e *notSeenError) Error() string {
	return e.message
}

// Unwrap returns the failure of every term.
func (e *notSeenError) Unwrap() []error {
	return e.failures
}

func (a *TheseAction) seeNone(actor *screenplay.Actor) error {
	for _, t := range a.terms {
		if err := actor.AttemptsTo(&termNotMet{term: t}); err != nil {
			return err
		}
	}

	return nil
}

// termNotMet sees that a single term is not met. It is how the actor checks
// each term of NoneOfTheseMet, so that each of them is narrated on its own
// line.
type termNotMet struct {
	term *term.Term
}

// String describes the action.
func (a *termNotMet) String() string {
	return fmt.Sprintf("see that \"%s\" is not met", a.term)
}

// PerformAs performs the task or the action as the provided actor.
func (a *termNotMet) PerformAs(actor *screenplay.Actor) error {
	verdict, err := a.term.JudgedBy(actor)
	if err != nil {
		return err
	}

	if verdict.Met {
		return fmt.Errorf("%s failed to see that none of these terms was met: \"%s\" is met: got %s",
			actor.Name(), a.term, verdict.Got)
	}

	return nil
}

// TheseAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*TheseAction)(nil)

// ThatAllOfThe starts an assertion that the answers to all the questions
// satisfy the same resolutions. Every question must share the same answer
// type; use AllOfTheseMet when they differ.
//
//	see.ThatAllOfThe(profileWidget, loginForm).Should(contain.Substring("Adam"))
func ThatAllOfThe[T any](questions ...screenplay.Question[T]) *TheSubjects[T] {
	return &TheSubjects[T]{quantifier: all, constructor: "see.ThatAllOfThe", questions: questions}
}

// ThatAnyOfThe starts an assertion that the answer to at least one of the
// questions satisfies the resolutions. Every question must share the same
// answer type; use AtLeastOneOfTheseMet when they differ.
//
//	see.ThatAnyOfThe(profileWidget, loginForm).Should(contain.Substring("Adam"))
func ThatAnyOfThe[T any](questions ...screenplay.Question[T]) *TheSubjects[T] {
	return &TheSubjects[T]{quantifier: atLeastOne, constructor: "see.ThatAnyOfThe", questions: questions}
}

// ThatNoneOfThe starts an assertion that the answer to none of the questions
// satisfies the resolutions. Every question must share the same answer type;
// use NoneOfTheseMet when they differ.
//
//	see.ThatNoneOfThe(profileWidget, loginForm).Should(contain.Substring("error"))
func ThatNoneOfThe[T any](questions ...screenplay.Question[T]) *TheSubjects[T] {
	return &TheSubjects[T]{quantifier: none, constructor: "see.ThatNoneOfThe", questions: questions}
}

// TheSubjects are the questions see.ThatAllOfThe, see.ThatAnyOfThe and
// see.ThatNoneOfThe are about, waiting for their verb.
type TheSubjects[T any] struct {
	quantifier  quantifier
	constructor string
	questions   []screenplay.Question[T]
}

// Should sees that the answers satisfy every resolution. It reads "should ...".
func (s *TheSubjects[T]) Should(resolutions ...screenplay.Resolution[T]) *TheseAction {
	return s.see(sentence.Should, resolutions)
}

// ShouldBe sees that the answers satisfy every resolution. It reads "should be ...".
func (s *TheSubjects[T]) ShouldBe(resolutions ...screenplay.Resolution[T]) *TheseAction {
	return s.see(sentence.ShouldBe, resolutions)
}

// ShouldHave sees that the answers satisfy every resolution. It reads "should have ...".
func (s *TheSubjects[T]) ShouldHave(resolutions ...screenplay.Resolution[T]) *TheseAction {
	return s.see(sentence.ShouldHave, resolutions)
}

// ShouldNot sees that the answers satisfy none of the resolutions. It reads "should not ...".
func (s *TheSubjects[T]) ShouldNot(resolutions ...screenplay.Resolution[T]) *TheseAction {
	return s.see(sentence.ShouldNot, resolutions)
}

// ShouldNotBe sees that the answers satisfy none of the resolutions. It reads "should not be ...".
func (s *TheSubjects[T]) ShouldNotBe(resolutions ...screenplay.Resolution[T]) *TheseAction {
	return s.see(sentence.ShouldNotBe, resolutions)
}

// ShouldNotHave sees that the answers satisfy none of the resolutions. It reads "should not have ...".
func (s *TheSubjects[T]) ShouldNotHave(resolutions ...screenplay.Resolution[T]) *TheseAction {
	return s.see(sentence.ShouldNotHave, resolutions)
}

// ShouldAnyOf sees that the answers satisfy at least one resolution. It reads "should be any of: ...".
func (s *TheSubjects[T]) ShouldAnyOf(resolutions ...screenplay.Resolution[T]) *TheseAction {
	return s.see(sentence.ShouldAnyOf, resolutions)
}

// see creates the action seeing that the answers satisfy the resolutions.
func (s *TheSubjects[T]) see(verb sentence.Verb, resolutions []screenplay.Resolution[T]) *TheseAction {
	clause := sentence.NewClause(verb, resolutions)

	terms := make([]*term.Term, len(s.questions))
	for i, question := range s.questions {
		terms[i] = sentence.NewTerm(question, clause)
	}

	return these(s.quantifier, s.constructor, terms)
}

// isAMistake reports whether an error comes from a sentence or a resolution
// that was built wrongly, rather than from a term that is not met.
func isAMistake(err error) bool {
	return errors.Is(err, term.ErrInvalidTerm) || errors.Is(err, screenplay.ErrInvalidResolution)
}
