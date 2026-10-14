// Package exitcode holds the question that answers with the exit code of the result of a command:
//
//	err := anActor.Should(see.The(exitcode.Of(last.Of(question.Responses()))).Should(be.EqualTo(0)))
package exitcode
