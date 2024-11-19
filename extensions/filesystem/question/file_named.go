package question

import (
	"fmt"

	"github.com/grandper/go-screenplay/extensions/filesystem/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

// FileNamed returns a file with the provided name.
func FileNamed(filename string) *FileNamedQuestion {
	return &FileNamedQuestion{
		filename: filename,
	}
}

// FileNamedQuestion returns a file with the provided name.
type FileNamedQuestion struct {
	filename string
}

// AnsweredBy returns the answer that an actor provided to the question.
func (q *FileNamedQuestion) AnsweredBy(theActor *screenplay.Actor) (any, error) {
	useTheFileSystem, err := screenplay.UseAbilityTo[*ability.UseTheFileSystemAbility]().Of(theActor)
	if err != nil {
		return nil, err
	}
	exists, err := useTheFileSystem.FileExists(q.filename)
	if err != nil {
		return nil, err
	}
	return newFile(q.filename, exists), nil
}

// String describes the question.
func (q *FileNamedQuestion) String() string {
	return fmt.Sprintf("file named '%s'", q.filename)
}

// Ensure FileNamedQuestion implements the screenplay.Question interface.
var _ screenplay.Question = (*FileNamedQuestion)(nil)
