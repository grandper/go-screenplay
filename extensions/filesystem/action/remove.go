package action

import (
	"fmt"

	"github.com/grandper/go-screenplay/v2/extensions/filesystem/ability"
	"github.com/grandper/go-screenplay/v2/internal/verdict"
	"github.com/grandper/go-screenplay/v2/screenplay"
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
	return fmt.Sprintf("remove the file %s", verdict.Describe(a.filename))
}

// PerformAs removes the file.
func (a *RemoveTheFileAction) PerformAs(actor *screenplay.Actor) error {
	useTheFileSystem := screenplay.UseAnAbilityOf(actor).That(ability.UseTheFileSystem())
	if a.filename == "" {
		return ErrFilenameNotProvided
	}
	exists, err := useTheFileSystem.ToTellIfTheFileExists(a.filename)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("file %s %w", verdict.Describe(a.filename), ErrFileDoesNotExist)
	}
	return useTheFileSystem.ToRemoveFile(a.filename)
}

// RemoveTheFileAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*RemoveTheFileAction)(nil)

// RemoveTheDirectoryAction is an action that removes a directory.
type RemoveTheDirectoryAction struct {
	directoryName string
}

// String describes the action.
func (a *RemoveTheDirectoryAction) String() string {
	return fmt.Sprintf("remove the directory %s", verdict.Describe(a.directoryName))
}

// PerformAs removes the directory.
func (a *RemoveTheDirectoryAction) PerformAs(actor *screenplay.Actor) error {
	useTheFileSystem := screenplay.UseAnAbilityOf(actor).That(ability.UseTheFileSystem())
	if a.directoryName == "" {
		return ErrDirectoryNameNotProvided
	}
	exists, err := useTheFileSystem.ToTellIfTheDirectoryExists(a.directoryName)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("directory %s %w", verdict.Describe(a.directoryName), ErrDirectoryDoesNotExist)
	}
	return useTheFileSystem.ToRemoveDirectory(a.directoryName)
}

// RemoveTheDirectoryAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*RemoveTheDirectoryAction)(nil)
