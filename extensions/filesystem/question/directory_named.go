package question

import (
	"fmt"

	"github.com/grandper/go-screenplay/v2/extensions/filesystem/internal/verdict"

	impl "github.com/grandper/go-screenplay/v2/extensions/filesystem/internal/ability"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// DirectoryNamed returns a directory with the provided name.
func DirectoryNamed(directoryPath string) screenplay.Question[*Directory] {
	return &directoryNamedQuestion{
		directoryPath: directoryPath,
	}
}

// directoryNamedQuestion returns a directory with the provided name.
type directoryNamedQuestion struct {
	directoryPath string
}

// AnsweredBy returns the answer that an actor provided to the question.
func (q *directoryNamedQuestion) AnsweredBy(theActor *screenplay.Actor) (*Directory, error) {
	useTheFileSystem := screenplay.UseAnAbilityOf(theActor).That(impl.UseTheFileSystem())
	exists, err := useTheFileSystem.ToTellIfTheDirectoryExists(q.directoryPath)
	if err != nil {
		return nil, err
	}
	return NewDirectory(q.directoryPath, exists), nil
}

// String describes the question.
func (q *directoryNamedQuestion) String() string {
	return fmt.Sprintf("directory named %s", verdict.Describe(q.directoryPath))
}

// directoryNamedQuestion implements the screenplay.Question interface.
var _ screenplay.Question[*Directory] = (*directoryNamedQuestion)(nil)
