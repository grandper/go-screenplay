package screenplay

// Cast provisions actors with abilities when they enter the stage.
type Cast interface {
	Prepare(actor *Actor)
}

// CastOfStandardActors creates a cast with no predefined abilities.
func CastOfStandardActors() Cast {
	return &standardCast{}
}

// CastWhereEveryoneCan creates a cast where every actor
// is given the listed abilities.
func CastWhereEveryoneCan(abilities ...Ability) Cast {
	return &standardCast{abilities: abilities}
}

type standardCast struct {
	abilities []Ability
}

func (c *standardCast) Prepare(actor *Actor) {
	actor.WhoCan(c.abilities...)
}

// CastFunc adapts a function to the Cast interface.
// This is useful for custom actor provisioning:
//
//	SetTheStage(CastFunc(func(actor *Actor) {
//	    actor.WhoCan(browseTheWeb)
//	}))
type CastFunc func(actor *Actor)

// Prepare calls the underlying function.
func (f CastFunc) Prepare(actor *Actor) {
	f(actor)
}
