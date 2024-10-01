package contains

import (
	"errors"
	"fmt"
	"strings"

	"github.com/grandper/go-screenplay/screenplay"
)

// TheText creates a matcher to tell if a string contains a given text.
func TheText(text string) *TheTextResolution {
	return &TheTextResolution{
		text: text,
	}
}

// TheTextResolution is a matcher to tell if a string contains a given text.
type TheTextResolution struct {
	text string
}

// Resolve creates a matcher to make an assertion.
func (r *TheTextResolution) Resolve() screenplay.Matcher {
	return func(obj any) (bool, error) {
		objValue, ok := obj.(string)
		if !ok {
			return false, errors.New("the object should be a string")
		}

		if strings.Contains(objValue, r.text) {
			return true, nil
		}

		return false, nil
	}
}

// String describe the resolution's expectation.
func (r *TheTextResolution) String() string {
	return fmt.Sprintf("containing the text %s", r.text)
}

// TheTextResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*TheTextResolution)(nil)
