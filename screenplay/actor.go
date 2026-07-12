package screenplay

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"
)

var (
	ErrActorMissingAbility = errors.New("actor does not have the required ability")
	ErrWrongAbilityType    = errors.New("ability has the wrong type")
)

// Actor represents the end user.
//
// An Actor is safe for concurrent use: its memory, abilities and cleanup tasks
// are guarded by a mutex so that actions running in parallel (for example with
// action.Concurrently or action.Asynchronously) can call Remember, Recall,
// Forget or use abilities without racing.
type Actor struct {
	*actorCore
	muted bool // a view's own flag; the actor itself is never muted
}

// actorCore holds the state shared by every view of an actor. A view (see
// Muted) shares this core by pointer, so all views see the same memory,
// abilities, cleanup tasks and production while each may decide on its own
// whether to narrate. This lets a decorator mute a scoped run without ever
// touching state other goroutines read.
type actorCore struct {
	name                    string
	ctx                     context.Context
	production              *Production
	mu                      sync.RWMutex
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

// WithContext sets the context the actor uses for cancellation and deadlines.
func (a *Actor) WithContext(ctx context.Context) *Actor {
	a.ctx = ctx
	return a
}

// narrator returns the microphone the actor narrates through, never nil: its
// production's narrator, or a silent one when it has no production, so it is
// always safe to narrate.
func (a *Actor) narrator() *Narrator {
	if a.production != nil {
		if narrator := a.production.Narrator(); narrator != nil {
			return narrator
		}
	}

	return &Narrator{}
}

// Timeout returns the timeout the actor's production configures for actions that
// wait on something, falling back to DefaultTimeout when the actor has no
// production.
func (a *Actor) Timeout() time.Duration {
	if a.production != nil {
		return a.production.Timeout()
	}

	return DefaultTimeout
}

// Polling returns the polling interval the actor's production configures for
// actions that wait on something, falling back to DefaultPolling when the actor
// has no production.
func (a *Actor) Polling() time.Duration {
	if a.production != nil {
		return a.production.Polling()
	}

	return DefaultPolling
}

// narrates reports whether the actor should narrate its next step: it has an
// active narrator and is either not muted or forced to narrate anyway by a
// production configured with WithForceAllNarration.
func (a *Actor) narrates() bool {
	if a.muted && !(a.production != nil && a.production.forceAllNarration) {
		return false
	}

	return a.narrator().active()
}

// Muted returns a view of the actor whose narration is muted: the steps it
// performs emit nothing, unless its production forces all narration. The view
// shares the actor's memory, abilities and cleanup tasks, so muting never
// touches the actor itself nor any concurrent branch. It is how action.Silently
// and question.Silently suppress narration.
func (a *Actor) Muted() *Actor {
	return &Actor{
		actorCore: a.actorCore,
		muted:     true,
	}
}

// Remember stores a value in the actor's memory.
// When the value is a Question, the actor answers it and stores the answer instead.
// If the question fails to be answered, nil is stored for the key.
func (a *Actor) Remember(key string, value any) {
	if question, isAQuestion := value.(Question); isAQuestion {
		// AnsweredBy is called without holding the lock: it may itself call
		// Recall or Remember on the same actor, which would deadlock otherwise.
		answer, err := question.AnsweredBy(a)
		if err != nil {
			a.store(key, nil)
			return
		}

		a.store(key, answer)
		return
	}

	a.store(key, value)
}

// store writes a value into the actor's memory under the actor's lock.
func (a *Actor) store(key string, value any) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.memory[key] = value
}

// Recall retrieves a value from the actor's memory.
func (a *Actor) Recall(key string) any {
	a.mu.RLock()
	defer a.mu.RUnlock()

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
	a.mu.Lock()
	defer a.mu.Unlock()

	delete(a.memory, key)
}

// WhoCan defines abilities that the actor can use.
func (a *Actor) WhoCan(abilities ...Ability) *Actor {
	a.mu.Lock()
	defer a.mu.Unlock()

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

	_, ok := a.ability(str)

	return ok
}

