package standard_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/extensions/cli/ability"
	"github.com/grandper/go-screenplay/v2/extensions/cli/action"
	"github.com/grandper/go-screenplay/v2/extensions/cli/question"
	"github.com/grandper/go-screenplay/v2/extensions/cli/question/standard"
	"github.com/grandper/go-screenplay/v2/fixture"
	"github.com/grandper/go-screenplay/v2/question/last"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestErrorOfQuestion(t *testing.T) {
	t.Run("returns the standard error of a result", func(t *testing.T) {
		theActor := screenplay.ActorNamed("Adam").WhoCan(ability.RunCLICommands())
		require.NoError(t, theActor.AttemptsTo(action.RunTheCommand("echo", "Hello World")))

		answer, err := standard.ErrorOf(last.Of(question.Responses())).AnsweredBy(theActor)

		require.NoError(t, err)
		assert.Equal(t, []byte(""), answer)
	})

	t.Run("fails when the answer is not a result", func(t *testing.T) {
		notAResult := fixture.NewFakeQuestion("not a result", "hello")

		answer, err := standard.ErrorOf(notAResult).AnsweredBy(screenplay.ActorNamed("Adam"))

		require.ErrorIs(t, err, standard.ErrAnswerIsNotAResult)
		assert.Nil(t, answer)
	})

	t.Run("fails when the wrapped question fails", func(t *testing.T) {
		failingQuestion := fixture.NewFailingFakeQuestion(
			"failing question",
			errors.New("failed to get the answer"),
		)

		answer, err := standard.ErrorOf(failingQuestion).AnsweredBy(screenplay.ActorNamed("Adam"))

		require.Error(t, err)
		assert.Nil(t, answer)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		q := standard.ErrorOf(question.Responses())
		assert.Equal(t, "standard error of responses", q.String())
	})
}
