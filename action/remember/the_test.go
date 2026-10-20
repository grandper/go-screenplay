package remember_test

import (
	"errors"
	"log" //nolint:depguard // the example narrates through the standard logger
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action/remember"
	"github.com/grandper/go-screenplay/v2/narrator/logadapter"
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
		require.EqualError(
			t,
			adam.AttemptsTo(remember.The(missingStatusCode).As("statusCode")),
			"failed to get the status code",
		)
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

func ExampleThe() {
	production := screenplay.NewProduction(
		screenplay.WithNarrator(screenplay.NewNarrator(logadapter.New(log.New(os.Stdout, "", 0)))))
	adam := production.SetTheStage(screenplay.CastOfStandardActors()).ActorNamed("Adam")
	statusCode := fixture.NewFakeQuestion("status code", 200)

	_ = adam.AttemptsTo(remember.The(statusCode).As("status code"))

	// Output:
	// Adam attempts to remember the status code as "status code"
}
