package question

import (
	"errors"
	"fmt"
	"os"

	"github.com/grandper/go-screenplay/extensions/cli/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

var (
	// ErrEnvironmentVariableNotFound is returned when the environment variable is not found.
	ErrEnvironmentVariableNotFound = errors.New("environment variable not found")
)

// EnvironmentVariableNamed represents an environment variable.
func EnvironmentVariableNamed(environmentVariableName string) *EnvironmentVariableQuestion {
	return &EnvironmentVariableQuestion{
		environmentVariableName: environmentVariableName,
	}
}

// EnvironmentVariableQuestion represents an environment variable data.
type EnvironmentVariableQuestion struct {
	environmentVariableName string
}

// AnsweredBy returns the answer that an actor provided to the question.
func (evq *EnvironmentVariableQuestion) AnsweredBy(theActor *screenplay.Actor) (any, error) {
	_, err := screenplay.UseAbilityTo[*ability.RunCLICommandsAbility]().Of(theActor)
	if err != nil {
		return nil, err
	}
	value, exist := os.LookupEnv(evq.environmentVariableName)
	if !exist {
		return nil, ErrEnvironmentVariableNotFound
	}
	return value, nil
}

// String describes the question.
func (evq *EnvironmentVariableQuestion) String() string {
	return fmt.Sprintf("environment variable named '%s'", evq.environmentVariableName)
}
