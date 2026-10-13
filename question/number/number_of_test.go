package number_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/fixture"
	"github.com/grandper/go-screenplay/v2/question/number"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestNumberOfQuestion(t *testing.T) {
	t.Parallel()

	adam := screenplay.ActorNamed("Adam")
	colors := fixture.NewFakeQuestion("colors", []string{"red", "green", "blue"})
	numbers := fixture.NewFakeQuestion("numbers", []int{1, 2, 3})
	headers := fixture.NewFakeQuestion("headers", map[string]string{"Content-Type": "application/json"})
	emptyList := fixture.NewFakeQuestion("empty list", []string{})
	emptyMap := fixture.NewFakeQuestion("empty map", map[string]int{})
	notACollection := fixture.NewFakeQuestion("not a collection", "hello")
	nilAnswer := fixture.NewFakeQuestion[any]("nil answer", nil)
	failingQuestion := fixture.NewFailingFakeQuestion[[]string](
		"failing question",
		errors.New("failed to get the answer"),
	)

	t.Run("returns the number of elements in a slice of strings", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(number.Of(colors))

		require.NoError(t, err)
		assert.Equal(t, 3, answer)
	})

	t.Run("returns the number of elements in a slice of ints", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(number.Of(numbers))

		require.NoError(t, err)
		assert.Equal(t, 3, answer)
	})

	t.Run("returns the number of entries in a map", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(number.Of(headers))

		require.NoError(t, err)
		assert.Equal(t, 1, answer)
	})

	t.Run("returns zero for an empty slice", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(number.Of(emptyList))

		require.NoError(t, err)
		assert.Equal(t, 0, answer)
	})

	t.Run("returns zero for an empty map", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(number.Of(emptyMap))

		require.NoError(t, err)
		assert.Equal(t, 0, answer)
	})

	t.Run("fails when the answer is not a slice or a map", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(number.Of(notACollection))

		require.ErrorIs(t, err, number.ErrAnswerIsNotACollection)
		assert.Zero(t, answer)
	})

	t.Run("fails when the answer is nil", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(number.Of(nilAnswer))

		require.ErrorIs(t, err, number.ErrAnswerIsNotACollection)
		assert.Zero(t, answer)
	})

	t.Run("fails when the wrapped question fails", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(number.Of(failingQuestion))

		require.Error(t, err)
		assert.Zero(t, answer)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		t.Parallel()

		question := number.Of(colors)

		assert.Equal(t, "number of colors", question.String())
	})

	t.Run("implements the question interface", func(t *testing.T) {
		t.Parallel()

		assert.Implements(t, (*screenplay.Question[int])(nil), number.Of(colors))
	})
}
