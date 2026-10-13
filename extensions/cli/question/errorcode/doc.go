// Package errorcode holds the question that answers with the error code of the result of a command:
//
//	err := anActor.Should(see.The(errorcode.Of(last.Of(question.Responses()))).Should(be.EqualTo(0)))
package errorcode
