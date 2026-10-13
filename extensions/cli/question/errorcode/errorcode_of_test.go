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

	t.Run("fails when the answer is a nil result", func(t *testing.T) {
		nilResult := fixture.NewFakeQuestion[*ability.Result]("nil result", nil)

		answer, err := errorcode.Of(nilResult).AnsweredBy(screenplay.ActorNamed("Adam"))

		require.ErrorIs(t, err, errorcode.ErrAnswerIsNil)
		assert.Zero(t, answer)
	})

	t.Run("fails when the wrapped question fails", func(t *testing.T) {
		failingQuestion := fixture.NewFailingFakeQuestion[*ability.Result](
			"failing question",
			errors.New("failed to get the answer"),
		)

		answer, err := errorcode.Of(failingQuestion).AnsweredBy(screenplay.ActorNamed("Adam"))

		require.Error(t, err)
		assert.Zero(t, answer)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		q := errorcode.Of(last.Of(question.Responses()))
		assert.Equal(t, "error code of last of responses", q.String())
	})

	t.Run("implements the question interface", func(t *testing.T) {
		assert.Implements(t, (*screenplay.Question[int])(nil), errorcode.Of(last.Of(question.Responses())))
	})
}
