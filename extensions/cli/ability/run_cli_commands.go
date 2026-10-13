package ability

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/grandper/go-screenplay/v2/screenplay"
)

var (
	// ErrNoCommandCurrentlyRunning is returned when there is no command currently running.
	ErrNoCommandCurrentlyRunning = errors.New("no command currently running")
)

// waitDelay bounds the time Wait spends on a program that does not exit after
// its context is canceled, or that exits leaving its output pipes open.
const waitDelay = 5 * time.Second

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
	// done is closed once the result of the current interactive command has
	// been recorded, so Forget can wait for it.
	done chan struct{}
}

// Type types the content in the input of an interactive command.
func (a *RunCLICommandsAbility) Type(input string) error {
	// The pipe is taken under the lock and written to after releasing it: a
	// write may block, and the goroutine recording the program's result needs
	// the lock to finish.
	a.mutex.RLock()
	stdin := a.stdin
	a.mutex.RUnlock()

	if stdin == nil {
		return ErrNoCommandCurrentlyRunning
	}
	_, err := io.WriteString(stdin, input)
	return err
}

// Run runs a CLI command.
func (a *RunCLICommandsAbility) Run(ctx context.Context, cmd *Command) error {
	c := exec.CommandContext(ctx, cmd.Program, cmd.Args...) // #nosec G204 -- cmd is controlled by the test framework
	c.WaitDelay = waitDelay
	if cmd.Dir != "" {
		c.Dir = cmd.Dir
	}
	// The child inherits the parent's environment (exec does so when Env is nil)
	// and the command's own variables are added on top of it. A later entry wins
	// over an inherited one of the same name, so a command variable overrides
	// the parent's.
	if len(cmd.Env) > 0 {
		env := os.Environ()
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
			err:      err,
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

	done := make(chan struct{})
	a.mutex.Lock()
	a.stdin = stdinPipe
	a.currentCmd = c
	a.done = done
	a.mutex.Unlock()

	if err = c.Start(); err != nil {
		// The program never ran: nothing to wait for, nothing to record.
		_ = stdinPipe.Close()
		a.mutex.Lock()
		a.stdin, a.currentCmd, a.done = nil, nil, nil
		a.mutex.Unlock()

		return err
	}

	go func() {
		defer close(done)
		waitErr := c.Wait()
		a.mutex.Lock()
		defer a.mutex.Unlock()
		result := &Result{
			exitCode: c.ProcessState.ExitCode(),
			stdOut:   outBuf.Bytes(),
			stdErr:   errBuf.Bytes(),
			err:      waitErr,
		}
		a.stdin = nil
		a.currentCmd = nil
		a.done = nil
		a.responses = append(a.responses, result)
	}()

	return nil
}

// Responses returns the previous responses.
func (a *RunCLICommandsAbility) Responses() []*Result {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.responses
}

// Forget cleans up the ability: it ends the interactive command still running,
// if any. The command's input is closed first, which is the end of file a
// program reading its input waits for; a program still alive after that is
// killed. Forget returns once the command's result has been recorded, so the
// responses are complete when the actor exits.
// The ability cannot be used after Forget() has been called.
func (a *RunCLICommandsAbility) Forget() error {
	a.mutex.RLock()
	stdin, currentCmd, done := a.stdin, a.currentCmd, a.done
	a.mutex.RUnlock()

	if currentCmd == nil {
		return nil
	}

	// Closing the input may fail when the program already exited: harmless.
	_ = stdin.Close()

	select {
	case <-done:
		return nil
	case <-time.After(waitDelay):
	}

	if currentCmd.Process != nil {
		// Kill fails when the program exited in the meantime: harmless.
		_ = currentCmd.Process.Kill()
	}
	<-done

	return nil
}

// RunCLICommandsAbility implements the screenplay.Ability interface.
var _ screenplay.Ability = (*RunCLICommandsAbility)(nil)
