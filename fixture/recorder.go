package fixture

import (
	"sync"

	"github.com/grandper/go-screenplay/screenplay"
)

// Recorder is a narration adapter that records every event it receives.
// It is meant to be used in tests to assert on the narrative a scenario
// produces. It is safe for concurrent use.
type Recorder struct {
	mutex  sync.Mutex
	events []screenplay.Event
}

// NewRecorder creates a new recorder.
func NewRecorder() *Recorder {
	return &Recorder{}
}

// Narrate records the event.
func (r *Recorder) Narrate(event screenplay.Event) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.events = append(r.events, event)
}

// Events returns a copy of the recorded events.
func (r *Recorder) Events() []screenplay.Event {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	events := make([]screenplay.Event, len(r.events))
	copy(events, r.events)

	return events
}

// Messages returns the messages of the recorded events, in order.
func (r *Recorder) Messages() []string {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	messages := make([]string, len(r.events))
	for i, event := range r.events {
		messages[i] = event.Message
	}

	return messages
}

// Recorder implements the screenplay.Adapter interface.
var _ screenplay.Adapter = (*Recorder)(nil)
