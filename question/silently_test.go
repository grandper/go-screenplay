package question_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/fixture"
	"github.com/grandper/go-screenplay/v2/question"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestSilentlyQuestion(t *testing.T) {
	t.Run("implements the stringer interface", func(t *testing.T) {
		silent := question.Silently(fixture.NewFakeQuestion("the token", "secret"))
		assert.Equal(t, "a secret", silent.String())
	})

	t.Run("answers with the answer of the wrapped question", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")

		answer, err := question.Silently(fixture.NewFakeQuestion("the token", "secret")).AnsweredBy(adam)
		require.NoError(t, err)
		assert.Equal(t, "secret", answer)
	})

	t.Run("propagates the error of the wrapped question", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		failing := fixture.NewFailingFakeQuestion("the token", assert.AnError)

		_, err := question.Silently(failing).AnsweredBy(adam)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("mutes the narration produced while the question is answered", func(t *testing.T) {
		recorder := fixture.NewRecorder()
		adam := narratedActor(recorder)

		// A question whose answer is itself obtained by asking another question,
		// which normally narrates an aside and its answer.
		inner := fixture.NewFakeQuestion("the inner value", "secret")
		chatty := question.FromFunc("the token", func(theActor *screenplay.Actor) (any, error) {
			return theActor.Sees(inner)
		})

		answer, err := adam.Sees(question.Silently(chatty))
		require.NoError(t, err)
		assert.Equal(t, "secret", answer)

		// Only the outer aside is narrated; the inner question stays muted.
		events := recorder.Events()
		require.Len(t, events, 2)
		assert.Equal(t, "Adam asks for a secret", events[0].Message)
		assert.Equal(t, "Adam asks for a secret", events[1].Message)
	})

	t.Run("never narrates the answer", func(t *testing.T) {
		recorder := fixture.NewRecorder()
		adam := narratedActor(recorder)

		answer, err := adam.Sees(question.Silently(fixture.NewFakeQuestion("the token", "tok-123")))
		require.NoError(t, err)
		assert.Equal(t, "tok-123", answer)

		for _, event := range recorder.Events() {
			assert.Nil(t, event.Answer)
			assert.NotContains(t, event.Message, "tok-123")
		}
	})

	t.Run("restores the narrator once the question is answered", func(t *testing.T) {
		recorder := fixture.NewRecorder()
		adam := narratedActor(recorder)

		_, err := adam.Sees(question.Silently(fixture.NewFakeQuestion("the token", "secret")))
		require.NoError(t, err)

		_, err = adam.Sees(fixture.NewFakeQuestion("the greeting", "hello"))
		require.NoError(t, err)

		assert.Contains(t, recorder.Messages(), "Adam asks for the greeting")
	})

	t.Run("is neutralized by a production forcing all narration", func(t *testing.T) {
		recorder := fixture.NewRecorder()
		production := screenplay.NewProduction(
			screenplay.WithNarrator(screenplay.NewNarrator(recorder)),
			screenplay.WithForceAllNarration(),
		)
		adam := production.SetTheStage(screenplay.CastOfStandardActors()).ActorNamed("Adam")

		inner := fixture.NewFakeQuestion("the inner value", "secret")
		chatty := question.FromFunc("the token", func(theActor *screenplay.Actor) (any, error) {
			return theActor.Sees(inner)
		})

		answer, err := adam.Sees(question.Silently(chatty))
		require.NoError(t, err)
		assert.Equal(t, "secret", answer)

		// The inner question narrates again because the production forces it.
		assert.Contains(t, recorder.Messages(), "Adam asks for the inner value")
	})

	t.Run("implements the question interface", func(t *testing.T) {
		assert.Implements(
			t,
			(*screenplay.Question)(nil),
			question.Silently(fixture.NewFakeQuestion("the token", "secret")),
		)
	})
}

// narratedActor casts an actor from a production whose narrator records through
// the given recorder, the way a real program wires narration.
func narratedActor(recorder *fixture.Recorder) *screenplay.Actor {
	production := screenplay.NewProduction(
		screenplay.WithNarrator(screenplay.NewNarrator(recorder)))

	return production.SetTheStage(screenplay.CastOfStandardActors()).ActorNamed("Adam")
}
