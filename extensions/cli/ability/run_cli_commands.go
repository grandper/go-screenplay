package ability

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os/exec"
	"sync"

	"github.com/grandper/go-screenplay/screenplay"
)

var (
	// ErrNoCommandCurrentlyRunning is returned when there is no command currently running.
	ErrNoCommandCurrentlyRunning = errors.New("no command currently running")
)

// RunCLICommands represents the ability of an actor to run CLI commands.
func RunCLICommands() *RunCLICommandsAbility {
	return &RunCLICommandsAbility{
		responses: nil,
	}
}

// RunCLICommandsAbility represents the ability of an actor to run CLI commands.
type RunCLICommandsAbility struct {
	mutex      sync.RWMutex
	responses  []*Result
	currentCmd *exec.Cmd
	stdin      io.WriteCloser
}

// Type types the content in the input of an interactive command.
func (a *RunCLICommandsAbility) Type(input string) error {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	if a.stdin == nil {
		return ErrNoCommandCurrentlyRunning
	}
	_, err := io.WriteString(a.stdin, input)
	return err
}

// Run runs a CLI command.
func (a *RunCLICommandsAbility) Run(ctx context.Context, cmd *Command) error {
	c := exec.CommandContext(ctx, cmd.Program, cmd.Args...) // #nosec G204 -- cmd is controlled by the test framework
	if cmd.Dir != "" {
		c.Dir = cmd.Dir
	}
	if cmd.Env != nil {
		env := make([]string, 0, len(cmd.Env))
		for k, v := range cmd.Env {
			env = append(env, k+"="+v)
		}
		c.Env = env
	}

	var outBuf, errBuf bytes.Buffer
	c.Stdout = &outBuf
	c.Stderr = &errBuf

	if !cmd.Interactive {
		err := c.Run()
		result := &Result{
			exitCode: c.ProcessState.ExitCode(),
			stdOut:   outBuf.Bytes(),
			stdErr:   errBuf.Bytes(),
		}
		a.mutex.Lock()
		defer a.mutex.Unlock()
		a.responses = append(a.responses, result)
		return err
	}

	stdinPipe, err := c.StdinPipe()
	if err != nil {
		return err
	}
	a.mutex.Lock()
	a.stdin = stdinPipe
	a.currentCmd = c
	a.mutex.Unlock()

	err = c.Start()

	go func() {
		_ = c.Wait()
		a.mutex.Lock()
		defer a.mutex.Unlock()
		result := &Result{
			exitCode: c.ProcessState.ExitCode(),
			stdOut:   outBuf.Bytes(),
			stdErr:   errBuf.Bytes(),
		}
		a.stdin = nil
		a.currentCmd = nil
		a.responses = append(a.responses, result)
	}()

	return err
}

// Responses returns the previous responses.
func (a *RunCLICommandsAbility) Responses() []*Result {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.responses
}

// Forget clean up the ability.
// The ability cannot be used after Forget() has been called.
// This method is used, e.g., to close connections to databases,
// deleting data, closing client cleanly.
func (a *RunCLICommandsAbility) Forget() error {
	return nil
}

// Ensure RunCLICommandsAbility implements the Ability interface.
var _ screenplay.Ability = (*RunCLICommandsAbility)(nil)
