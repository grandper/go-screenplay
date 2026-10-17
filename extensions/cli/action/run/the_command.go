package run

import (
	"fmt"
	"strings"

	"github.com/grandper/go-screenplay/v2/internal/verdict"

	"github.com/grandper/go-screenplay/v2/extensions/cli/ability"
	"github.com/grandper/go-screenplay/v2/extensions/cli/action"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// TheCommand runs a given command with arguments:
//
//	run.TheCommand("echo", "Hello World")
func TheCommand(name string, args ...string) *TheCommandAction {
	return &TheCommandAction{
		name:        name,
		args:        args,
		env:         map[string]string{},
		workingDir:  "",
		interactive: false,
	}
}

// TheCommandAction is an action to run a given command with arguments.
type TheCommandAction struct {
	name        string
	args        []string
	env         map[string]string
	workingDir  string
	interactive bool
}

// InTheWorkingDirectory indicates that the command should be run in a given working directory.
func (a *TheCommandAction) InTheWorkingDirectory(workingDir string) *TheCommandAction {
	a.workingDir = workingDir
	return a
}

// Interactively indicates that the command should be run interactively.
func (a *TheCommandAction) Interactively() *TheCommandAction {
	a.interactive = true
	return a
}

// WithEnv sets environment variables for the command.
func (a *TheCommandAction) WithEnv(env map[string]string) *TheCommandAction {
	a.env = env
	return a
}

// WithEnvVar adds a single environment variable for the command.
func (a *TheCommandAction) WithEnvVar(key, value string) *TheCommandAction {
	if a.env == nil {
		a.env = map[string]string{}
	}
	a.env[key] = value
	return a
}

// String describes the action.
func (a *TheCommandAction) String() string {
	str := fmt.Sprintf("run the command %s", verdict.Describe(strings.Join(append([]string{a.name}, a.args...), " ")))
	if a.interactive {
		str += " interactively"
	}
	if a.workingDir != "" {
		str += fmt.Sprintf(" in the working directory %s", verdict.Describe(a.workingDir))
	}
	return str
}

// PerformAs performs the task or the action as the provided actor.
func (a *TheCommandAction) PerformAs(actor *screenplay.Actor) error {
	if a.name == "" {
		return action.ErrEmptyCommandName
	}

	cmd := &ability.Command{
		Dir:         a.workingDir,
		Program:     a.name,
		Args:        a.args,
		Env:         a.env,
		Interactive: a.interactive,
	}

	return screenplay.UseAnAbilityOf(actor).That(ability.RunCLICommands()).ToRun(actor.Context(), cmd)
}

// TheCommandAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*TheCommandAction)(nil)
