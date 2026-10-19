package question

import (
	"errors"
	"fmt"

	"github.com/grandper/go-monad"

	impl "github.com/grandper/go-screenplay/v2/extensions/cli/internal/ability"
	"github.com/grandper/go-screenplay/v2/extensions/cli/internal/verdict"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

var (
	// ErrEnvironmentVariableNotFound is returned when the environment variable is not found.
	ErrEnvironmentVariableNotFound = errors.New("environment variable not found")
)

// EnvironmentVariableNamed represents an environment variable.
func EnvironmentVariableNamed(environmentVariableName string) screenplay.Question[string] {
	return &environmentVariableNamedQuestion{
		environmentVariableName: environmentVariableName,
	}
}

// environmentVariableNamedQuestion represents an environment variable data.
type environmentVariableNamedQuestion struct {
	environmentVariableName string
}

// AnsweredBy returns the answer that an actor provided to the question.
func (q *environmentVariableNamedQuestion) AnsweredBy(theActor *screenplay.Actor) (string, error) {
	if _, err := screenplay.UseAnAbilityOf(theActor).That(impl.RunCLICommands()).ToRetrieveResponses(); err != nil {
		return "", err
	}
	return monad.GetStringFromEnv(q.environmentVariableName).OrElseError(ErrEnvironmentVariableNotFound)
}

// String describes the question.
func (q *environmentVariableNamedQuestion) String() string {
	return fmt.Sprintf("environment variable named %s", verdict.Describe(q.environmentVariableName))
}

// environmentVariableNamedQuestion implements the screenplay.Question interface.
var _ screenplay.Question[string] = (*environmentVariableNamedQuestion)(nil)
