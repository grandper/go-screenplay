// Package action holds the errors of the actions that drive command-line
// programs. The actions live in packages named after their verb:
//
//	run.TheCommand("echo", "Hello World")                 extensions/cli/action/run
//	enter.TheText("Alice in Wonderland").AndPressEnter()  extensions/cli/action/enter
//
// For example:
//
//	err := anActor.AttemptsTo(run.TheCommand("echo", "Hello World"))
package action
