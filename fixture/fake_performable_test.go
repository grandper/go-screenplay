package fixture_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/fixture"
)

func TestFakePerformable(t *testing.T) {
	t.Parallel()

	t.Run("can be created", func(t *testing.T) {
		t.Parallel()

		fakePerformable := fixture.NewFakePerformable("fakePerformable action", nil)
		assert.NotNil(t, fakePerformable)
	})

	t.Run("returns the provided description", func(t *testing.T) {
		t.Parallel()

		fakePerformable := fixture.NewFakePerformable("fakePerformable action", nil)
		assert.Equal(t, "fakePerformable action", fakePerformable.String())
	})

	t.Run("succeeds when no error is provided", func(t *testing.T) {
		t.Parallel()

		fakePerformable := fixture.NewFakePerformable("fakePerformable action", nil)
		assert.NoError(t, fakePerformable.PerformAs(nil))
	})

	t.Run("fails when an error is provided", func(t *testing.T) {
		t.Parallel()

		fakePerformable := fixture.NewFakePerformable("fakePerformable action", assert.AnError)
		assert.ErrorIs(t, fakePerformable.PerformAs(nil), assert.AnError)
	})
}
