package sentence_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/fixture"
	"github.com/grandper/go-screenplay/v2/internal/sentence"
	"github.com/grandper/go-screenplay/v2/resolution/be"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestTerm(t *testing.T) {
	t.Parallel()

	adam := screenplay.ActorNamed("Adam")
	statusCode := fixture.NewFakeQuestion("status code", 404)
	shouldBeOK := sentence.NewClause(sentence.Should, []screenplay.Resolution[int]{be.EqualTo(200)})

	t.Run("describes itself with the article before the question", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "the status code should be equal to 200", sentence.NewTerm(statusCode, shouldBeOK).String())
	})

	t.Run("is met when the answer satisfies its clause", func(t *testing.T) {
		t.Parallel()

		ok := fixture.NewFakeQuestion("status code", 200)

		require.NoError(t, sentence.NewTerm(ok, shouldBeOK).MetBy(adam))
	})

	t.Run("explains why it is not met", func(t *testing.T) {
		t.Parallel()

		err := sentence.NewTerm(statusCode, shouldBeOK).MetBy(adam)

		var unmet *sentence.UnmetError
		require.ErrorAs(t, err, &unmet)
		assert.EqualError(t, err, "the status code should be equal to 200: got 404")
	})

	t.Run("fails with the error of a question that cannot be answered", func(t *testing.T) {
		t.Parallel()

		errNoResponse := errors.New("no response received")
		missing := fixture.NewFailingFakeQuestion[int]("status code", errNoResponse)

		_, err := sentence.NewTerm(missing, shouldBeOK).JudgedBy(adam)

		require.ErrorIs(t, err, errNoResponse)
	})

	t.Run("remembers the mistake of a nil question", func(t *testing.T) {
		t.Parallel()

		term := sentence.NewTerm(nil, shouldBeOK)

		require.ErrorIs(t, term.Mistake(), sentence.ErrInvalidTerm)
		assert.Equal(t, "the <nil> should be equal to 200", term.String())
		require.ErrorIs(t, term.MetBy(adam), sentence.ErrInvalidTerm)
	})

	t.Run("remembers the mistake of its clause", func(t *testing.T) {
		t.Parallel()

		term := sentence.NewTerm(statusCode, sentence.NewClause[int](sentence.Should, nil))

		require.ErrorIs(t, term.Mistake(), sentence.ErrInvalidTerm)
		require.ErrorIs(t, term.MetBy(adam), sentence.ErrInvalidTerm)
	})

	t.Run("has no mistake when it was built correctly", func(t *testing.T) {
		t.Parallel()

		assert.NoError(t, sentence.NewTerm(statusCode, shouldBeOK).Mistake())
	})
}
