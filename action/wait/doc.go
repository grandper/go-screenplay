// Package wait holds the actions that make an actor wait.
//
//	wait.UntilThe(question).HappensTo(...)       waits until an answer happens to be as stated
//	wait.UntilTheEnterKeyIsPressed()             waits for the person debugging the test
//	wait.UntilAnInputIsProvidedBy(reader)        waits until a line is read from a reader
//	wait.For(10).Seconds().Because(reason)       waits for a given time
//
// For example:
//
//	err := actor.AttemptsTo(wait.UntilThe(logs).HappenTo(contain.Substring("ready")))
//
// reads "wait until the logs happen to contain substring "ready"", which is
// also what the actor narrates.
package wait
