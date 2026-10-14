package sentence_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/internal/sentence"
	"github.com/grandper/go-screenplay/v2/resolution/be"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestClause(t *testing.T) {
	t.Parallel()

	t.Run("each verb reads as it is written", func(t *testing.T) {
		t.Parallel()

		expected := map[sentence.Verb]string{
			sentence.Should:        "should",
			sentence.ShouldBe:      "should be",
			sentence.ShouldHave:    "should have",
			sentence.ShouldNot:     "should not",
			sentence.ShouldNotBe:   "should not be",
			sentence.ShouldNotHave: "should not have",
			sentence.ShouldAnyOf:   "should be any of:",
		}
		for verb, words := range expected {
			assert.Equal(t, words, verb.String())
		}
	})

	t.Run("an expectation reads after its subject", func(t *testing.T) {
		t.Parallel()

		inRange := []screenplay.Resolution[int]{be.AtLeast(200), be.LessThan(300)}

		assert.Equal(t, "should be at least 200, and be less than 300",
			sentence.NewClause(sentence.Should, inRange).String())
		assert.Equal(t, "should not be at least 200, nor be less than 300",
			sentence.NewClause(sentence.ShouldNot, inRange).String())
		assert.Equal(t, "should be any of: equal to 200, equal to 201",
			sentence.NewClause(sentence.ShouldAnyOf,
				[]screenplay.Resolution[int]{be.EqualTo(200), be.EqualTo(201)}).String())
	})

	t.Run("a condition reads as a fact about a singular or a plural subject", func(t *testing.T) {
		t.Parallel()

		ok := []screenplay.Resolution[int]{be.EqualTo(200)}

		assert.Equal(t, "happens to be equal to 200", sentence.NewCondition(sentence.Should, false, ok).String())
		assert.Equal(t, "happen to be equal to 200", sentence.NewCondition(sentence.Should, true, ok).String())
		assert.Equal(t, "does not happen to be equal to 200",
			sentence.NewCondition(sentence.ShouldNot, false, ok).String())
		assert.Equal(t, "do not happen to be equal to 200",
			sentence.NewCondition(sentence.ShouldNot, true, ok).String())
	})

	t.Run("judges a value against every resolution of Should", func(t *testing.T) {
		t.Parallel()

		clause := sentence.NewClause(sentence.Should, []screenplay.Resolution[int]{be.AtLeast(200), be.LessThan(300)})

		met, err := clause.Judge("the status code", 204)
		require.NoError(t, err)
		assert.True(t, met.Met)
		assert.Equal(t, "204", met.Got)

		unmet, err := clause.Judge("the status code", 404)
		require.NoError(t, err)
		assert.False(t, unmet.Met)
		assert.Equal(t, "the status code should be less than 300: got 404", unmet.Failure)
		require.Error(t, unmet.Reason)
	})

	t.Run("judges a value against none of the resolutions of ShouldNot", func(t *testing.T) {
		t.Parallel()

		clause := sentence.NewClause(sentence.ShouldNot, []screenplay.Resolution[int]{be.EqualTo(500), be.EqualTo(503)})

		met, err := clause.Judge("the status code", 200)
		require.NoError(t, err)
		assert.True(t, met.Met)

		unmet, err := clause.Judge("the status code", 503)
		require.NoError(t, err)
		assert.False(t, unmet.Met)
		assert.Equal(t, "the status code should not be equal to 503: got 503", unmet.Failure)
	})

	t.Run("judges a value against at least one resolution of ShouldAnyOf", func(t *testing.T) {
		t.Parallel()

		clause := sentence.NewClause(sentence.ShouldAnyOf, []screenplay.Resolution[int]{be.EqualTo(200), be.EqualTo(201)})

		met, err := clause.Judge("the status code", 201)
		require.NoError(t, err)
		assert.True(t, met.Met)

		unmet, err := clause.Judge("the status code", 404)
		require.NoError(t, err)
		assert.False(t, unmet.Met)
		assert.Equal(t, "the status code should be any of: equal to 200, equal to 201: got 404", unmet.Failure)
	})

	t.Run("fails with ErrInvalidTerm when it has no resolution", func(t *testing.T) {
		t.Parallel()

		_, err := sentence.NewClause[int](sentence.ShouldBe, nil).Judge("the status code", 200)

		require.ErrorIs(t, err, sentence.ErrInvalidTerm)
		assert.EqualError(t, err, "invalid term: ShouldBe needs at least one resolution")
	})

	t.Run("fails with ErrInvalidTerm when a resolution is nil", func(t *testing.T) {
		t.Parallel()

		_, err := sentence.NewCondition(sentence.Should, false, []screenplay.Resolution[int]{be.EqualTo(200), nil}).
			Judge("the status code", 200)

		require.ErrorIs(t, err, sentence.ErrInvalidTerm)
		assert.EqualError(t, err, "invalid term: the resolution at position 1 given to HappensTo is nil")
	})

	t.Run("fails with ErrInvalidResolution when a resolution was built wrongly", func(t *testing.T) {
		t.Parallel()

		clause := sentence.NewClause(sentence.Should, []screenplay.Resolution[float64]{be.CloseTo(3.0, -1)})

		_, err := clause.Judge("the price", 3.0)

		require.ErrorIs(t, err, screenplay.ErrInvalidResolution)
	})
}
