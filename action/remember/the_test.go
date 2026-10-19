package remember_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action/remember"
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/grandper/go-screenplay/v2/screenplay/fixture"
)

func TestThe(t *testing.T) {
	adam := screenplay.ActorNamed("Adam")
	statusCode := fixture.NewFakeQuestion("status code", 200)
	missingStatusCode := fixture.NewFailingFakeQuestion[int]("status code", errors.New("failed to get the status code"))
	optionalData := fixture.NewFakeQuestion[any]("optional data", nil)

	t.Run("remembers the answer to the question under the given key", func(t *testing.T) {
		require.NoError(t, adam.AttemptsTo(remember.The(statusCode).As("statusCode")))
		assert.Equal(t, 200, adam.Recall[int]("statusCode").OrElse(0))
	})

	t.Run("fails when the question fails", func(t *testing.T) {
		require.Error(t, adam.AttemptsTo(remember.The(missingStatusCode).As("statusCode")))
	})

	t.Run("fails when the answer is nil", func(t *testing.T) {
		require.ErrorIs(t, adam.AttemptsTo(remember.The(optionalData).As("data")), remember.ErrAnswerIsNil)
	})

	t.Run("remembers a nil answer when allowing nil", func(t *testing.T) {
		require.NoError(t, adam.AttemptsTo(remember.The(optionalData).As("data").AllowingNil()))
		assert.True(t, adam.Recall[any]("data").IsEmpty())
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		action1 := remember.The(statusCode).As("statusCode")
		assert.Equal(t, "remember the status code as \"statusCode\"", action1.String())

		action2 := remember.The(optionalData).As("data").AllowingNil()
		assert.Equal(t, "remember the optional data as \"data\" allowing nil", action2.String())
	})
}
