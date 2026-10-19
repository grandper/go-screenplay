// Package fileadapter writes go-screenplay narration to a file, as the same
// indented tree the log adapter renders.
package fileadapter

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/grandper/go-screenplay/v2/narrator/logadapter"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// ErrConflictingOptions is returned by New when it receives both
// WithOverwriting and WithAppending.
var ErrConflictingOptions = errors.New("conflicting options")

// settings collects the options passed to New.
type settings struct {
	overwriting bool
	appending   bool
	timestamps  bool
}

// Option configures the Adapter created by New.
type Option func(*settings)

// WithOverwriting empties the file when it is opened. It is the default.
func WithOverwriting() Option {
	return func(s *settings) { s.overwriting = true }
}

// WithAppending adds the narration to the end of the file, keeping what it
// already holds.
func WithAppending() Option {
	return func(s *settings) { s.appending = true }
}

// WithTimestamps prefixes each line with the date and time.
func WithTimestamps() Option {
	return func(s *settings) { s.timestamps = true }
}

// Adapter writes narration to a file.
type Adapter struct {
	file     *os.File
	writer   *errorRecorder
	renderer *logadapter.Adapter
}

// New opens the file at path and returns an Adapter writing to it. The file is
// created if it does not exist, and overwritten unless WithAppending is given.
func New(path string, options ...Option) (*Adapter, error) {
	var s settings
	for _, option := range options {
		option(&s)
	}

	if s.overwriting && s.appending {
		return nil, fmt.Errorf("%w: overwriting and appending to %s", ErrConflictingOptions, path)
	}

	flags := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	if s.appending {
		flags = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	}

	file, err := os.OpenFile(path, flags, 0o644) //nolint:gosec // a narration file is meant to be read by anyone
	if err != nil {
		return nil, err
	}

	logFlags := 0
	if s.timestamps {
		logFlags = log.LstdFlags
	}

	writer := &errorRecorder{writer: file}

	return &Adapter{
		file:     file,
		writer:   writer,
		renderer: logadapter.New(log.New(writer, "", logFlags)),
	}, nil
}

// Narrate implements screenplay.Adapter.
func (a *Adapter) Narrate(event screenplay.Event) {
	a.renderer.Narrate(event)
}

// Close closes the file. It returns the first error met while writing the
// narration, or the error met while closing the file.
func (a *Adapter) Close() error {
	return errors.Join(a.writer.err(), a.file.Close())
}

// errorRecorder keeps the first error met while writing, which log.Logger
// would otherwise discard.
type errorRecorder struct {
	writer   io.Writer
	mutex    sync.Mutex
	firstErr error
}

// Write implements io.Writer.
func (r *errorRecorder) Write(p []byte) (int, error) {
	n, err := r.writer.Write(p)
	if err != nil {
		r.mutex.Lock()
		if r.firstErr == nil {
			r.firstErr = err
		}
		r.mutex.Unlock()
	}

	return n, err
}

// err returns the first error met while writing.
func (r *errorRecorder) err() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.firstErr
}

// Adapter implements the screenplay.Adapter interface.
var _ screenplay.Adapter = (*Adapter)(nil)
