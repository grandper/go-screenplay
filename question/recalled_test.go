package question_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action"
	"github.com/grandper/go-screenplay/v2/fixture"
	"github.com/grandper/go-screenplay/v2/question"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestRecalled(t *testing.T) {
	adam := screenplay.ActorNamed("Adam")

	t.Run("answers with the value remembered under the given key", func(t *testing.T) {
		adam.Remember("statusCode", 200)

		answer, err := adam.AsksFor(question.Recalled[int]("statusCode"))
		require.NoError(t, err)
		assert.Equal(t, 200, answer)
	})

	t.Run("answers with nil when nothing is remembered under the given key", func(t *testing.T) {
		answer, err := adam.AsksFor(question.Recalled[any]("unknown"))
		require.NoError(t, err)
		assert.Nil(t, answer)
	})

	t.Run("answers with the zero value when nothing is remembered under the given key", func(t *testing.T) {
		answer, err := adam.AsksFor(question.Recalled[int]("unknown"))
		require.NoError(t, err)
		assert.Equal(t, 0, answer)
	})

	t.Run("fails when the remembered value is not of the expected type", func(t *testing.T) {
		adam.Remember("statusCode", 200)

		answer, err := adam.AsksFor(question.Recalled[string]("statusCode"))
		require.ErrorIs(t, err, question.ErrRecalledValueHasWrongType)
		assert.Empty(t, answer)
	})

	t.Run("answers with the latest remembered value", func(t *testing.T) {
		adam.Remember("counter", 1)
		adam.Remember("counter", 2)

		answer, err := adam.AsksFor(question.Recalled[int]("counter"))
		require.NoError(t, err)
		assert.Equal(t, 2, answer)
	})

	t.Run("is the counterpart of the Remember action", func(t *testing.T) {
		userID := fixture.NewFakeQuestion("user id", 1234)
		require.NoError(t, adam.AttemptsTo(action.Remember(userID).As("userID")))

		answer, err := adam.AsksFor(question.Recalled[int]("userID"))
		require.NoError(t, err)
		assert.Equal(t, 1234, answer)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "recalled 'statusCode'", question.Recalled[int]("statusCode").String())
	})

	t.Run("implements the question interface", func(t *testing.T) {
		assert.Implements(t, (*screenplay.Question[int])(nil), question.Recalled[int]("statusCode"))
	})
}
