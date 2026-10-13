package end

import (
	"strings"

	"github.com/grandper/go-screenplay/v2/internal/verdict"
	"github.com/grandper/go-screenplay/v2/resolution"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// With creates a resolution to tell if a text ends with a given suffix.
func With(suffix string) screenplay.Resolution[string] {
	return resolution.FromFunc(
		"end with "+verdict.Describe(suffix),
		func(text string) error {
			return verdict.Of(strings.HasSuffix(text, suffix), text)
		},
	)
}
