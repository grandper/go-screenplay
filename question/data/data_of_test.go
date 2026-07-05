package data_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/fixture"
	"github.com/grandper/go-screenplay/question/data"
	"github.com/grandper/go-screenplay/screenplay"
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
	pointerMessage := fixture.NewFakeQuestion("pointer message", &message{Subject: "created", Data: "created"})
	noData := fixture.NewFakeQuestion("no data", withoutData{Subject: "empty"})
	notAStruct := fixture.NewFakeQuestion("not a struct", "hello")
	nilAnswer := fixture.NewFakeQuestion("nil answer", nil)
	failingQuestion := fixture.NewFailingFakeQuestion(
		"failing question",
		errors.New("failed to get the answer"),
	)

	t.Run("returns the Data field of a struct", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(data.Of(textMessage))

		require.NoError(t, err)
		assert.Equal(t, "OK", answer)
	})

	t.Run("returns the Data field whatever its type", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(data.Of(numberMessage))

		require.NoError(t, err)
		assert.Equal(t, 2, answer)
	})

	t.Run("returns the Data field of a pointer to a struct", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(data.Of(pointerMessage))

		require.NoError(t, err)
		assert.Equal(t, "created", answer)
	})

	t.Run("fails when the struct has no Data field", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(data.Of(noData))

		require.ErrorIs(t, err, data.ErrFieldNotFound)
		assert.Nil(t, answer)
	})

	t.Run("fails when the answer is not a struct", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(data.Of(notAStruct))

		require.ErrorIs(t, err, data.ErrAnswerIsNotAStruct)
		assert.Nil(t, answer)
	})

	t.Run("fails when the answer is nil", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(data.Of(nilAnswer))

		require.ErrorIs(t, err, data.ErrAnswerIsNotAStruct)
		assert.Nil(t, answer)
	})

	t.Run("fails when the wrapped question fails", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(data.Of(failingQuestion))

		require.Error(t, err)
		assert.Nil(t, answer)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		t.Parallel()

		question := data.Of(textMessage)

		assert.Equal(t, "data of message", question.String())
	})
}
