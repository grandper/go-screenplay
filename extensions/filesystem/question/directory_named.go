package question

import (
	"fmt"

	"github.com/grandper/go-screenplay/extensions/filesystem/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

// DirectoryNamed returns a directory with the provided name.
func DirectoryNamed(directoryPath string) *DirectoryNamedQuestion {
	return &DirectoryNamedQuestion{
		directoryPath: directoryPath,
	}
}

// DirectoryNamedQuestion returns a directory with the provided name.
type DirectoryNamedQuestion struct {
	directoryPath string
}

// AnsweredBy returns the answer that an actor provided to the question.
func (q *DirectoryNamedQuestion) AnsweredBy(theActor *screenplay.Actor) (any, error) {
	useTheFileSystem, err := screenplay.UseAbilityTo[*ability.UseTheFileSystemAbility]().Of(theActor)
	if err != nil {
		return nil, err
	}
	exists, err := useTheFileSystem.DirectoryExists(q.directoryPath)
	if err != nil {
		return nil, err
	}
	return newDirectory(q.directoryPath, exists), nil
}

// String describes the question.
func (q *DirectoryNamedQuestion) String() string {
	return fmt.Sprintf("directory named '%s'", q.directoryPath)
}

// Ensure DirectoryNamedQuestion implements the screenplay.Question interface.
var _ screenplay.Question = (*DirectoryNamedQuestion)(nil)
