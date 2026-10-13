// Package be holds the resolutions that complete "should be ...": equality,
// comparisons, emptiness and booleans. They are bare infinitives and take the
// verb Should (or ShouldNot, or ShouldAnyOf):
//
//	err := actor.Should(see.The(statusCode).Should(be.EqualTo(200)))
//	err := actor.Should(see.The(total).Should(be.Between(1, 100)))
//	err := actor.Should(see.The(customers).ShouldNot(be.Empty[[]string]()))
//	err := actor.Should(see.The(statusCode).ShouldAnyOf(be.EqualTo(200), be.EqualTo(201)))
//
// The first one reads "see the status code should be equal to 200".
//
// The resolutions built from a value inspect the values of its type:
// be.EqualTo(200) is a resolution on int. The ones built without a value take
// the inspected type explicitly, as be.Empty[[]string]() and be.Zero[int]() do.
package be
