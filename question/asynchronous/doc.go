// Package asynchronous holds the question about the actions launched with
// proceed.AsynchronouslyTo.
//
//	err := actor.Should(see.The(asynchronous.Errors()).Should(be.Empty[[]error]()))
//
// reads "see the asynchronous errors should be empty"; answering it waits for
// every pending asynchronous action to complete.
package asynchronous
