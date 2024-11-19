package action

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/grandper/go-screenplay/extensions/filesystem/ability"
	"github.com/grandper/go-screenplay/screenplay"
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
		filename: "",
		reader:   strings.NewReader(text),
	}
}

// AppendTheContent adds the content to the end of the file without altering the existing content.
func AppendTheContent(reader io.Reader) *AppendToAction {
	return &AppendToAction{
		filename: "",
		reader:   reader,
	}
}

// AppendTheBytes adds the bytes to the end of the file without altering the existing content.
func AppendTheBytes(b []byte) *AppendToAction {
	return &AppendToAction{
		filename: "",
		reader:   bytes.NewReader(b),
	}
}

// AppendToAction is an action to append content to a file.
type AppendToAction struct {
	filename string
	reader   io.Reader
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
	if a.reader == nil {
		return ErrReaderNotProvided
	}
	useTheFileSystem, err := screenplay.UseAbilityTo[*ability.UseTheFileSystemAbility]().Of(actor)
	if err != nil {
		return err
	}
	return useTheFileSystem.AppendTo(a.filename, a.reader)
}

// Ensure AppendToAction implements the screenplay.Action interface.
var _ screenplay.Performable = (*AppendToAction)(nil)
