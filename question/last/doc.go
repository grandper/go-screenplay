// Package last holds the question that answers with the last element of the
// slice another question answers with:
//
//	err := actor.Should(see.The(status.CodeOf(last.Of(Responses()))).Should(be.EqualTo(200)))
//
// The answer is typed as the element of the slice, and wrapping a question
// that does not answer with a slice is a compile-time error.
package last
