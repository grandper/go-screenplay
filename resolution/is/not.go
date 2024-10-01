package is

import (
	"fmt"

	"github.com/grandper/go-screenplay/screenplay"
)

// Not creates a matcher to tell if the inverse of a resolution.
func Not(resolution screenplay.Resolution) *NotResolution {
	return &NotResolution{
		resolution: resolution,
	}
}

// NotResolution is a matcher to tell if the inverse of a resolution.
type NotResolution struct {
	resolution screenplay.Resolution
}

// Resolve creates a matcher to make an assertion.
func (r *NotResolution) Resolve() screenplay.Matcher {
	resolver := r.resolution.Resolve()

	return func(obj any) (bool, error) {
		result, err := resolver(obj)
		if err != nil {
			return false, err
		}

		return !result, nil
	}
}

// String describe the resolution's expectation.
func (r *NotResolution) String() string {
	return fmt.Sprintf("not %v", r.resolution.String())
}

// NotResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*NotResolution)(nil)
