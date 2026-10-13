package action

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/grandper/go-screenplay/v2/extensions/filesystem/ability"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

var (
	// ErrFilenameNotProvided is returned when no filename is provided.
	ErrFilenameNotProvided = errors.New("filename must be provided")

	// ErrReaderNotProvided is returned when no reader is provided.
	ErrReaderNotProvided = errors.New("reader must be provided")
)

// AppendTheText adds the text to the end of the file without altering the existing content.
func AppendTheText(text string) *AppendToAction {
	return &AppendToAction{
		content:    []byte(text),
		hasContent: true,
	}
}

// AppendTheContent adds the content to the end of the file without altering the
// existing content. The reader is buffered once, so the action can be performed
// more than once.
func AppendTheContent(reader io.Reader) *AppendToAction {
	if reader == nil {
		return &AppendToAction{}
	}

	content, err := readTheContent(reader)

	return &AppendToAction{
		content:    content,
		hasContent: true,
		contentErr: err,
	}
}

// AppendTheBytes adds the bytes to the end of the file without altering the existing content.
func AppendTheBytes(b []byte) *AppendToAction {
	return &AppendToAction{
		content:    bytes.Clone(b),
		hasContent: true,
	}
}

// AppendToAction is an action to append content to a file.
type AppendToAction struct {
	filename   string
	content    []byte
	hasContent bool
	contentErr error
}

// To sets the filename to which the content will be appended.
func (a *AppendToAction) To(filename string) *AppendToAction {
	a.filename = filename
	return a
}

// String describes the action.
func (a *AppendToAction) String() string {
	return fmt.Sprintf("append content to the file %s", a.filename)
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
	useTheFileSystem, err := screenplay.UseAbilityTo[*ability.UseTheFileSystemAbility]().Of(actor)
	if err != nil {
		return err
	}
	return useTheFileSystem.AppendTo(a.filename, bytes.NewReader(a.content))
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
		return nil, fmt.Errorf("the content exceeds %d MiB", screenplay.DefaultMaxBufferedSize>>20)
	}
	return content, nil
}

// AppendToAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*AppendToAction)(nil)
