package see_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/action/see"
	"github.com/grandper/go-screenplay/fixture"
	"github.com/grandper/go-screenplay/resolution/contains"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestSeeThatAnyOfThe(t *testing.T) {
	t.Parallel()

	adam := screenplay.ActorNamed("Adam")

	t.Run("succeeds when at least one of the questions matches the resolution", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, adam.AttemptsTo(see.ThatAnyOfThe(
			fixture.NewFakeQuestion("widget", "goodbye Adam"),
			fixture.NewFakeQuestion("widget", "hello Eve"),
		)(contains.TheText("hello"))))
	})

	t.Run("succeeds when there is no question to see", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, adam.AttemptsTo(see.ThatAnyOfThe()(contains.TheText("hello"))))
	})

	t.Run("fails when none of the questions match the resolution", func(t *testing.T) {
		t.Parallel()
		require.Error(t, adam.AttemptsTo(see.ThatAnyOfThe(
			fixture.NewFakeQuestion("widget", "goodbye Adam"),
			fixture.NewFakeQuestion("widget", "farewell Eve"),
		)(contains.TheText("hello"))))
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		t.Parallel()

		action := see.ThatAnyOfThe(
			fixture.NewFakeQuestion("widget", "hello Adam"),
			fixture.NewFakeQuestion("widget", "hello Eve"),
		)(contains.TheText("hello"))

		assert.Equal(
			t,
			"see if the widget is containing the text hello, or see if the widget is containing the text hello",
			action.String(),
		)
	})
}
