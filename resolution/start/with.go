package start

import (
	"errors"
	"fmt"
	"strings"

	"github.com/grandper/go-screenplay/screenplay"
)

// With creates a matcher to tell if a string starts with a given substring.
func With(prefix string) *WithResolution {
	return &WithResolution{
		prefix: prefix,
	}
}

// WithResolution is a matcher to tell if a string starts with a given substring.
type WithResolution struct {
	prefix string
}

// Resolve creates a matcher to make an assertion.
func (r *WithResolution) Resolve() screenplay.Matcher {
	return func(obj any) (bool, error) {
		objValue, ok := obj.(string)
		if !ok {
			return false, errors.New("the object should be a string")
		}

		if strings.HasPrefix(objValue, r.prefix) {
			return true, nil
		}

		return false, nil
	}
}

// String describes the resolution's expectation.
func (r *WithResolution) String() string {
	return fmt.Sprintf("starting with %s", r.prefix)
}

// WithResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*WithResolution)(nil)
