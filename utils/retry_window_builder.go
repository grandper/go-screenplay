package utils

// RetryWindowBuilder gives any builder a fluent, natural-language API to
// configure the two durations of a RetryWindow: how long the actor keeps trying
// (Total) and how often it tries (Interval).
//
// It is generic over the parent builder type T. The timeout methods (For and its
// aliases) and the polling methods (Polling and its aliases) each return a
// TimeFrameBuilder pointed at the relevant field of the window, so the caller can
// pick a unit (Milliseconds, Seconds, ...) and keep the fluent chain flowing back
// to the parent. Embed a *RetryWindowBuilder[T] into your builder to expose this
// vocabulary without re-declaring every alias.
type RetryWindowBuilder[T any] struct {
	parent *T
	window *RetryWindow
}

// NewRetryWindowBuilder creates a RetryWindowBuilder that writes into the given
// window and returns parent from the unit methods so the fluent chain can
// continue. parent is typically the builder that embeds the returned value.
func NewRetryWindowBuilder[T any](parent *T, window *RetryWindow) *RetryWindowBuilder[T] {
	return &RetryWindowBuilder[T]{
		parent: parent,
		window: window,
	}
}

// For sets the time during which the actor keeps on trying.
func (b *RetryWindowBuilder[T]) For(amount int) *TimeFrameBuilder[T] {
	return NewTimeFrameBuilder(b.parent, &b.window.Total).For(amount)
}

// TryingFor sets the time during which the actor keeps on trying. It is an alias for For.
func (b *RetryWindowBuilder[T]) TryingFor(amount int) *TimeFrameBuilder[T] {
	return b.For(amount)
}

// TryingForNoLongerThan sets the time during which the actor keeps on trying. It is an alias for For.
func (b *RetryWindowBuilder[T]) TryingForNoLongerThan(amount int) *TimeFrameBuilder[T] {
	return b.For(amount)
}

// WaitingFor sets the time during which the actor keeps on trying. It is an alias for For.
func (b *RetryWindowBuilder[T]) WaitingFor(amount int) *TimeFrameBuilder[T] {
	return b.For(amount)
}

// Polling sets the polling frequency.
func (b *RetryWindowBuilder[T]) Polling(amount int) *TimeFrameBuilder[T] {
	return NewTimeFrameBuilder(b.parent, &b.window.Interval).For(amount)
}

// PollingEvery sets the polling frequency. It is an alias for Polling.
func (b *RetryWindowBuilder[T]) PollingEvery(amount int) *TimeFrameBuilder[T] {
	return b.Polling(amount)
}

// TryingEvery sets the polling frequency. It is an alias for Polling.
func (b *RetryWindowBuilder[T]) TryingEvery(amount int) *TimeFrameBuilder[T] {
	return b.Polling(amount)
}
