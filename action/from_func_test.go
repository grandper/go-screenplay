package action_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/action"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestFromFunc(t *testing.T) {
	adam := screenplay.ActorNamed("Adam")

	doNothing := action.FromFunc("do nothing", func(_ *screenplay.Actor) error {
		return nil
	})
	assert.Implements(t, (*screenplay.Performable)(nil), doNothing)
	assert.Equal(t, "do nothing", doNothing.String())
	require.NoError(t, adam.AttemptsTo(doNothing))

	failToDoSomething := action.FromFunc("fail to do something", func(_ *screenplay.Actor) error {
		return assert.AnError
	})
	assert.Implements(t, (*screenplay.Performable)(nil), failToDoSomething)
	assert.Equal(t, "fail to do something", failToDoSomething.String())
	require.Error(t, adam.AttemptsTo(failToDoSomething))
}
