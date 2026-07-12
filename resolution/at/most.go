package at

import (
	"cmp"
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/screenplay"
)

// Most starts a matcher that tells if a number is within a given delta from a
// reference number. Complete it with AwayFrom: at.Most(delta).AwayFrom(number).
func Most[T cmp.Ordered](delta T) *MostBuilder[T] {
	return &MostBuilder[T]{
		delta: delta,
	}
}

// MostBuilder holds the delta of an at.Most matcher while it waits for the
// reference number supplied through AwayFrom.
type MostBuilder[T cmp.Ordered] struct {
	delta T
}

// AwayFrom completes the matcher with the number the value is compared against.
func (b *MostBuilder[T]) AwayFrom(number T) *MostResolution[T] {
	return &MostResolution[T]{
		number: number,
		delta:  b.delta,
	}
}

// MostResolution is a matcher to tell if a number is within a given delta from a number.
type MostResolution[T cmp.Ordered] struct {
	number T
	delta  T
}

// Resolve creates a matcher to make an assertion.
func (r *MostResolution[T]) Resolve() screenplay.Matcher {
	return func(obj any) (bool, error) {
		objValue, ok := obj.(T)
		if !ok {
			return false, errors.New("the object should be a comparable type")
		}

		if objValue <= r.delta+r.number && r.number <= objValue+r.delta {
			return true, nil
		}

		return false, nil
	}
}

// String describe the resolution's expectation.
func (r *MostResolution[T]) String() string {
	return fmt.Sprintf("at most %v away from %v", r.delta, r.number)
}

// MostResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*MostResolution[int])(nil)
