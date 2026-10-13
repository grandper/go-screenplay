// Package standard holds the questions that answer with the standard output and the standard error
// of the result of a command:
//
//	err := anActor.Should(
//		see.The(standard.OutputOf(last.Of(question.Responses()))).Should(contain.TheBytes([]byte("Hello World"))),
//	)
package standard
