package ability

import (
	"bytes"
	"context"
	"errors"
	"fmt"
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

// ToType types the content in the input of an interactive command. On a nil
// receiver, the actor does not have the ability, it fails with
// screenplay.ErrActorMissingAbility.
func (a *RunCLICommandsAbility) ToType(input string) error {
	if a == nil {
		return screenplay.ErrActorMissingAbility
	}
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

// ToRun runs a CLI command. On a nil receiver, the actor does not have the
// ability, it fails with screenplay.ErrActorMissingAbility.
func (a *RunCLICommandsAbility) ToRun(ctx context.Context, cmd *Command) error { //nolint:funlen
	if a == nil {
		return screenplay.ErrActorMissingAbility
	}
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

	var outBuf, errBuf boundedBuffer
	c.Stdout = &outBuf
	c.Stderr = &errBuf

	if !cmd.Interactive {
		runErr := c.Run()
		err := runErr
		if err == nil {
			err = outputWithinBounds(&outBuf, &errBuf)
		}
		result := &Result{
			exitCode: c.ProcessState.ExitCode(),
			stdOut:   outBuf.Bytes(),
			stdErr:   errBuf.Bytes(),
			err:      runErr,
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
		if waitErr == nil {
			waitErr = outputWithinBounds(&outBuf, &errBuf)
		}
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

// ToRetrieveResponses returns the previous responses. On a nil receiver, the actor does
// not have the ability, it fails with screenplay.ErrActorMissingAbility.
func (a *RunCLICommandsAbility) ToRetrieveResponses() ([]*Result, error) {
	if a == nil {
		return nil, screenplay.ErrActorMissingAbility
	}
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.responses, nil
}

// Forget cleans up the ability: it ends the interactive command still running,
// if any. The command's input is closed first, which is the end of file a
// program reading its input waits for; a program still alive after that is
// killed. Forget returns once the command's result has been recorded, so the
// responses are complete when the actor exits.
// The ability cannot be used after Forget() has been called. On a nil receiver,
// the actor does not have the ability, it fails with
// screenplay.ErrActorMissingAbility.
func (a *RunCLICommandsAbility) Forget() error {
	if a == nil {
		return screenplay.ErrActorMissingAbility
	}
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

// boundedBuffer collects the output of a command up to
// screenplay.DefaultMaxBufferedSize. Beyond that the output is dropped and the
// overflow is remembered, so the command keeps running and its result reports
// the overflow, while the test process never exhausts its memory.
type boundedBuffer struct {
	buffer     bytes.Buffer
	overflowed bool
}

// Write implements io.Writer; it never fails and always accounts for the whole
// input, so the command's pipes stay drained.
func (b *boundedBuffer) Write(p []byte) (int, error) {
	written := len(p)
	room := screenplay.DefaultMaxBufferedSize - b.buffer.Len()
	if len(p) > room {
		b.overflowed = true
		p = p[:max(room, 0)]
	}
	b.buffer.Write(p)
	return written, nil
}

// Bytes returns the output collected so far, never nil.
func (b *boundedBuffer) Bytes() []byte {
	if b.buffer.Len() == 0 {
		return []byte{}
	}
	return b.buffer.Bytes()
}

// outputWithinBounds fails when the standard output or error of the command
// exceeded screenplay.DefaultMaxBufferedSize.
func outputWithinBounds(stdout, stderr *boundedBuffer) error {
	if stdout.overflowed || stderr.overflowed {
		return fmt.Errorf(
			"the output of the command exceeds %d MiB",
			screenplay.DefaultMaxBufferedSize>>20, //nolint:mnd // shifting by 20 converts bytes to MiB
		)
	}
	return nil
}
