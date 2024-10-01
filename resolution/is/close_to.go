package is

import (
	"cmp"
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/screenplay"
)

// CloseTo creates a matcher to tell if a number is within a given delta from a number.
func CloseTo[T cmp.Ordered](number, delta T) *CloseToResolution[T] {
	return &CloseToResolution[T]{
		number: number,
		delta:  delta,
	}
}

// CloseToResolution is a matcher to tell if a number is within a given delta from a number.
type CloseToResolution[T cmp.Ordered] struct {
	number T
	delta  T
}

// Resolve creates a matcher to make an assertion.
func (r *CloseToResolution[T]) Resolve() screenplay.Matcher {
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
func (r *CloseToResolution[T]) String() string {
	return fmt.Sprintf("at most %v away from %v", r.delta, r.number)
}

// CloseToResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*CloseToResolution[int])(nil)
