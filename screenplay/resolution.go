package screenplay

// Matcher is a function that evaluates an object against some
// criteria and returns the result and an error if one occurred.
type Matcher func(obj any) (bool, error)

// MatcherFactory is a function that creates a Matcher.
type MatcherFactory func() Matcher

// Resolution defines an interface for creating matchers.
type Resolution interface {
	// Resolve creates a matcher to make an assertion.
	Resolve() Matcher
	// String describes the resolution's expectation.
	String() string
}

// ResolutionToSeeThatTheObject creates a resolution with the given description and matcher factory.
func ResolutionToSeeThatTheObject(description string, matcherFactory MatcherFactory) Resolution {
	return &anonymousResolution{
		description:    description,
		matcherFactory: matcherFactory,
	}
}

type anonymousResolution struct {
	description    string
	matcherFactory MatcherFactory
}

// Resolve creates a matcher to make an assertion.
func (ar *anonymousResolution) Resolve() Matcher {
	return ar.matcherFactory()
}

// String describes the resolution's expectation.
func (ar *anonymousResolution) String() string {
	return ar.description
}

// Ensure that anonymousResolution implements the Resolution interface.
var _ Resolution = (*anonymousResolution)(nil)
