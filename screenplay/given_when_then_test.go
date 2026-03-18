package screenplay_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/fixture"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestGivenWhenThenHelpers(t *testing.T) {
	t.Parallel()

	t.Run("Given returns the same actor pointer", func(t *testing.T) {
		t.Parallel()
		adam := screenplay.ActorNamed("Adam")
		assert.Same(t, adam, screenplay.Given(adam))
	})

	t.Run("When returns the same actor pointer", func(t *testing.T) {
		t.Parallel()
		adam := screenplay.ActorNamed("Adam")
		assert.Same(t, adam, screenplay.When(adam))
	})

	t.Run("Then returns the same actor pointer", func(t *testing.T) {
		t.Parallel()
		adam := screenplay.ActorNamed("Adam")
		assert.Same(t, adam, screenplay.Then(adam))
	})

	t.Run("And returns the same actor pointer", func(t *testing.T) {
		t.Parallel()
		adam := screenplay.ActorNamed("Adam")
		assert.Same(t, adam, screenplay.And(adam))
	})

	t.Run("supports method chaining on the returned value", func(t *testing.T) {
		t.Parallel()
		adam := screenplay.ActorNamed("Adam")
		task := fixture.NewFakePerformable("some task", nil)
		require.NoError(t, screenplay.Given(adam).AttemptsTo(task))
	})
}
