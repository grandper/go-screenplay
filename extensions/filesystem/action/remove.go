package action

import (
	"fmt"

	"github.com/grandper/go-screenplay/extensions/filesystem/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

// Remove creates a RemoveBuilder to build a Remove action.
func Remove() *RemoveBuilder {
	return &RemoveBuilder{}
}

// RemoveBuilder is used to build a Remove action.
type RemoveBuilder struct{}

// TheFile creates a RemoveTheFileAction.
func (b *RemoveBuilder) TheFile(fileName string) *RemoveTheFileAction {
	return &RemoveTheFileAction{filename: fileName}
}

// TheDirectory creates a RemoveTheDirectoryAction.
func (b *RemoveBuilder) TheDirectory(directoryName string) *RemoveTheDirectoryAction {
	return &RemoveTheDirectoryAction{directoryName: directoryName}
}

// RemoveTheFileAction is an action that removes a file.
type RemoveTheFileAction struct {
	filename string
}

// String describes the action.
func (a *RemoveTheFileAction) String() string {
	return fmt.Sprintf("remove the file '%s'", a.filename)
}

// PerformAs removes the file or directory.
func (a *RemoveTheFileAction) PerformAs(actor *screenplay.Actor) error {
	useTheFileSystem, err := screenplay.UseAbilityTo[*ability.UseTheFileSystemAbility]().Of(actor)
	if err != nil {
		return err
	}
	if a.filename == "" {
		return ErrFilenameNotProvided
	}
	exists, err := useTheFileSystem.FileExists(a.filename)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("file '%s' does not exist", a.filename)
	}
	return useTheFileSystem.RemoveTheFile(a.filename)
}

// Ensure RemoveTheFileAction implements the screenplay.Action interface.
var _ screenplay.Performable = (*RemoveTheFileAction)(nil)

// RemoveTheDirectoryAction is an action that removes a directory.
type RemoveTheDirectoryAction struct {
	directoryName string
}

// String describes the action.
func (a *RemoveTheDirectoryAction) String() string {
	return fmt.Sprintf("remove the directory '%s'", a.directoryName)
}

// PerformAs removes the file or directory.
func (a *RemoveTheDirectoryAction) PerformAs(actor *screenplay.Actor) error {
	useTheFileSystem, err := screenplay.UseAbilityTo[*ability.UseTheFileSystemAbility]().Of(actor)
	if err != nil {
		return err
	}
	if a.directoryName == "" {
		return ErrDirectoryNameNotProvided
	}
	exists, err := useTheFileSystem.DirectoryExists(a.directoryName)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("directory '%s' does not exist", a.directoryName)
	}
	return useTheFileSystem.RemoveDirectory(a.directoryName)
}

// Ensure RemoveTheDirectoryAction implements the screenplay.Action interface.
var _ screenplay.Performable = (*RemoveTheDirectoryAction)(nil)
