package slogadapter_test

import (
	"bytes"
	"context"
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

	t.Run("maps every screenplay level onto its slog level", func(t *testing.T) {
		t.Parallel()

		levels := map[screenplay.Level]string{
			screenplay.LevelDebug: "DEBUG",
			screenplay.LevelInfo:  "INFO",
			screenplay.LevelWarn:  "WARN",
			screenplay.LevelError: "ERROR",
			screenplay.Level(42):  "INFO",
		}

		for level, want := range levels {
			var buffer bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&buffer, &slog.HandlerOptions{Level: slog.LevelDebug}))
			adapter := slogadapter.New(logger)

			adapter.Narrate(screenplay.Event{
				Kind:    screenplay.KindAside,
				Level:   level,
				Phase:   screenplay.PhaseBegin,
				Message: "an aside",
			})

			record := decode(t, buffer.Bytes())
			assert.Equal(t, want, record["level"])
		}
	})

	t.Run("logs through the configured context", func(t *testing.T) {
		t.Parallel()

		type key struct{}
		var buffer bytes.Buffer
		handler := &contextHandler{inner: slog.NewJSONHandler(&buffer, nil), key: key{}}
		adapter := slogadapter.New(slog.New(handler)).
			WithContext(context.WithValue(context.Background(), key{}, "trace-42"))

		adapter.Narrate(screenplay.Event{
			Kind:    screenplay.KindBeat,
			Level:   screenplay.LevelInfo,
			Phase:   screenplay.PhaseBegin,
			Message: "Wanda sends a GET request",
		})

		record := decode(t, buffer.Bytes())
		assert.Equal(t, "trace-42", record["trace"])
	})
}

// contextHandler is a slog.Handler that copies a value from the logging context
// onto every record, so a test can prove the adapter forwards its context.
type contextHandler struct {
	inner slog.Handler
	key   any
}

// Enabled implements slog.Handler.
func (h *contextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

// Handle copies the context value onto the record before delegating.
func (h *contextHandler) Handle(ctx context.Context, record slog.Record) error {
	if value := ctx.Value(h.key); value != nil {
		record.AddAttrs(slog.Any("trace", value))
	}

	return h.inner.Handle(ctx, record)
}

// WithAttrs implements slog.Handler.
func (h *contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &contextHandler{inner: h.inner.WithAttrs(attrs), key: h.key}
}

// WithGroup implements slog.Handler.
func (h *contextHandler) WithGroup(name string) slog.Handler {
	return &contextHandler{inner: h.inner.WithGroup(name), key: h.key}
}

// contextHandler implements the slog.Handler interface.
var _ slog.Handler = (*contextHandler)(nil)

func decode(t *testing.T, data []byte) map[string]any {
	t.Helper()

	var record map[string]any
	require.NoError(t, json.Unmarshal(data, &record))

	return record
}
