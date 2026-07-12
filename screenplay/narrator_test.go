package screenplay_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/fixture"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestNarrator(t *testing.T) {
	t.Parallel()

	t.Run("the zero value is a valid, silent narrator", func(t *testing.T) {
		t.Parallel()

		var narrator screenplay.Narrator

		assert.NotPanics(t, func() {
			narrator.WhispersTheAside("Adam", "just a comment")
			_ = narrator.StatesTheFact("Adam", "does something", func() error { return nil })
		})
	})

	t.Run("a nil narrator still runs the beat", func(t *testing.T) {
		t.Parallel()

		var narrator *screenplay.Narrator
		performed := false

		err := narrator.StatesTheFact("Adam", "does something", func() error {
			performed = true
			return nil
		})

		require.NoError(t, err)
		assert.True(t, performed)
	})

	t.Run("fans out every event to every attached adapter", func(t *testing.T) {
		t.Parallel()

		first := fixture.NewRecorder()
		second := fixture.NewRecorder()
		narrator := screenplay.NewNarrator(first)
		narrator.AttachAdapter(second)

		narrator.WhispersTheAside("Adam", "hello")

		assert.Equal(t, []string{"hello"}, first.Messages())
		assert.Equal(t, []string{"hello"}, second.Messages())
	})

	t.Run("a beat announces itself, runs, then reports its outcome", func(t *testing.T) {
		t.Parallel()

		recorder := fixture.NewRecorder()
		narrator := screenplay.NewNarrator(recorder)

		err := narrator.StatesTheFact("Adam", "Adam does something", func() error { return nil })
		require.NoError(t, err)

		events := recorder.Events()
		require.Len(t, events, 2)

		assert.Equal(t, screenplay.PhaseBegin, events[0].Phase)
		assert.Equal(t, screenplay.KindBeat, events[0].Kind)
		assert.Equal(t, screenplay.LevelInfo, events[0].Level)

		assert.Equal(t, screenplay.PhaseEnd, events[1].Phase)
		assert.Equal(t, screenplay.LevelInfo, events[1].Level)
		assert.NoError(t, events[1].Err)
	})

	t.Run("a failing beat is reported at the error level", func(t *testing.T) {
		t.Parallel()

		recorder := fixture.NewRecorder()
		narrator := screenplay.NewNarrator(recorder)
		boom := errors.New("boom")

		err := narrator.StatesTheFact("Adam", "Adam does something", func() error { return boom })
		require.ErrorIs(t, err, boom)

		events := recorder.Events()
		require.Len(t, events, 2)
		assert.Equal(t, screenplay.LevelError, events[1].Level)
		assert.ErrorIs(t, events[1].Err, boom)
	})

	t.Run("nested beats increase the depth", func(t *testing.T) {
		t.Parallel()

		recorder := fixture.NewRecorder()
		narrator := screenplay.NewNarrator(recorder)

		err := narrator.StatesTheFact("Adam", "outer", func() error {
			return narrator.StatesTheFact("Adam", "inner", func() error { return nil })
		})
		require.NoError(t, err)

		events := recorder.Events()
		require.Len(t, events, 4)
		assert.Equal(t, 0, events[0].Depth, "outer begin")
		assert.Equal(t, 1, events[1].Depth, "inner begin")
		assert.Equal(t, 1, events[2].Depth, "inner end")
		assert.Equal(t, 0, events[3].Depth, "outer end")
	})
}

func TestKindString(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "act", screenplay.KindAct.String())
	assert.Equal(t, "scene", screenplay.KindScene.String())
	assert.Equal(t, "beat", screenplay.KindBeat.String())
	assert.Equal(t, "aside", screenplay.KindAside.String())
	assert.Equal(t, "unknown", screenplay.Kind(42).String())
}

func TestActorNarration(t *testing.T) {
	t.Parallel()

	t.Run("actors are silent by default", func(t *testing.T) {
		t.Parallel()

		adam := screenplay.ActorNamed("Adam")

		assert.NotPanics(t, func() {
			_ = adam.AttemptsTo(fixture.NewFakePerformable("does something", nil))
		})
	})

	t.Run("narrates every attempted action", func(t *testing.T) {
		t.Parallel()

		recorder := fixture.NewRecorder()
		adam := narratedActor(recorder)

		require.NoError(t, adam.AttemptsTo(fixture.NewFakePerformable("waves", nil)))

		events := recorder.Events()
		require.Len(t, events, 2)
		assert.Equal(t, "Adam waves", events[0].Message)
		assert.Equal(t, "Adam", events[0].Actor)
	})

	t.Run("narrates the steps of a task as nested beats", func(t *testing.T) {
		t.Parallel()

		recorder := fixture.NewRecorder()
		adam := narratedActor(recorder)

		task := screenplay.TaskWhere("greet everyone",
			fixture.NewFakePerformable("says hello", nil),
			fixture.NewFakePerformable("waves", nil),
		)
		require.NoError(t, adam.AttemptsTo(task))

		events := recorder.Events()
		require.Len(t, events, 6)

		assert.Equal(t, "Adam greet everyone", events[0].Message)
		assert.Equal(t, 0, events[0].Depth)
		assert.Equal(t, "Adam says hello", events[1].Message)
		assert.Equal(t, 1, events[1].Depth)
		assert.Equal(t, "Adam waves", events[3].Message)
		assert.Equal(t, 1, events[3].Depth)
	})

	t.Run("narrates a failing action at the error level", func(t *testing.T) {
		t.Parallel()

		recorder := fixture.NewRecorder()
		adam := narratedActor(recorder)
		boom := errors.New("boom")

		require.Error(t, adam.AttemptsTo(fixture.NewFakePerformable("stumbles", boom)))

		events := recorder.Events()
		require.Len(t, events, 2)
		assert.Equal(t, screenplay.PhaseEnd, events[1].Phase)
		assert.Equal(t, screenplay.LevelError, events[1].Level)
		assert.ErrorIs(t, events[1].Err, boom)
	})

	t.Run("narrates questions and their answers", func(t *testing.T) {
		t.Parallel()

		recorder := fixture.NewRecorder()
		adam := narratedActor(recorder)

		answer, err := adam.Sees(fixture.NewFakeQuestion("the greeting", "hello world"))
		require.NoError(t, err)
		assert.Equal(t, "hello world", answer)

		events := recorder.Events()
		require.Len(t, events, 2)

		assert.Equal(t, screenplay.KindAside, events[0].Kind)
		assert.Equal(t, "Adam asks for the greeting", events[0].Message)

		assert.Equal(t, screenplay.PhaseEnd, events[1].Phase)
		assert.Equal(t, "hello world", events[1].Answer)
	})
}

// narratedActor casts an actor from a production whose narrator records through
// the given recorder, the way a real program wires narration.
func narratedActor(recorder *fixture.Recorder) *screenplay.Actor {
	production := screenplay.NewProduction(
		screenplay.WithNarrator(screenplay.NewNarrator(recorder)))

	return production.SetTheStage(screenplay.CastOfStandardActors()).ActorNamed("Adam")
}
