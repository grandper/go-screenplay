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
	// ErrDirectoryNameNotProvided is returned when the directory name is not provided.
	ErrDirectoryNameNotProvided = errors.New("directory name not provided")
)

// Create creates a CreateBuild to build a Create action.
func Create() *CreateBuilder {
	return &CreateBuilder{}
}

// CreateBuilder is used to build a Create action.
type CreateBuilder struct{}

// TheFile creates a CreateTheFileAction.
func (b *CreateBuilder) TheFile(fileName string) *CreateTheFileAction {
	return &CreateTheFileAction{filename: fileName}
}

// TheTemporaryFile creates a CreateTheTemporaryFileAction.
func (b *CreateBuilder) TheTemporaryFile(fileName string) *CreateTheFileAction {
	return &CreateTheFileAction{
		filename:  fileName,
		temporary: true,
	}
}

// TheDirectory creates a CreateTheDirectoryAction.
func (b *CreateBuilder) TheDirectory(directoryName string) *CreateTheDirectoryAction {
	return &CreateTheDirectoryAction{directoryName: directoryName}
}

// TheTemporaryDirectory creates a CreateTheTemporaryDirectoryAction.
func (b *CreateBuilder) TheTemporaryDirectory(directoryName string) *CreateTheDirectoryAction {
	return &CreateTheDirectoryAction{
		directoryName: directoryName,
		temporary:     true,
	}
}

// CreateTheFileAction is an action to create a file.
type CreateTheFileAction struct {
	filename    string
	content     io.Reader
	temporary   bool
	placeholder *string
}

// Containing sets the content to write into the file.
func (a *CreateTheFileAction) Containing(content io.Reader) *CreateTheFileAction {
	a.content = content
	return a
}

// ContainingBytes sets the content to write into the file.
func (a *CreateTheFileAction) ContainingBytes(b []byte) *CreateTheFileAction {
	a.content = bytes.NewReader(b)
	return a
}

// ContainingTheText sets the content to write into the file.
func (a *CreateTheFileAction) ContainingTheText(text string) *CreateTheFileAction {
	a.content = strings.NewReader(text)
	return a
}

// AndSaveNameTo saves the created filename into the provided placeholder.
func (a *CreateTheFileAction) AndSaveNameTo(placeholder *string) *CreateTheFileAction {
	a.placeholder = placeholder
	return a
}

// String describes the action.
func (a *CreateTheFileAction) String() string {
	if a.temporary {
		return fmt.Sprintf("create the temporary file %s", a.filename)
	}
	return fmt.Sprintf("create the file %s", a.filename)
}

// PerformAs performs the task or the action as the provided actor.
func (a *CreateTheFileAction) PerformAs(actor *screenplay.Actor) error {
	if a.filename == "" {
		return ErrFilenameNotProvided
	}
	useTheFileSystem, err := screenplay.UseAbilityTo[*ability.UseTheFileSystemAbility]().Of(actor)
	if err != nil {
		return err
	}
	if !a.temporary {
		if a.content != nil {
			return useTheFileSystem.CreateFileWithContent(a.filename, a.content)
		}
		return useTheFileSystem.CreateFile(a.filename)
	}
	if a.content != nil {
		name, errCreate := useTheFileSystem.CreateTemporaryFileWithContent(a.filename, a.content)
		if errCreate != nil {
			return errCreate
		}
		if a.placeholder != nil {
			*a.placeholder = name
		}
		return nil
	}
	name, err := useTheFileSystem.CreateTemporaryFile(a.filename)
	if err != nil {
		return err
	}
	if a.placeholder != nil {
		*a.placeholder = name
	}
	return nil
}

// Ensure CreateTheFileAction implements the screenplay.Action interface.
var _ screenplay.Performable = (*CreateTheFileAction)(nil)

// CreateTheDirectoryAction is an action to create a directory.
type CreateTheDirectoryAction struct {
	directoryName string
	temporary     bool
	placeholder   *string
}

// AndSaveNameTo saves the created directory name into the provided placeholder.
func (a *CreateTheDirectoryAction) AndSaveNameTo(placeholder *string) *CreateTheDirectoryAction {
	a.placeholder = placeholder
	return a
}

// String describes the action.
func (a *CreateTheDirectoryAction) String() string {
	if a.temporary {
		return fmt.Sprintf("create the temporary directory %s", a.directoryName)
	}
	return fmt.Sprintf("create the directory %s", a.directoryName)
}

const fileMode = 0755

// PerformAs performs the task or the action as the provided actor.
func (a *CreateTheDirectoryAction) PerformAs(actor *screenplay.Actor) error {
	if a.directoryName == "" {
		return ErrDirectoryNameNotProvided
	}
	useTheFileSystem, err := screenplay.UseAbilityTo[*ability.UseTheFileSystemAbility]().Of(actor)
	if err != nil {
		return err
	}
	if !a.temporary {
		return useTheFileSystem.CreateDirectory(a.directoryName, fileMode)
	}
	name, err := useTheFileSystem.CreateTemporaryDirectory(a.directoryName)
	if err != nil {
		return err
	}
	if a.placeholder != nil {
		*a.placeholder = name
	}
	return nil
}

// Ensure CreateTheDirectoryAction implements the screenplay.Action interface.
var _ screenplay.Performable = (*CreateTheDirectoryAction)(nil)
