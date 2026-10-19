// Package sized measures the values that have a length: channels, maps,
// slices, arrays, strings, and pointers to those.
package sized

import (
	"fmt"
	"reflect"
)

// Check returns an error when values of type T can never have a length. It is
// called when a resolution about a length is built, so that a resolution that
// cannot hold is reported as such rather than as a value that does not
// satisfy it. An interface type is accepted: its values are measured when
// inspected.
func Check[T any]() error {
	kind := reflect.TypeFor[T]()
	if kind.Kind() == reflect.Pointer {
		kind = kind.Elem()
	}

	if kind.Kind() == reflect.Interface || hasLength(kind.Kind()) {
		return nil
	}

	return fmt.Errorf("values of type %s have no length", reflect.TypeFor[T]())
}

// Len returns the length of a value, and whether it has one. A nil value and
// a nil pointer have the length zero.
func Len(value any) (int, bool) {
	if value == nil {
		return 0, true
	}

	reflected := reflect.ValueOf(value)
	if reflected.Kind() == reflect.Pointer {
		if reflected.IsNil() {
			return 0, true
		}

		reflected = reflected.Elem()
	}

	if !hasLength(reflected.Kind()) {
		return 0, false
	}

	return reflected.Len(), true
}

func hasLength(kind reflect.Kind) bool {
	switch kind { //nolint:exhaustive // only the kinds with a length are listed
	case reflect.Chan, reflect.Map, reflect.Slice, reflect.Array, reflect.String:
		return true
	default:
		return false
	}
}
