package screenplay_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/screenplay"
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

			performTesting := performTestingAbility{}
			checkErrors := checkErrorsAbility{}

			cast := screenplay.CastWhereEveryoneCan(performTesting, checkErrors)
			adam := screenplay.ActorNamed("Adam")
			cast.Prepare(adam)

			assert.Equal(t, 2, adam.NumAbilities())
			assert.True(t, adam.HasAbilityTo(performTesting))
			assert.True(t, adam.HasAbilityTo(checkErrors))
		})

		t.Run("prepares multiple actors with the same abilities", func(t *testing.T) {
			t.Parallel()

			performTesting := performTestingAbility{}

			cast := screenplay.CastWhereEveryoneCan(performTesting)
			adam := screenplay.ActorNamed("Adam")
			bob := screenplay.ActorNamed("Bob")
			cast.Prepare(adam)
			cast.Prepare(bob)

			assert.True(t, adam.HasAbilityTo(performTesting))
			assert.True(t, bob.HasAbilityTo(performTesting))
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
			assert.Equal(t, "tester", adam.Recall("role"))
		})
	})
}
