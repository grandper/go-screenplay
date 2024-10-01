package is

import (
	"cmp"
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/screenplay"
)

// GreaterThan creates a matcher to tell if a number is greater than a given number.
func GreaterThan[T cmp.Ordered](number T) *GreaterThanResolution[T] {
	return &GreaterThanResolution[T]{
		number: number,
	}
}

// GreaterThanResolution is a matcher to tell if a number is greater than a given number.
type GreaterThanResolution[T cmp.Ordered] struct {
	number T
}

// Resolve creates a matcher to make an assertion.
func (r *GreaterThanResolution[T]) Resolve() screenplay.Matcher {
	return func(obj any) (bool, error) {
		objValue, ok := obj.(T)
		if !ok {
			return false, errors.New("the object should be a comparable type")
		}

		if objValue > r.number {
			return true, nil
		}

		return false, nil
	}
}

// String describe the resolution's expectation.
func (r *GreaterThanResolution[T]) String() string {
	return fmt.Sprintf("greater than %v", r.number)
}

// GreaterThanResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*GreaterThanResolution[int])(nil)
