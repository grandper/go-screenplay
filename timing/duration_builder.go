// Package timing provides small, reusable building blocks for describing
// durations and retry windows, shared across the go-screenplay actions,
// questions, and resolutions.
package timing

import (
	"fmt"
	"time"
)

// DurationBuilder builds a duration by combining an amount and a unit and
// writes the resulting time.Duration into the value pointed to by duration.
//
// It is generic over the parent builder type T. The amount is supplied fluently
// through the For/During methods and their wording aliases. The caller then picks
// a unit (Milliseconds, Seconds, ...), which writes amount*unit through the
// duration pointer and returns the parent so the fluent chain can continue.
// Pointing several builders at different fields is how a parent can describe more
// than one duration, for example how long to keep trying and how long to wait
// between tries.
//
// The unit label tracks the amount: String reports a singular unit for a single
// time unit (for example "1 second") and a plural one otherwise ("20 seconds").
type DurationBuilder[T any] struct {
	parent   *T
	amount   int
	unit     string
	duration *time.Duration
}

// NewDurationBuilder creates a DurationBuilder. The computed duration is written
// into the value pointed to by duration and parent is returned by the unit
// methods. Set the amount with For/During or one of their aliases before picking a
// unit, for example NewDurationBuilder(parent, &duration).For(30).Seconds().
func NewDurationBuilder[T any](parent *T, duration *time.Duration) *DurationBuilder[T] {
	return &DurationBuilder[T]{
		parent:   parent,
		amount:   0,
		unit:     "",
		duration: duration,
	}
}

// String describes the duration.
func (b *DurationBuilder[T]) String() string {
	return fmt.Sprintf("%d %s", b.amount, b.unit)
}

// For sets the number of time units.
func (b *DurationBuilder[T]) For(amount int) *DurationBuilder[T] {
	b.amount = amount

	return b
}

// During sets the number of time units. It is an alias for For.
func (b *DurationBuilder[T]) During(amount int) *DurationBuilder[T] {
	return b.For(amount)
}

// TryingFor sets the number of time units. It is an alias for For.
func (b *DurationBuilder[T]) TryingFor(amount int) *DurationBuilder[T] {
	return b.For(amount)
}

// TryingForNoLongerThan sets the number of time units. It is an alias for For.
func (b *DurationBuilder[T]) TryingForNoLongerThan(amount int) *DurationBuilder[T] {
	return b.For(amount)
}

// WaitingFor sets the number of time units. It is an alias for For.
func (b *DurationBuilder[T]) WaitingFor(amount int) *DurationBuilder[T] {
	return b.For(amount)
}

// Every sets the number of time units. It is an alias for For, phrased for
// interval/frequency APIs (for example a polling period).
func (b *DurationBuilder[T]) Every(amount int) *DurationBuilder[T] {
	return b.For(amount)
}

// PollingEvery sets the number of time units. It is an alias for For.
func (b *DurationBuilder[T]) PollingEvery(amount int) *DurationBuilder[T] {
	return b.For(amount)
}

// TryingEvery sets the number of time units. It is an alias for For.
func (b *DurationBuilder[T]) TryingEvery(amount int) *DurationBuilder[T] {
	return b.For(amount)
}

// Milliseconds sets the duration in milliseconds and returns the parent.
func (b *DurationBuilder[T]) Milliseconds() *T {
	return b.resolve("millisecond", "milliseconds", time.Millisecond)
}

// Millisecond sets the duration in milliseconds. It is an alias for Milliseconds.
func (b *DurationBuilder[T]) Millisecond() *T { return b.Milliseconds() }

// Seconds sets the duration in seconds and returns the parent.
func (b *DurationBuilder[T]) Seconds() *T {
	return b.resolve("second", "seconds", time.Second)
}

// Second sets the duration in seconds. It is an alias for Seconds.
func (b *DurationBuilder[T]) Second() *T { return b.Seconds() }

// Minutes sets the duration in minutes and returns the parent.
func (b *DurationBuilder[T]) Minutes() *T {
	return b.resolve("minute", "minutes", time.Minute)
}

// Minute sets the duration in minutes. It is an alias for Minutes.
func (b *DurationBuilder[T]) Minute() *T { return b.Minutes() }

// resolve records the unit label, picking the singular or plural form to match
// the amount, writes amount*unit through the duration pointer, and returns the
// parent.
func (b *DurationBuilder[T]) resolve(singular, plural string, magnitude time.Duration) *T {
	b.unit = singularOrPlural(b.amount, singular, plural)

	if b.duration != nil {
		*b.duration = time.Duration(b.amount) * magnitude
	}

	return b.parent
}

// singularOrPlural returns the plural form when count denotes more than one unit
// and the singular form otherwise.
func singularOrPlural(count int, singular, plural string) string {
	if count > 1 {
		return plural
	}

	return singular
}
