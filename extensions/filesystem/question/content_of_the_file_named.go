package question

import (
	"fmt"

	"github.com/grandper/go-screenplay/v2/extensions/filesystem/internal/verdict"

	impl "github.com/grandper/go-screenplay/v2/extensions/filesystem/internal/ability"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// ContentOfTheFileNamed returns the content of the file with the provided name.
func ContentOfTheFileNamed(filename string) screenplay.Question[[]byte] {
	return &contentOfTheFileNamedQuestion{
		filename: filename,
	}
}

// contentOfTheFileNamedQuestion returns the content of a file.
type contentOfTheFileNamedQuestion struct {
	filename string
}

// AnsweredBy returns the answer that an actor provided to the question.
func (q *contentOfTheFileNamedQuestion) AnsweredBy(theActor *screenplay.Actor) ([]byte, error) {
	return screenplay.UseAnAbilityOf(theActor).That(impl.UseTheFileSystem()).ToRead(q.filename)
}

// String describes the question.
func (q *contentOfTheFileNamedQuestion) String() string {
	return fmt.Sprintf("content of the file named %s", verdict.Describe(q.filename))
}

// contentOfTheFileNamedQuestion implements the screenplay.Question interface.
var _ screenplay.Question[[]byte] = (*contentOfTheFileNamedQuestion)(nil)
