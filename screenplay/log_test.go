package screenplay_test

import (
	"bytes"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/fixture"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestLogAction(t *testing.T) {
	adam := screenplay.ActorNamed("Adam")
	formField := fixture.NewFakeQuestion("form field", "hello world")
	missingFormField := fixture.NewFailingFakeQuestion("form field", errors.New("failed to get the field content"))

	t.Run("should log the answer of the question", func(t *testing.T) {
		require.NoError(t, adam.AttemptsTo(screenplay.Log(formField)))
	})

	t.Run("fails when the question fails", func(t *testing.T) {
		require.Error(t, adam.AttemptsTo(screenplay.Log(missingFormField)))
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		action := screenplay.Log(formField)
		assert.Equal(t, "log the form field", action.String())
	})

	t.Run("whispers the answer as an aside", func(t *testing.T) {
		recorder := fixture.NewRecorder()
		narratedAdam := narratedActor(recorder)

		require.NoError(t, narratedAdam.AttemptsTo(screenplay.Log(formField)))

		// The action is narrated as a beat, with the value whispered underneath
		// as a nested aside — never through an "asks for" line.
		events := recorder.Events()
		require.Len(t, events, 3)
		assert.Equal(t, "Adam log the form field", events[0].Message)
		assert.Equal(t, screenplay.KindAside, events[1].Kind)
		assert.Equal(t, "=> hello world", events[1].Message)
		assert.Equal(t, 1, events[1].Depth)
	})

	t.Run("produces no output when no adapter is attached", func(t *testing.T) {
		var buffer bytes.Buffer

		// Capture the default slog logger, the sink the action historically
		// wrote to, so the test fails if any output is produced without an
		// adapter attached.
		previousDefault := slog.Default()
		slog.SetDefault(slog.New(slog.NewTextHandler(&buffer, nil)))
		t.Cleanup(func() { slog.SetDefault(previousDefault) })

		silentActor := screenplay.ActorNamed("Adam") // no narrator, no adapter

		require.NoError(t, silentActor.AttemptsTo(screenplay.Log(formField)))
		assert.Empty(t, buffer.String())
	})
}
