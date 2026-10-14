package question

import (
	"fmt"

	"github.com/grandper/go-screenplay/v2/internal/verdict"

	"github.com/grandper/go-screenplay/v2/extensions/filesystem/ability"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// ContentOfTheFileNamed returns the content of the file with the provided name.
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
func (q *ContentOfTheFileNamedQuestion) AnsweredBy(theActor *screenplay.Actor) ([]byte, error) {
	return screenplay.UseAnAbilityOf(theActor).That(ability.UseTheFileSystem()).ToRead(q.filename)
}

// String describes the question.
func (q *ContentOfTheFileNamedQuestion) String() string {
	return fmt.Sprintf("content of the file named %s", verdict.Describe(q.filename))
}

// ContentOfTheFileNamedQuestion implements the screenplay.Question interface.
var _ screenplay.Question[[]byte] = (*ContentOfTheFileNamedQuestion)(nil)
