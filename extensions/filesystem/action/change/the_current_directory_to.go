package change

import (
	"fmt"

	"github.com/grandper/go-screenplay/v2/extensions/filesystem/action"
	impl "github.com/grandper/go-screenplay/v2/extensions/filesystem/internal/ability"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/internal/verdict"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// TheCurrentDirectoryTo changes the current directory:
//
//	change.TheCurrentDirectoryTo("MyDocuments")
func TheCurrentDirectoryTo(directory string) screenplay.Performable {
	return &theCurrentDirectoryToAction{
		directory: directory,
	}
}

// theCurrentDirectoryToAction is an action that changes the current directory.
type theCurrentDirectoryToAction struct {
	directory string
}

// String describes the action.
func (a *theCurrentDirectoryToAction) String() string {
	return fmt.Sprintf("change the current directory to %s", verdict.Describe(a.directory))
}

// PerformAs changes the current directory.
func (a *theCurrentDirectoryToAction) PerformAs(actor *screenplay.Actor) error {
	if a.directory == "" {
		return action.ErrDirectoryNameNotProvided
	}
	return screenplay.UseAnAbilityOf(actor).That(impl.UseTheFileSystem()).ToChangeDirectory(a.directory)
}

// theCurrentDirectoryToAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*theCurrentDirectoryToAction)(nil)
