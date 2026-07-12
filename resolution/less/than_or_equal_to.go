package less

import (
	"cmp"
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/screenplay"
)

// ThanOrEqualTo creates a matcher to tell if a number is less than or equal to a given number.
func ThanOrEqualTo[T cmp.Ordered](number T) *ThanOrEqualToResolution[T] {
	return &ThanOrEqualToResolution[T]{
		number: number,
	}
}

// ThanOrEqualToResolution is a matcher to tell if a number is less than or equal to a given number.
type ThanOrEqualToResolution[T cmp.Ordered] struct {
	number T
}

// Resolve creates a matcher to make an assertion.
func (r *ThanOrEqualToResolution[T]) Resolve() screenplay.Matcher {
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
func (r *ThanOrEqualToResolution[T]) String() string {
	return fmt.Sprintf("less than or equal to %v", r.number)
}

// ThanOrEqualToResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*ThanOrEqualToResolution[int])(nil)
