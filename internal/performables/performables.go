// Package performables holds the helpers the actions that take a list of
// performables share.
package performables

import (
	"strings"

	"github.com/grandper/go-screenplay/v2/screenplay"
)

// WithoutNil drops the nil items of a list of performables, as the actor
// does, so that a nil written by mistake in a test never panics.
func WithoutNil(performables []screenplay.Performable) []screenplay.Performable {
	kept := make([]screenplay.Performable, 0, len(performables))
	for _, performable := range performables {
		if performable != nil {
			kept = append(kept, performable)
		}
	}

	return kept
}

// Describe describes a list of performables, separated by commas.
func Describe(performables []screenplay.Performable) string {
	descriptions := make([]string, 0, len(performables))
	for _, performable := range performables {
		descriptions = append(descriptions, performable.String())
	}

	return strings.Join(descriptions, ", ")
}
