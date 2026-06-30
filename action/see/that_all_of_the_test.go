package see_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/action/see"
	"github.com/grandper/go-screenplay/fixture"
	"github.com/grandper/go-screenplay/resolution/contains"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestSeeThatAllOfThe(t *testing.T) {
	t.Parallel()

	adam := screenplay.ActorNamed("Adam")

	t.Run("succeeds when all of the questions match the resolution", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, adam.AttemptsTo(see.ThatAllOfThe(
			fixture.NewFakeQuestion("widget", "hello Adam"),
			fixture.NewFakeQuestion("widget", "hello Eve"),
		)(contains.TheText("hello"))))
	})

	t.Run("succeeds when there is no question to see", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, adam.AttemptsTo(see.ThatAllOfThe()(contains.TheText("hello"))))
	})

	t.Run("fails when one of the questions does not match the resolution", func(t *testing.T) {
		t.Parallel()
		require.Error(t, adam.AttemptsTo(see.ThatAllOfThe(
			fixture.NewFakeQuestion("widget", "hello Adam"),
			fixture.NewFakeQuestion("widget", "goodbye Eve"),
		)(contains.TheText("hello"))))
	})

	t.Run("fails when the actor fails to answer one of the questions", func(t *testing.T) {
		t.Parallel()
		require.Error(t, adam.AttemptsTo(see.ThatAllOfThe(
			fixture.NewFakeQuestion("widget", "hello Adam"),
			fixture.NewFailingFakeQuestion("widget", errors.New("failed to get the widget")),
		)(contains.TheText("hello"))))
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		t.Parallel()

		action := see.ThatAllOfThe(
			fixture.NewFakeQuestion("widget", "hello Adam"),
			fixture.NewFakeQuestion("widget", "hello Eve"),
		)(contains.TheText("hello"))

		assert.Equal(
			t,
			"see if the widget is containing the text hello, and see if the widget is containing the text hello",
			action.String(),
		)
	})
}
