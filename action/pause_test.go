package action_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/action"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestPauseAction(t *testing.T) {
	adam := screenplay.ActorNamed("Adam")

	t.Run("should PauseActionBuilder for a given time", func(t *testing.T) {
		start := time.Now()
		require.NoError(t, adam.AttemptsTo(action.PauseFor(20).Milliseconds().Because("it is a test")))
		assert.WithinDuration(t, start.Add(20*time.Millisecond), time.Now(), 5*time.Millisecond)
	})

	t.Run("fails if no reason is provided for the PauseActionBuilder", func(t *testing.T) {
		err := adam.AttemptsTo(action.PauseFor(10).Seconds())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "without a reason")
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		action1 := action.PauseFor(1).Second().Because("it is a test")
		assert.Equal(t, "PauseActionBuilder for 1 second because it is a test", action1.String())

		action2 := action.PauseFor(1).Seconds().Because("it is a test")
		assert.Equal(t, "PauseActionBuilder for 1 second because it is a test", action2.String())

		action3 := action.PauseFor(20).Milliseconds().Because("it is a test")
		assert.Equal(t, "PauseActionBuilder for 20 milliseconds because it is a test", action3.String())

		action4 := action.PauseFor(20).Millisecond().Because("it is a test")
		assert.Equal(t, "PauseActionBuilder for 20 milliseconds because it is a test", action4.String())
	})
}
