// Package length holds the resolutions about the length of a value: a channel,
// a map, a slice, an array, a string, or a pointer to one. They are nouns and
// take the verb ShouldHave (or ShouldNotHave):
//
//	err := actor.Should(see.The(customers).ShouldHave(length.Exactly[[]string](2)))
//	err := actor.Should(see.The(password).ShouldHave(length.AtLeast[string](12)))
//
// The first one reads "see the customers should have length exactly 2".
//
// Go cannot infer the type of the measured value from a length, so it is given
// explicitly as the type argument. A resolution is invalid, and fails with
// screenplay.ErrInvalidResolution when it inspects, when values of that type
// can never have a length, when a bound is negative, or when the bounds are
// out of order.
package length
