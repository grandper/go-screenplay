package write

import (
	"io"

	impl "github.com/grandper/go-screenplay/v2/extensions/filesystem/action/internal/write"
)

// TheText writes the text at the end of a file, without altering its existing content:
//
//	write.TheText("Hello World").AtTheEndOf("notes.txt")
func TheText(text string) *impl.TheContentAction {
	return impl.TheText(text)
}

// TheContent writes the content of a reader at the end of a file, without
// altering its existing content. The reader is buffered once, so the action can be performed
// more than once.
func TheContent(reader io.Reader) *impl.TheContentAction {
	return impl.TheContent(reader)
}

// TheBytes writes the bytes at the end of a file, without altering its existing content.
func TheBytes(b []byte) *impl.TheContentAction {
	return impl.TheBytes(b)
}
