package overwrite

import (
	"bytes"
	"fmt"
	"io"

	"github.com/grandper/go-screenplay/v2/extensions/filesystem/ability"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/action"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/internal/content"
	"github.com/grandper/go-screenplay/v2/internal/verdict"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// TheFile replaces the entire content of an existing file with new data:
//
//	overwrite.TheFile("notes.txt").WithTheText("Hello World")
func TheFile(filename string) *TheFileAction {
	return &TheFileAction{
		filename: filename,
	}
}

// TheFileAction is an action that overwrites a file.
type TheFileAction struct {
	filename    string
	content     []byte
	contentKind string // "text", "bytes" or "content", as the method reads
	hasContent  bool
	contentErr  error
}

// WithTheContent sets the content to write into the file. The reader is
// buffered once, so the action can be performed more than once.
func (a *TheFileAction) WithTheContent(reader io.Reader) *TheFileAction {
	if reader == nil {
		return a
	}

	a.content, a.contentErr = content.Read(reader)
	a.contentKind = "content"
	a.hasContent = true
	return a
}

// WithTheBytes sets the content to write into the file.
func (a *TheFileAction) WithTheBytes(b []byte) *TheFileAction {
	a.content = bytes.Clone(b)
	a.contentKind = "bytes"
	a.hasContent = true
	return a
}

// WithTheText sets the content to write into the file.
func (a *TheFileAction) WithTheText(text string) *TheFileAction {
	a.content = []byte(text)
	a.contentKind = "text"
	a.hasContent = true
	return a
}

// String describes the action.
func (a *TheFileAction) String() string {
	if !a.hasContent {
		return fmt.Sprintf("overwrite the file %s", verdict.Describe(a.filename))
	}
	return fmt.Sprintf("overwrite the file %s with %s",
		verdict.Describe(a.filename), content.Describe(a.contentKind, a.content, a.contentErr))
}

// PerformAs performs the task or the action as the provided actor.
func (a *TheFileAction) PerformAs(actor *screenplay.Actor) error {
	if a.filename == "" {
		return action.ErrFilenameNotProvided
	}
	if !a.hasContent {
		return action.ErrContentNotProvided
	}
	if a.contentErr != nil {
		return fmt.Errorf("failed to read the content: %w", a.contentErr)
	}
	return screenplay.UseAnAbilityOf(actor).
		That(ability.UseTheFileSystem()).
		ToOverwriteFileWithContent(a.filename, bytes.NewReader(a.content))
}

// TheFileAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*TheFileAction)(nil)
