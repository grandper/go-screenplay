// Package keep holds the action that retries another action until it
// succeeds.
//
//	err := actor.AttemptsTo(
//		keep.TryingTo(CancelTheOrder()).For(5).Seconds().PollingEvery(500).Milliseconds(),
//	)
//
// reads "keep trying to cancel the order for 5 seconds polling every 500
// milliseconds", which is also what the actor narrates. When it never
// succeeds, it fails with "Adam kept trying to cancel the order for 5s (10
// tries) without success: ...".
package keep
