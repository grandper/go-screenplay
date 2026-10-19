package create

import impl "github.com/grandper/go-screenplay/v2/extensions/filesystem/action/internal/create"

// TheFile creates a file:
//
//	create.TheFile("notes.txt").ContainingTheText("Hello World")
func TheFile(fileName string) *impl.TheFileAction {
	return impl.TheFile(fileName)
}

// TheTemporaryFile creates a temporary file, whose name is built from the
// pattern given.
func TheTemporaryFile(fileName string) *impl.TheFileAction {
	return impl.TheTemporaryFile(fileName)
}

// TheDirectory creates a directory.
func TheDirectory(directoryName string) *impl.TheDirectoryAction {
	return impl.TheDirectory(directoryName)
}

// TheTemporaryDirectory creates a temporary directory, whose name is built
// from the pattern given.
func TheTemporaryDirectory(directoryName string) *impl.TheDirectoryAction {
	return impl.TheTemporaryDirectory(directoryName)
}
