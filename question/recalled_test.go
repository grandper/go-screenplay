package question_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/action"
	"github.com/grandper/go-screenplay/fixture"
	"github.com/grandper/go-screenplay/question"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestRecalled(t *testing.T) {
	adam := screenplay.ActorNamed("Adam")

	t.Run("answers with the value remembered under the given key", func(t *testing.T) {
		adam.Remember("statusCode", 200)

		answer, err := adam.AsksFor(question.Recalled("statusCode"))
		require.NoError(t, err)
		assert.Equal(t, 200, answer)
	})

	t.Run("answers with nil when nothing is remembered under the given key", func(t *testing.T) {
		answer, err := adam.AsksFor(question.Recalled("unknown"))
		require.NoError(t, err)
		assert.Nil(t, answer)
	})

	t.Run("answers with the latest remembered value", func(t *testing.T) {
		adam.Remember("counter", 1)
		adam.Remember("counter", 2)

		answer, err := adam.AsksFor(question.Recalled("counter"))
		require.NoError(t, err)
		assert.Equal(t, 2, answer)
	})

	t.Run("is the counterpart of the Remember action", func(t *testing.T) {
		userID := fixture.NewFakeQuestion("user id", 1234)
		require.NoError(t, adam.AttemptsTo(action.Remember(userID).As("userID")))

		answer, err := adam.AsksFor(question.Recalled("userID"))
		require.NoError(t, err)
		assert.Equal(t, 1234, answer)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "recalled 'statusCode'", question.Recalled("statusCode").String())
	})

	t.Run("implements the question interface", func(t *testing.T) {
		assert.Implements(t, (*screenplay.Question)(nil), question.Recalled("statusCode"))
	})
}
