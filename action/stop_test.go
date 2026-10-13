package action_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action"
	"github.com/grandper/go-screenplay/v2/fixture"
	"github.com/grandper/go-screenplay/v2/question"
	"github.com/grandper/go-screenplay/v2/resolution/be"
	"github.com/grandper/go-screenplay/v2/resolution/match"
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/grandper/go-screenplay/v2/term"
)

func TestStopAction(t *testing.T) {
	adam := screenplay.ActorNamed("Adam")
	formField := fixture.NewFakeQuestion("form field", "")
	failingFormField := fixture.NewFailingFakeQuestion[string](
		"form field",
		errors.New("failed to find the form field"),
	)

	t.Run("should stop until the actor press enter", func(t *testing.T) {
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
		require.NoError(t, adam.AttemptsTo(action.Stop().UntilAnInputIsProvidedBy(r)))
	})

	t.Run("should stop until a question is answered", func(t *testing.T) {
		formFieldWillChange := fixture.NewFakeQuestion("form field", "")

		go func() {
			time.Sleep(5 * time.Millisecond)

			formFieldWillChange.AnswerWith("hello world")
		}()

		require.NoError(
			t,
			adam.AttemptsTo(action.Stop().UntilThe(formFieldWillChange).HappensTo(be.EqualTo("hello world"))),
		)
	})

	t.Run("fails if the attempt continuously failed", func(t *testing.T) {
		err := adam.AttemptsTo(action.Stop().UntilThe(formField).HappensTo(be.EqualTo("hello world")))
		require.Error(t, err)
		require.EqualError(t, err, fmt.Sprintf(
			`Adam stopped for %s without seeing that the form field happens to be equal to "hello world": got ""`,
			screenplay.DefaultTimeout,
		))
	})

	t.Run("fails when the question fails", func(t *testing.T) {
		require.Error(t, adam.AttemptsTo(action.Stop().UntilThe(failingFormField).HappensTo(be.EqualTo("hello world"))))
	})

	t.Run("fails without waiting when the condition is built wrongly", func(t *testing.T) {
		started := time.Now()

		require.ErrorIs(t, adam.AttemptsTo(action.Stop().UntilThe(formField).HappensTo()), term.ErrInvalidTerm)
		require.ErrorIs(t,
			adam.AttemptsTo(action.Stop().UntilThe(formField).DoesNotHappenTo(match.PatternString("hello("))),
			screenplay.ErrInvalidResolution,
		)
		assert.Less(t, time.Since(started), screenplay.DefaultTimeout)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		action1 := action.Stop()
		assert.Equal(t, "stop until the 'enter' key is pressed", action1.String())

		action2 := action.Stop().UntilThe(formField).HappensTo(be.EqualTo("finished"))
		assert.Equal(t, `stop until the form field happens to be equal to "finished"`, action2.String())
	})

	t.Run("fails when the reader is closed", func(t *testing.T) {
		r, _, err := os.Pipe()
		require.NoError(t, err)
		require.NoError(t, r.Close())
		require.Error(t, adam.AttemptsTo(action.Stop().UntilAnInputIsProvidedBy(r)))
	})

	t.Run("stop action honors the actor context", func(t *testing.T) {
		// A reader that never delivers anything, like a stdin nobody types into.
		never, _ := io.Pipe()

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		adamWithADeadline := screenplay.ActorNamed("Adam").WithContext(ctx)

		start := time.Now()
		err := adamWithADeadline.AttemptsTo(action.Stop().UntilAnInputIsProvidedBy(never))
		require.ErrorIs(t, err, context.DeadlineExceeded)
		assert.Less(t, time.Since(start), 5*time.Second)
		assert.EqualError(t, err, "Adam stopped waiting for an input: context deadline exceeded")
	})

	t.Run("fails when the condition is met only after the timeout", func(t *testing.T) {
		production := screenplay.NewProduction(
			screenplay.WithTimeout(30*time.Millisecond),
			screenplay.WithPolling(12*time.Millisecond),
		)
		bob := production.SetTheStage(screenplay.CastOfStandardActors()).ActorNamed("Bob")

		// False during the whole window, true right after it: the last judgment,
		// made after the timeout to explain the failure, sees it met.
		start := time.Now()
		lateCondition := question.FromFunc("late condition", func(_ *screenplay.Actor) (bool, error) {
			return time.Since(start) > 30*time.Millisecond, nil
		})

		err := bob.AttemptsTo(action.Stop().UntilThe(lateCondition).HappensTo(be.True()))
		require.Error(t, err)
		assert.EqualError(
			t,
			err,
			"Bob stopped for 30ms without seeing that the late condition happens to be true in time",
		)
	})
}
