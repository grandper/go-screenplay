package contains

import (
	"fmt"
	"reflect"

	"github.com/grandper/go-screenplay/screenplay"
)

// TheValue creates a matcher to tell if a map or a slice contains a given value.
func TheValue[T any](value T) *TheValueResolution[T] {
	return &TheValueResolution[T]{
		value: value,
	}
}

// TheValueResolution is a matcher to tell if a number is greater than a given number.
type TheValueResolution[T any] struct {
	value T
}

// Resolve creates a matcher to make an assertion.
func (r *TheValueResolution[T]) Resolve() screenplay.Matcher {
	return func(obj any) (bool, error) {
		return hasValue(obj, r.value)
	}
}

func hasValue(obj any, value any) (bool, error) {
	if obj == nil {
		return false, nil
	}

	objValue := reflect.ValueOf(obj)

	switch objValue.Kind() { //nolint:exhaustive // we handle only the cases we need
	case reflect.Chan:
		continueScanning := true
		var chanValue reflect.Value
		var found bool

		for continueScanning {
			chanValue, continueScanning = objValue.Recv()
			if reflect.DeepEqual(chanValue.Interface(), value) {
				found = true
			}
		}

		return found, nil
	case reflect.Map:
		iter := objValue.MapRange()
		for iter.Next() {
			mapValue := iter.Value()
			if reflect.DeepEqual(mapValue.Interface(), value) {
				return true, nil
			}
		}

		return false, nil
	case reflect.Slice, reflect.Array:
		for i := range objValue.Len() {
			sliceValue := objValue.Index(i)
			if reflect.DeepEqual(sliceValue.Interface(), value) {
				return true, nil
			}
		}

		return false, nil
	case reflect.Ptr:
		if objValue.IsNil() {
			return false, nil
		}

		ref := objValue.Elem().Interface()

		return hasValue(ref, value)
	default:
		return reflect.DeepEqual(obj, value), nil
	}
}

// String describe the resolution's expectation.
func (r *TheValueResolution[T]) String() string {
	return fmt.Sprintf("containing the value %v", r.value)
}

// TheValueResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*TheValueResolution[int])(nil)
