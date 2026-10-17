package change

import (
	"fmt"

	"github.com/grandper/go-screenplay/v2/extensions/filesystem/ability"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/action"
	"github.com/grandper/go-screenplay/v2/internal/verdict"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// TheCurrentDirectoryTo changes the current directory:
//
//	change.TheCurrentDirectoryTo("MyDocuments")
func TheCurrentDirectoryTo(directory string) *TheCurrentDirectoryToAction {
	return &TheCurrentDirectoryToAction{
		directory: directory,
	}
}

// TheCurrentDirectoryToAction is an action that changes the current directory.
type TheCurrentDirectoryToAction struct {
	directory string
}

// String describes the action.
func (a *TheCurrentDirectoryToAction) String() string {
	return fmt.Sprintf("change the current directory to %s", verdict.Describe(a.directory))
}

// PerformAs changes the current directory.
func (a *TheCurrentDirectoryToAction) PerformAs(actor *screenplay.Actor) error {
	if a.directory == "" {
		return action.ErrDirectoryNameNotProvided
	}
	return screenplay.UseAnAbilityOf(actor).That(ability.UseTheFileSystem()).ToChangeDirectory(a.directory)
}

// TheCurrentDirectoryToAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*TheCurrentDirectoryToAction)(nil)
