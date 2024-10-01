package resolution

import (
	"errors"
	"fmt"
	"strings"

	"github.com/grandper/go-screenplay/screenplay"
)

// StartsWith creates a matcher to tell if a string starts with a given substring.
func StartsWith(prefix string) *StartsWithResolution {
	return &StartsWithResolution{
		prefix: prefix,
	}
}

// StartsWithResolution is a matcher to tell if a string starts with a given substring.
type StartsWithResolution struct {
	prefix string
}

// Resolve creates a matcher to make an assertion.
func (r *StartsWithResolution) Resolve() screenplay.Matcher {
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

// String describe the resolution's expectation.
func (r *StartsWithResolution) String() string {
	return fmt.Sprintf("starting with %s", r.prefix)
}

// StartsWithResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*StartsWithResolution)(nil)
