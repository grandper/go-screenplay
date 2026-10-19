package screenplay

import "reflect"

// UseAnAbilityOf starts using an ability of the actor inside an action or a
// question; it is written in PerformAs or AnsweredBy, never in a test:
//
//	err := screenplay.UseAnAbilityOf(theActor).That(UseTerminalCommands()).Run("ls")
func UseAnAbilityOf(actor *Actor) UseAnAbilityBuilder {
	return UseAnAbilityBuilder{actor: actor}
}

// UseAnAbilityBuilder holds the actor given to UseAnAbilityOf and waits for the
// ability asked for.
type UseAnAbilityBuilder struct {
	actor *Actor
}

// That returns the actor's own instance of the ability of type A. The argument
// is a pointer to an initialized ability, used only for its type: an ability is
// always a pointer, which the constraint enforces at compile time.
//
// It never fails. When the actor is nil, has no ability of type A, or the
// ability stored under that type's name is not an A, it returns nil. Every
// method of an ability guards a nil receiver and returns
// ErrActorMissingAbility, so the operation chained after That reports the error.
func (b UseAnAbilityBuilder) That[T any, A interface {
	*T
	Ability
}](ability A) A {
	if b.actor == nil {
		return nil
	}

	found, ok := b.actor.ability(abilityName(reflect.TypeOf(ability)))
	if !ok {
		return nil
	}

	theAbility, isAnA := found.(A)
	if !isAnA {
		return nil
	}

	return theAbility
}
