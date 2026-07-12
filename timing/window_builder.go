package timing

// WindowBuilder gives any builder a fluent, natural-language API to
// configure the two durations of a Window: how long the actor keeps trying
// (Total) and how often it tries (Interval).
//
// It is generic over the parent builder type T. The timeout methods (For and its
// aliases) and the polling methods (Polling and its aliases) each return a
// DurationBuilder pointed at the relevant field of the window, so the caller can
// pick a unit (Milliseconds, Seconds, ...) and keep the fluent chain flowing back
// to the parent. Embed a *WindowBuilder[T] into your builder to expose this
// vocabulary without re-declaring every alias.
type WindowBuilder[T any] struct {
	parent *T
	window *Window
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
	return NewDurationBuilder(b.parent, &b.window.Total).For(amount)
}

// TryingFor sets the time during which the actor keeps on trying. It is an alias for For.
func (b *WindowBuilder[T]) TryingFor(amount int) *DurationBuilder[T] {
	return b.For(amount)
}

// TryingForNoLongerThan sets the time during which the actor keeps on trying. It is an alias for For.
func (b *WindowBuilder[T]) TryingForNoLongerThan(amount int) *DurationBuilder[T] {
	return b.For(amount)
}

// WaitingFor sets the time during which the actor keeps on trying. It is an alias for For.
func (b *WindowBuilder[T]) WaitingFor(amount int) *DurationBuilder[T] {
	return b.For(amount)
}

// Polling sets the polling frequency.
func (b *WindowBuilder[T]) Polling(amount int) *DurationBuilder[T] {
	return NewDurationBuilder(b.parent, &b.window.Interval).For(amount)
}

// PollingEvery sets the polling frequency. It is an alias for Polling.
func (b *WindowBuilder[T]) PollingEvery(amount int) *DurationBuilder[T] {
	return b.Polling(amount)
}

// TryingEvery sets the polling frequency. It is an alias for Polling.
func (b *WindowBuilder[T]) TryingEvery(amount int) *DurationBuilder[T] {
	return b.Polling(amount)
}
