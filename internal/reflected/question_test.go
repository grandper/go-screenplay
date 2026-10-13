package reflected_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/internal/reflected"
)

type actor struct{ name string }

type statusCode struct{ err error }

func (q statusCode) AnsweredBy(_ *actor) (int, error) {
	if q.err != nil {
		return 0, q.err
	}

	return 200, nil
}

func (q statusCode) String() string { return "status code" }

type noDescription struct{}

func (noDescription) AnsweredBy(_ *actor) (int, error) { return 0, nil }

type wrongActor struct{}

func (wrongActor) AnsweredBy(_ string) (int, error) { return 0, nil }
func (wrongActor) String() string                   { return "wrong actor" }

type noError struct{}

func (noError) AnsweredBy(_ *actor) int { return 0 }
func (noError) String() string          { return "no error" }

func TestQuestionOf(t *testing.T) {
	t.Parallel()

	actorType := reflect.TypeFor[*actor]()

	t.Run("recognizes a question whatever its answer type", func(t *testing.T) {
		t.Parallel()

		question, isAQuestion := reflected.QuestionOf(statusCode{}, actorType)
		require.True(t, isAQuestion)
		assert.Equal(t, "status code", question.String())

		answer, err := question.AnsweredBy(&actor{name: "Adam"})
		require.NoError(t, err)
		assert.Equal(t, 200, answer)
	})

	t.Run("returns the error of a question that fails to be answered", func(t *testing.T) {
		t.Parallel()

		errNoResponse := errors.New("no response")
		question, isAQuestion := reflected.QuestionOf(statusCode{err: errNoResponse}, actorType)
		require.True(t, isAQuestion)

		answer, err := question.AnsweredBy(&actor{name: "Adam"})
		require.ErrorIs(t, err, errNoResponse)
		assert.Nil(t, answer)
	})

	t.Run("does not recognize a value that is not a question", func(t *testing.T) {
		t.Parallel()

		for _, value := range []any{nil, 42, "question", noDescription{}, wrongActor{}, noError{}} {
			_, isAQuestion := reflected.QuestionOf(value, actorType)
			assert.False(t, isAQuestion, "%T", value)
		}
	})
}
