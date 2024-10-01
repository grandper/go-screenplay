package action_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/action"
	"github.com/grandper/go-screenplay/fixture"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestLogAction(t *testing.T) {
	adam := screenplay.ActorNamed("Adam")
	formField := fixture.NewFakeQuestion("form field", "hello world")
	missingFormField := fixture.NewFailingFakeQuestion("form field", errors.New("failed to get the field content"))

	t.Run("should log the answer of the question", func(t *testing.T) {
		require.NoError(t, adam.AttemptsTo(action.Log(formField)))
	})

	t.Run("fails when the question fails", func(t *testing.T) {
		require.Error(t, adam.AttemptsTo(action.Log(missingFormField)))
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		action1 := action.Log(formField)
		assert.Equal(t, "log the form field", action1.String())
	})
}
