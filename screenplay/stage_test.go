package screenplay_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestStage(t *testing.T) {
	t.Parallel()

	t.Run("can be set with a cast", func(t *testing.T) {
		t.Parallel()

		stage := screenplay.SetTheStage(screenplay.CastOfStandardActors())
		assert.NotNil(t, stage)
	})

	t.Run("the actor called", func(t *testing.T) {
		t.Run("creates a new actor when called for the first time", func(t *testing.T) {
			t.Parallel()

			stage := screenplay.SetTheStage(screenplay.CastOfStandardActors())
			adam := stage.ActorNamed("Adam")

			assert.Equal(t, "Adam", adam.Name())
		})

		t.Run("returns the same actor when called with the same name", func(t *testing.T) {
			t.Parallel()

			stage := screenplay.SetTheStage(screenplay.CastOfStandardActors())
			adam := stage.ActorNamed("Adam")
			sameAdam := stage.ActorNamed("Adam")

			assert.Same(t, adam, sameAdam)
		})

		t.Run("matches actor names case-insensitively", func(t *testing.T) {
			t.Parallel()

			stage := screenplay.SetTheStage(screenplay.CastOfStandardActors())
			adam := stage.ActorNamed("Adam")
			sameAdam := stage.ActorNamed("adam")

			assert.Same(t, adam, sameAdam)
		})

		t.Run("prepares the actor through the cast", func(t *testing.T) {
			t.Parallel()

			stage := screenplay.SetTheStage(screenplay.CastWhereEveryoneCan(newPerformTestingAbility))
			adam := stage.ActorNamed("Adam")

			assert.True(t, adam.HasAbilityTo(performTestingAbility{}))
		})

		t.Run("places the actor in the spotlight", func(t *testing.T) {
			t.Parallel()

			stage := screenplay.SetTheStage(screenplay.CastOfStandardActors())
			adam := stage.ActorNamed("Adam")

			spotlight, err := stage.TheActorInTheSpotlight()
			require.NoError(t, err)
			assert.Same(t, adam, spotlight)
		})

		t.Run("moves the spotlight to the last actor called", func(t *testing.T) {
			t.Parallel()

			stage := screenplay.SetTheStage(screenplay.CastOfStandardActors())
			stage.ActorNamed("Adam")
			bob := stage.ActorNamed("Bob")

			spotlight, err := stage.TheActorInTheSpotlight()
			require.NoError(t, err)
			assert.Same(t, bob, spotlight)
		})
	})

	t.Run("the actor in the spotlight", func(t *testing.T) {
		t.Run("returns an error when no actor is on stage", func(t *testing.T) {
			t.Parallel()

			stage := screenplay.SetTheStage(screenplay.CastOfStandardActors())

			_, err := stage.TheActorInTheSpotlight()
			assert.ErrorIs(t, err, screenplay.ErrNoActorInTheSpotlight)
		})
	})

	t.Run("an actor is on stage", func(t *testing.T) {
		t.Run("returns false when no actor has been called", func(t *testing.T) {
			t.Parallel()

			stage := screenplay.SetTheStage(screenplay.CastOfStandardActors())
			assert.False(t, stage.AnActorIsOnStage())
		})

		t.Run("returns true after an actor has been called", func(t *testing.T) {
			t.Parallel()

			stage := screenplay.SetTheStage(screenplay.CastOfStandardActors())
			stage.ActorNamed("Adam")
			assert.True(t, stage.AnActorIsOnStage())
		})
	})

	t.Run("draw the curtain", func(t *testing.T) {
		t.Run("makes all actors exit the stage", func(t *testing.T) {
			t.Parallel()

			stage := screenplay.SetTheStage(screenplay.CastWhereEveryoneCan(newPerformTestingAbility))
			adam := stage.ActorNamed("Adam")
			stage.ActorNamed("Bob")

			require.NoError(t, stage.DrawTheCurtain())

			assert.Equal(t, 0, adam.NumAbilities())
			assert.False(t, stage.AnActorIsOnStage())
		})

		t.Run("clears the spotlight", func(t *testing.T) {
			t.Parallel()

			stage := screenplay.SetTheStage(screenplay.CastOfStandardActors())
			stage.ActorNamed("Adam")

			require.NoError(t, stage.DrawTheCurtain())

			_, err := stage.TheActorInTheSpotlight()
			assert.ErrorIs(t, err, screenplay.ErrNoActorInTheSpotlight)
		})

		t.Run("allows new actors to enter after the curtain was drawn", func(t *testing.T) {
			t.Parallel()

			stage := screenplay.SetTheStage(screenplay.CastOfStandardActors())
			stage.ActorNamed("Adam")
			require.NoError(t, stage.DrawTheCurtain())

			bob := stage.ActorNamed("Bob")
			spotlight, err := stage.TheActorInTheSpotlight()
			require.NoError(t, err)
			assert.Same(t, bob, spotlight)
		})
	})

	t.Run("stage is safe for concurrent use", func(t *testing.T) {
		// Run with -race.
		stage := screenplay.SetTheStage(screenplay.CastOfStandardActors())

		var group sync.WaitGroup
		for range 8 {
			group.Add(1)
			go func() {
				defer group.Done()
				stage.ActorNamed("Adam")
				stage.ActorNamed("Wanda")
				_, _ = stage.TheActorInTheSpotlight()
				_ = stage.AnActorIsOnStage()
			}()
		}
		group.Wait()

		adam := stage.ActorNamed("Adam")
		again := stage.ActorNamed("adam")
		assert.Same(t, adam, again, "one actor per name, whatever the goroutine that created it")
		require.NoError(t, stage.DrawTheCurtain())
		assert.False(t, stage.AnActorIsOnStage())
	})
}
