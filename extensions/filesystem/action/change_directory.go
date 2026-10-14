package action

import (
	"fmt"

	"github.com/grandper/go-screenplay/v2/extensions/filesystem/ability"
	"github.com/grandper/go-screenplay/v2/internal/verdict"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// ChangeDirectoryTo creates an action to change the current directory.
func ChangeDirectoryTo(directory string) *ChangeDirectoryAction {
	return &ChangeDirectoryAction{
		directory: directory,
	}
}

// ChangeDirectoryAction is an action to change the current directory.
type ChangeDirectoryAction struct {
	directory string
}

// String describes the action.
func (a *ChangeDirectoryAction) String() string {
	return fmt.Sprintf("change directory to %s", verdict.Describe(a.directory))
}

// PerformAs changes the current directory.
func (a *ChangeDirectoryAction) PerformAs(actor *screenplay.Actor) error {
	if a.directory == "" {
		return ErrDirectoryNameNotProvided
	}
	return screenplay.UseAnAbilityOf(actor).That(ability.UseTheFileSystem()).ToChangeDirectory(a.directory)
}

// ChangeDirectoryAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*ChangeDirectoryAction)(nil)
