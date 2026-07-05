package resolution

import (
	"github.com/grandper/go-screenplay/screenplay"
)

// FromFunc creates a resolution with the given description and matcher factory.
func FromFunc(description string, matcherFactory screenplay.MatcherFactory) screenplay.Resolution {
	return &funcResolution{
		description:    description,
		matcherFactory: matcherFactory,
	}
}

type funcResolution struct {
	description    string
	matcherFactory screenplay.MatcherFactory
}

// Resolve creates a matcher to make an assertion.
func (r *funcResolution) Resolve() screenplay.Matcher {
	return r.matcherFactory()
}

// String describes the resolution's expectation.
func (r *funcResolution) String() string {
	return r.description
}

// funcResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*funcResolution)(nil)
