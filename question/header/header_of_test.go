package header_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/fixture"
	"github.com/grandper/go-screenplay/v2/question"
	"github.com/grandper/go-screenplay/v2/question/header"
	"github.com/grandper/go-screenplay/v2/screenplay"
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
	pointerResponse := fixture.NewFakeQuestion("pointer response", &response{Status: 200, Header: "text/plain"})
	noHeader := fixture.NewFakeQuestion("no header", withoutHeader{Status: 200})
	notAStruct := fixture.NewFakeQuestion("not a struct", "hello")
	nilAnswer := fixture.NewFakeQuestion[any]("nil answer", nil)
	failingQuestion := fixture.NewFailingFakeQuestion[response](
		"failing question",
		errors.New("failed to get the answer"),
	)

	t.Run("returns the Header field of a struct", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(header.Of[string](textResponse))

		require.NoError(t, err)
		assert.Equal(t, "application/json", answer)
	})

	t.Run("returns the Header field typed, ready to use without assertion", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(header.Of[string](textResponse))

		require.NoError(t, err)
		assert.Equal(t, strings.ToUpper("application/json"), strings.ToUpper(answer))
	})

	t.Run("returns the Header field whatever its type", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(header.Of[int](numberResponse))

		require.NoError(t, err)
		assert.Equal(t, 2, answer)
	})

	t.Run("returns the Header field as any when the type is not specified", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(header.Of[any](numberResponse))

		require.NoError(t, err)
		assert.Equal(t, 2, answer)
	})

	t.Run("returns the Header field of a pointer to a struct", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(header.Of[string](pointerResponse))

		require.NoError(t, err)
		assert.Equal(t, "text/plain", answer)
	})

	t.Run("fails when the Header field is not of the expected type", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(header.Of[int](textResponse))

		require.ErrorIs(t, err, question.ErrFieldHasWrongType)
		assert.Zero(t, answer)
	})

	t.Run("fails when the struct has no Header field", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(header.Of[any](noHeader))

		require.ErrorIs(t, err, question.ErrFieldNotFound)
		assert.Nil(t, answer)
	})

	t.Run("fails when the answer is not a struct", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(header.Of[any](notAStruct))

		require.ErrorIs(t, err, question.ErrAnswerIsNotAStruct)
		assert.Nil(t, answer)
	})

	t.Run("fails when the answer is nil", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(header.Of[any](nilAnswer))

		require.ErrorIs(t, err, question.ErrAnswerIsNotAStruct)
		assert.Nil(t, answer)
	})

	t.Run("fails when the wrapped question fails", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(header.Of[string](failingQuestion))

		require.Error(t, err)
		assert.Empty(t, answer)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		t.Parallel()

		question := header.Of[string](textResponse)

		assert.Equal(t, "header of response", question.String())
	})

	t.Run("implements the question interface", func(t *testing.T) {
		t.Parallel()

		assert.Implements(t, (*screenplay.Question[string])(nil), header.Of[string](textResponse))
	})
}
