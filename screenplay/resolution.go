package screenplay

import "errors"

// ErrInvalidResolution is wrapped by the error a resolution returns from
// Inspect when it was built wrongly: an invalid pattern, an impossible bound, a
// type whose values it cannot inspect. Such a resolution can never be judged,
// so every verb fails on it, the negative ones included, and the error tells a
// mistake in the test apart from a value that is not satisfying.
//
// The built-in resolutions never panic: a resolution built wrongly is reported
// through this error when the actor performs the action that uses it. A
// resolution of your own does the same by returning an error wrapping
// ErrInvalidResolution from Inspect.
var ErrInvalidResolution = errors.New("invalid resolution")

// Resolution represents the expectation a value of type T should satisfy. T is
// the type of the values it inspects, so that a resolution only pairs with the
// questions answering with that type.
//
// An error returned by Inspect means that the value does not satisfy the
// resolution, which is what lets a verb such as ShouldNot invert the outcome.
// The one exception is an error wrapping ErrInvalidResolution, which reports a
// resolution that was built wrongly and makes every verb fail.
type Resolution[T any] interface {
	// Inspect returns nil when the value satisfies the resolution, and an
	// error explaining why when it does not.
	Inspect(value T) error

	// String describes the resolution's expectation, written to follow
	// "should": "be equal to 200", "contain substring "hello"".
	String() string
}
