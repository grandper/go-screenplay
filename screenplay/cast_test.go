package screenplay_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestCast(t *testing.T) {
	t.Run("of standard actors", func(t *testing.T) {
		t.Run("prepares an actor with no abilities", func(t *testing.T) {
			t.Parallel()

			cast := screenplay.CastOfStandardActors()
			adam := screenplay.ActorNamed("Adam")
			cast.Prepare(adam)

			assert.Equal(t, 0, adam.NumAbilities())
		})
	})

	t.Run("where everyone can use abilities", func(t *testing.T) {
		t.Run("prepares an actor with the listed abilities", func(t *testing.T) {
			t.Parallel()

			cast := screenplay.CastWhereEveryoneCan(newPerformTestingAbility).And(newCheckErrorsAbility)
			adam := screenplay.ActorNamed("Adam")
			cast.Prepare(adam)

			assert.Equal(t, 2, adam.NumAbilities())
			assert.True(t, adam.HasAbilityTo(performTestingAbility{}))
			assert.True(t, adam.HasAbilityTo(checkErrorsAbility{}))
		})

		t.Run("prepares multiple actors with the same abilities", func(t *testing.T) {
			t.Parallel()

			cast := screenplay.CastWhereEveryoneCan(newPerformTestingAbility)
			adam := screenplay.ActorNamed("Adam")
			bob := screenplay.ActorNamed("Bob")
			cast.Prepare(adam)
			cast.Prepare(bob)

			assert.True(t, adam.HasAbilityTo(performTestingAbility{}))
			assert.True(t, bob.HasAbilityTo(performTestingAbility{}))
		})

		t.Run("gives each actor its own instance of every ability", func(t *testing.T) {
			t.Parallel()

			cast := screenplay.CastWhereEveryoneCan(newCountingAbility)
			adam := screenplay.ActorNamed("Adam")
			bob := screenplay.ActorNamed("Bob")
			cast.Prepare(adam)
			cast.Prepare(bob)

			adamCounts, err := screenplay.UseAbilityTo[*countingAbility]().Of(adam)
			require.NoError(t, err)
			bobCounts, err := screenplay.UseAbilityTo[*countingAbility]().Of(bob)
			require.NoError(t, err)
			assert.NotSame(t, adamCounts, bobCounts)

			// What Adam records stays Adam's, and Adam's exit forgets Adam's ability only.
			adamCounts.count++
			assert.Equal(t, 0, bobCounts.count)
			require.NoError(t, adam.Exit())
			assert.True(t, adamCounts.forgotten)
			assert.False(t, bobCounts.forgotten)
		})
	})

	t.Run("from a function", func(t *testing.T) {
		t.Run("prepares an actor using custom logic", func(t *testing.T) {
			t.Parallel()

			cast := screenplay.CastFunc(func(actor *screenplay.Actor) {
				actor.WhoCan(performTestingAbility{})
				actor.Remember("role", "tester")
			})

			adam := screenplay.ActorNamed("Adam")
			cast.Prepare(adam)

			assert.Equal(t, 1, adam.NumAbilities())
			assert.Equal(t, "tester", adam.Recall[string]("role").OrElse(""))
		})
	})
}

func newPerformTestingAbility() performTestingAbility { return performTestingAbility{} }

func newCheckErrorsAbility() checkErrorsAbility { return checkErrorsAbility{} }

// countingAbility records what its actor did, to show that each actor owns its
// instance.
type countingAbility struct {
	count     int
	forgotten bool
}

func newCountingAbility() *countingAbility { return &countingAbility{} }

func (c *countingAbility) Forget() error {
	c.forgotten = true

	return nil
}
