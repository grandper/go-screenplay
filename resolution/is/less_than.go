package is

import (
	"cmp"
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/screenplay"
)

// LessThan creates a matcher to tell if a number is less than a given number.
func LessThan[T cmp.Ordered](number T) *LessThanResolution[T] {
	return &LessThanResolution[T]{
		number: number,
	}
}

// LessThanResolution is a matcher to tell if a number is less than a given number.
type LessThanResolution[T cmp.Ordered] struct {
	number T
}

// Resolve creates a matcher to make an assertion.
func (r *LessThanResolution[T]) Resolve() screenplay.Matcher {
	return func(obj any) (bool, error) {
		objValue, ok := obj.(T)
		if !ok {
			return false, errors.New("the object should be a comparable type")
		}

		if objValue < r.number {
			return true, nil
		}

		return false, nil
	}
}

// String describe the resolution's expectation.
func (r *LessThanResolution[T]) String() string {
	return fmt.Sprintf("less than %v", r.number)
}

// LessThanResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*LessThanResolution[int])(nil)
