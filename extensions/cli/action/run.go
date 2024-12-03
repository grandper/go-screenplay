package action

import (
	"errors"
	"fmt"
	"strings"

	"github.com/grandper/go-screenplay/extensions/cli/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

var (
	// ErrEmptyCommandName is returned when the command name is empty.
	ErrEmptyCommandName = errors.New("the command name cannot be empty")
)

// RunTheCommand runs a given command with arguments.
func RunTheCommand(name string, args ...string) *RunTheCommandAction {
	return &RunTheCommandAction{
		name:        name,
		args:        args,
		env:         map[string]string{},
		workingDir:  "",
		interactive: false,
	}
}

// RunTheCommandAction is an action to run a given command with arguments.
type RunTheCommandAction struct {
	name        string
	args        []string
	env         map[string]string
	workingDir  string
	interactive bool
}

// InTheWorkingDirectory indicates that the command should be run in a given working directory.
func (a *RunTheCommandAction) InTheWorkingDirectory(workingDir string) *RunTheCommandAction {
	a.workingDir = workingDir
	return a
}

// Interactively indicates that the command should be run interactively.
func (a *RunTheCommandAction) Interactively() *RunTheCommandAction {
	a.interactive = true
	return a
}

// WithEnv sets environment variables for the command.
func (a *RunTheCommandAction) WithEnv(env map[string]string) *RunTheCommandAction {
	a.env = env
	return a
}

// WithEnvVar adds a single environment variable for the command.
func (a *RunTheCommandAction) WithEnvVar(key, value string) *RunTheCommandAction {
	a.env[key] = value
	return a
}

// String describes the action.
func (a *RunTheCommandAction) String() string {
	str := fmt.Sprintf("run the command '%s %s'", a.name, strings.Join(a.args, " "))
	if a.interactive {
		str += " interactively"
	}
	if a.workingDir != "" {
		str += fmt.Sprintf(" in the working directory %s", a.workingDir)
	}
	return str
}

// PerformAs performs the task or the action as the provided actor.
func (a *RunTheCommandAction) PerformAs(actor *screenplay.Actor) error {
	if a.name == "" {
		return ErrEmptyCommandName
	}

	runCLICommands, err := screenplay.UseAbilityTo[*ability.RunCLICommandsAbility]().Of(actor)
	if err != nil {
		return err
	}

	cmd := &ability.Command{
		Dir:         a.workingDir,
		Program:     a.name,
		Args:        a.args,
		Env:         a.env,
		Interactive: a.interactive,
	}

	return runCLICommands.Run(actor.Context(), cmd)
}

// RunTheCommandAction implements the screenplay.Performable interface.
var _ screenplay.Performable = &RunTheCommandAction{}
