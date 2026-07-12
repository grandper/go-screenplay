package action_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/action"
	"github.com/grandper/go-screenplay/fixture"
	"github.com/grandper/go-screenplay/screenplay"
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
		assert.Contains(t, err.Error(), "over 0.100000 seconds")
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
		assert.Contains(t, err.Error(), "over 0.100000 seconds")
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		action1 := action.Eventually(openTheHomePage)
		assert.Equal(t, "eventually open the home page", action1.String())
	})
}
