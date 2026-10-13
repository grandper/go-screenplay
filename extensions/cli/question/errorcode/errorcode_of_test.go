package errorcode_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/extensions/cli/ability"
	"github.com/grandper/go-screenplay/v2/extensions/cli/action"
	"github.com/grandper/go-screenplay/v2/extensions/cli/question"
	"github.com/grandper/go-screenplay/v2/extensions/cli/question/errorcode"
	"github.com/grandper/go-screenplay/v2/fixture"
	"github.com/grandper/go-screenplay/v2/question/last"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestErrorCodeOfQuestion(t *testing.T) {
	t.Run("returns the error code of a result", func(t *testing.T) {
		theActor := screenplay.ActorNamed("Adam").WhoCan(ability.RunCLICommands())
		require.NoError(t, theActor.AttemptsTo(action.RunTheCommand("echo", "Hello World")))

		answer, err := errorcode.Of(last.Of(question.Responses())).AnsweredBy(theActor)

		require.NoError(t, err)
		assert.Equal(t, 0, answer)
	})

	t.Run("fails when the answer is not a result", func(t *testing.T) {
		notAResult := fixture.NewFakeQuestion("not a result", "hello")

		answer, err := errorcode.Of(notAResult).AnsweredBy(screenplay.ActorNamed("Adam"))

		require.ErrorIs(t, err, errorcode.ErrAnswerIsNotAResult)
		assert.Nil(t, answer)
	})

	t.Run("fails when the wrapped question fails", func(t *testing.T) {
		failingQuestion := fixture.NewFailingFakeQuestion(
			"failing question",
			errors.New("failed to get the answer"),
		)

		answer, err := errorcode.Of(failingQuestion).AnsweredBy(screenplay.ActorNamed("Adam"))

		require.Error(t, err)
		assert.Nil(t, answer)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		q := errorcode.Of(question.Responses())
		assert.Equal(t, "error code of responses", q.String())
	})
}
