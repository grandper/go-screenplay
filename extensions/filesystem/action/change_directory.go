package action

import (
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/extensions/filesystem/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

var (
	ErrDirectoryNotProvided = errors.New("directory must be provided")
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
	return fmt.Sprintf("change the directory to %s", a.directory)
}

// PerformAs changes the current directory.
func (a *ChangeDirectoryAction) PerformAs(actor *screenplay.Actor) error {
	if a.directory == "" {
		return ErrDirectoryNotProvided
	}
	useTheFileSystem, err := screenplay.UseAbilityTo[*ability.UseTheFileSystemAbility]().Of(actor)
	if err != nil {
		return err
	}
	return useTheFileSystem.ChangeDirectory(a.directory)
}

// Ensure ChangeDirectoryAction implements the screenplay.Action interface.
var _ screenplay.Performable = (*ChangeDirectoryAction)(nil)
