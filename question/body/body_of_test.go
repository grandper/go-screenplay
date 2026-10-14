package body_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/fixture"
	"github.com/grandper/go-screenplay/v2/question"
	"github.com/grandper/go-screenplay/v2/question/body"
	"github.com/grandper/go-screenplay/v2/screenplay"
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
	pointerResponse := fixture.NewFakeQuestion("pointer response", &response{Status: 200, Body: "created"})
	noBody := fixture.NewFakeQuestion("no body", withoutBody{Status: 200})
	notAStruct := fixture.NewFakeQuestion("not a struct", "hello")
	nilAnswer := fixture.NewFakeQuestion[any]("nil answer", nil)
	failingQuestion := fixture.NewFailingFakeQuestion[response](
		"failing question",
		errors.New("failed to get the answer"),
	)

	t.Run("returns the Body field of a struct", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(body.Of[string](textResponse))

		require.NoError(t, err)
		assert.Equal(t, "OK", answer)
	})

	t.Run("returns the Body field typed, ready to use without assertion", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(body.Of[string](textResponse))

		require.NoError(t, err)
		assert.Equal(t, strings.ToUpper("OK"), strings.ToUpper(answer))
	})

	t.Run("returns the Body field whatever its type", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(body.Of[int](numberResponse))

		require.NoError(t, err)
		assert.Equal(t, 2, answer)
	})

	t.Run("returns the Body field as any when the type is not specified", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(body.Of[any](numberResponse))

		require.NoError(t, err)
		assert.Equal(t, 2, answer)
	})

	t.Run("returns the Body field of a pointer to a struct", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(body.Of[string](pointerResponse))

		require.NoError(t, err)
		assert.Equal(t, "created", answer)
	})

	t.Run("fails when the Body field is not of the expected type", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(body.Of[int](textResponse))

		require.ErrorIs(t, err, question.ErrFieldHasWrongType)
		assert.Zero(t, answer)
	})

	t.Run("fails when the struct has no Body field", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(body.Of[any](noBody))

		require.ErrorIs(t, err, question.ErrFieldNotFound)
		assert.Nil(t, answer)
	})

	t.Run("fails when the answer is not a struct", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(body.Of[any](notAStruct))

		require.ErrorIs(t, err, question.ErrAnswerIsNotAStruct)
		assert.Nil(t, answer)
	})

	t.Run("fails when the answer is nil", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(body.Of[any](nilAnswer))

		require.ErrorIs(t, err, question.ErrAnswerIsNotAStruct)
		assert.Nil(t, answer)
	})

	t.Run("fails when the wrapped question fails", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(body.Of[string](failingQuestion))

		require.Error(t, err)
		assert.Empty(t, answer)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		t.Parallel()

		question := body.Of[string](textResponse)

		assert.Equal(t, "body of response", question.String())
	})

	t.Run("implements the question interface", func(t *testing.T) {
		t.Parallel()

		assert.Implements(t, (*screenplay.Question[string])(nil), body.Of[string](textResponse))
	})
}
