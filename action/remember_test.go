package action_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action"
	"github.com/grandper/go-screenplay/v2/fixture"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestRememberAction(t *testing.T) {
	adam := screenplay.ActorNamed("Adam")
	statusCode := fixture.NewFakeQuestion("status code", 200)
	missingStatusCode := fixture.NewFailingFakeQuestion("status code", errors.New("failed to get the status code"))
	optionalData := fixture.NewFakeQuestion("optional data", nil)

	t.Run("remembers the answer to the question under the given key", func(t *testing.T) {
		require.NoError(t, adam.AttemptsTo(action.Remember(statusCode).As("statusCode")))
		assert.Equal(t, 200, adam.Recall[int]("statusCode").OrElse(0))
	})

	t.Run("fails when the question fails", func(t *testing.T) {
		require.Error(t, adam.AttemptsTo(action.Remember(missingStatusCode).As("statusCode")))
	})

	t.Run("fails when the answer is nil", func(t *testing.T) {
		require.ErrorIs(t, adam.AttemptsTo(action.Remember(optionalData).As("data")), action.ErrAnswerIsNil)
	})

	t.Run("remembers a nil answer when allowing nil", func(t *testing.T) {
		require.NoError(t, adam.AttemptsTo(action.Remember(optionalData).As("data").AllowingNil()))
		assert.True(t, adam.Recall[any]("data").IsEmpty())
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		action1 := action.Remember(statusCode).As("statusCode")
		assert.Equal(t, "remember the status code as 'statusCode'", action1.String())
	})
}
