package remove

import (
	"fmt"

	"github.com/grandper/go-screenplay/v2/extensions/filesystem/ability"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/action"
	"github.com/grandper/go-screenplay/v2/internal/verdict"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// TheFile removes a file:
//
//	remove.TheFile("draft.txt")
func TheFile(fileName string) *TheFileAction {
	return &TheFileAction{filename: fileName}
}

// TheDirectory removes a directory.
func TheDirectory(directoryName string) *TheDirectoryAction {
	return &TheDirectoryAction{directoryName: directoryName}
}

// TheFileAction is an action that removes a file.
type TheFileAction struct {
	filename string
}

// String describes the action.
func (a *TheFileAction) String() string {
	return fmt.Sprintf("remove the file %s", verdict.Describe(a.filename))
}

// PerformAs removes the file.
func (a *TheFileAction) PerformAs(actor *screenplay.Actor) error {
	useTheFileSystem := screenplay.UseAnAbilityOf(actor).That(ability.UseTheFileSystem())
	if a.filename == "" {
		return action.ErrFilenameNotProvided
	}
	exists, err := useTheFileSystem.ToTellIfTheFileExists(a.filename)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("file %s %w", verdict.Describe(a.filename), action.ErrFileDoesNotExist)
	}
	return useTheFileSystem.ToRemoveFile(a.filename)
}

// TheFileAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*TheFileAction)(nil)

// TheDirectoryAction is an action that removes a directory.
type TheDirectoryAction struct {
	directoryName string
}

// String describes the action.
func (a *TheDirectoryAction) String() string {
	return fmt.Sprintf("remove the directory %s", verdict.Describe(a.directoryName))
}

// PerformAs removes the directory.
func (a *TheDirectoryAction) PerformAs(actor *screenplay.Actor) error {
	useTheFileSystem := screenplay.UseAnAbilityOf(actor).That(ability.UseTheFileSystem())
	if a.directoryName == "" {
		return action.ErrDirectoryNameNotProvided
	}
	exists, err := useTheFileSystem.ToTellIfTheDirectoryExists(a.directoryName)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("directory %s %w", verdict.Describe(a.directoryName), action.ErrDirectoryDoesNotExist)
	}
	return useTheFileSystem.ToRemoveDirectory(a.directoryName)
}

// TheDirectoryAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*TheDirectoryAction)(nil)
