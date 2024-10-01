package resolution

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/grandper/go-screenplay/screenplay"
)

// Matches creates a matcher to tell if a string satisfies a regular expression.
func Matches(regex *regexp.Regexp) *MatchesResolution {
	return &MatchesResolution{
		regex: regex,
	}
}

// MatchesRegexString is a matcher to tell if a string satisfies a regular expression string.
func MatchesRegexString(regexStr string) *MatchesResolution {
	return &MatchesResolution{
		regex: regexp.MustCompile(regexStr),
	}
}

// MatchesResolution creates a matcher to tell if a string satisfies a regular expression.
type MatchesResolution struct {
	regex *regexp.Regexp
}

// Resolve creates a matcher to make an assertion.
func (r *MatchesResolution) Resolve() screenplay.Matcher {
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

// String describe the resolution's expectation.
func (r *MatchesResolution) String() string {
	return fmt.Sprintf("text matching the pattern %s", r.regex)
}

// MatchesResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*MatchesResolution)(nil)
