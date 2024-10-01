package contains

import (
	"fmt"
	"reflect"

	"github.com/grandper/go-screenplay/screenplay"
)

// TheItem creates a matcher to tell if an iterable contains a given item.
func TheItem[T any](item T) *TheItemResolution[T] {
	return &TheItemResolution[T]{
		item: item,
	}
}

// TheItemResolution is a matcher to tell if an iterable contains a given item.
type TheItemResolution[T any] struct {
	item T
}

// Resolve creates a matcher to make an assertion.
func (r *TheItemResolution[T]) Resolve() screenplay.Matcher {
	return func(obj any) (bool, error) {
		return hasItem(obj, r.item)
	}
}

func hasItem(obj any, value any) (bool, error) {
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

		return hasItem(ref, value)
	default:
		return reflect.DeepEqual(obj, value), nil
	}
}

// String describe the resolution's expectation.
func (r *TheItemResolution[T]) String() string {
	return fmt.Sprintf("containing the item %v", r.item)
}

// TheItemResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*TheItemResolution[int])(nil)
