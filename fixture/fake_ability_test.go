package fixture_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/v2/fixture"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestFakeAbility(t *testing.T) {
	t.Parallel()

	t.Run("can be forgotten", func(t *testing.T) {
		t.Parallel()

		assert.NoError(t, fixture.NewFakeAbility().Forget())
	})

	t.Run("can be given to an actor", func(t *testing.T) {
		t.Parallel()

		adam := screenplay.ActorNamed("Adam").WhoCan(fixture.NewFakeAbility())
		assert.True(t, adam.HasAbilityTo(fixture.NewFakeAbility()))
	})
}
