// Package logadapter renders go-screenplay narration through the standard
// library "log" package.
package logadapter

import (
	"log"
	"strings"

	"github.com/grandper/go-screenplay/screenplay"
)

// Adapter writes narration to a *log.Logger.
type Adapter struct {
	logger *log.Logger
	indent string
}

// New returns an Adapter writing to the given logger. If logger is nil,
// log.Default() is used. Each nesting level is indented by four spaces.
func New(logger *log.Logger) *Adapter {
	if logger == nil {
		logger = log.Default()
	}

	return &Adapter{logger: logger, indent: "    "}
}

// WithIndent overrides the per-level indentation string.
func (a *Adapter) WithIndent(indent string) *Adapter {
	a.indent = indent
	return a
}

// Narrate implements screenplay.Adapter.
func (a *Adapter) Narrate(event screenplay.Event) {
	// Announce steps on the way in; only report the way out when it failed or
	// produced an answer, to avoid doubling every line.
	if event.Phase == screenplay.PhaseEnd && event.Err == nil && event.Answer == nil {
		return
	}

	pad := strings.Repeat(a.indent, event.Depth)

	switch {
	case event.Err != nil:
		a.logger.Printf("%s✗ %s: %v", pad, event.Message, event.Err)
	case event.Answer != nil:
		a.logger.Printf("%s=> %v", pad, event.Answer)
	default:
		a.logger.Printf("%s%s", pad, event.Message)
	}
}

// Adapter implements the screenplay.Adapter interface.
var _ screenplay.Adapter = (*Adapter)(nil)
