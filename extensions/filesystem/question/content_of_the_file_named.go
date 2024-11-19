package question

import (
	"fmt"

	"github.com/grandper/go-screenplay/extensions/filesystem/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

// ContentOfTheFileNamed returns a file with the provided name.
func ContentOfTheFileNamed(filename string) *ContentOfTheFileNamedQuestion {
	return &ContentOfTheFileNamedQuestion{
		filename: filename,
	}
}

// ContentOfTheFileNamedQuestion returns the content of a file.
type ContentOfTheFileNamedQuestion struct {
	filename string
}

// AnsweredBy returns the answer that an actor provided to the question.
func (q *ContentOfTheFileNamedQuestion) AnsweredBy(theActor *screenplay.Actor) (any, error) {
	useTheFileSystem, err := screenplay.UseAbilityTo[*ability.UseTheFileSystemAbility]().Of(theActor)
	if err != nil {
		return nil, err
	}
	content, err := useTheFileSystem.Read(q.filename)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// String describes the question.
func (q *ContentOfTheFileNamedQuestion) String() string {
	return fmt.Sprintf("content of the file named '%s'", q.filename)
}

// Ensure FileNamedQuestion implements the screenplay.Question interface.
var _ screenplay.Question = (*ContentOfTheFileNamedQuestion)(nil)
