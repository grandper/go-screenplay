package ability_test

import (
	"context"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/extensions/cli/ability"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestRunCLICommands(t *testing.T) {
	ctx := context.Background()

	t.Run("can execute a command", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			runCLICommands := ability.RunCLICommands()
			err := runCLICommands.ToRun(ctx, &ability.Command{
				Program: "echo",
				Args:    []string{"Hello, World!"},
			})
			require.NoError(t, err)

			responses, err := runCLICommands.ToRetrieveResponses()
			require.NoError(t, err)
			require.Len(t, responses, 1)
			require.Equal(t, []byte("Hello, World!\n"), responses[0].StdOut())
			require.Empty(t, responses[0].StdErr())
			require.Equal(t, 0, responses[0].ExitCode())
		})

		t.Run("in a workdir", func(t *testing.T) {
			tmpDir := t.TempDir()
			const testDir = "testdata"
			runCLICommands := ability.RunCLICommands()
			require.NoError(t, runCLICommands.ToRun(ctx, &ability.Command{
				Program: "mkdir",
				Args:    []string{testDir},
				Dir:     tmpDir,
			}))
			assert.DirExists(t, filepath.Join(tmpDir, testDir))
		})

		t.Run("with environment variable", func(t *testing.T) {
			runCLICommands := ability.RunCLICommands()
			err := runCLICommands.ToRun(ctx, &ability.Command{
				Program: "printenv",
				Args:    []string{"FOO"},
				Env:     map[string]string{"FOO": "bar"},
			})
			require.NoError(t, err)

			responses, err := runCLICommands.ToRetrieveResponses()
			require.NoError(t, err)
			require.Len(t, responses, 1)
			assert.Equal(t, []byte("bar\n"), responses[0].StdOut())
			assert.Empty(t, responses[0].StdErr())
			assert.Equal(t, 0, responses[0].ExitCode())
		})

		t.Run("fail when the command does not exist", func(t *testing.T) {
			runCLICommands := ability.RunCLICommands()
			err := runCLICommands.ToRun(ctx, &ability.Command{
				Program: "nonexistentcommand",
			})
			require.Error(t, err)

			responses, err := runCLICommands.ToRetrieveResponses()
			require.NoError(t, err)
			require.Len(t, responses, 1)
			require.Empty(t, responses[0].StdOut())
			require.Empty(t, responses[0].StdErr())
			require.Equal(t, -1, responses[0].ExitCode())
		})
	})

	t.Run("can execute an interactive command", func(t *testing.T) {
		t.Run("even though no input is provided", func(t *testing.T) {
			runCLICommands := ability.RunCLICommands()
			err := runCLICommands.ToRun(ctx, &ability.Command{
				Program:     "echo",
				Args:        []string{"Hello, World!"},
				Interactive: true,
			})
			require.NoError(t, err)

			time.Sleep(time.Second)

			responses, err := runCLICommands.ToRetrieveResponses()
			require.NoError(t, err)
			require.Len(t, responses, 1)
			require.Equal(t, []byte("Hello, World!\n"), responses[0].StdOut())
			require.Empty(t, responses[0].StdErr())
			require.Equal(t, 0, responses[0].ExitCode())
		})

		t.Run("and provide input", func(t *testing.T) {
			runCLICommands := ability.RunCLICommands()
			err := runCLICommands.ToRun(ctx, &ability.Command{
				Program:     "sh",
				Args:        []string{"-c", "read -p 'Enter your name: ' name; echo Hello, $name"},
				Interactive: true,
			})
			require.NoError(t, err)

			require.NoError(t, runCLICommands.ToType("Alice\n"))

			time.Sleep(time.Second)

			responses, err := runCLICommands.ToRetrieveResponses()
			require.NoError(t, err)
			require.Len(t, responses, 1)
			require.Equal(t, []byte("Hello, Alice\n"), responses[0].StdOut())
			require.Empty(t, responses[0].StdErr())
			require.Equal(t, 0, responses[0].ExitCode())
		})

		t.Run("fail to type when no command in executed", func(t *testing.T) {
			runCLICommands := ability.RunCLICommands()
			require.ErrorIs(t, runCLICommands.ToType("Some input\n"), ability.ErrNoCommandCurrentlyRunning)
		})
	})

	t.Run("should be able to forget", func(t *testing.T) {
		runCLICommands := ability.RunCLICommands()
		require.NoError(t, runCLICommands.Forget())
	})

	t.Run("captures the Wait error in the result for interactive commands", func(t *testing.T) {
		t.Run("is nil when the command succeeds", func(t *testing.T) {
			runCLICommands := ability.RunCLICommands()
			require.NoError(t, runCLICommands.ToRun(ctx, &ability.Command{
				Program:     "sh",
				Args:        []string{"-c", "exit 0"},
				Interactive: true,
			}))

			time.Sleep(200 * time.Millisecond)

			responses, err := runCLICommands.ToRetrieveResponses()
			require.NoError(t, err)
			require.Len(t, responses, 1)
			assert.NoError(t, responses[0].Err())
		})

		t.Run("is non-nil when the command exits with a non-zero code", func(t *testing.T) {
			runCLICommands := ability.RunCLICommands()
			require.NoError(t, runCLICommands.ToRun(ctx, &ability.Command{
				Program:     "sh",
				Args:        []string{"-c", "exit 42"},
				Interactive: true,
			}))

			time.Sleep(200 * time.Millisecond)

			responses, err := runCLICommands.ToRetrieveResponses()
			require.NoError(t, err)
			require.Len(t, responses, 1)
			require.Error(t, responses[0].Err())
			assert.Equal(t, 42, responses[0].ExitCode())
		})
	})

	t.Run("Type is safe to call concurrently with a finishing interactive command", func(t *testing.T) {
		runCLICommands := ability.RunCLICommands()
		err := runCLICommands.ToRun(ctx, &ability.Command{
			Program:     "sh",
			Args:        []string{"-c", "read name; echo Hello, $name"},
			Interactive: true,
		})
		require.NoError(t, err)

		// Write input and let the command finish while a concurrent reader checks stdin.
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			_ = runCLICommands.ToType("Alice\n")
		}()

		go func() {
			defer wg.Done()
			// Calling ToType concurrently races with the goroutine that sets stdin=nil on exit.
			_ = runCLICommands.ToType("Bob\n")
		}()

		wg.Wait()
	})

	t.Run("Forget ends the interactive command still running", func(t *testing.T) {
		t.Run("by closing its input, which a program reading it waits for", func(t *testing.T) {
			runCLICommands := ability.RunCLICommands()
			require.NoError(t, runCLICommands.ToRun(ctx, &ability.Command{
				Program:     "cat",
				Interactive: true,
			}))
			require.NoError(t, runCLICommands.ToType("hello"))

			start := time.Now()
			require.NoError(t, runCLICommands.Forget())
			assert.Less(t, time.Since(start), 4*time.Second, "cat exits on end of file, before the kill")

			responses, err := runCLICommands.ToRetrieveResponses()
			require.NoError(t, err)
			require.Len(t, responses, 1)
			assert.Equal(t, []byte("hello"), responses[0].StdOut())
			assert.Equal(t, 0, responses[0].ExitCode())
			require.NoError(t, responses[0].Err())

			require.ErrorIs(t, runCLICommands.ToType("again"), ability.ErrNoCommandCurrentlyRunning)
		})

		t.Run("by killing a program that ignores the end of its input", func(t *testing.T) {
			runCLICommands := ability.RunCLICommands()
			require.NoError(t, runCLICommands.ToRun(ctx, &ability.Command{
				Program:     "sh",
				Args:        []string{"-c", "trap '' PIPE; while true; do sleep 1; done"},
				Interactive: true,
			}))

			require.NoError(t, runCLICommands.Forget())

			responses, err := runCLICommands.ToRetrieveResponses()
			require.NoError(t, err)
			require.Len(t, responses, 1)
			assert.Error(t, responses[0].Err(), "the result records the kill")
		})

		t.Run("and is a no-op when no command is running", func(t *testing.T) {
			runCLICommands := ability.RunCLICommands()
			require.NoError(t, runCLICommands.Forget())
			require.NoError(t, runCLICommands.Forget())
		})
	})

	t.Run("reports an output larger than the buffered size", func(t *testing.T) {
		runCLICommands := ability.RunCLICommands()
		err := runCLICommands.ToRun(ctx, &ability.Command{
			Program: "head",
			Args:    []string{"-c", strconv.Itoa(screenplay.DefaultMaxBufferedSize + 1), "/dev/zero"},
		})
		require.EqualError(t, err, "the output of the command exceeds 64 MiB")

		responses, err := runCLICommands.ToRetrieveResponses()
		require.NoError(t, err)
		require.Len(t, responses, 1)
		assert.Len(t, responses[0].StdOut(), screenplay.DefaultMaxBufferedSize, "the output is kept up to the bound")
		assert.Equal(t, 0, responses[0].ExitCode(), "the command itself ran to completion")
	})

	t.Run("records the error of a command exiting with a non-zero code", func(t *testing.T) {
		runCLICommands := ability.RunCLICommands()
		err := runCLICommands.ToRun(ctx, &ability.Command{Program: "sh", Args: []string{"-c", "exit 3"}})
		require.Error(t, err)

		responses, err := runCLICommands.ToRetrieveResponses()
		require.NoError(t, err)
		require.Len(t, responses, 1)
		assert.Equal(t, 3, responses[0].ExitCode())
		var exitErr *exec.ExitError
		require.ErrorAs(t, responses[0].Err(), &exitErr)
	})

	t.Run("records nothing when an interactive command cannot start", func(t *testing.T) {
		runCLICommands := ability.RunCLICommands()
		err := runCLICommands.ToRun(ctx, &ability.Command{Program: "/nonexistent/program", Interactive: true})
		require.Error(t, err)

		responses, err := runCLICommands.ToRetrieveResponses()
		require.NoError(t, err)
		assert.Empty(t, responses)
		require.ErrorIs(t, runCLICommands.ToType("hello"), ability.ErrNoCommandCurrentlyRunning)
		require.NoError(t, runCLICommands.Forget())
	})

	t.Run("fails every operation on a nil receiver", func(t *testing.T) {
		var runCLICommands *ability.RunCLICommandsAbility
		require.ErrorIs(
			t,
			runCLICommands.ToRun(ctx, &ability.Command{Program: "echo"}),
			screenplay.ErrActorMissingAbility,
		)
		require.ErrorIs(t, runCLICommands.ToType("hello"), screenplay.ErrActorMissingAbility)
		_, err := runCLICommands.ToRetrieveResponses()
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		assert.ErrorIs(t, runCLICommands.Forget(), screenplay.ErrActorMissingAbility)
	})
}
