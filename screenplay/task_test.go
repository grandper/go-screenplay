package screenplay_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/action"
	"github.com/grandper/go-screenplay/fixture"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestTask(t *testing.T) {
	adam := screenplay.ActorNamed("Adam")
	formField := fixture.NewFakeQuestion("form field", "hello world")
	missingFormField := fixture.NewFailingFakeQuestion("form field", errors.New("failed to get the field content"))

	t.Run("can be created", func(t *testing.T) {
		t.Run("with a description", func(t *testing.T) {
			t.Parallel()

			doNothing := screenplay.TaskWhere("do nothing")
			assert.Equal(t, "do nothing", doNothing.String())
			assert.NoError(t, adam.AttemptsTo(doNothing))
		})

		t.Run("with a description and successful steps", func(t *testing.T) {
			t.Parallel()

			logFormField := screenplay.TaskWhere("log the form field",
				action.Log(formField))
			assert.Equal(t, "log the form field", logFormField.String())
			assert.NoError(t, adam.AttemptsTo(logFormField))
		})

		t.Run("with a description and failing steps", func(t *testing.T) {
			t.Parallel()

			logFormField := screenplay.TaskWhere("log the form field",
				action.Log(missingFormField))
			assert.Equal(t, "log the form field", logFormField.String())
			assert.Error(t, adam.AttemptsTo(logFormField))
		})
	})
}
