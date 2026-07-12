package length

import (
	"fmt"
	"reflect"

	"github.com/grandper/go-screenplay/screenplay"
)

// Of creates a matcher to tell if a collection has a specific length.
func Of(length int) *OfResolution {
	return &OfResolution{
		length: length,
	}
}

// OfResolution is a matcher to tell if a collection has a specific length.
type OfResolution struct {
	length int
}

// Resolve creates a matcher to make an assertion.
func (r *OfResolution) Resolve() screenplay.Matcher {
	return func(obj any) (bool, error) {
		if lengthOf(obj) != r.length {
			return false, nil
		}

		return true, nil
	}
}

func lengthOf(obj any) int {
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

		return lengthOf(ref)
	case reflect.String:
		return objValue.Len()
	default:
		return 1
	}
}

// String describe the resolution's expectation.
func (r *OfResolution) String() string {
	return fmt.Sprintf("%s long", r.numItemsString())
}

func (r *OfResolution) numItemsString() string {
	if r.length > 1 {
		return fmt.Sprintf("%d items", r.length)
	}

	return fmt.Sprintf("%d item", r.length)
}

// OfResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*OfResolution)(nil)
