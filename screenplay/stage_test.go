package screenplay_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/screenplay"
)

func TestStage(t *testing.T) {
	t.Run("can be set with a cast", func(t *testing.T) {
		t.Parallel()

		stage := screenplay.SetTheStage(screenplay.CastOfStandardActors())
		assert.NotNil(t, stage)
	})

	t.Run("the actor called", func(t *testing.T) {
		t.Run("creates a new actor when called for the first time", func(t *testing.T) {
			t.Parallel()

			stage := screenplay.SetTheStage(screenplay.CastOfStandardActors())
			adam := stage.TheActorCalled("Adam")

			assert.Equal(t, "Adam", adam.Name())
		})

		t.Run("returns the same actor when called with the same name", func(t *testing.T) {
			t.Parallel()

			stage := screenplay.SetTheStage(screenplay.CastOfStandardActors())
			adam := stage.TheActorCalled("Adam")
			sameAdam := stage.TheActorCalled("Adam")

			assert.Same(t, adam, sameAdam)
		})

		t.Run("matches actor names case-insensitively", func(t *testing.T) {
			t.Parallel()

			stage := screenplay.SetTheStage(screenplay.CastOfStandardActors())
			adam := stage.TheActorCalled("Adam")
			sameAdam := stage.TheActorCalled("adam")

			assert.Same(t, adam, sameAdam)
		})

		t.Run("prepares the actor through the cast", func(t *testing.T) {
			t.Parallel()

			performTesting := performTestingAbility{}
			stage := screenplay.SetTheStage(screenplay.CastWhereEveryoneCan(performTesting))
			adam := stage.TheActorCalled("Adam")

			assert.True(t, adam.HasAbilityTo(performTesting))
		})

		t.Run("places the actor in the spotlight", func(t *testing.T) {
			t.Parallel()

			stage := screenplay.SetTheStage(screenplay.CastOfStandardActors())
			adam := stage.TheActorCalled("Adam")

			spotlight, err := stage.TheActorInTheSpotlight()
			require.NoError(t, err)
			assert.Same(t, adam, spotlight)
		})

		t.Run("moves the spotlight to the last actor called", func(t *testing.T) {
			t.Parallel()

			stage := screenplay.SetTheStage(screenplay.CastOfStandardActors())
			stage.TheActorCalled("Adam")
			bob := stage.TheActorCalled("Bob")

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
			stage.TheActorCalled("Adam")
			assert.True(t, stage.AnActorIsOnStage())
		})
	})

	t.Run("draw the curtain", func(t *testing.T) {
		t.Run("makes all actors exit the stage", func(t *testing.T) {
			t.Parallel()

			performTesting := performTestingAbility{}
			stage := screenplay.SetTheStage(screenplay.CastWhereEveryoneCan(performTesting))
			adam := stage.TheActorCalled("Adam")
			stage.TheActorCalled("Bob")

			require.NoError(t, stage.DrawTheCurtain())

			assert.Equal(t, 0, adam.NumAbilities())
			assert.False(t, stage.AnActorIsOnStage())
		})

		t.Run("clears the spotlight", func(t *testing.T) {
			t.Parallel()

			stage := screenplay.SetTheStage(screenplay.CastOfStandardActors())
			stage.TheActorCalled("Adam")

			require.NoError(t, stage.DrawTheCurtain())

			_, err := stage.TheActorInTheSpotlight()
			assert.ErrorIs(t, err, screenplay.ErrNoActorInTheSpotlight)
		})

		t.Run("allows new actors to enter after the curtain was drawn", func(t *testing.T) {
			t.Parallel()

			stage := screenplay.SetTheStage(screenplay.CastOfStandardActors())
			stage.TheActorCalled("Adam")
			require.NoError(t, stage.DrawTheCurtain())

			bob := stage.TheActorCalled("Bob")
			spotlight, err := stage.TheActorInTheSpotlight()
			require.NoError(t, err)
			assert.Same(t, bob, spotlight)
		})
	})
}
