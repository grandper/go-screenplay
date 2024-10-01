package screenplay

// Ability represents an ability of an actor.
// Abilities can be forgotten.
type Ability interface {
	// Forget clean up the ability.
	// The ability cannot be used after Forget() has been called.
	// This method is used, e.g., to close connections to databases,
	// deleting data, closing client cleanly.
	Forget() error
}
