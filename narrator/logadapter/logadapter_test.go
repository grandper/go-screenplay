package logadapter_test

import (
	"bytes"
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/narrator/logadapter"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestLogAdapter(t *testing.T) {
	t.Parallel()

	t.Run("announces a beat on the way in", func(t *testing.T) {
		t.Parallel()

		var buffer bytes.Buffer
		adapter := logadapter.New(log.New(&buffer, "", 0))

		adapter.Narrate(screenplay.Event{
			Kind:    screenplay.KindBeat,
			Phase:   screenplay.PhaseBegin,
			Message: "Wanda sends a GET request to /birds",
		})

		assert.Equal(t, "Wanda sends a GET request to /birds\n", buffer.String())
	})

	t.Run("stays silent on a successful end phase", func(t *testing.T) {
		t.Parallel()

		var buffer bytes.Buffer
		adapter := logadapter.New(log.New(&buffer, "", 0))

		adapter.Narrate(screenplay.Event{Kind: screenplay.KindBeat, Phase: screenplay.PhaseEnd})

		assert.Empty(t, buffer.String())
	})

	t.Run("indents according to the depth", func(t *testing.T) {
		t.Parallel()

		var buffer bytes.Buffer
		adapter := logadapter.New(log.New(&buffer, "", 0))

		adapter.Narrate(screenplay.Event{
			Kind:    screenplay.KindBeat,
			Phase:   screenplay.PhaseBegin,
			Depth:   2,
			Message: "sets the Accept header",
		})

		assert.Equal(t, "        sets the Accept header\n", buffer.String())
	})

	t.Run("renders an answer", func(t *testing.T) {
		t.Parallel()

		var buffer bytes.Buffer
		adapter := logadapter.New(log.New(&buffer, "", 0))

		adapter.Narrate(screenplay.Event{
			Kind:   screenplay.KindBeat,
			Phase:  screenplay.PhaseEnd,
			Answer: 200,
		})

		assert.Equal(t, "=> 200\n", buffer.String())
	})

	t.Run("renders a failure", func(t *testing.T) {
		t.Parallel()

		var buffer bytes.Buffer
		adapter := logadapter.New(log.New(&buffer, "", 0))

		adapter.Narrate(screenplay.Event{
			Kind:    screenplay.KindBeat,
			Phase:   screenplay.PhaseEnd,
			Message: "Wanda sends a GET request",
			Err:     errors.New("connection refused"),
		})

		assert.Equal(t, "✗ Wanda sends a GET request: connection refused\n", buffer.String())
	})

	t.Run("honors a custom indentation", func(t *testing.T) {
		t.Parallel()

		var buffer bytes.Buffer
		adapter := logadapter.New(log.New(&buffer, "", 0)).WithIndent(">>")

		adapter.Narrate(screenplay.Event{
			Kind:    screenplay.KindBeat,
			Phase:   screenplay.PhaseBegin,
			Depth:   1,
			Message: "step",
		})

		assert.Equal(t, ">>step\n", buffer.String())
	})
}
