package empty

import (
	"errors"
	"reflect"

	"github.com/grandper/go-screenplay/screenplay"
)

// ErrObjectShouldBeACollection is returned when the object is not a collection.
var ErrObjectShouldBeACollection = errors.New("the object should be a collection")

// Collection creates a matcher to tell if a collection is empty. A collection is
// a channel, a map, a slice, an array or a string (pointers to those are
// dereferenced).
func Collection() *CollectionResolution {
	return &CollectionResolution{}
}

// CollectionResolution is a matcher to tell if a collection is empty.
type CollectionResolution struct{}

// Resolve creates a matcher to make an assertion.
func (r *CollectionResolution) Resolve() screenplay.Matcher {
	return isEmptyCollection
}

func isEmptyCollection(obj any) (bool, error) {
	if obj == nil {
		return true, nil
	}

	objValue := reflect.ValueOf(obj)
	switch objValue.Kind() { //nolint:exhaustive // we handle only the collection kinds
	case reflect.Chan, reflect.Map, reflect.Slice, reflect.Array, reflect.String:
		return objValue.Len() == 0, nil
	case reflect.Ptr:
		if objValue.IsNil() {
			return true, nil
		}
		return isEmptyCollection(objValue.Elem().Interface())
	default:
		return false, ErrObjectShouldBeACollection
	}
}

// String describe the resolution's expectation.
func (r *CollectionResolution) String() string {
	return "an empty collection"
}

// CollectionResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*CollectionResolution)(nil)

// Value creates a matcher to tell if a value equals its zero value. A nil object,
// a nil pointer, the false boolean, the number zero and a struct whose fields are
// all zero are all considered empty (pointers are dereferenced).
func Value() *ValueResolution {
	return &ValueResolution{}
}

// ValueResolution is a matcher to tell if a value equals its zero value.
type ValueResolution struct{}

// Resolve creates a matcher to make an assertion.
func (r *ValueResolution) Resolve() screenplay.Matcher {
	return isEmptyValue
}

func isEmptyValue(obj any) (bool, error) {
	if obj == nil {
		return true, nil
	}

	objValue := reflect.ValueOf(obj)
	if objValue.Kind() == reflect.Ptr {
		if objValue.IsNil() {
			return true, nil
		}
		return isEmptyValue(objValue.Elem().Interface())
	}

	return objValue.IsZero(), nil
}

// String describe the resolution's expectation.
func (r *ValueResolution) String() string {
	return "an empty value"
}

// ValueResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*ValueResolution)(nil)
