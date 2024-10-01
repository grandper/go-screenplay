package screenplay

// Matcher is a function that evaluates an object against some
// criteria and returns the result and an error if one occurred.
type Matcher func(obj any) (bool, error)

// Resolution defines an interface for creating matchers.
type Resolution interface {
	// Resolve creates a matcher to make an assertion.
	Resolve() Matcher
	// String describes the resolution's expectation.
	String() string
}
