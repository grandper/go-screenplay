package create

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

// TheFile creates a file:
//
//	create.TheFile("notes.txt").ContainingTheText("Hello World")
func TheFile(fileName string) *TheFileAction {
	return &TheFileAction{filename: fileName}
}

// TheTemporaryFile creates a temporary file, whose name is built from the
// pattern given.
func TheTemporaryFile(fileName string) *TheFileAction {
	return &TheFileAction{
		filename:  fileName,
		temporary: true,
	}
}

// TheDirectory creates a directory.
func TheDirectory(directoryName string) *TheDirectoryAction {
	return &TheDirectoryAction{directoryName: directoryName}
}

// TheTemporaryDirectory creates a temporary directory, whose name is built
// from the pattern given.
func TheTemporaryDirectory(directoryName string) *TheDirectoryAction {
	return &TheDirectoryAction{
		directoryName: directoryName,
		temporary:     true,
	}
}

// TheFileAction is an action to create a file.
type TheFileAction struct {
	filename    string
	content     []byte
	contentKind string // "text", "bytes" or "content", as the method reads
	hasContent  bool
	contentErr  error
	temporary   bool
	placeholder *string
}

// Containing sets the content to write into the file. The reader is buffered
// once, so the action can be performed more than once.
func (a *TheFileAction) Containing(reader io.Reader) *TheFileAction {
	if reader == nil {
		return a
	}

	a.content, a.contentErr = content.Read(reader)
	a.contentKind = "content"
	a.hasContent = true
	return a
}

// ContainingBytes sets the content to write into the file.
func (a *TheFileAction) ContainingBytes(b []byte) *TheFileAction {
	a.content = bytes.Clone(b)
	a.contentKind = "bytes"
	a.hasContent = true
	return a
}

// ContainingTheText sets the content to write into the file.
func (a *TheFileAction) ContainingTheText(text string) *TheFileAction {
	a.content = []byte(text)
	a.contentKind = "text"
	a.hasContent = true
	return a
}

// AndSaveNameTo saves the created filename into the provided placeholder.
func (a *TheFileAction) AndSaveNameTo(placeholder *string) *TheFileAction {
	a.placeholder = placeholder
	return a
}

// String describes the action.
func (a *TheFileAction) String() string {
	description := fmt.Sprintf("create the file %s", verdict.Describe(a.filename))
	if a.temporary {
		description = fmt.Sprintf("create the temporary file %s", verdict.Describe(a.filename))
	}
	if a.hasContent {
		description += " containing " + content.Describe(a.contentKind, a.content, a.contentErr)
	}
	return description
}

// PerformAs performs the task or the action as the provided actor.
func (a *TheFileAction) PerformAs(actor *screenplay.Actor) error {
	if a.filename == "" {
		return action.ErrFilenameNotProvided
	}
	if a.contentErr != nil {
		return fmt.Errorf("failed to read the content: %w", a.contentErr)
	}
	useTheFileSystem := screenplay.UseAnAbilityOf(actor).That(impl.UseTheFileSystem())
	if !a.temporary {
		if a.hasContent {
			return useTheFileSystem.ToCreateFileWithContent(a.filename, bytes.NewReader(a.content))
		}
		return useTheFileSystem.ToCreateFile(a.filename)
	}
	if a.hasContent {
		name, errCreate := useTheFileSystem.ToCreateTemporaryFileWithContent(a.filename, bytes.NewReader(a.content))
		if errCreate != nil {
			return errCreate
		}
		if a.placeholder != nil {
			*a.placeholder = name
		}
		return nil
	}
	name, err := useTheFileSystem.ToCreateTemporaryFile(a.filename)
	if err != nil {
		return err
	}
	if a.placeholder != nil {
		*a.placeholder = name
	}
	return nil
}

// TheFileAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*TheFileAction)(nil)

// TheDirectoryAction is an action to create a directory.
type TheDirectoryAction struct {
	directoryName string
	temporary     bool
	placeholder   *string
}

// AndSaveNameTo saves the created directory name into the provided placeholder.
func (a *TheDirectoryAction) AndSaveNameTo(placeholder *string) *TheDirectoryAction {
	a.placeholder = placeholder
	return a
}

// String describes the action.
func (a *TheDirectoryAction) String() string {
	if a.temporary {
		return fmt.Sprintf("create the temporary directory %s", verdict.Describe(a.directoryName))
	}
	return fmt.Sprintf("create the directory %s", verdict.Describe(a.directoryName))
}

const fileMode = 0755

// PerformAs performs the task or the action as the provided actor.
func (a *TheDirectoryAction) PerformAs(actor *screenplay.Actor) error {
	if a.directoryName == "" {
		return action.ErrDirectoryNameNotProvided
	}
	useTheFileSystem := screenplay.UseAnAbilityOf(actor).That(impl.UseTheFileSystem())
	if !a.temporary {
		return useTheFileSystem.ToCreateDirectory(a.directoryName, fileMode)
	}
	name, err := useTheFileSystem.ToCreateTemporaryDirectory(a.directoryName)
	if err != nil {
		return err
	}
	if a.placeholder != nil {
		*a.placeholder = name
	}
	return nil
}

// TheDirectoryAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*TheDirectoryAction)(nil)
