// Package status holds the question that answers with the status code of an HTTP response:
//
//	err := anActor.Should(see.The(status.CodeOf(last.Of(question.Responses()))).Should(be.EqualTo(200)))
package status
