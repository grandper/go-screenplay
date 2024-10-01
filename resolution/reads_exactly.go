package resolution

import (
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/screenplay"
)

// ReadsExactly creates a matcher to tell if a string match an exact text.
func ReadsExactly(text string) *ReadsExactlyResolution {
	return &ReadsExactlyResolution{
		text: text,
	}
}

// ReadsExactlyResolution is a matcher to tell if a string match an exact text.
type ReadsExactlyResolution struct {
	text string
}

// Resolve creates a matcher to make an assertion.
func (r *ReadsExactlyResolution) Resolve() screenplay.Matcher {
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

// String describe the resolution's expectation.
func (r *ReadsExactlyResolution) String() string {
	return fmt.Sprintf("reading exactly '%s'", r.text)
}

// ReadsExactlyResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*ReadsExactlyResolution)(nil)
