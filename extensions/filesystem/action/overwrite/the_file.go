package overwrite

import impl "github.com/grandper/go-screenplay/v2/extensions/filesystem/action/internal/overwrite"

// TheFile replaces the entire content of an existing file with new data:
//
//	overwrite.TheFile("notes.txt").WithTheText("Hello World")
func TheFile(filename string) *impl.TheFileAction {
	return impl.TheFile(filename)
}
