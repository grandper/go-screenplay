package question_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/question"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestFromFunc(t *testing.T) {
	adam := screenplay.ActorNamed("Adam")

	theAnswer := question.FromFunc("the answer", func(_ *screenplay.Actor) (any, error) {
		return 42, nil
	})
	assert.Implements(t, (*screenplay.Question)(nil), theAnswer)
	assert.Equal(t, "the answer", theAnswer.String())
	answer, err := theAnswer.AnsweredBy(adam)
	require.NoError(t, err)
	assert.Equal(t, 42, answer)

	failToAnswer := question.FromFunc("fail to answer", func(_ *screenplay.Actor) (any, error) {
		return nil, assert.AnError
	})
	assert.Implements(t, (*screenplay.Question)(nil), failToAnswer)
	assert.Equal(t, "fail to answer", failToAnswer.String())
	_, err = failToAnswer.AnsweredBy(adam)
	require.ErrorIs(t, err, assert.AnError)
}
