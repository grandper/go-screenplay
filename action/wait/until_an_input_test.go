package wait_test

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action/wait"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestUntilAnInput(t *testing.T) {
	adam := screenplay.ActorNamed("Adam")

	t.Run("waits until an input is provided", func(t *testing.T) {
		r, w, err := os.Pipe()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, r.Close())
		}()

		input := []byte("\n")
		go func() {
			_, err = w.Write(input)
			assert.NoError(t, err)
			assert.NoError(t, w.Close())
		}()
		require.NoError(t, adam.AttemptsTo(wait.UntilAnInputIsProvidedBy(r)))
	})

	t.Run("fails when the reader is closed", func(t *testing.T) {
		r, _, err := os.Pipe()
		require.NoError(t, err)
		require.NoError(t, r.Close())
		require.Error(t, adam.AttemptsTo(wait.UntilAnInputIsProvidedBy(r)))
	})

	t.Run("honors the actor context", func(t *testing.T) {
		// A reader that never delivers anything, like a stdin nobody types into.
		never, _ := io.Pipe()

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		adamWithADeadline := screenplay.ActorNamed("Adam").WithContext(ctx)

		start := time.Now()
		err := adamWithADeadline.AttemptsTo(wait.UntilAnInputIsProvidedBy(never))
		require.ErrorIs(t, err, context.DeadlineExceeded)
		assert.Less(t, time.Since(start), 5*time.Second)
		assert.EqualError(t, err, "Adam stopped waiting for an input: context deadline exceeded")
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "wait until the enter key is pressed", wait.UntilTheEnterKeyIsPressed().String())

		r, _, err := os.Pipe()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, r.Close())
		}()
		assert.Equal(t, "wait until an input is provided", wait.UntilAnInputIsProvidedBy(r).String())
	})
}
