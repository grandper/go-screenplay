package data_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/fixture"
	"github.com/grandper/go-screenplay/v2/question/data"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestDataOfQuestion(t *testing.T) {
	t.Parallel()

	adam := screenplay.ActorNamed("Adam")

	type message struct {
		Subject string
		Data    string
	}
	type counter struct{ Data int }
	type withoutData struct{ Subject string }

	textMessage := fixture.NewFakeQuestion("message", message{Subject: "greetings", Data: "OK"})
	numberMessage := fixture.NewFakeQuestion("counter", counter{Data: 2})
	pointerMessage := fixture.NewFakeQuestion("pointer message", &message{Subject: "greetings", Data: "created"})
	noData := fixture.NewFakeQuestion("no data", withoutData{Subject: "greetings"})
	notAStruct := fixture.NewFakeQuestion("not a struct", "hello")
	nilAnswer := fixture.NewFakeQuestion[any]("nil answer", nil)
	failingQuestion := fixture.NewFailingFakeQuestion[message](
		"failing question",
		errors.New("failed to get the answer"),
	)

	t.Run("returns the Data field of a struct", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(data.Of[string](textMessage))

		require.NoError(t, err)
		assert.Equal(t, "OK", answer)
	})

	t.Run("returns the Data field typed, ready to use without assertion", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(data.Of[string](textMessage))

		require.NoError(t, err)
		assert.Equal(t, strings.ToUpper("OK"), strings.ToUpper(answer))
	})

	t.Run("returns the Data field whatever its type", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(data.Of[int](numberMessage))

		require.NoError(t, err)
		assert.Equal(t, 2, answer)
	})

	t.Run("returns the Data field as any when the type is not specified", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(data.Of[any](numberMessage))

		require.NoError(t, err)
		assert.Equal(t, 2, answer)
	})

	t.Run("returns the Data field of a pointer to a struct", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(data.Of[string](pointerMessage))

		require.NoError(t, err)
		assert.Equal(t, "created", answer)
	})

	t.Run("fails when the Data field is not of the expected type", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(data.Of[int](textMessage))

		require.ErrorIs(t, err, data.ErrFieldHasWrongType)
		assert.Zero(t, answer)
	})

	t.Run("fails when the struct has no Data field", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(data.Of[any](noData))

		require.ErrorIs(t, err, data.ErrFieldNotFound)
		assert.Nil(t, answer)
	})

	t.Run("fails when the answer is not a struct", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(data.Of[any](notAStruct))

		require.ErrorIs(t, err, data.ErrAnswerIsNotAStruct)
		assert.Nil(t, answer)
	})

	t.Run("fails when the answer is nil", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(data.Of[any](nilAnswer))

		require.ErrorIs(t, err, data.ErrAnswerIsNotAStruct)
		assert.Nil(t, answer)
	})

	t.Run("fails when the wrapped question fails", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(data.Of[string](failingQuestion))

		require.Error(t, err)
		assert.Empty(t, answer)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		t.Parallel()

		question := data.Of[string](textMessage)

		assert.Equal(t, "data of message", question.String())
	})

	t.Run("implements the question interface", func(t *testing.T) {
		t.Parallel()

		assert.Implements(t, (*screenplay.Question[string])(nil), data.Of[string](textMessage))
	})
}
