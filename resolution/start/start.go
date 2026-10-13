package start

import (
	"strings"

	"github.com/grandper/go-screenplay/v2/internal/verdict"
	"github.com/grandper/go-screenplay/v2/resolution"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// With creates a resolution to tell if a text starts with a given prefix.
func With(prefix string) screenplay.Resolution[string] {
	return resolution.FromFunc(
		"start with "+verdict.Describe(prefix),
		func(text string) error {
			return verdict.Of(strings.HasPrefix(text, prefix), text)
		},
	)
}
