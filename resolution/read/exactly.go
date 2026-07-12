package read

import (
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/screenplay"
)

// Exactly creates a matcher to tell if a string match an exact text.
func Exactly(text string) *ExactlyResolution {
	return &ExactlyResolution{
		text: text,
	}
}

// ExactlyResolution is a matcher to tell if a string match an exact text.
type ExactlyResolution struct {
	text string
}

// Resolve creates a matcher to make an assertion.
func (r *ExactlyResolution) Resolve() screenplay.Matcher {
	return func(obj any) (bool, error) {
		objValue, ok := obj.(string)
		if !ok {
			return false, errors.New("the object should be a string")
		}

		if objValue == r.text {
			return true, nil
		}

		return false, nil
	}
}

// String describes the resolution's expectation.
func (r *ExactlyResolution) String() string {
	return fmt.Sprintf("reading exactly '%s'", r.text)
}

// ExactlyResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*ExactlyResolution)(nil)
