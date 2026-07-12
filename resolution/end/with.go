package end

import (
	"errors"
	"strings"

	"github.com/grandper/go-screenplay/screenplay"
)

// ErrObjectTypeShouldBeString is returned when the object type is not a string.
var ErrObjectTypeShouldBeString = errors.New("the object type should be a string")

// With creates a matcher to tell if a string ends with a given substring.
func With(suffix string) *WithResolution {
	return &WithResolution{
		suffix: suffix,
	}
}

// WithResolution is a matcher to tell if a string ends with a given substring.
type WithResolution struct {
	suffix string
}

// Resolve creates a matcher to make an assertion.
func (r *WithResolution) Resolve() screenplay.Matcher {
	return func(obj any) (bool, error) {
		objValue, ok := obj.(string)
		if !ok {
			return false, ErrObjectTypeShouldBeString
		}

		if strings.HasSuffix(objValue, r.suffix) {
			return true, nil
		}

		return false, nil
	}
}

// String describes the resolution's expectation.
func (r *WithResolution) String() string {
	return "ending with " + r.suffix
}

// WithResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*WithResolution)(nil)
