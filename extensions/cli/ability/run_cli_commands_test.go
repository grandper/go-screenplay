package ability_test

import (
	"context"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/extensions/cli/ability"
)

func TestRunCLICommands(t *testing.T) {
	ctx := context.Background()

	t.Run("can execute a command", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			runCLICommands := ability.RunCLICommands()
			err := runCLICommands.Run(ctx, &ability.Command{
				Program: "echo",
				Args:    []string{"Hello, World!"},
			})
			require.NoError(t, err)

			responses := runCLICommands.Responses()
			require.Len(t, responses, 1)
			require.Equal(t, []byte("Hello, World!\n"), responses[0].StdOut())
			require.Empty(t, responses[0].StdErr())
			require.Equal(t, 0, responses[0].ExitCode())
		})

		t.Run("in a workdir", func(t *testing.T) {
			tmpDir := t.TempDir()
			const testDir = "testdata"
			runCLICommands := ability.RunCLICommands()
			require.NoError(t, runCLICommands.Run(ctx, &ability.Command{
				Program: "mkdir",
				Args:    []string{testDir},
				Dir:     tmpDir,
			}))
			assert.DirExists(t, filepath.Join(tmpDir, testDir))
		})

		t.Run("with environment variable", func(t *testing.T) {
			runCLICommands := ability.RunCLICommands()
			err := runCLICommands.Run(ctx, &ability.Command{
				Program: "printenv",
				Args:    []string{"FOO"},
				Env:     map[string]string{"FOO": "bar"},
			})
			require.NoError(t, err)

			responses := runCLICommands.Responses()
			require.Len(t, responses, 1)
			assert.Equal(t, []byte("bar\n"), responses[0].StdOut())
			assert.Empty(t, responses[0].StdErr())
			assert.Equal(t, 0, responses[0].ExitCode())
		})

		t.Run("fail when the command does not exist", func(t *testing.T) {
			runCLICommands := ability.RunCLICommands()
			err := runCLICommands.Run(ctx, &ability.Command{
				Program: "nonexistentcommand",
			})
			require.Error(t, err)

			responses := runCLICommands.Responses()
			require.Len(t, responses, 1)
			require.Empty(t, responses[0].StdOut())
			require.Empty(t, responses[0].StdErr())
			require.Equal(t, -1, responses[0].ExitCode())
		})
	})

	t.Run("can execute an interactive command", func(t *testing.T) {
		t.Run("even though no input is provided", func(t *testing.T) {
			runCLICommands := ability.RunCLICommands()
			err := runCLICommands.Run(ctx, &ability.Command{
				Program:     "echo",
				Args:        []string{"Hello, World!"},
				Interactive: true,
			})
			require.NoError(t, err)

			time.Sleep(time.Second)

			responses := runCLICommands.Responses()
			require.Len(t, responses, 1)
			require.Equal(t, []byte("Hello, World!\n"), responses[0].StdOut())
			require.Empty(t, responses[0].StdErr())
			require.Equal(t, 0, responses[0].ExitCode())
		})

		t.Run("and provide input", func(t *testing.T) {
			runCLICommands := ability.RunCLICommands()
			err := runCLICommands.Run(ctx, &ability.Command{
				Program:     "sh",
				Args:        []string{"-c", "read -p 'Enter your name: ' name; echo Hello, $name"},
				Interactive: true,
			})
			require.NoError(t, err)

			require.NoError(t, runCLICommands.Type("Alice\n"))

			time.Sleep(time.Second)

			responses := runCLICommands.Responses()
			require.Len(t, responses, 1)
			require.Equal(t, []byte("Hello, Alice\n"), responses[0].StdOut())
			require.Empty(t, responses[0].StdErr())
			require.Equal(t, 0, responses[0].ExitCode())
		})

		t.Run("fail to type when no command in executed", func(t *testing.T) {
			runCLICommands := ability.RunCLICommands()
			require.ErrorIs(t, runCLICommands.Type("Some input\n"), ability.ErrNoCommandCurrentlyRunning)
		})
	})

	t.Run("should be able to forget", func(t *testing.T) {
		runCLICommands := ability.RunCLICommands()
		require.NoError(t, runCLICommands.Forget())
	})

	t.Run("Type is safe to call concurrently with a finishing interactive command", func(t *testing.T) {
		runCLICommands := ability.RunCLICommands()
		err := runCLICommands.Run(ctx, &ability.Command{
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
			_ = runCLICommands.Type("Alice\n")
		}()

		go func() {
			defer wg.Done()
			// Calling Type concurrently races with the goroutine that sets stdin=nil on exit.
			_ = runCLICommands.Type("Bob\n")
		}()

		wg.Wait()
	})
}
