package proceed_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action"
	"github.com/grandper/go-screenplay/v2/action/proceed"
	"github.com/grandper/go-screenplay/v2/fixture"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestSilentlyTo(t *testing.T) {
	t.Run("implements the stringer interface", func(t *testing.T) {
		silent := proceed.SilentlyTo(fixture.NewFakePerformable("log in", nil))
		assert.Equal(t, "proceed silently to perform an action", silent.String())
	})

	t.Run("performs the wrapped performable", func(t *testing.T) {
		performed := false
		performable := action.FromFunc("log in", func(_ *screenplay.Actor) error {
			performed = true
			return nil
		})

		adam := screenplay.ActorNamed("Adam")
		require.NoError(t, adam.AttemptsTo(proceed.SilentlyTo(performable)))
		assert.True(t, performed)
	})

	t.Run("propagates the error of the wrapped performable", func(t *testing.T) {
		boom := errors.New("boom")
		adam := screenplay.ActorNamed("Adam")

		err := adam.AttemptsTo(proceed.SilentlyTo(fixture.NewFakePerformable("stumbles", boom)))
		require.ErrorIs(t, err, boom)
	})

	t.Run("mutes the narration produced while the performable runs", func(t *testing.T) {
		recorder := fixture.NewRecorder()
		adam := narratedActor(recorder)

		task := screenplay.TaskWhere("log in",
			fixture.NewFakePerformable("enter the password", nil),
			fixture.NewFakePerformable("submit the form", nil),
		)
		require.NoError(t, adam.AttemptsTo(proceed.SilentlyTo(task)))

		events := recorder.Events()
		require.Len(t, events, 2)
		assert.Equal(t, "Adam attempts to proceed silently to perform an action", events[0].Message)
		assert.Equal(t, "Adam attempts to proceed silently to perform an action", events[1].Message)
	})

	t.Run("never narrates the description of the wrapped performable", func(t *testing.T) {
		recorder := fixture.NewRecorder()
		adam := narratedActor(recorder)

		require.NoError(
			t,
			adam.AttemptsTo(proceed.SilentlyTo(fixture.NewFakePerformable("type the password 'hunter2'", nil))),
		)

		for _, message := range recorder.Messages() {
			assert.NotContains(t, message, "hunter2")
		}
	})

	t.Run("restores the narrator once the performable is done", func(t *testing.T) {
		recorder := fixture.NewRecorder()
		adam := narratedActor(recorder)

		require.NoError(t, adam.AttemptsTo(
			proceed.SilentlyTo(fixture.NewFakePerformable("log in", nil)),
			fixture.NewFakePerformable("wave", nil),
		))

		assert.Contains(t, recorder.Messages(), "Adam attempts to wave")
	})

	t.Run("is neutralized by a production forcing all narration", func(t *testing.T) {
		recorder := fixture.NewRecorder()
		production := screenplay.NewProduction(
			screenplay.WithNarrator(screenplay.NewNarrator(recorder)),
			screenplay.WithForceAllNarration(),
		)
		adam := production.SetTheStage(screenplay.CastOfStandardActors()).ActorNamed("Adam")

		task := screenplay.TaskWhere("log in",
			fixture.NewFakePerformable("enter the password", nil),
			fixture.NewFakePerformable("submit the form", nil),
		)
		require.NoError(t, adam.AttemptsTo(proceed.SilentlyTo(task)))

		// The silenced steps narrate again because the production forces it.
		assert.Contains(t, recorder.Messages(), "Adam attempts to enter the password")
		assert.Contains(t, recorder.Messages(), "Adam attempts to submit the form")
	})

	t.Run("mutes only its own view under ConcurrentlyTo", func(t *testing.T) {
		recorder := fixture.NewRecorder()
		adam := narratedActor(recorder)

		login := screenplay.TaskWhere("log in",
			fixture.NewFakePerformable("enter the password", nil),
			fixture.NewFakePerformable("submit the form", nil),
		)
		greeting := screenplay.TaskWhere("greet", fixture.NewFakePerformable("wave", nil))

		// Each Silently mutes its own view of the actor, never the actor itself,
		// so the concurrent greeting branch still narrates while the login steps
		// stay silent. Run with -race to catch a regression.
		err := adam.AttemptsTo(proceed.ConcurrentlyTo(
			proceed.SilentlyTo(login),
			greeting,
			proceed.SilentlyTo(login),
		))
		require.NoError(t, err)

		messages := recorder.Messages()
		assert.Contains(t, messages, "Adam attempts to wave")
		assert.NotContains(t, messages, "Adam attempts to enter the password")
	})

	t.Run("implements the performable interface", func(t *testing.T) {
		assert.Implements(
			t,
			(*screenplay.Performable)(nil),
			proceed.SilentlyTo(fixture.NewFakePerformable("log in", nil)),
		)
	})

	t.Run("fails instead of panicking when the performable is nil", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		err := adam.AttemptsTo(proceed.SilentlyTo(nil))
		require.EqualError(t, err, "the performable given to SilentlyTo is nil")
	})
}

// narratedActor casts an actor from a production whose narrator records through
// the given recorder, the way a real program wires narration.
func narratedActor(recorder *fixture.Recorder) *screenplay.Actor {
	production := screenplay.NewProduction(
		screenplay.WithNarrator(screenplay.NewNarrator(recorder)))

	return production.SetTheStage(screenplay.CastOfStandardActors()).ActorNamed("Adam")
}
