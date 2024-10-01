package action_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/action"
	"github.com/grandper/go-screenplay/fixture"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestEventuallyAction(t *testing.T) {
	adam := screenplay.ActorNamed("Adam")
	openTheHomePage := fixture.NewFakePerformable("open the home page", nil)
	openTheHomePageButFailed := fixture.NewFakePerformable(
		"open the home page",
		errors.New("the actor failed to perform the task"),
	)

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
		require.Error(
			t,
			adam.AttemptsTo(
				action.Eventually(openTheHomePageButFailed).For(100).Milliseconds().PollingEvery(10).Milliseconds(),
			),
		)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		action1 := action.Eventually(openTheHomePage)
		assert.Equal(t, "eventually open the home page", action1.String())
	})

	t.Run("support alternative wordings", func(t *testing.T) {
		wording1 := action.Eventually(openTheHomePageButFailed).For(100).Milliseconds().Polling(10).Milliseconds()

		wording2 := action.Eventually(openTheHomePageButFailed).TryingFor(100).Milliseconds().Polling(10).Milliseconds()
		assert.Equal(t, wording1, wording2)

		wording3 := action.Eventually(openTheHomePageButFailed).
			TryingForNoLongerThan(100).
			Milliseconds().
			Polling(10).
			Milliseconds()
		assert.Equal(t, wording1, wording3)

		wording4 := action.Eventually(openTheHomePageButFailed).
			WaitingFor(100).
			Milliseconds().
			Polling(10).
			Milliseconds()
		assert.Equal(t, wording1, wording4)

		wording5 := action.Eventually(openTheHomePageButFailed).For(100).Milliseconds().PollingEvery(10).Milliseconds()
		assert.Equal(t, wording1, wording5)

		wording6 := action.Eventually(openTheHomePageButFailed).For(100).Milliseconds().TryingEvery(10).Milliseconds()
		assert.Equal(t, wording1, wording6)
	})

	t.Run("support unit wording flexibility", func(t *testing.T) {
		wording1 := action.Eventually(openTheHomePageButFailed).For(100).Millisecond().Polling(10).Millisecond()
		wording2 := action.Eventually(openTheHomePageButFailed).For(100).Milliseconds().Polling(10).Milliseconds()
		assert.Equal(t, wording1, wording2)

		wording1 = action.Eventually(openTheHomePageButFailed).For(100).Second().Polling(10).Second()
		wording2 = action.Eventually(openTheHomePageButFailed).For(100).Seconds().Polling(10).Seconds()
		assert.Equal(t, wording1, wording2)
	})
}
