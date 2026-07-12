package equal

import (
	"fmt"
	"reflect"

	"github.com/grandper/go-screenplay/screenplay"
)

// To creates a matcher to tell if an object equals a given object.
func To[T any](value T) *ToResolution[T] {
	return &ToResolution[T]{
		value: value,
	}
}

// ToResolution is a matcher to tell if an object equals a given object.
type ToResolution[T any] struct {
	value T
}

// Resolve creates a matcher to make an assertion.
func (r *ToResolution[T]) Resolve() screenplay.Matcher {
	return func(obj any) (bool, error) {
		if reflect.DeepEqual(r.value, obj) {
			return true, nil
		}

		return false, nil
	}
}

// String describe the resolution's expectation.
func (r *ToResolution[T]) String() string {
	return fmt.Sprintf("equal to %v", r.value)
}

// ToResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*ToResolution[fmt.Stringer])(nil)
