package question

import (
	"errors"
	"fmt"
	"os"

	"github.com/grandper/go-screenplay/v2/extensions/cli/ability"
	"github.com/grandper/go-screenplay/v2/internal/verdict"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

var (
	// ErrEnvironmentVariableNotFound is returned when the environment variable is not found.
	ErrEnvironmentVariableNotFound = errors.New("environment variable not found")
)

// EnvironmentVariableNamed represents an environment variable.
func EnvironmentVariableNamed(environmentVariableName string) *EnvironmentVariableNamedQuestion {
	return &EnvironmentVariableNamedQuestion{
		environmentVariableName: environmentVariableName,
	}
}

// EnvironmentVariableNamedQuestion represents an environment variable data.
type EnvironmentVariableNamedQuestion struct {
	environmentVariableName string
}

// AnsweredBy returns the answer that an actor provided to the question.
func (q *EnvironmentVariableNamedQuestion) AnsweredBy(theActor *screenplay.Actor) (string, error) {
	if _, err := screenplay.UseAnAbilityOf(theActor).That(ability.RunCLICommands()).ToRetrieveResponses(); err != nil {
		return "", err
	}
	value, exist := os.LookupEnv(q.environmentVariableName)
	if !exist {
		return "", ErrEnvironmentVariableNotFound
	}
	return value, nil
}

// String describes the question.
func (q *EnvironmentVariableNamedQuestion) String() string {
	return fmt.Sprintf("environment variable named %s", verdict.Describe(q.environmentVariableName))
}

// EnvironmentVariableNamedQuestion implements the screenplay.Question interface.
var _ screenplay.Question[string] = (*EnvironmentVariableNamedQuestion)(nil)
