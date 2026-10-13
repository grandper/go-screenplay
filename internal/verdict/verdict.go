// Package verdict holds the small helpers the built-in resolutions share to
// describe values, to report the value that does not satisfy them, and to
// report a resolution that was built wrongly.
package verdict

import (
	"fmt"
	"reflect"

	"github.com/grandper/go-screenplay/v2/screenplay"
)

// Describe renders a value as it reads in a description or in a failure:
// texts and bytes are quoted, everything else uses its default format.
func Describe(value any) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("%q", v)
	case []byte:
		return fmt.Sprintf("%q", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// Of returns nil when the value satisfies a resolution, and otherwise the
// error every built-in resolution reports: the value it got.
func Of(satisfied bool, value any) error {
	if satisfied {
		return nil
	}

	return fmt.Errorf("got %s", Describe(value))
}

// Equal compares two values of a comparable type. When T is an interface type,
// such as any, a value whose dynamic type is not comparable (a slice, a map, a
// func) would make == panic; Equal reports it as an invalid resolution instead,
// so that the mistake fails when the actor performs the action.
func Equal[T comparable](value, expected T) (bool, error) {
	if reflect.TypeFor[T]().Kind() == reflect.Interface {
		for _, v := range []T{value, expected} {
			if rv := reflect.ValueOf(v); rv.IsValid() && !rv.Comparable() {
				return false, fmt.Errorf("%w: values of type %T cannot be compared", screenplay.ErrInvalidResolution, v)
			}
		}
	}

	return value == expected, nil
}

// Invalid creates the resolution a constructor returns when it is given
// something it cannot build a resolution from. Its Inspect always fails with
// an error wrapping screenplay.ErrInvalidResolution and the reason, so that
// the mistake is reported when the actor performs the action, whatever the
// verb the resolution is used with.
func Invalid[T any](description string, reason error) screenplay.Resolution[T] {
	return &invalidResolution[T]{
		description: description,
		err:         fmt.Errorf("%w %q: %w", screenplay.ErrInvalidResolution, description, reason),
	}
}

type invalidResolution[T any] struct {
	description string
	err         error
}

// Inspect always fails: the resolution was built wrongly.
func (r *invalidResolution[T]) Inspect(_ T) error {
	return r.err
}

// String describes the resolution's expectation.
func (r *invalidResolution[T]) String() string {
	return r.description
}
