package match

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/grandper/go-screenplay/screenplay"
)

// TheRegex creates a matcher to tell if a string satisfies a regular expression.
func TheRegex(regex *regexp.Regexp) *TheRegexResolution {
	return &TheRegexResolution{
		regex: regex,
	}
}

// TheRegexString creates a matcher to tell if a string satisfies a regular expression string.
func TheRegexString(regexStr string) *TheRegexResolution {
	return &TheRegexResolution{
		regex: regexp.MustCompile(regexStr),
	}
}

// TheRegexResolution is a matcher to tell if a string satisfies a regular expression.
type TheRegexResolution struct {
	regex *regexp.Regexp
}

// Resolve creates a matcher to make an assertion.
func (r *TheRegexResolution) Resolve() screenplay.Matcher {
	return func(obj any) (bool, error) {
		objValue, ok := obj.(string)
		if !ok {
			return false, errors.New("the object should be a string")
		}

		if r.regex.MatchString(objValue) {
			return true, nil
		}

		return false, nil
	}
}

// String describes the resolution's expectation.
func (r *TheRegexResolution) String() string {
	return fmt.Sprintf("text matching the pattern %s", r.regex)
}

// TheRegexResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*TheRegexResolution)(nil)
