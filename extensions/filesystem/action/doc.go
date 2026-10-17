// Package action holds the errors of the actions that drive the file system.
// The actions live in packages named after their verb:
//
//	create.TheFile("notes.txt").ContainingTheText("Hello World")  extensions/filesystem/action/create
//	remove.TheFile("draft.txt")                                   extensions/filesystem/action/remove
//	write.TheText("Hello World").AtTheEndOf("notes.txt")          extensions/filesystem/action/write
//	overwrite.TheFile("notes.txt").WithTheText("Hello World")     extensions/filesystem/action/overwrite
//	change.TheCurrentDirectoryTo("MyDocuments")                   extensions/filesystem/action/change
//
// For example:
//
//	err := anActor.AttemptsTo(create.TheFile("notes.txt").ContainingTheText("Hello World"))
package action
