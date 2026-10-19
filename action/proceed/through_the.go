package proceed

import (
	"fmt"
	"strings"

	"github.com/grandper/go-screenplay/v2/internal/verdict"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// ThroughThe starts an action that goes through the items the answer to a
// question holds. Close it with the body performed for each item:
//
//	proceed.ThroughThe(temporaryFiles).AndForEach(remove.TheFile)
func ThroughThe[T any](question screenplay.Question[[]T]) *ThroughTheBuilder[T] {
	if question == nil {
		return &ThroughTheBuilder[T]{
			description: "proceed through the <nil>",
			items: func(*screenplay.Actor) ([]T, error) {
				return nil, fmt.Errorf("%w: the question is nil", screenplay.ErrInvalidResolution)
			},
		}
	}

	return &ThroughTheBuilder[T]{
		description: fmt.Sprintf("proceed through the %s", question),
		items:       question.AnsweredBy,
	}
}

// Through starts an action that goes through the items written in the
// statement. Close it with the body performed for each item:
//
//	proceed.Through(endpoints...).AndForEach(CheckTheEndpoint)
func Through[T any](items ...T) *ThroughTheBuilder[T] {
	descriptions := make([]string, 0, len(items))
	for _, item := range items {
		descriptions = append(descriptions, verdict.Describe(item))
	}

	return &ThroughTheBuilder[T]{
		description: strings.TrimSpace("proceed through " + strings.Join(descriptions, ", ")),
		items:       func(*screenplay.Actor) ([]T, error) { return items, nil },
	}
}

// ThroughTheBuilder holds the items of a proceed.ThroughThe or proceed.Through
// action, and waits for the body performed for each of them.
type ThroughTheBuilder[T any] struct {
	description string
	items       func(theActor *screenplay.Actor) ([]T, error)
}

// AndForEach sets the body performed for each item: a function from the item
// to an action. A function that returns a concrete action works without a
// wrapper:
//
//	proceed.ThroughThe(temporaryFiles).AndForEach(remove.TheFile)
func (b *ThroughTheBuilder[T]) AndForEach[P screenplay.Performable](body func(item T) P) *ThroughTheAction[T] {
	action := &ThroughTheAction[T]{description: b.description, items: b.items}
	if body != nil {
		action.body = func(item T) screenplay.Performable { return body(item) }
	}

	return action
}

// ThroughTheAction is an action that performs its body for each item, in
// order, and stops at the first failure.
type ThroughTheAction[T any] struct {
	description string
	items       func(theActor *screenplay.Actor) ([]T, error)
	body        func(item T) screenplay.Performable
}

// String describes the action: "proceed through the temporary files".
func (a *ThroughTheAction[T]) String() string {
	return a.description
}

// PerformAs answers the question, then performs the body for each item. A
// failure ends the loop with an error that names the pass and the item.
func (a *ThroughTheAction[T]) PerformAs(theActor *screenplay.Actor) error {
	if a.body == nil {
		return fmt.Errorf("%s failed to %s: %w: the function given to AndForEach is nil",
			theActor.Name(), a, screenplay.ErrInvalidResolution)
	}

	items, err := a.items(theActor)
	if err != nil {
		return fmt.Errorf("%s failed to %s: %w", theActor.Name(), a, err)
	}

	for index, item := range items {
		if err = theActor.AttemptsTo(a.body(item)); err != nil {
			return fmt.Errorf("%s failed to %s on the %s pass, for %s: %w",
				theActor.Name(), a, ordinalOf(index+1), verdict.Describe(item), err)
		}
	}

	return nil
}

// ThroughTheAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*ThroughTheAction[string])(nil)
