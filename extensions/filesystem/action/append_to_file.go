package action

import (
	"bytes"
	"fmt"
	"io"

	"github.com/grandper/go-screenplay/v2/extensions/filesystem/ability"
	"github.com/grandper/go-screenplay/v2/internal/verdict"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// AppendTheText adds the text to the end of the file without altering the existing content.
func AppendTheText(text string) *AppendToAction {
	return &AppendToAction{
		content:     []byte(text),
		contentKind: "text",
		hasContent:  true,
	}
}

// AppendTheContent adds the content to the end of the file without altering the
// existing content. The reader is buffered once, so the action can be performed
// more than once.
func AppendTheContent(reader io.Reader) *AppendToAction {
	if reader == nil {
		return &AppendToAction{contentKind: "content"}
	}

	content, err := readTheContent(reader)

	return &AppendToAction{
		content:     content,
		contentKind: "content",
		hasContent:  true,
		contentErr:  err,
	}
}

// AppendTheBytes adds the bytes to the end of the file without altering the existing content.
func AppendTheBytes(b []byte) *AppendToAction {
	return &AppendToAction{
		content:     bytes.Clone(b),
		contentKind: "bytes",
		hasContent:  true,
	}
}

// AppendToAction is an action to append content to a file.
type AppendToAction struct {
	filename    string
	content     []byte
	contentKind string // "text", "bytes" or "content", as the constructor reads
	hasContent  bool
	contentErr  error
}

// To sets the filename to which the content will be appended.
func (a *AppendToAction) To(filename string) *AppendToAction {
	a.filename = filename
	return a
}

// String describes the action.
func (a *AppendToAction) String() string {
	return fmt.Sprintf("append %s to %s",
		describeContent(a.contentKind, a.content, a.contentErr), verdict.Describe(a.filename))
}

// PerformAs performs the task or the action as the provided actor.
func (a *AppendToAction) PerformAs(actor *screenplay.Actor) error {
	if a.filename == "" {
		return ErrFilenameNotProvided
	}
	if !a.hasContent {
		return ErrReaderNotProvided
	}
	if a.contentErr != nil {
		return fmt.Errorf("failed to read the content: %w", a.contentErr)
	}
	return screenplay.UseAnAbilityOf(actor).
		That(ability.UseTheFileSystem()).
		ToAppendTo(a.filename, bytes.NewReader(a.content))
}

// readTheContent buffers the content of a reader once, up to
// screenplay.DefaultMaxBufferedSize, so that an action holding it can be
// performed more than once: a reader is consumed by its first use.
func readTheContent(reader io.Reader) ([]byte, error) {
	content, err := io.ReadAll(io.LimitReader(reader, screenplay.DefaultMaxBufferedSize+1))
	if err != nil {
		return nil, err
	}
	if len(content) > screenplay.DefaultMaxBufferedSize {
		return nil, fmt.Errorf(
			"the content exceeds %d MiB",
			screenplay.DefaultMaxBufferedSize>>20, //nolint:mnd // shifting by 20 converts bytes to MiB
		)
	}
	return content, nil
}

// describeContent describes the content of an action as its constructor reads,
// followed by the content itself: the text "Hello", the bytes "Hello". A
// content that could not be read is described without its value.
func describeContent(kind string, content []byte, err error) string {
	if err != nil {
		return "the " + kind
	}

	return fmt.Sprintf("the %s %s", kind, verdict.Describe(content))
}

// AppendToAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*AppendToAction)(nil)
