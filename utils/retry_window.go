package utils

import "time"

// RetryWindow captures the two durations that bound a retry loop: the total time
// during which the actor keeps trying (for how long we repeat) and the interval
// the actor waits between two tries (how much time between every trial).
type RetryWindow struct {
	Total    time.Duration // for how long we keep repeating
	Interval time.Duration // how much time between every trial
}

// NewRetryWindow returns a RetryWindow with the given total time and interval
// between tries.
func NewRetryWindow(total, interval time.Duration) RetryWindow {
	return RetryWindow{
		Total:    total,
		Interval: interval,
	}
}

// Valid reports whether the interval between tries is not larger than the total
// time during which the actor keeps trying.
func (w RetryWindow) Valid() bool {
	return w.Interval <= w.Total
}
