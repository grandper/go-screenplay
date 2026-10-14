// Package fixture holds the fakes to test code that receives the answers of
// the file system questions, without touching the file system:
//
//	notes := fixture.CreateFakeExistingFile("notes.txt")
//	build := fixture.CreateFakeMissingDirectory("build")
package fixture
