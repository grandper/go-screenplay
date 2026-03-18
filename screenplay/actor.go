package screenplay

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrActorMissingAbility = errors.New("actor does not have the required ability")
	ErrWrongAbilityType    = errors.New("ability has the wrong type")
)

// Actor represents the end user.
type Actor struct {
	name                    string
	ctx                     context.Context
	abilities               map[string]Ability
	orderedCleanupTasks     []Performable
	independentCleanupTasks []Performable
	memory                  map[string]any
}

// Name returns the name of the actor.
func (a *Actor) Name() string {
	return a.name
}

// Context returns the context of the actor.
func (a *Actor) Context() context.Context {
	if a.ctx == nil {
		return context.Background()
	}

	return a.ctx
}

func (a *Actor) WithContext(ctx context.Context) *Actor {
	a.ctx = ctx
	return a
}

// Remember stores a value in the actor's memory.
func (a *Actor) Remember(key string, value any) {
	a.memory[key] = value
}

// Recall retrieves a value from the actor's memory.
func (a *Actor) Recall(key string) any {
	return a.memory[key]
}

// Share starts sharing a value from the actor's memory.
// Usage: actor.Share("mykey").With(anotherActor).
func (a *Actor) Share(key string) *ShareAction {
	return &ShareAction{
		source: a,
		key:    key,
	}
}

// ShareAction holds the context for sharing a memory value between actors.
type ShareAction struct {
	source *Actor
	key    string
}

// With completes the share by copying the value into the target actor's memory.
func (s *ShareAction) With(target *Actor) {
	value := s.source.Recall(s.key)
	if value == nil {
		return
	}
	target.Remember(s.key, value)
}

// Forget removes a value from the actor's memory.
func (a *Actor) Forget(key string) {
	delete(a.memory, key)
}

// WhoCan defines abilities that the actor can use.
func (a *Actor) WhoCan(abilities ...Ability) *Actor {
	for _, ability := range abilities {
		str := abilityStr(ability)
		a.abilities[str] = ability
	}

	return a
}

// Can defines an abilities that the actor can use.
func (a *Actor) Can(abilities ...Ability) *Actor {
	return a.WhoCan(abilities...)
}

// HasAbilityTo returns whether an actor has an ability or not.
func (a *Actor) HasAbilityTo(ability Ability) bool {
	str := abilityStr(ability)
	_, ok := a.abilities[str]

	return ok
}

// UseAbilityTo is used to access the ability of an actor.
// It is typically used as follows:
// ability, err := UseAbilityTo[BrowseTheWeb]().Of(Adam) .
func UseAbilityTo[A Ability]() UseAbility[A] {
	return UseAbility[A]{}
}

// UseAbility is a way to extract an ability of an actor.
type UseAbility[A Ability] struct{}

// Of specifies which actor the ability is extracted from.
func (ae UseAbility[A]) Of(actor *Actor) (A, error) {
	var ability A

	str := abilityStr(ability)
	if v, found := actor.abilities[str]; found {
		var ok bool
		if ability, ok = v.(A); ok {
			return ability, nil
		}

		return ability, fmt.Errorf("%w: the ability '%s' learned by '%s' has the wrong type",
			ErrWrongAbilityType, str, actor.Name())
	}

	return ability, fmt.Errorf(
		"%w: actor '%s' does not have the ability '%s'",
		ErrActorMissingAbility,
		actor.Name(),
		str,
	)
}

func abilityStr(a Ability) string {
	t := reflect.TypeOf(a)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name()
}

// NumAbilities returns the number of abilities of the actor.
func (a *Actor) NumAbilities() int {
	return len(a.abilities)
}

// HasOrderedCleanupTasks assigns one or more tasks to an actor.
// Those tasks will be performed when the actor exit the stage.
// Ordered cleanup tasks will be performed in order.
// When a task fails, the subsequent tasks won't be done.
func (a *Actor) HasOrderedCleanupTasks(tasks ...Performable) {
	a.orderedCleanupTasks = append(a.orderedCleanupTasks, tasks...)
}

// WithOrderedCleanupTasks is an alias for HasOrderedCleanupTasks.
func (a *Actor) WithOrderedCleanupTasks(tasks ...Performable) {
	a.HasOrderedCleanupTasks(tasks...)
}

// HasIndependentCleanupTasks assigns one or more tasks to an actor.
// Those tasks will be performed when the actor exit the stage.
// Ordered cleanup tasks will all be performed even if some of them failed.
// When a task fails, the subsequent tasks won't be done.
func (a *Actor) HasIndependentCleanupTasks(tasks ...Performable) {
	a.independentCleanupTasks = append(a.independentCleanupTasks, tasks...)
}

