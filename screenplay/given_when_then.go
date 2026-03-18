package screenplay

// Given provides a convenient way to introduce "Given" statement in tests.
func Given(actor *Actor) *Actor {
	return actor
}

// When provides a convenient way to introduce "When" statement in tests.
func When(actor *Actor) *Actor {
	return actor
}

// Then provides a convenient way to introduce "Then" statement in tests.
func Then(actor *Actor) *Actor {
	return actor
}

// And can be used to avoid repeating given, when, or then.
func And(actor *Actor) *Actor {
	return actor
}
