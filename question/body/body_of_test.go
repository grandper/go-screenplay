package body_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/fixture"
	"github.com/grandper/go-screenplay/question/body"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestBodyOfQuestion(t *testing.T) {
	t.Parallel()

	adam := screenplay.ActorNamed("Adam")

	type response struct {
		Status int
		Body   string
	}
	type counter struct{ Body int }
	type withoutBody struct{ Status int }

	textResponse := fixture.NewFakeQuestion("response", response{Status: 200, Body: "OK"})
	numberResponse := fixture.NewFakeQuestion("counter", counter{Body: 2})
	pointerResponse := fixture.NewFakeQuestion("pointer response", &response{Status: 201, Body: "created"})
	noBody := fixture.NewFakeQuestion("no body", withoutBody{Status: 500})
	notAStruct := fixture.NewFakeQuestion("not a struct", "hello")
	nilAnswer := fixture.NewFakeQuestion("nil answer", nil)
	failingQuestion := fixture.NewFailingFakeQuestion(
		"failing question",
		errors.New("failed to get the answer"),
	)

	t.Run("returns the Body field of a struct", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(body.Of(textResponse))

		require.NoError(t, err)
		assert.Equal(t, "OK", answer)
	})

	t.Run("returns the Body field whatever its type", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(body.Of(numberResponse))

		require.NoError(t, err)
		assert.Equal(t, 2, answer)
	})

	t.Run("returns the Body field of a pointer to a struct", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(body.Of(pointerResponse))

		require.NoError(t, err)
		assert.Equal(t, "created", answer)
	})

	t.Run("fails when the struct has no Body field", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(body.Of(noBody))

		require.ErrorIs(t, err, body.ErrFieldNotFound)
		assert.Nil(t, answer)
	})

	t.Run("fails when the answer is not a struct", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(body.Of(notAStruct))

		require.ErrorIs(t, err, body.ErrAnswerIsNotAStruct)
		assert.Nil(t, answer)
	})

	t.Run("fails when the answer is nil", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(body.Of(nilAnswer))

		require.ErrorIs(t, err, body.ErrAnswerIsNotAStruct)
		assert.Nil(t, answer)
	})

	t.Run("fails when the wrapped question fails", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(body.Of(failingQuestion))

		require.Error(t, err)
		assert.Nil(t, answer)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		t.Parallel()

		question := body.Of(textResponse)

		assert.Equal(t, "body of response", question.String())
	})
}
