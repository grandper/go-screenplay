package slogadapter_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/narrator/slogadapter"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestSlogAdapter(t *testing.T) {
	t.Parallel()

	t.Run("emits a structured record for a beat", func(t *testing.T) {
		t.Parallel()

		var buffer bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&buffer, nil))
		adapter := slogadapter.New(logger)

		adapter.Narrate(screenplay.Event{
			Kind:    screenplay.KindBeat,
			Level:   screenplay.LevelInfo,
			Phase:   screenplay.PhaseBegin,
			Depth:   1,
			Actor:   "Wanda",
			Message: "Wanda sends a GET request",
		})

		record := decode(t, buffer.Bytes())
		assert.Equal(t, "Wanda sends a GET request", record["msg"])
		assert.Equal(t, "INFO", record["level"])
		assert.Equal(t, "beat", record["kind"])
		assert.Equal(t, "Wanda", record["actor"])
		assert.InDelta(t, 1, record["depth"], 0)
	})

	t.Run("stays silent on a successful end phase", func(t *testing.T) {
		t.Parallel()

		var buffer bytes.Buffer
		adapter := slogadapter.New(slog.New(slog.NewJSONHandler(&buffer, nil)))

		adapter.Narrate(screenplay.Event{Kind: screenplay.KindBeat, Phase: screenplay.PhaseEnd})

		assert.Empty(t, buffer.String())
	})

	t.Run("records an answer", func(t *testing.T) {
		t.Parallel()

		var buffer bytes.Buffer
		adapter := slogadapter.New(slog.New(slog.NewJSONHandler(&buffer, nil)))

		adapter.Narrate(screenplay.Event{
			Kind:   screenplay.KindBeat,
			Level:  screenplay.LevelInfo,
			Phase:  screenplay.PhaseEnd,
			Answer: 200,
		})

		record := decode(t, buffer.Bytes())
		assert.InDelta(t, 200, record["answer"], 0)
	})

	t.Run("records a failure at the error level", func(t *testing.T) {
		t.Parallel()

		var buffer bytes.Buffer
		adapter := slogadapter.New(slog.New(slog.NewJSONHandler(&buffer, nil)))

		adapter.Narrate(screenplay.Event{
			Kind:    screenplay.KindBeat,
			Level:   screenplay.LevelError,
			Phase:   screenplay.PhaseEnd,
			Message: "Wanda sends a GET request",
			Err:     errors.New("connection refused"),
		})

		record := decode(t, buffer.Bytes())
		assert.Equal(t, "ERROR", record["level"])
		assert.Equal(t, "connection refused", record["error"])
	})
}

func decode(t *testing.T, data []byte) map[string]any {
	t.Helper()

	var record map[string]any
	require.NoError(t, json.Unmarshal(data, &record))

	return record
}
