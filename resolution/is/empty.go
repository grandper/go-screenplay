package is

import (
	"reflect"

	"github.com/grandper/go-screenplay/screenplay"
)

// Empty creates a matcher to tell if a collection is empty.
func Empty() *EmptyResolution {
	return &EmptyResolution{}
}

// EmptyResolution is a matcher to tell if a collection is empty.
type EmptyResolution struct{}

// Resolve creates a matcher to make an assertion.
func (r *EmptyResolution) Resolve() screenplay.Matcher {
	return isEmpty
}

//nolint:gocognit
func isEmpty(obj any) (bool, error) {
	if obj == nil {
		return true, nil
	}

	objValue := reflect.ValueOf(obj)
	switch objValue.Kind() {
	case reflect.Chan, reflect.Map, reflect.Slice, reflect.Array, reflect.String:
		if objValue.Len() == 0 {
			return true, nil
		}
	case reflect.Ptr:
		if objValue.IsNil() {
			return true, nil
		}
		return isEmpty(objValue.Elem().Interface())
	case reflect.Bool:
		if !objValue.Bool() {
			return true, nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if objValue.Int() == 0 {
			return true, nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if objValue.Uint() == 0 {
			return true, nil
		}
	case reflect.Float32, reflect.Float64:
		if objValue.Float() == 0 {
			return true, nil
		}
	case reflect.Complex64, reflect.Complex128:
		if objValue.Complex() == 0 {
			return true, nil
		}
	case reflect.Interface, reflect.Func, reflect.UnsafePointer:
		if objValue.IsNil() {
			return true, nil
		}
	case reflect.Invalid: // Invalid kind is considered empty
		return true, nil
	case reflect.Struct:
		return structIsEmpty(objValue)
	default:
		zero := reflect.Zero(objValue.Type())
		if reflect.DeepEqual(objValue.Interface(), zero.Interface()) {
			return true, nil
		}
	}

	return false, nil
}

func structIsEmpty(objValue reflect.Value) (bool, error) {
	for i := range objValue.NumField() {
		field := objValue.Field(i)
		if !field.IsZero() {
			return false, nil
		}
	}

	return true, nil
}

// String describe the resolution's expectation.
func (r *EmptyResolution) String() string {
	return "an empty collection"
}

// EmptyResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*EmptyResolution)(nil)
