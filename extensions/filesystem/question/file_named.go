package question

import (
	"fmt"

	"github.com/grandper/go-screenplay/v2/internal/verdict"

	"github.com/grandper/go-screenplay/v2/extensions/filesystem/ability"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// FileNamed returns a file with the provided name.
func FileNamed(filePath string) *FileNamedQuestion {
	return &FileNamedQuestion{
		filePath: filePath,
	}
}

// FileNamedQuestion returns a file with the provided name.
type FileNamedQuestion struct {
	filePath string
}

// AnsweredBy returns the answer that an actor provided to the question.
func (q *FileNamedQuestion) AnsweredBy(theActor *screenplay.Actor) (*File, error) {
	useTheFileSystem := screenplay.UseAnAbilityOf(theActor).That(ability.UseTheFileSystem())
	exists, err := useTheFileSystem.ToTellIfTheFileExists(q.filePath)
	if err != nil {
		return nil, err
	}
	return NewFile(q.filePath, exists), nil
}

// String describes the question.
func (q *FileNamedQuestion) String() string {
	return fmt.Sprintf("file named %s", verdict.Describe(q.filePath))
}

// FileNamedQuestion implements the screenplay.Question interface.
var _ screenplay.Question[*File] = (*FileNamedQuestion)(nil)
