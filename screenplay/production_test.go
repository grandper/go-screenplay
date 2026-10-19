package screenplay_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/grandper/go-screenplay/v2/screenplay/fixture"
)

func TestProduction(t *testing.T) {
	t.Parallel()

	t.Run("defaults the timeout, the polling interval and the buffered size", func(t *testing.T) {
		t.Parallel()

		production := screenplay.NewProduction()

		assert.Equal(t, screenplay.DefaultTimeout, production.Timeout())
		assert.Equal(t, screenplay.DefaultPolling, production.Polling())
		assert.Equal(t, screenplay.DefaultMaxBufferedSize, production.MaxBufferedSize())
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
			screenplay.WithMaxBufferedSize(8<<20),
			screenplay.WithForceAllNarration(),
		)

		assert.Same(t, narrator, production.Narrator())
		assert.Equal(t, 10*time.Second, production.Timeout())
		assert.Equal(t, 250*time.Millisecond, production.Polling())
		assert.Equal(t, 8<<20, production.MaxBufferedSize())
		assert.True(t, production.ForceAllNarration())
	})

	t.Run("sets a stage whose actors narrate through the narrator", func(t *testing.T) {
		t.Parallel()

		recorder := fixture.NewRecorder()
		production := screenplay.NewProduction(
			screenplay.WithNarrator(screenplay.NewNarrator(recorder)))

		stage := production.SetTheStage(screenplay.CastOfStandardActors())
		adam := stage.ActorNamed("Adam")

		require.NoError(t, adam.AttemptsTo(fixture.NewFakePerformable("wave", nil)))
		assert.Equal(t, []string{"Adam attempts to wave", "Adam attempts to wave"}, recorder.Messages())
	})

	t.Run("sets a stage even without a narrator", func(t *testing.T) {
		t.Parallel()

		stage := screenplay.NewProduction().SetTheStage(screenplay.CastOfStandardActors())
		adam := stage.ActorNamed("Adam")

		assert.Equal(t, "Adam", adam.Name())
	})
}
