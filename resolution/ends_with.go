package resolution

import (
	"errors"
	"strings"

	"github.com/grandper/go-screenplay/screenplay"
)

var (
	// ErrObjectTypeShouldBeString is returned when the object type is not a string.
	ErrObjectTypeShouldBeString = errors.New("the object type should be a string")
)

// EndsWith creates a matcher to tell if a string ends with a given substring.
func EndsWith(suffix string) *EndsWithResolution {
	return &EndsWithResolution{
		suffix: suffix,
	}
}

// EndsWithResolution is a matcher to tell if a string ends with a given substring.
type EndsWithResolution struct {
	suffix string
}

// Resolve creates a matcher to make an assertion.
func (r *EndsWithResolution) Resolve() screenplay.Matcher {
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

// String describe the resolution's expectation.
func (r *EndsWithResolution) String() string {
	return "ending with " + r.suffix
}

// EndsWithResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*EndsWithResolution)(nil)
