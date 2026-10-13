package testdata

import (
	"errors"

	"github.com/grandper/go-screenplay/v2/screenplay"
)

// errNotSatisfied is the error a fake resolution returns when it is built not
// to be satisfied.
var errNotSatisfied = errors.New("the fake resolution is not satisfied")

// FakeResolution is a fake resolution for test purpose. T is the type of the
// values it inspects.
type FakeResolution[T any] struct {
	satisfied   bool
	description string
}

// NewFakeResolution creates a new fake resolution that every value satisfies,
// or that none does. The inspected type cannot be inferred and must be
// provided, for example NewFakeResolution[string]("be anything", true).
func NewFakeResolution[T any](description string, satisfied bool) *FakeResolution[T] {
	return &FakeResolution[T]{
		satisfied:   satisfied,
		description: description,
	}
}

// Inspect returns nil when the fake resolution is built to be satisfied.
func (fr *FakeResolution[T]) Inspect(_ T) error {
	if fr.satisfied {
		return nil
	}

	return errNotSatisfied
}

// String describes the resolution's expectation.
func (fr *FakeResolution[T]) String() string {
	return fr.description
}

// FakeResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution[any] = (*FakeResolution[any])(nil)
