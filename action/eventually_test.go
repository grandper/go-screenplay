package action_test

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action"
	"github.com/grandper/go-screenplay/v2/fixture"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestEventuallyAction(t *testing.T) {
	adam := screenplay.ActorNamed("Adam")
	openTheHomePage := fixture.NewFakePerformable("open the home page", nil)

	t.Run("should perform the action when the action is fast enough", func(t *testing.T) {
		require.NoError(t, adam.AttemptsTo(action.Eventually(openTheHomePage)))
	})

	t.Run("fails if the polling is larger than a timeout", func(t *testing.T) {
		require.Error(
			t,
			adam.AttemptsTo(action.Eventually(openTheHomePage).PollingEvery(10).Seconds().For(1).Millisecond()),
		)
	})

	t.Run("fails if the polling is zero or negative instead of panicking", func(t *testing.T) {
		err := adam.AttemptsTo(action.Eventually(openTheHomePage).PollingEvery(-1).Seconds())
		require.ErrorIs(t, err, action.ErrPollingPeriodAndTimeoutMustBePositive)

		err = adam.AttemptsTo(action.Eventually(openTheHomePage).For(-1).Seconds())
		require.ErrorIs(t, err, action.ErrPollingPeriodAndTimeoutMustBePositive)
	})

	t.Run("fails when the underlying action fails", func(t *testing.T) {
		underlyingErr := errors.New("the actor failed to perform the task")
		failing := fixture.NewFakePerformable("open the home page", underlyingErr)

		err := adam.AttemptsTo(
			action.Eventually(failing).For(100).Milliseconds().PollingEvery(10).Milliseconds(),
		)
		require.Error(t, err)
		assert.ErrorIs(t, err, underlyingErr)
	})

	t.Run("deduplicates repeated errors in the timeout message", func(t *testing.T) {
		underlyingErr := errors.New("always the same error")
		failing := fixture.NewFakePerformable("failing action", underlyingErr)

		err := adam.AttemptsTo(
			action.Eventually(failing).For(100).Milliseconds().PollingEvery(10).Milliseconds(),
		)
		require.Error(t, err)

		// The error message should contain the unique error exactly once.
		errMsg := err.Error()
		first := strings.Index(errMsg, underlyingErr.Error())
		require.NotEqual(t, -1, first, "expected error message to contain the underlying error")
		second := strings.Index(errMsg[first+1:], underlyingErr.Error())
		assert.Equal(t, -1, second, "expected the underlying error to appear only once (deduplicated)")
	})

	t.Run("reads the timeout and polling from the production when not overridden", func(t *testing.T) {
		underlyingErr := errors.New("still failing")
		failing := fixture.NewFakePerformable("failing action", underlyingErr)

		production := screenplay.NewProduction(
			screenplay.WithTimeout(100*time.Millisecond),
			screenplay.WithPolling(10*time.Millisecond),
		)
		stage := production.SetTheStage(screenplay.CastOfStandardActors())
		bob := stage.ActorNamed("Bob")

		err := bob.AttemptsTo(action.Eventually(failing))
		require.Error(t, err)
		// The timeout in the message reflects the production's timeout (0.1s), not
		// the DefaultTimeout of 2s.
		assert.Contains(t, err.Error(), "over 100ms")
	})

	t.Run("explicit overrides win over the production timeout and polling", func(t *testing.T) {
		underlyingErr := errors.New("still failing")
		failing := fixture.NewFakePerformable("failing action", underlyingErr)

		production := screenplay.NewProduction(
			screenplay.WithTimeout(10*time.Second),
			screenplay.WithPolling(5*time.Second),
		)
		stage := production.SetTheStage(screenplay.CastOfStandardActors())
		bob := stage.ActorNamed("Bob")

		err := bob.AttemptsTo(
			action.Eventually(failing).For(100).Milliseconds().PollingEvery(10).Milliseconds(),
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "over 100ms")
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		action1 := action.Eventually(openTheHomePage)
		assert.Equal(t, "eventually open the home page", action1.String())
	})

	t.Run("eventually action concurrency", func(t *testing.T) {
		t.Run("one instance can be performed concurrently", func(t *testing.T) {
			// Run with -race: the state of a run is local to it, never on the action.
			var attempts atomic.Int32
			flaky := action.FromFunc("open the home page", func(_ *screenplay.Actor) error {
				if attempts.Add(1)%3 != 0 {
					return errors.New("not yet")
				}
				return nil
			})
			waitForTheHomePage := action.Eventually(flaky).For(2).Seconds().PollingEvery(10).Milliseconds()

			adam := screenplay.ActorNamed("Adam")
			require.NoError(t, adam.AttemptsTo(action.Concurrently(waitForTheHomePage, waitForTheHomePage, waitForTheHomePage)))
		})

		t.Run("stops when the actor's context ends", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()
			adam := screenplay.ActorNamed("Adam").WithContext(ctx)
			failing := fixture.NewFakePerformable("open the home page", errors.New("not yet"))

			start := time.Now()
			err := adam.AttemptsTo(action.Eventually(failing).For(5).Minutes().PollingEvery(10).Milliseconds())
			require.ErrorIs(t, err, context.DeadlineExceeded)
			assert.Less(t, time.Since(start), 5*time.Second)
			assert.EqualError(t, err, "an error occurred when Adam tried to eventually open the home page: context deadline exceeded")
		})
	})

	t.Run("returns within its window when the first attempt is slow", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		slow := action.FromFunc("open the home page", func(_ *screenplay.Actor) error {
			time.Sleep(300 * time.Millisecond)
			return errors.New("still failing")
		})

		start := time.Now()
		err := adam.AttemptsTo(action.Eventually(slow).For(50).Milliseconds().PollingEvery(10).Milliseconds())
		require.ErrorIs(t, err, action.ErrNoAttemptCompleted)
		assert.Less(t, time.Since(start), 250*time.Millisecond)
		assert.EqualError(t, err, "an error occurred when Adam tried to eventually open the home page 1 times over 50ms: no attempt completed within the timeout")
	})

	t.Run("reports the collected failures over the window", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		failing := fixture.NewFakePerformable("open the home page", errors.New("not yet"))

		err := adam.AttemptsTo(action.Eventually(failing).For(50).Milliseconds().PollingEvery(10).Milliseconds())
		require.Error(t, err)
		assert.Regexp(t, `^an error occurred when Adam tried to eventually open the home page \d+ times over 50ms: not yet$`, err.Error())
	})

	t.Run("waits a full polling interval after a failed attempt longer than the interval", func(t *testing.T) {
		var attempts atomic.Int32
		slow := action.FromFunc("open the home page", func(_ *screenplay.Actor) error {
			attempts.Add(1)
			time.Sleep(30 * time.Millisecond)
			return errors.New("not yet")
		})
		adam := screenplay.ActorNamed("Adam")

		// 30 ms attempts followed by 20 ms pauses: at most 4 attempts fit in 200 ms.
		_ = adam.AttemptsTo(action.Eventually(slow).For(200).Milliseconds().PollingEvery(20).Milliseconds())
		assert.LessOrEqual(t, attempts.Load(), int32(5))
	})

}
