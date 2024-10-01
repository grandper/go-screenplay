package is

import (
	"cmp"
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/screenplay"
)

// GreaterThanOrEqualTo creates a matcher to tell if a number is greater than or equal to a given number.
func GreaterThanOrEqualTo[T cmp.Ordered](number T) *GreaterThanOrEqualToResolution[T] {
	return &GreaterThanOrEqualToResolution[T]{
		number: number,
	}
}

// GreaterThanOrEqualToResolution is a matcher to tell if a number is greater than or equal to a given number.
type GreaterThanOrEqualToResolution[T cmp.Ordered] struct {
	number T
}

// Resolve creates a matcher to make an assertion.
func (r *GreaterThanOrEqualToResolution[T]) Resolve() screenplay.Matcher {
	return func(obj any) (bool, error) {
		objValue, ok := obj.(T)
		if !ok {
			return false, errors.New("the object should be a comparable type")
		}

		if objValue >= r.number {
			return true, nil
		}

		return false, nil
	}
}

// String describe the resolution's expectation.
func (r *GreaterThanOrEqualToResolution[T]) String() string {
	return fmt.Sprintf("greater than or equal to %v", r.number)
}

// GreaterThanOrEqualToResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*GreaterThanOrEqualToResolution[int])(nil)
