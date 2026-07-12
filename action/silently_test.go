package action_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/action"
	"github.com/grandper/go-screenplay/fixture"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestSilentlyAction(t *testing.T) {
	t.Run("implements the stringer interface", func(t *testing.T) {
		silent := action.Silently(fixture.NewFakePerformable("log in", nil))
		assert.Equal(t, "silently log in", silent.String())
	})

	t.Run("performs the wrapped performable", func(t *testing.T) {
		performed := false
		performable := action.FromFunc("log in", func(_ *screenplay.Actor) error {
			performed = true
			return nil
		})

		adam := screenplay.ActorNamed("Adam")
		require.NoError(t, adam.AttemptsTo(action.Silently(performable)))
		assert.True(t, performed)
	})

	t.Run("propagates the error of the wrapped performable", func(t *testing.T) {
		boom := errors.New("boom")
		adam := screenplay.ActorNamed("Adam")

		err := adam.AttemptsTo(action.Silently(fixture.NewFakePerformable("stumbles", boom)))
		require.ErrorIs(t, err, boom)
	})

	t.Run("mutes the narration produced while the performable runs", func(t *testing.T) {
		recorder := fixture.NewRecorder()
		adam := narratedActor(recorder)

		task := screenplay.TaskWhere("log in",
			fixture.NewFakePerformable("enters the password", nil),
			fixture.NewFakePerformable("submits the form", nil),
		)
		require.NoError(t, adam.AttemptsTo(action.Silently(task)))

		events := recorder.Events()
		require.Len(t, events, 2)
		assert.Equal(t, "Adam silently log in", events[0].Message)
		assert.Equal(t, "Adam silently log in", events[1].Message)
	})

	t.Run("restores the narrator once the performable is done", func(t *testing.T) {
		recorder := fixture.NewRecorder()
		adam := narratedActor(recorder)

		require.NoError(t, adam.AttemptsTo(
			action.Silently(fixture.NewFakePerformable("logs in", nil)),
			fixture.NewFakePerformable("waves", nil),
		))

		assert.Contains(t, recorder.Messages(), "Adam waves")
	})

	t.Run("is neutralized by a production forcing all narration", func(t *testing.T) {
		recorder := fixture.NewRecorder()
		production := screenplay.NewProduction(
			screenplay.WithNarrator(screenplay.NewNarrator(recorder)),
			screenplay.WithForceAllNarration(),
		)
		adam := production.SetTheStage(screenplay.CastOfStandardActors()).ActorNamed("Adam")

		task := screenplay.TaskWhere("log in",
			fixture.NewFakePerformable("enters the password", nil),
			fixture.NewFakePerformable("submits the form", nil),
		)
		require.NoError(t, adam.AttemptsTo(action.Silently(task)))

		// The silenced steps narrate again because the production forces it.
		assert.Contains(t, recorder.Messages(), "Adam enters the password")
		assert.Contains(t, recorder.Messages(), "Adam submits the form")
	})

	t.Run("mutes only its own view under Concurrently", func(t *testing.T) {
		recorder := fixture.NewRecorder()
		adam := narratedActor(recorder)

		login := screenplay.TaskWhere("log in",
			fixture.NewFakePerformable("enters the password", nil),
			fixture.NewFakePerformable("submits the form", nil),
		)
		greeting := screenplay.TaskWhere("greet", fixture.NewFakePerformable("waves", nil))

		// Each Silently mutes its own view of the actor, never the actor itself,
		// so the concurrent greeting branch still narrates while the login steps
		// stay silent. Run with -race to catch a regression.
		err := adam.AttemptsTo(action.Concurrently(
			action.Silently(login),
			greeting,
			action.Silently(login),
		))
		require.NoError(t, err)

		messages := recorder.Messages()
		assert.Contains(t, messages, "Adam waves")
		assert.NotContains(t, messages, "Adam enters the password")
	})

	t.Run("implements the performable interface", func(t *testing.T) {
		assert.Implements(t, (*screenplay.Performable)(nil), action.Silently(fixture.NewFakePerformable("log in", nil)))
	})
}

// narratedActor casts an actor from a production whose narrator records through
// the given recorder, the way a real program wires narration.
func narratedActor(recorder *fixture.Recorder) *screenplay.Actor {
	production := screenplay.NewProduction(
		screenplay.WithNarrator(screenplay.NewNarrator(recorder)))

	return production.SetTheStage(screenplay.CastOfStandardActors()).ActorNamed("Adam")
}
