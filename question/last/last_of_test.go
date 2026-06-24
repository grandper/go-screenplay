package last_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/fixture"
	"github.com/grandper/go-screenplay/question/last"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestLastOfQuestion(t *testing.T) {
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

	t.Run("returns the last element of a slice of strings", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(last.Of(colors))

		require.NoError(t, err)
		assert.Equal(t, "blue", answer)
	})

	t.Run("returns the last element of a slice of ints", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(last.Of(numbers))

		require.NoError(t, err)
		assert.Equal(t, 3, answer)
	})

	t.Run("returns the last element of a heterogeneous slice", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(last.Of(mixed))

		require.NoError(t, err)
		assert.Equal(t, true, answer)
	})

	t.Run("fails when the answer is an empty slice", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(last.Of(emptyList))

		require.ErrorIs(t, err, last.ErrAnswerIsEmpty)
		assert.Nil(t, answer)
	})

	t.Run("fails when the answer is not a slice", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(last.Of(notASlice))

		require.ErrorIs(t, err, last.ErrAnswerIsNotASlice)
		assert.Nil(t, answer)
	})

	t.Run("fails when the answer is nil", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(last.Of(nilAnswer))

		require.ErrorIs(t, err, last.ErrAnswerIsNotASlice)
		assert.Nil(t, answer)
	})

	t.Run("fails when the wrapped question fails", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(last.Of(failingQuestion))

		require.Error(t, err)
		assert.Nil(t, answer)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		t.Parallel()

		question := last.Of(colors)

		assert.Equal(t, "last of colors", question.String())
	})
}
