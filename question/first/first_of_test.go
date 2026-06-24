package first_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/fixture"
	"github.com/grandper/go-screenplay/question/first"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestFirstOfQuestion(t *testing.T) {
	t.Parallel()

	adam := screenplay.ActorNamed("Adam")
	colors := fixture.NewFakeQuestion("colors", []string{"red", "green", "blue"})
	numbers := fixture.NewFakeQuestion("numbers", []int{1, 2, 3})
	mixed := fixture.NewFakeQuestion("mixed", []any{"red", 1, true})
	emptyList := fixture.NewFakeQuestion("empty list", []string{})
	notASlice := fixture.NewFakeQuestion("not a slice", "hello")
	nilAnswer := fixture.NewFakeQuestion("nil answer", nil)
	failingQuestion := fixture.NewFailingFakeQuestion(
		"failing question",
		errors.New("failed to get the answer"),
	)

	t.Run("returns the first element of a slice of strings", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(first.Of(colors))

		require.NoError(t, err)
		assert.Equal(t, "red", answer)
	})

	t.Run("returns the first element of a slice of ints", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(first.Of(numbers))

		require.NoError(t, err)
		assert.Equal(t, 1, answer)
	})

	t.Run("returns the first element of a heterogeneous slice", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(first.Of(mixed))

		require.NoError(t, err)
		assert.Equal(t, "red", answer)
	})

	t.Run("fails when the answer is an empty slice", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(first.Of(emptyList))

		require.ErrorIs(t, err, first.ErrAnswerIsEmpty)
		assert.Nil(t, answer)
	})

	t.Run("fails when the answer is not a slice", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(first.Of(notASlice))

		require.ErrorIs(t, err, first.ErrAnswerIsNotASlice)
		assert.Nil(t, answer)
	})

	t.Run("fails when the answer is nil", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(first.Of(nilAnswer))

		require.ErrorIs(t, err, first.ErrAnswerIsNotASlice)
		assert.Nil(t, answer)
	})

	t.Run("fails when the wrapped question fails", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(first.Of(failingQuestion))

		require.Error(t, err)
		assert.Nil(t, answer)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		t.Parallel()

		question := first.Of(colors)

		assert.Equal(t, "first of colors", question.String())
	})
}
