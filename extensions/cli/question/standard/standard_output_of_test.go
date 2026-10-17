package standard_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/extensions/cli/ability"
	"github.com/grandper/go-screenplay/v2/extensions/cli/action/run"
	"github.com/grandper/go-screenplay/v2/extensions/cli/question"
	"github.com/grandper/go-screenplay/v2/extensions/cli/question/standard"
	"github.com/grandper/go-screenplay/v2/fixture"
	"github.com/grandper/go-screenplay/v2/question/last"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestOutputOfQuestion(t *testing.T) {
	t.Run("returns the standard output of a result", func(t *testing.T) {
		theActor := screenplay.ActorNamed("Adam").WhoCan(ability.RunCLICommands())
		require.NoError(t, theActor.AttemptsTo(run.TheCommand("echo", "Hello World")))

		answer, err := standard.OutputOf(last.Of(question.Responses())).AnsweredBy(theActor)

		require.NoError(t, err)
		assert.Equal(t, []byte("Hello World\n"), answer)
	})

	t.Run("fails when the answer is a nil result", func(t *testing.T) {
		nilResult := fixture.NewFakeQuestion[*ability.Result]("nil result", nil)

		answer, err := standard.OutputOf(nilResult).AnsweredBy(screenplay.ActorNamed("Adam"))

		require.ErrorIs(t, err, standard.ErrResultIsNil)
		assert.Nil(t, answer)
	})

	t.Run("fails when the wrapped question fails", func(t *testing.T) {
		failingQuestion := fixture.NewFailingFakeQuestion[*ability.Result](
			"failing question",
			errors.New("failed to get the answer"),
		)

		answer, err := standard.OutputOf(failingQuestion).AnsweredBy(screenplay.ActorNamed("Adam"))

		require.Error(t, err)
		assert.Nil(t, answer)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		q := standard.OutputOf(last.Of(question.Responses()))
		assert.Equal(t, "standard output of last of command results", q.String())
	})

	t.Run("implements the question interface", func(t *testing.T) {
		assert.Implements(t, (*screenplay.Question[[]byte])(nil), standard.OutputOf(last.Of(question.Responses())))
	})
}
