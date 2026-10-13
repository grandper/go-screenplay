package match

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/grandper/go-screenplay/v2/internal/verdict"
	"github.com/grandper/go-screenplay/v2/resolution"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// Pattern creates a resolution to tell if a text matches a regular expression.
// The resolution is invalid when the expression is nil.
func Pattern(pattern *regexp.Regexp) screenplay.Resolution[string] {
	if pattern == nil {
		return verdict.Invalid[string]("match pattern <nil>", errors.New("the regular expression is nil"))
	}

	return resolution.FromFunc(
		fmt.Sprintf("match pattern `%s`", pattern),
		func(text string) error {
			return verdict.Of(pattern.MatchString(text), text)
		},
	)
}

// PatternString creates a resolution to tell if a text matches a regular
// expression given as a text. The resolution is invalid when the expression
// does not compile.
func PatternString(pattern string) screenplay.Resolution[string] {
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return verdict.Invalid[string](fmt.Sprintf("match pattern `%s`", pattern), err)
	}

	return Pattern(compiled)
}
