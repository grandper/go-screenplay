// Package resolution holds the resolution to tell if a file or a directory exists. It is a bare
// infinitive and takes the verb Should (or ShouldNot):
//
//	err := anActor.Should(see.The(question.FileNamed("notes.txt")).Should(resolution.Exist[*question.File]()))
//
// reads "see the file named notes.txt should exist".
package resolution
