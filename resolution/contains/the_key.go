package contains

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/grandper/go-screenplay/screenplay"
)

// TheKey creates a matcher to tell if a map contains a given key.
func TheKey[T any](key T) *TheKeyResolution[T] {
	return &TheKeyResolution[T]{
		key: key,
	}
}

// TheKeyResolution is a matcher to tell if a number is greater than a given number.
type TheKeyResolution[T any] struct {
	key T
}

// Resolve creates a matcher to make an assertion.
func (r *TheKeyResolution[T]) Resolve() screenplay.Matcher {
	return func(obj any) (bool, error) {
		return hasKey(obj, r.key)
	}
}

func hasKey(obj any, key any) (bool, error) {
	if obj == nil {
		return false, nil
	}

	objValue := reflect.ValueOf(obj)

	switch objValue.Kind() { //nolint:exhaustive // we handle only the cases we need
	case reflect.Map:
		iter := objValue.MapRange()
		for iter.Next() {
			mapKey := iter.Key()
			if reflect.DeepEqual(mapKey.Interface(), key) {
				return true, nil
			}
		}

		return false, nil
	case reflect.Ptr:
		if objValue.IsNil() {
			return false, nil
		}

		ref := objValue.Elem().Interface()

		return hasKey(ref, key)
	default:
		return false, errors.New("the object should be a map")
	}
}

// String describe the resolution's expectation.
func (r *TheKeyResolution[T]) String() string {
	return fmt.Sprintf("containing the key %v", r.key)
}

// TheKeyResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*TheKeyResolution[int])(nil)
