// Package slogadapter renders go-screenplay narration through log/slog,
// preserving structure (actor, kind, depth, error) as attributes.
package slogadapter

import (
	"context"
	"log/slog"

	"github.com/grandper/go-screenplay/screenplay"
)

// Adapter writes narration to an *slog.Logger.
type Adapter struct {
	logger *slog.Logger
	ctx    context.Context
}

// New returns an Adapter writing to the given logger. If logger is nil,
// slog.Default() is used.
func New(logger *slog.Logger) *Adapter {
	if logger == nil {
		logger = slog.Default()
	}

	return &Adapter{logger: logger, ctx: context.Background()}
}

// WithContext sets the context passed to the logger (for trace propagation).
func (a *Adapter) WithContext(ctx context.Context) *Adapter {
	a.ctx = ctx
	return a
}

// Narrate implements screenplay.Adapter.
func (a *Adapter) Narrate(event screenplay.Event) {
	// One structured record per step (on begin) plus outcomes worth recording.
	if event.Phase == screenplay.PhaseEnd && event.Err == nil && event.Answer == nil {
		return
	}

	attrs := []slog.Attr{
		slog.String("kind", event.Kind.String()),
		slog.Int("depth", event.Depth),
	}

	if event.Actor != "" {
		attrs = append(attrs, slog.String("actor", event.Actor))
	}

	if event.Answer != nil {
		attrs = append(attrs, slog.Any("answer", event.Answer))
	}

	if event.Err != nil {
		attrs = append(attrs, slog.Any("error", event.Err))
	}

	a.logger.LogAttrs(a.ctx, mapLevel(event.Level), event.Message, attrs...)
}

// mapLevel maps a screenplay.Level onto an slog.Level.
func mapLevel(level screenplay.Level) slog.Level {
	switch level {
	case screenplay.LevelDebug:
		return slog.LevelDebug
	case screenplay.LevelInfo:
		return slog.LevelInfo
	case screenplay.LevelWarn:
		return slog.LevelWarn
	case screenplay.LevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Adapter implements the screenplay.Adapter interface.
var _ screenplay.Adapter = (*Adapter)(nil)
