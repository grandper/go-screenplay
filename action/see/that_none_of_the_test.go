package see_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/action/see"
	"github.com/grandper/go-screenplay/fixture"
	"github.com/grandper/go-screenplay/resolution/contain"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestSeeThatNoneOfThe(t *testing.T) {
	t.Parallel()

	adam := screenplay.ActorNamed("Adam")

	t.Run("succeeds when none of the questions match the resolution", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, adam.AttemptsTo(see.ThatNoneOfThe(
			fixture.NewFakeQuestion("widget", "goodbye Adam"),
			fixture.NewFakeQuestion("widget", "farewell Eve"),
		)(contain.TheText("hello"))))
	})

	t.Run("succeeds when there is no question to see", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, adam.AttemptsTo(see.ThatNoneOfThe()(contain.TheText("hello"))))
	})

	t.Run("fails when one of the questions matches the resolution", func(t *testing.T) {
		t.Parallel()
		require.Error(t, adam.AttemptsTo(see.ThatNoneOfThe(
			fixture.NewFakeQuestion("widget", "hello Adam"),
			fixture.NewFakeQuestion("widget", "farewell Eve"),
		)(contain.TheText("hello"))))
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		t.Parallel()

		action := see.ThatNoneOfThe(
			fixture.NewFakeQuestion("widget", "goodbye Adam"),
			fixture.NewFakeQuestion("widget", "farewell Eve"),
		)(contain.TheText("hello"))

		assert.Equal(
			t,
			"see if the widget is not containing the text hello, and see if the widget is not containing the text hello",
			action.String(),
		)
	})
}