// ability returns the ability stored under the given name under the actor's lock.
func (a *Actor) ability(name string) (Ability, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	ability, ok := a.abilities[name]

	return ability, ok
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
	if v, found := actor.ability(str); found {
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
	a.mu.RLock()
	defer a.mu.RUnlock()

	return len(a.abilities)
}

// HasOrderedCleanupTasks assigns one or more tasks to an actor.
// Those tasks will be performed when the actor exit the stage.
// Ordered cleanup tasks will be performed in order.
// When a task fails, the subsequent tasks won't be done.
func (a *Actor) HasOrderedCleanupTasks(tasks ...Performable) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.orderedCleanupTasks = append(a.orderedCleanupTasks, tasks...)
}

// WithOrderedCleanupTasks is an alias for HasOrderedCleanupTasks.
func (a *Actor) WithOrderedCleanupTasks(tasks ...Performable) {
	a.HasOrderedCleanupTasks(tasks...)
}

// HasIndependentCleanupTasks assigns one or more tasks to an actor.
// Those tasks will be performed when the actor exit the stage.
// Independent cleanup tasks will all be performed even if some of them failed.
func (a *Actor) HasIndependentCleanupTasks(tasks ...Performable) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.independentCleanupTasks = append(a.independentCleanupTasks, tasks...)
}

// WithIndependentCleanupTasks is an alias for HasIndependentCleanupTasks.
func (a *Actor) WithIndependentCleanupTasks(tasks ...Performable) {
	a.HasIndependentCleanupTasks(tasks...)
}

// AttemptsTo makes the actor perform a list of actions and return
// an error when the first action failed.
// Nil actions are ignored so that optional or conditionally built actions can be
// passed without guarding against nil.
// Aliases:
//
//	WasAbleTo, Does, Did, Will, TriesTo, TriedTo, Should, Shall
func (a *Actor) AttemptsTo(actions ...Performable) error {
	for _, action := range actions {
		if action == nil {
			continue
		}

		// When the actor is not narrating (no adapter, or muted without a
		// production forcing narration) the action is performed directly, without
		// paying for its description (String may be costly or, for some actions,
		// have side effects).
		if !a.narrates() {
			if err := action.PerformAs(a); err != nil {
				return err
			}

			continue
		}

		message := fmt.Sprintf("%s %s", a.name, action.String())

		err := a.narrator().StatesTheFact(a.name, message, func() error {
			return action.PerformAs(a)
		})
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
	if !a.narrates() {
		return question.AnsweredBy(a)
	}

	narrator := a.narrator()
	message := fmt.Sprintf("%s asks for %s", a.name, question.String())
	narrator.WhispersTheAside(a.name, message)

	answer, err := question.AnsweredBy(a)
	narrator.RevealsTheAnswer(a.name, message, answer, err)

	return answer, err
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
	independentErr := a.cleansUpIndependentTasks()
	orderedErr := a.cleansUpOrderedTasks()
	return errors.Join(independentErr, orderedErr)
}

func (a *Actor) cleansUpIndependentTasks() error {
	// The tasks are snapshotted under the lock and performed without holding it,
	// as a task may call back into the actor (Remember, Recall, ...).
	a.mu.Lock()
	tasks := a.independentCleanupTasks
	a.independentCleanupTasks = []Performable{}
	a.mu.Unlock()

	var errs []error

	for _, task := range tasks {
		err := task.PerformAs(a)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (a *Actor) cleansUpOrderedTasks() error {
	a.mu.Lock()
	tasks := a.orderedCleanupTasks
	a.orderedCleanupTasks = []Performable{}
	a.mu.Unlock()

	for _, task := range tasks {
		err := task.PerformAs(a)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *Actor) forgetsAbilities() error {
	a.mu.RLock()
	abilities := make([]Ability, 0, len(a.abilities))
	for _, ability := range a.abilities {
		abilities = append(abilities, ability)
	}
	a.mu.RUnlock()

	for _, ability := range abilities {
		err := ability.Forget()
		if err != nil {
			return err
		}
	}

	a.mu.Lock()
	a.abilities = map[string]Ability{}
	a.mu.Unlock()

	return nil
}

// ActorNamed creates a new actor with the provided name.
func ActorNamed(name string) *Actor {
	return &Actor{
		actorCore: &actorCore{
			name:                    name,
			abilities:               map[string]Ability{},
			orderedCleanupTasks:     []Performable{},
			independentCleanupTasks: []Performable{},
			memory:                  map[string]any{},
		},
	}
}
