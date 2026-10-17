// Package condition holds the condition that guards a branch of the flow
// actions: proceed.To(...).IfThe(...) and choose.To(...).WhenThe(...).
package condition

import (
	"github.com/grandper/go-screenplay/v2/internal/sentence"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// Condition reports whether a guarded branch should run, and describes itself.
// Keyword introduces the description: "if", "when" or "unless".
type Condition struct {
	Keyword     string
	Description string
	holds       func(theActor *screenplay.Actor) (bool, error)
}

// Boolean builds a condition from a boolean expression, introduced by "if".
func Boolean(cond bool) Condition {
	return Condition{
		Keyword:     "if",
		Description: "the condition holds",
		holds:       func(*screenplay.Actor) (bool, error) { return cond, nil },
	}
}

// Term builds a condition, introduced by "if", that holds when a term is met.
func Term(t *sentence.Term) Condition {
	return Condition{
		Keyword:     "if",
		Description: t.String(),
		holds: func(actor *screenplay.Actor) (bool, error) {
			verdict, err := t.JudgedBy(actor)
			if err != nil {
				return false, err
			}

			return verdict.Met, nil
		},
	}
}

// IntroducedBy returns the same condition, introduced by another keyword.
func (c Condition) IntroducedBy(keyword string) Condition {
	c.Keyword = keyword

	return c
}

// Negated returns the condition, introduced by "unless", that holds when the
// receiver does not.
func (c Condition) Negated() Condition {
	return Condition{
		Keyword:     "unless",
		Description: c.Description,
		holds: func(theActor *screenplay.Actor) (bool, error) {
			held, err := c.holds(theActor)

			return !held, err
		},
	}
}

// HoldsFor reports whether the condition holds for an actor.
func (c Condition) HoldsFor(theActor *screenplay.Actor) (bool, error) {
	return c.holds(theActor)
}
