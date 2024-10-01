package resolution

import (
	"fmt"
	"reflect"

	"github.com/grandper/go-screenplay/screenplay"
)

// HasLength creates a matcher to tell if a collection has a specific length.
func HasLength(length int) *HasLengthResolution {
	return &HasLengthResolution{
		length: length,
	}
}

// HasLengthResolution is a matcher to tell if a collection has a specific length.
type HasLengthResolution struct {
	length int
}

// Resolve creates a matcher to make an assertion.
func (r *HasLengthResolution) Resolve() screenplay.Matcher {
	return func(obj any) (bool, error) {
		if length(obj) != r.length {
			return false, nil
		}

		return true, nil
	}
}

func length(obj any) int {
	if obj == nil {
		return 0
	}

	objValue := reflect.ValueOf(obj)

	switch objValue.Kind() { //nolint:exhaustive // we handle only the cases we need
	case reflect.Chan, reflect.Map, reflect.Slice, reflect.Array:
		return objValue.Len()
	case reflect.Ptr:
		if objValue.IsNil() {
			return 0
		}
		ref := objValue.Elem().Interface()

		return length(ref)
	case reflect.String:
		return objValue.Len()
	default:
		return 1
	}
}

// String describe the resolution's expectation.
func (r *HasLengthResolution) String() string {
	return fmt.Sprintf("%s long", r.numItemsString())
}

func (r *HasLengthResolution) numItemsString() string {
	if r.length > 1 {
		return fmt.Sprintf("%d items", r.length)
	}

	return fmt.Sprintf("%d item", r.length)
}

// HasLengthResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*HasLengthResolution)(nil)
