package question

import (
	"fmt"

	"github.com/grandper/go-screenplay/v2/extensions/filesystem/internal/verdict"

	impl "github.com/grandper/go-screenplay/v2/extensions/filesystem/internal/ability"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// FileNamed returns a file with the provided name.
func FileNamed(filePath string) screenplay.Question[*File] {
	return &fileNamedQuestion{
		filePath: filePath,
	}
}

// fileNamedQuestion returns a file with the provided name.
type fileNamedQuestion struct {
	filePath string
}

// AnsweredBy returns the answer that an actor provided to the question.
func (q *fileNamedQuestion) AnsweredBy(theActor *screenplay.Actor) (*File, error) {
	useTheFileSystem := screenplay.UseAnAbilityOf(theActor).That(impl.UseTheFileSystem())
	exists, err := useTheFileSystem.ToTellIfTheFileExists(q.filePath)
	if err != nil {
		return nil, err
	}
	return NewFile(q.filePath, exists), nil
}

// String describes the question.
func (q *fileNamedQuestion) String() string {
	return fmt.Sprintf("file named %s", verdict.Describe(q.filePath))
}

// fileNamedQuestion implements the screenplay.Question interface.
var _ screenplay.Question[*File] = (*fileNamedQuestion)(nil)
