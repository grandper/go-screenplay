package is

import (
	"cmp"
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/screenplay"
)

// LessThanOrEqualTo creates a matcher to tell if a number is less than or equal to a given number.
func LessThanOrEqualTo[T cmp.Ordered](number T) *LessThanOrEqualToResolution[T] {
	return &LessThanOrEqualToResolution[T]{
		number: number,
	}
}

// LessThanOrEqualToResolution is a matcher to tell if a number is less than or equal to a given number.
type LessThanOrEqualToResolution[T cmp.Ordered] struct {
	number T
}

// Resolve creates a matcher to make an assertion.
func (r *LessThanOrEqualToResolution[T]) Resolve() screenplay.Matcher {
	return func(obj any) (bool, error) {
		objValue, ok := obj.(T)
		if !ok {
			return false, errors.New("the object should be a comparable type")
		}

		if objValue <= r.number {
			return true, nil
		}

		return false, nil
	}
}

// String describe the resolution's expectation.
func (r *LessThanOrEqualToResolution[T]) String() string {
	return fmt.Sprintf("less than or equal to %v", r.number)
}

// LessThanOrEqualToResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*LessThanOrEqualToResolution[int])(nil)
