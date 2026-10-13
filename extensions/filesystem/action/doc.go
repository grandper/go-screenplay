// Package action holds the actions to drive the file system: creating, appending to, overwriting
// and removing files and directories, and changing the working directory.
//
//	err := anActor.AttemptsTo(action.Create().TheFile("notes.txt").ContainingTheText("Hello World"))
package action
