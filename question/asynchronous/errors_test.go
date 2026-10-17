package asynchronous_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action"
	"github.com/grandper/go-screenplay/v2/action/proceed"
	"github.com/grandper/go-screenplay/v2/action/remember"
	"github.com/grandper/go-screenplay/v2/question/asynchronous"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestErrors(t *testing.T) {
	t.Run("is empty when nothing was launched asynchronously", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")

		errs, err := asynchronous.Errors().AnsweredBy(adam)
		require.NoError(t, err)
		assert.Empty(t, errs)
	})

	t.Run("waits for the asynchronous actions and returns their errors", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		rebuildTheIndex := action.FromFunc("rebuild the index", func(_ *screenplay.Actor) error {
			return errors.New("the server is down")
		})

		require.NoError(t, adam.AttemptsTo(proceed.AsynchronouslyTo(rebuildTheIndex)))

		errs, err := asynchronous.Errors().AnsweredBy(adam)
		require.NoError(t, err)
		require.Len(t, errs, 1)
		assert.EqualError(t, errs[0], "rebuild the index: the server is down")
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "asynchronous errors", asynchronous.Errors().String())
		// The sentence builders add the article once.
		assert.Equal(
			t,
			"remember the asynchronous errors as \"errs\"",
			remember.The(asynchronous.Errors()).As("errs").String(),
		)
	})
}
