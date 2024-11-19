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
	// ErrContentNotProvided is returned when the content is not provided.
	ErrContentNotProvided = errors.New("content not provided")
)

// OverwriteTo explicitly replaces the entire content of an existing file with new data.
func OverwriteTo(filename string) *OverwriteToAction {
	return &OverwriteToAction{
		filename: filename,
	}
}

// OverwriteToAction is an action to overwrite a file.
type OverwriteToAction struct {
	filename string
	content  io.Reader
}

// WithTheContent sets the content to write into the file.
func (a *OverwriteToAction) WithTheContent(content io.Reader) *OverwriteToAction {
	a.content = content
	return a
}

// WithTheBytes sets the content to write into the file.
func (a *OverwriteToAction) WithTheBytes(b []byte) *OverwriteToAction {
	a.content = bytes.NewReader(b)
	return a
}

// WithTheText sets the content to write into the file.
func (a *OverwriteToAction) WithTheText(text string) *OverwriteToAction {
	a.content = strings.NewReader(text)
	return a
}

// String describes the action.
func (a *OverwriteToAction) String() string {
	return fmt.Sprintf("overwrite the file %s", a.filename)
}

// PerformAs performs the task or the action as the provided actor.
func (a *OverwriteToAction) PerformAs(actor *screenplay.Actor) error {
	if a.filename == "" {
		return ErrFilenameNotProvided
	}
	if a.content == nil {
		return ErrContentNotProvided
	}
	useTheFileSystem, err := screenplay.UseAbilityTo[*ability.UseTheFileSystemAbility]().Of(actor)
	if err != nil {
		return err
	}
	return useTheFileSystem.OverwriteFileWithContent(a.filename, a.content)
}

// Ensure OverwriteToAction implements the screenplay.Action interface.
var _ screenplay.Performable = (*OverwriteToAction)(nil)
