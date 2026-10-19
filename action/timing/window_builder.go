package timing

import "strings"

// WindowBuilder gives any builder a fluent, natural-language API to
// configure the two durations of a Window: how long the actor keeps trying
// (Total) and how often it tries (Interval).
//
// It is generic over the parent builder type T. The timeout methods (For and its
// aliases) and the polling methods (Polling and its aliases) each return a
// DurationBuilder pointed at the relevant field of the window, so the caller can
// pick a unit (Milliseconds, Seconds, ...) and keep the fluent chain flowing back
// to the parent. Embed a *WindowBuilder[T] into your builder to expose this
// vocabulary without re-declaring every alias. String describes the window as
// it was written, so the parent can narrate it.
type WindowBuilder[T any] struct {
	parent  *T
	window  *Window
	written []writtenDuration[T]
}

// writtenDuration is one duration of the window, with the words of the method
// that set it: "polling every" and "500 milliseconds".
type writtenDuration[T any] struct {
	words    string
	duration *DurationBuilder[T]
}

// NewWindowBuilder creates a WindowBuilder that writes into the given
// window and returns parent from the unit methods so the fluent chain can
// continue. parent is typically the builder that embeds the returned value.
func NewWindowBuilder[T any](parent *T, window *Window) *WindowBuilder[T] {
	return &WindowBuilder[T]{
		parent: parent,
		window: window,
	}
}

// For sets the time during which the actor keeps on trying.
func (b *WindowBuilder[T]) For(amount int) *DurationBuilder[T] {
	return b.total("for", amount)
}

// TryingFor sets the time during which the actor keeps on trying. It behaves like For and is described "trying for".
func (b *WindowBuilder[T]) TryingFor(amount int) *DurationBuilder[T] {
	return b.total("trying for", amount)
}

// TryingForNoLongerThan sets the time during which the actor keeps on trying. It behaves like For and is described "trying for no longer than".
func (b *WindowBuilder[T]) TryingForNoLongerThan(amount int) *DurationBuilder[T] {
	return b.total("trying for no longer than", amount)
}

// WaitingFor sets the time during which the actor keeps on trying. It behaves like For and is described "waiting for".
func (b *WindowBuilder[T]) WaitingFor(amount int) *DurationBuilder[T] {
	return b.total("waiting for", amount)
}

// Polling sets the polling frequency.
func (b *WindowBuilder[T]) Polling(amount int) *DurationBuilder[T] {
	return b.interval("polling", amount)
}

// PollingEvery sets the polling frequency. It behaves like Polling and is described "polling every".
func (b *WindowBuilder[T]) PollingEvery(amount int) *DurationBuilder[T] {
	return b.interval("polling every", amount)
}

// TryingEvery sets the polling frequency. It behaves like Polling and is described "trying every".
func (b *WindowBuilder[T]) TryingEvery(amount int) *DurationBuilder[T] {
	return b.interval("trying every", amount)
}

// String describes the window as it was written, in call order: "for 5
// seconds polling every 500 milliseconds". It is empty when no duration was
// set.
func (b *WindowBuilder[T]) String() string {
	descriptions := make([]string, 0, len(b.written))
	for _, written := range b.written {
		descriptions = append(descriptions, written.words+" "+written.duration.String())
	}

	return strings.Join(descriptions, " ")
}

// total sets the time during which the actor keeps on trying, and records the
// words of the method that set it.
func (b *WindowBuilder[T]) total(words string, amount int) *DurationBuilder[T] {
	return b.records(words, NewDurationBuilder(b.parent, &b.window.Total).For(amount))
}

// interval sets the polling frequency, and records the words of the method
// that set it.
func (b *WindowBuilder[T]) interval(words string, amount int) *DurationBuilder[T] {
	return b.records(words, NewDurationBuilder(b.parent, &b.window.Interval).For(amount))
}

// records keeps a duration of the window with the words of the method that set it.
func (b *WindowBuilder[T]) records(words string, duration *DurationBuilder[T]) *DurationBuilder[T] {
	b.written = append(b.written, writtenDuration[T]{words: words, duration: duration})

	return duration
}
