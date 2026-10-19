package last_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/question"
	"github.com/grandper/go-screenplay/v2/question/last"
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/grandper/go-screenplay/v2/screenplay/fixture"
)

func TestLastOfQuestion(t *testing.T) {
	t.Parallel()

	adam := screenplay.ActorNamed("Adam")
	colors := fixture.NewFakeQuestion("colors", []string{"red", "green", "blue"})
	numbers := fixture.NewFakeQuestion("numbers", []int{1, 2, 3})
	mixed := fixture.NewFakeQuestion("mixed", []any{"red", 1, true})
	emptyList := fixture.NewFakeQuestion("empty list", []string{})
	failingQuestion := fixture.NewFailingFakeQuestion[[]string](
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
		require.ErrorIs(t, err, question.ErrAnswerIsEmpty)
		assert.Empty(t, answer)
	})

	t.Run("fails when the wrapped question fails", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(last.Of(failingQuestion))

		require.Error(t, err)
		assert.Empty(t, answer)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		t.Parallel()

		question := last.Of(colors)

		assert.Equal(t, "last of colors", question.String())
	})

	t.Run("implements the question interface", func(t *testing.T) {
		t.Parallel()

		assert.Implements(t, (*screenplay.Question[string])(nil), last.Of(colors))
	})
}
