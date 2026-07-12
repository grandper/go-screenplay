package screenplay_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/fixture"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestProduction(t *testing.T) {
	t.Parallel()

	t.Run("defaults the timeout and the polling interval", func(t *testing.T) {
		t.Parallel()

		production := screenplay.NewProduction()

		assert.Equal(t, screenplay.DefaultTimeout, production.Timeout())
		assert.Equal(t, screenplay.DefaultPolling, production.Polling())
		assert.Equal(t, screenplay.DefaultForceAllNarration, production.ForceAllNarration())
		assert.Nil(t, production.Narrator())
	})

	t.Run("can be configured with options", func(t *testing.T) {
		t.Parallel()

		narrator := screenplay.NewNarrator()
		production := screenplay.NewProduction(
			screenplay.WithNarrator(narrator),
			screenplay.WithTimeout(10*time.Second),
			screenplay.WithPolling(250*time.Millisecond),
			screenplay.WithForceAllNarration(),
		)

		assert.Same(t, narrator, production.Narrator())
		assert.Equal(t, 10*time.Second, production.Timeout())
		assert.Equal(t, 250*time.Millisecond, production.Polling())
		assert.True(t, production.ForceAllNarration())
	})

	t.Run("sets a stage whose actors narrate through the narrator", func(t *testing.T) {
		t.Parallel()

		recorder := fixture.NewRecorder()
		production := screenplay.NewProduction(
			screenplay.WithNarrator(screenplay.NewNarrator(recorder)))

		stage := production.SetTheStage(screenplay.CastOfStandardActors())
		adam := stage.ActorNamed("Adam")

		require.NoError(t, adam.AttemptsTo(fixture.NewFakePerformable("waves", nil)))
		assert.Equal(t, []string{"Adam waves", "Adam waves"}, recorder.Messages())
	})

	t.Run("sets a stage even without a narrator", func(t *testing.T) {
		t.Parallel()

		stage := screenplay.NewProduction().SetTheStage(screenplay.CastOfStandardActors())
		adam := stage.ActorNamed("Adam")

		assert.Equal(t, "Adam", adam.Name())
	})
}
