package contains

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/grandper/go-screenplay/screenplay"
)

// TheEntry creates a matcher to tell if a map contains a given key.
func TheEntry[T any, V any](key T, value V) *TheEntryResolution[T, V] {
	return &TheEntryResolution[T, V]{
		key:   key,
		value: value,
	}
}

// TheEntryResolution is a matcher to tell if a number is greater than a given number.
type TheEntryResolution[T any, V any] struct {
	key   T
	value V
}

// Resolve creates a matcher to make an assertion.
func (r *TheEntryResolution[T, V]) Resolve() screenplay.Matcher {
	return func(obj any) (bool, error) {
		return hasEntry(obj, r.key, r.value)
	}
}

func hasEntry(obj any, key any, value any) (bool, error) {
	if obj == nil {
		return false, nil
	}

	objValue := reflect.ValueOf(obj)

	switch objValue.Kind() { //nolint:exhaustive // we handle only the cases we need
	case reflect.Map:
		iter := objValue.MapRange()
		for iter.Next() {
			mapKey := iter.Key()
			mapValue := iter.Value()

			if reflect.DeepEqual(mapKey.Interface(), key) &&
				reflect.DeepEqual(mapValue.Interface(), value) {
				return true, nil
			}
		}

		return false, nil
	case reflect.Ptr:
		if objValue.IsNil() {
			return false, nil
		}

		ref := objValue.Elem().Interface()

		return hasEntry(ref, key, value)
	default:
		return false, errors.New("the object should be a map")
	}
}

// String describe the resolution's expectation.
func (r *TheEntryResolution[T, V]) String() string {
	return fmt.Sprintf("containing the entry [%v: %v]", r.key, r.value)
}

// TheEntryResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*TheEntryResolution[int, int])(nil)
