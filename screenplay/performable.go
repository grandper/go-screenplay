package screenplay

// Performable is either a task or an action that can be performed by an actor.
type Performable interface {
	// PerformAs performs the task or the action as the provided actor.
	PerformAs(actor *Actor) error
	// String describes the performable.
	String() string
}
