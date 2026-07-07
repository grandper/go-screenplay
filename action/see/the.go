package see

import (
	"fmt"

	"github.com/grandper/go-screenplay/screenplay"
)

// The starts an assertion about the answer to a question. The verb (and its negation) is
// chosen with a method: Is, IsNot, Are, AreNot, Does, DoesNot, Do, DoNot, Has, HasNot,
// Have, HaveNot, Had, HadNot, Was, WasNot, Were, WereNot.
func The(question screenplay.Question) *Assertion {
	return &Assertion{
		question: question,
	}
}

// Assertion picks the verb and polarity of an assertion about a question's answer.
type Assertion struct {
	question screenplay.Question
}

// Is asserts that the answer matches the resolution.
func (a *Assertion) Is(resolution screenplay.Resolution) *TheAction {
	return a.check("is", false, resolution)
}

// IsNot asserts that the answer does not match the resolution.
func (a *Assertion) IsNot(resolution screenplay.Resolution) *TheAction {
	return a.check("is not", true, resolution)
}

// Are asserts that the answer matches the resolution.
func (a *Assertion) Are(resolution screenplay.Resolution) *TheAction {
	return a.check("are", false, resolution)
}

// AreNot asserts that the answer does not match the resolution.
func (a *Assertion) AreNot(resolution screenplay.Resolution) *TheAction {
	return a.check("are not", true, resolution)
}

// Does asserts that the answer matches the resolution.
func (a *Assertion) Does(resolution screenplay.Resolution) *TheAction {
	return a.check("does", false, resolution)
}

// DoesNot asserts that the answer does not match the resolution.
func (a *Assertion) DoesNot(resolution screenplay.Resolution) *TheAction {
	return a.check("does not", true, resolution)
}

// Do asserts that the answer matches the resolution.
func (a *Assertion) Do(resolution screenplay.Resolution) *TheAction {
	return a.check("do", false, resolution)
}

// DoNot asserts that the answer does not match the resolution.
func (a *Assertion) DoNot(resolution screenplay.Resolution) *TheAction {
	return a.check("do not", true, resolution)
}

// Has asserts that the answer matches the resolution.
func (a *Assertion) Has(resolution screenplay.Resolution) *TheAction {
	return a.check("has", false, resolution)
}

// HasNot asserts that the answer does not match the resolution.
func (a *Assertion) HasNot(resolution screenplay.Resolution) *TheAction {
	return a.check("has not", true, resolution)
}

// Have asserts that the answer matches the resolution.
func (a *Assertion) Have(resolution screenplay.Resolution) *TheAction {
	return a.check("have", false, resolution)
}

// HaveNot asserts that the answer does not match the resolution.
func (a *Assertion) HaveNot(resolution screenplay.Resolution) *TheAction {
	return a.check("have not", true, resolution)
}

// Had asserts that the answer matches the resolution.
func (a *Assertion) Had(resolution screenplay.Resolution) *TheAction {
	return a.check("had", false, resolution)
}

// HadNot asserts that the answer does not match the resolution.
func (a *Assertion) HadNot(resolution screenplay.Resolution) *TheAction {
	return a.check("had not", true, resolution)
}

// Was asserts that the answer matches the resolution.
func (a *Assertion) Was(resolution screenplay.Resolution) *TheAction {
	return a.check("was", false, resolution)
}

// WasNot asserts that the answer does not match the resolution.
func (a *Assertion) WasNot(resolution screenplay.Resolution) *TheAction {
	return a.check("was not", true, resolution)
}

// Were asserts that the answer matches the resolution.
func (a *Assertion) Were(resolution screenplay.Resolution) *TheAction {
	return a.check("were", false, resolution)
}

// WereNot asserts that the answer does not match the resolution.
func (a *Assertion) WereNot(resolution screenplay.Resolution) *TheAction {
	return a.check("were not", true, resolution)
}

func (a *Assertion) check(verb string, negated bool, resolution screenplay.Resolution) *TheAction {
	return &TheAction{
		question:   a.question,
		verb:       verb,
		negated:    negated,
		resolution: resolution,
	}
}

// TheAction sees whether the answer to a question matches the resolution (or, when negated,
// does not match it).
type TheAction struct {
	question   screenplay.Question
	verb       string
	negated    bool
	resolution screenplay.Resolution
}

// String describes the action.
func (a *TheAction) String() string {
	return fmt.Sprintf("see if the %s %s %s", a.question, a.verb, a.resolution)
}

// PerformAs performs the task or the action as the provided actor.
func (a *TheAction) PerformAs(actor *screenplay.Actor) error {
	value, err := a.question.AnsweredBy(actor)
	if err != nil {
		return err
	}

	matched, err := a.resolution.Resolve()(value)
	if err != nil {
		return fmt.Errorf("an error occurred when %s attempted to see %s: %w", actor.Name(), a.question, err)
	}

	if matched == a.negated {
		return fmt.Errorf("%s failed to see that the %s %s %s", actor.Name(), a.question, a.verb, a.resolution)
	}

	return nil
}

// TheAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*TheAction)(nil)
