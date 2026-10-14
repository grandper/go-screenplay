package action

import (
	"bytes"
	"fmt"
	"io"

	"github.com/grandper/go-screenplay/v2/extensions/filesystem/ability"
	"github.com/grandper/go-screenplay/v2/internal/verdict"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// OverwriteTo explicitly replaces the entire content of an existing file with new data.
func OverwriteTo(filename string) *OverwriteToAction {
	return &OverwriteToAction{
		filename: filename,
	}
}

// OverwriteToAction is an action to overwrite a file.
type OverwriteToAction struct {
	filename    string
	content     []byte
	contentKind string // "text", "bytes" or "content", as the method reads
	hasContent  bool
	contentErr  error
}

// WithTheContent sets the content to write into the file. The reader is
// buffered once, so the action can be performed more than once.
func (a *OverwriteToAction) WithTheContent(content io.Reader) *OverwriteToAction {
	if content == nil {
		return a
	}

	a.content, a.contentErr = readTheContent(content)
	a.contentKind = "content"
	a.hasContent = true
	return a
}

// WithTheBytes sets the content to write into the file.
func (a *OverwriteToAction) WithTheBytes(b []byte) *OverwriteToAction {
	a.content = bytes.Clone(b)
	a.contentKind = "bytes"
	a.hasContent = true
	return a
}

// WithTheText sets the content to write into the file.
func (a *OverwriteToAction) WithTheText(text string) *OverwriteToAction {
	a.content = []byte(text)
	a.contentKind = "text"
	a.hasContent = true
	return a
}

// String describes the action.
func (a *OverwriteToAction) String() string {
	if !a.hasContent {
		return fmt.Sprintf("overwrite %s", verdict.Describe(a.filename))
	}
	return fmt.Sprintf("overwrite %s with %s",
		verdict.Describe(a.filename), describeContent(a.contentKind, a.content, a.contentErr))
}

// PerformAs performs the task or the action as the provided actor.
func (a *OverwriteToAction) PerformAs(actor *screenplay.Actor) error {
	if a.filename == "" {
		return ErrFilenameNotProvided
	}
	if !a.hasContent {
		return ErrContentNotProvided
	}
	if a.contentErr != nil {
		return fmt.Errorf("failed to read the content: %w", a.contentErr)
	}
	return screenplay.UseAnAbilityOf(actor).
		That(ability.UseTheFileSystem()).
		ToOverwriteFileWithContent(a.filename, bytes.NewReader(a.content))
}

// OverwriteToAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*OverwriteToAction)(nil)
