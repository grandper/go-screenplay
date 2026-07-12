package less

import (
	"cmp"
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/screenplay"
)

// Than creates a matcher to tell if a number is less than a given number.
func Than[T cmp.Ordered](number T) *ThanResolution[T] {
	return &ThanResolution[T]{
		number: number,
	}
}

// ThanResolution is a matcher to tell if a number is less than a given number.
type ThanResolution[T cmp.Ordered] struct {
	number T
}

// Resolve creates a matcher to make an assertion.
func (r *ThanResolution[T]) Resolve() screenplay.Matcher {
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
func (r *ThanResolution[T]) String() string {
	return fmt.Sprintf("less than %v", r.number)
}

// ThanResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*ThanResolution[int])(nil)