// WithIndependentCleanupTasks is an alias for HasIndependentCleanupTasks.
func (a *Actor) WithIndependentCleanupTasks(tasks ...Performable) {
	a.HasIndependentCleanupTasks(tasks...)
}

// AttemptsTo makes the actor perform a list of actions and return
// an error when the first action failed.
// Aliases:
//
//	WasAbleTo, Does, Did, Will, TriesTo, TriedTo, Should, Shall
func (a *Actor) AttemptsTo(actions ...Performable) error {
	for _, action := range actions {
		err := action.PerformAs(a)
		if err != nil {
			return err
		}
	}

	return nil
}

// WasAbleTo performs a list of actions and return an error when the
// first action failed.
// This method is an alias to 'AttemptsTo'.
func (a *Actor) WasAbleTo(actions ...Performable) error {
	return a.AttemptsTo(actions...)
}

// Does performs a list of actions and return an error when the
// first action failed.
// This method is an alias to 'AttemptsTo'.
func (a *Actor) Does(actions ...Performable) error {
	return a.AttemptsTo(actions...)
}

// Did performs a list of actions and return an error when the
// first action failed.
// This method is an alias to 'AttemptsTo'.
func (a *Actor) Did(actions ...Performable) error {
	return a.AttemptsTo(actions...)
}

// Will performs a list of actions and return an error when the
// first action failed.
// This method is an alias to 'AttemptsTo'.
func (a *Actor) Will(actions ...Performable) error {
	return a.AttemptsTo(actions...)
}

// TriesTo performs a list of actions and return an error when the
// first action failed.
// This method is an alias to 'AttemptsTo'.
func (a *Actor) TriesTo(actions ...Performable) error {
	return a.AttemptsTo(actions...)
}

// TriedTo performs a list of actions and return an error when the
// first action failed.
// This method is an alias to 'AttemptsTo'.
func (a *Actor) TriedTo(actions ...Performable) error {
	return a.AttemptsTo(actions...)
}

// Tries performs a list of actions and return an error when the
// first action failed.
// This method is an alias to 'AttemptsTo'.
func (a *Actor) Tries(actions ...Performable) error {
	return a.AttemptsTo(actions...)
}

// Tried performs a list of actions and return an error when the
// first action failed.
// This method is an alias to 'AttemptsTo'.
func (a *Actor) Tried(actions ...Performable) error {
	return a.AttemptsTo(actions...)
}

// Shall performs a list of actions and return an error when the
// first action failed.
// This method is an alias to 'AttemptsTo'.
func (a *Actor) Shall(actions ...Performable) error {
	return a.AttemptsTo(actions...)
}

// Should performs a list of actions and return an error when the
// first action failed.
// This method is an alias to 'AttemptsTo'.
func (a *Actor) Should(actions ...Performable) error {
	return a.AttemptsTo(actions...)
}

// AsksFor asks the given question.
func (a *Actor) AsksFor(question Question) (any, error) {
	return question.AnsweredBy(a)
}

// Sees asks a question about what the actor sees on the screen.
func (a *Actor) Sees(question Question) (any, error) {
	return a.AsksFor(question)
}

// Exit makes the actor exit the stage.
// The actor will perform all his clean-up tasks and forget
// all of his abilities.
func (a *Actor) Exit() error {
	err := a.cleansUp()
	if err != nil {
		return err
	}

	return a.forgetsAbilities()
}

func (a *Actor) cleansUp() error {
	err := a.cleansUpIndependentTasks()
	if err != nil {
		return err
	}

	return a.cleansUpOrderedTasks()
}

func (a *Actor) cleansUpIndependentTasks() error {
	var errs []error

	for _, task := range a.independentCleanupTasks {
		err := task.PerformAs(a)
		if err != nil {
			errs = append(errs, err)
		}
	}

	a.independentCleanupTasks = []Performable{}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (a *Actor) cleansUpOrderedTasks() error {
	for _, task := range a.orderedCleanupTasks {
		err := task.PerformAs(a)
		if err != nil {
			return err
		}
	}

	a.orderedCleanupTasks = []Performable{}

	return nil
}

func (a *Actor) forgetsAbilities() error {
	for _, ability := range a.abilities {
		err := ability.Forget()
		if err != nil {
			return err
		}
	}

	a.abilities = map[string]Ability{}

	return nil
}

// ActorNamed creates a new actor with the provided name.
func ActorNamed(name string) *Actor {
	return &Actor{
		name:                    name,
		abilities:               map[string]Ability{},
		orderedCleanupTasks:     []Performable{},
		independentCleanupTasks: []Performable{},
		memory:                  map[string]any{},
	}
}
