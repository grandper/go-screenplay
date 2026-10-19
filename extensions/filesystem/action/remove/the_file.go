package remove

import (
	"fmt"

	"github.com/grandper/go-screenplay/v2/extensions/filesystem/action"
	impl "github.com/grandper/go-screenplay/v2/extensions/filesystem/internal/ability"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/internal/verdict"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// TheFile removes a file:
//
//	remove.TheFile("draft.txt")
func TheFile(fileName string) screenplay.Performable {
	return &theFileAction{filename: fileName}
}

// TheDirectory removes a directory.
func TheDirectory(directoryName string) screenplay.Performable {
	return &theDirectoryAction{directoryName: directoryName}
}

// theFileAction is an action that removes a file.
type theFileAction struct {
	filename string
}

// String describes the action.
func (a *theFileAction) String() string {
	return fmt.Sprintf("remove the file %s", verdict.Describe(a.filename))
}

// PerformAs removes the file.
func (a *theFileAction) PerformAs(actor *screenplay.Actor) error {
	useTheFileSystem := screenplay.UseAnAbilityOf(actor).That(impl.UseTheFileSystem())
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

// theFileAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*theFileAction)(nil)

// theDirectoryAction is an action that removes a directory.
type theDirectoryAction struct {
	directoryName string
}

// String describes the action.
func (a *theDirectoryAction) String() string {
	return fmt.Sprintf("remove the directory %s", verdict.Describe(a.directoryName))
}

// PerformAs removes the directory.
func (a *theDirectoryAction) PerformAs(actor *screenplay.Actor) error {
	useTheFileSystem := screenplay.UseAnAbilityOf(actor).That(impl.UseTheFileSystem())
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

// theDirectoryAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*theDirectoryAction)(nil)
