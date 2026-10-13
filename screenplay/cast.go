package screenplay

// Cast provisions actors with abilities when they enter the stage.
type Cast interface {
	Prepare(actor *Actor)
}

// CastOfStandardActors creates a cast with no predefined abilities.
func CastOfStandardActors() Cast {
	return &standardCast{}
}

// CastWhereEveryoneCan creates a cast where every actor is given its own
// instance of the ability built by newAbility, so that what an actor records
// through an ability (HTTP responses, command results, a current directory)
// stays its own, and an actor's Exit forgets its abilities only. Add the other
// abilities with And:
//
//	SetTheStage(CastWhereEveryoneCan(MakeHTTPRequests).And(UseDatabase))
func CastWhereEveryoneCan[A Ability](newAbility func() A) *standardCast {
	return (&standardCast{}).And(newAbility)
}

// standardCast builds, for each actor it prepares, a fresh instance of every
// ability it was given a constructor for.
type standardCast struct {
	newAbilities []func() Ability
}

// And adds an ability every actor is given its own instance of. It reads
// "cast where everyone can make HTTP requests and use database".
func (c *standardCast) And[A Ability](newAbility func() A) *standardCast {
	c.newAbilities = append(c.newAbilities, func() Ability { return newAbility() })

	return c
}

// Prepare gives the actor a fresh instance of every ability of the cast.
func (c *standardCast) Prepare(actor *Actor) {
	for _, newAbility := range c.newAbilities {
		actor.WhoCan(newAbility())
	}
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
