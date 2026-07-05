package header_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/fixture"
	"github.com/grandper/go-screenplay/question/header"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestHeaderOfQuestion(t *testing.T) {
	t.Parallel()

	adam := screenplay.ActorNamed("Adam")

	type response struct {
		Status int
		Header string
	}
	type counter struct{ Header int }
	type withoutHeader struct{ Status int }

	textResponse := fixture.NewFakeQuestion("response", response{Status: 200, Header: "application/json"})
	numberResponse := fixture.NewFakeQuestion("counter", counter{Header: 2})
	pointerResponse := fixture.NewFakeQuestion("pointer response", &response{Status: 201, Header: "text/plain"})
	noHeader := fixture.NewFakeQuestion("no header", withoutHeader{Status: 500})
	notAStruct := fixture.NewFakeQuestion("not a struct", "hello")
	nilAnswer := fixture.NewFakeQuestion("nil answer", nil)
	failingQuestion := fixture.NewFailingFakeQuestion(
		"failing question",
		errors.New("failed to get the answer"),
	)

	t.Run("returns the Header field of a struct", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(header.Of(textResponse))

		require.NoError(t, err)
		assert.Equal(t, "application/json", answer)
	})

	t.Run("returns the Header field whatever its type", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(header.Of(numberResponse))

		require.NoError(t, err)
		assert.Equal(t, 2, answer)
	})

	t.Run("returns the Header field of a pointer to a struct", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(header.Of(pointerResponse))

		require.NoError(t, err)
		assert.Equal(t, "text/plain", answer)
	})

	t.Run("fails when the struct has no Header field", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(header.Of(noHeader))

		require.ErrorIs(t, err, header.ErrFieldNotFound)
		assert.Nil(t, answer)
	})

	t.Run("fails when the answer is not a struct", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(header.Of(notAStruct))

		require.ErrorIs(t, err, header.ErrAnswerIsNotAStruct)
		assert.Nil(t, answer)
	})

	t.Run("fails when the answer is nil", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(header.Of(nilAnswer))

		require.ErrorIs(t, err, header.ErrAnswerIsNotAStruct)
		assert.Nil(t, answer)
	})

	t.Run("fails when the wrapped question fails", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(header.Of(failingQuestion))

		require.Error(t, err)
		assert.Nil(t, answer)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		t.Parallel()

		question := header.Of(textResponse)

		assert.Equal(t, "header of response", question.String())
	})
}
