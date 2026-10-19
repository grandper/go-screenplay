package write

import (
	"bytes"
	"fmt"
	"io"

	"github.com/grandper/go-screenplay/v2/extensions/filesystem/action"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/action/internal/content"
	impl "github.com/grandper/go-screenplay/v2/extensions/filesystem/internal/ability"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/internal/verdict"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// TheText writes the text at the end of a file, without altering its existing content:
//
//	write.TheText("Hello World").AtTheEndOf("notes.txt")
func TheText(text string) *TheContentAction {
	return &TheContentAction{
		content:     []byte(text),
		contentKind: "text",
		hasContent:  true,
	}
}

// TheContent writes the content of a reader at the end of a file, without
// altering its existing content. The reader is buffered once, so the action can be performed
// more than once.
func TheContent(reader io.Reader) *TheContentAction {
	if reader == nil {
		return &TheContentAction{contentKind: "content"}
	}

	content, err := content.Read(reader)

	return &TheContentAction{
		content:     content,
		contentKind: "content",
		hasContent:  true,
		contentErr:  err,
	}
}

// TheBytes writes the bytes at the end of a file, without altering its existing content.
func TheBytes(b []byte) *TheContentAction {
	return &TheContentAction{
		content:     bytes.Clone(b),
		contentKind: "bytes",
		hasContent:  true,
	}
}

// TheContentAction is an action that writes content at the end of a file.
type TheContentAction struct {
	filename    string
	content     []byte
	contentKind string // "text", "bytes" or "content", as the constructor reads
	hasContent  bool
	contentErr  error
}

// AtTheEndOf sets the file at the end of which the content is written.
func (a *TheContentAction) AtTheEndOf(filename string) *TheContentAction {
	a.filename = filename
	return a
}

// String describes the action.
func (a *TheContentAction) String() string {
	return fmt.Sprintf("write %s at the end of %s",
		content.Describe(a.contentKind, a.content, a.contentErr), verdict.Describe(a.filename))
}

// PerformAs performs the task or the action as the provided actor.
func (a *TheContentAction) PerformAs(actor *screenplay.Actor) error {
	if a.filename == "" {
		return action.ErrFilenameNotProvided
	}
	if !a.hasContent {
		return action.ErrReaderNotProvided
	}
	if a.contentErr != nil {
		return fmt.Errorf("failed to read the content: %w", a.contentErr)
	}
	return screenplay.UseAnAbilityOf(actor).
		That(impl.UseTheFileSystem()).
		ToAppendTo(a.filename, bytes.NewReader(a.content))
}

// TheContentAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*TheContentAction)(nil)
