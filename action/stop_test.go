package action_test

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/action"
	"github.com/grandper/go-screenplay/fixture"
	"github.com/grandper/go-screenplay/resolution/is"
	"github.com/grandper/go-screenplay/resolution/testdata"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestStopAction(t *testing.T) {
	adam := screenplay.ActorNamed("Adam")
	formField := fixture.NewFakeQuestion("form field", "")
	failingFormField := fixture.NewFailingFakeQuestion("form field", errors.New("failed to find the form field"))
	isEqualFailed := testdata.NewFailingResolution("is equal", errors.New("resolution failed"))

	t.Run("should stop until the actor press enter", func(t *testing.T) {
		r, w, err := os.Pipe()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, r.Close())
		}()

		input := []byte("\n")
		go func() {
			_, err = w.Write(input)
			assert.NoError(t, err)
			assert.NoError(t, w.Close())
		}()
		require.NoError(t, adam.AttemptsTo(action.Stop().UntilAnInputIsProvidedBy(r)))
	})

	t.Run("should stop until a question is answered", func(t *testing.T) {
		formFieldWillChange := fixture.NewFakeQuestion("form field", "")

		go func() {
			time.Sleep(5 * time.Millisecond)

			formFieldWillChange.AnswerWith("hello world")
		}()

		require.NoError(t, adam.AttemptsTo(action.Stop().UntilThe(formFieldWillChange, is.EqualTo("hello world"))))
	})

	t.Run("fails if the attempt continuously failed", func(t *testing.T) {
		err := adam.AttemptsTo(action.Stop().UntilThe(formField, is.EqualTo("hello world")))
		require.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("stopped for %s,", screenplay.DefaultTimeout))
	})

	t.Run("fails when the question fails", func(t *testing.T) {
		require.Error(t, adam.AttemptsTo(action.Stop().UntilThe(failingFormField, is.EqualTo("hello world"))))
	})

	t.Run("fails when the resolution fails", func(t *testing.T) {
		require.Error(t, adam.AttemptsTo(action.Stop().UntilThe(formField, isEqualFailed)))
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		action1 := action.Stop()
		assert.Equal(t, "stop until the 'enter' key is pressed", action1.String())

		action2 := action.Stop().UntilThe(formField, is.EqualTo("finished"))
		assert.Equal(t, "stop until the form field is equal to finished", action2.String())
	})

	t.Run("fails when the reader is closed", func(t *testing.T) {
		r, _, err := os.Pipe()
		require.NoError(t, err)
		require.NoError(t, r.Close())
		require.Error(t, adam.AttemptsTo(action.Stop().UntilAnInputIsProvidedBy(r)))
	})
}
