package see_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/action/see"
	"github.com/grandper/go-screenplay/fixture"
	"github.com/grandper/go-screenplay/resolution/is"
	"github.com/grandper/go-screenplay/resolution/testdata"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestSeeTheAction(t *testing.T) {
	t.Parallel()

	adam := screenplay.ActorNamed("Adam")
	formField := fixture.NewFakeQuestion("form field", "hello world")
	missingFormField := fixture.NewFailingFakeQuestion("form field", errors.New("failed to get the field content"))
	isEqualButFails := testdata.NewFailingResolution(
		"equal to hello world",
		errors.New("failed to match the content of the field"),
	)

	t.Run("should see something", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, adam.AttemptsTo(see.The(formField, is.EqualTo("hello world"))))
	})

	t.Run("fails when there is nothing to see", func(t *testing.T) {
		t.Parallel()
		require.Error(t, adam.AttemptsTo(see.The(formField, is.EqualTo("hello everybody"))))
	})

	t.Run("fails when the actor fails to answer the question", func(t *testing.T) {
		t.Parallel()
		require.Error(t, adam.AttemptsTo(see.The(missingFormField, is.EqualTo("hello everybody"))))
	})

	t.Run("fails when the resolution fails", func(t *testing.T) {
		t.Parallel()
		require.Error(t, adam.AttemptsTo(see.The(formField, isEqualButFails)))
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		t.Parallel()

		action := see.The(formField, is.EqualTo("hello world"))

		assert.Equal(t, "see if the form field is equal to hello world", action.String())
	})
}
