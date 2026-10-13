// Package first holds the question that answers with the first element of the
// slice another question answers with:
//
//	err := actor.Should(see.The(first.Of(items.InThe(TodoList))).Should(be.EqualTo("Write code")))
//
// The answer is typed as the element of the slice, and wrapping a question
// that does not answer with a slice is a compile-time error.
package first
