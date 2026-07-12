package timing

import "time"

// Window captures the two durations that bound a retry loop: the total time
// during which the actor keeps trying (for how long we repeat) and the interval
// the actor waits between two tries (how much time between every trial).
type Window struct {
	Total    time.Duration // for how long we keep repeating
	Interval time.Duration // how much time between every trial
}

// NewWindow returns a Window with the given total time and interval
// between tries.
func NewWindow(total, interval time.Duration) Window {
	return Window{
		Total:    total,
		Interval: interval,
	}
}

// Valid reports whether the interval between tries is not larger than the total
// time during which the actor keeps trying.
func (w Window) Valid() bool {
	return w.Interval <= w.Total
}
