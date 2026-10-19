package wait_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action/see/term"
	"github.com/grandper/go-screenplay/v2/action/wait"
	"github.com/grandper/go-screenplay/v2/question"
	"github.com/grandper/go-screenplay/v2/resolution/be"
	"github.com/grandper/go-screenplay/v2/resolution/match"
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/grandper/go-screenplay/v2/screenplay/fixture"
)

func TestUntilThe(t *testing.T) {
	adam := screenplay.ActorNamed("Adam")
	formField := fixture.NewFakeQuestion("form field", "")
	failingFormField := fixture.NewFailingFakeQuestion[string](
		"form field",
		errors.New("failed to find the form field"),
	)

	t.Run("waits until a question is answered", func(t *testing.T) {
		formFieldWillChange := fixture.NewFakeQuestion("form field", "")

		go func() {
			time.Sleep(5 * time.Millisecond)

			formFieldWillChange.AnswerWith("hello world")
		}()

		require.NoError(
			t,
			adam.AttemptsTo(wait.UntilThe(formFieldWillChange).HappensTo(be.EqualTo("hello world"))),
		)
	})

	t.Run("fails if the attempt continuously failed", func(t *testing.T) {
		err := adam.AttemptsTo(wait.UntilThe(formField).HappensTo(be.EqualTo("hello world")))
		require.Error(t, err)
		require.EqualError(t, err, fmt.Sprintf(
			`Adam waited for %s without seeing that the form field happens to be equal to "hello world": got ""`,
			screenplay.DefaultTimeout,
		))
	})

	t.Run("fails when the question fails", func(t *testing.T) {
		require.Error(t, adam.AttemptsTo(wait.UntilThe(failingFormField).HappensTo(be.EqualTo("hello world"))))
	})

	t.Run("fails without waiting when the condition is built wrongly", func(t *testing.T) {
		started := time.Now()

		require.ErrorIs(t, adam.AttemptsTo(wait.UntilThe(formField).HappensTo()), term.ErrInvalidTerm)
		require.ErrorIs(t,
			adam.AttemptsTo(wait.UntilThe(formField).DoesNotHappenTo(match.PatternString("hello("))),
			screenplay.ErrInvalidResolution,
		)
		assert.Less(t, time.Since(started), screenplay.DefaultTimeout)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		action := wait.UntilThe(formField).HappensTo(be.EqualTo("finished"))
		assert.Equal(t, `wait until the form field happens to be equal to "finished"`, action.String())
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

		err := bob.AttemptsTo(wait.UntilThe(lateCondition).HappensTo(be.True()))
		require.Error(t, err)
		assert.EqualError(
			t,
			err,
			"Bob waited for 30ms without seeing that the late condition happens to be true in time",
		)
	})
}
