package screenplay_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/fixture"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// usingAnAbility is an action that asks for an ability and hands what That
// returned to the test through got, the way an action reaches an ability.
type usingAnAbility[T any, A interface {
	*T
	screenplay.Ability
}] struct {
	ability A
	got     *A
}

// PerformAs asks the actor for the ability and keeps what That returned.
func (u usingAnAbility[T, A]) PerformAs(actor *screenplay.Actor) error {
	*u.got = screenplay.UseAnAbilityOf(actor).That(u.ability)

	return nil
}

// String describes the action.
func (u usingAnAbility[T, A]) String() string {
	return "use an ability"
}

// landingAbility is a pointer ability whose only operation guards a nil
// receiver, as every ability must.
type landingAbility struct {
	landed bool
}

// Land is the operation of the ability; on a nil receiver it reports that the
// actor does not have the ability.
func (l *landingAbility) Land() error {
	if l == nil {
		return screenplay.ErrActorMissingAbility
	}
	l.landed = true

	return nil
}

// Forget cleans up the ability; there is nothing to clean up.
func (l *landingAbility) Forget() error {
	return nil
}

// FakeAbility has the name of fixture.FakeAbility, in another package.
type FakeAbility struct{}

// Forget cleans up the ability; there is nothing to clean up.
func (*FakeAbility) Forget() error {
	return nil
}

func TestUseAnAbility(t *testing.T) {
	t.Parallel()

	t.Run("returns the actor's own instance", func(t *testing.T) {
		t.Parallel()

		landing := &landingAbility{}
		adam := screenplay.ActorNamed("Adam").WhoCan(landing)

		var got *landingAbility
		require.NoError(
			t,
			adam.AttemptsTo(usingAnAbility[landingAbility, *landingAbility]{ability: &landingAbility{}, got: &got}),
		)
		assert.Same(t, landing, got)
	})

	t.Run("returns nil when the actor has no such ability, and the operation reports it", func(t *testing.T) {
		t.Parallel()

		adam := screenplay.ActorNamed("Adam")

		var got *landingAbility
		require.NoError(
			t,
			adam.AttemptsTo(usingAnAbility[landingAbility, *landingAbility]{ability: &landingAbility{}, got: &got}),
		)
		assert.Nil(t, got)
		assert.ErrorIs(t, got.Land(), screenplay.ErrActorMissingAbility)
	})

	t.Run("returns nil for a nil actor without panicking", func(t *testing.T) {
		t.Parallel()

		assert.Nil(t, screenplay.UseAnAbilityOf(nil).That(&landingAbility{}))
	})

	t.Run("returns the actor's ability for a nil pointer argument without panicking", func(t *testing.T) {
		t.Parallel()

		landing := &landingAbility{}
		adam := screenplay.ActorNamed("Adam").WhoCan(landing)

		var none *landingAbility
		var got *landingAbility
		require.NoError(t, adam.AttemptsTo(usingAnAbility[landingAbility, *landingAbility]{ability: none, got: &got}))
		assert.Same(t, landing, got)
	})

	t.Run("returns nil when the ability of that name is not an A", func(t *testing.T) {
		t.Parallel()

		adam := screenplay.ActorNamed("Adam").WhoCan(performTestingAbility{})

		var got *performTestingAbility
		require.NoError(
			t,
			adam.AttemptsTo(
				usingAnAbility[performTestingAbility, *performTestingAbility]{
					ability: &performTestingAbility{},
					got:     &got,
				},
			),
		)
		assert.Nil(t, got)
	})

	t.Run("tells apart two abilities of one name from two packages", func(t *testing.T) {
		t.Parallel()

		adam := screenplay.ActorNamed("Adam").WhoCan(fixture.NewFakeAbility(), &FakeAbility{})
		assert.Equal(t, 2, adam.NumAbilities())

		var got *FakeAbility
		require.NoError(
			t,
			adam.AttemptsTo(usingAnAbility[FakeAbility, *FakeAbility]{ability: &FakeAbility{}, got: &got}),
		)
		assert.NotNil(t, got)
	})
}
