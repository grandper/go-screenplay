package is

import (
	"fmt"
	"reflect"

	"github.com/grandper/go-screenplay/screenplay"
)

// EqualTo creates a matcher to tell if an object equals a given object.
func EqualTo[T any](value T) *EqualToResolution[T] {
	return &EqualToResolution[T]{
		value: value,
	}
}

// EqualToResolution is a matcher to tell if an object equals a given object.
type EqualToResolution[T any] struct {
	value T
}

// Resolve creates a matcher to make an assertion.
func (r *EqualToResolution[T]) Resolve() screenplay.Matcher {
	return func(obj any) (bool, error) {
		if reflect.DeepEqual(r.value, obj) {
			return true, nil
		}

		return false, nil
	}
}

// String describe the resolution's expectation.
func (r *EqualToResolution[T]) String() string {
	return fmt.Sprintf("equal to %v", r.value)
}

// EqualToResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*EqualToResolution[fmt.Stringer])(nil)
