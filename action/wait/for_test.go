package wait_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action/wait"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestFor(t *testing.T) {
	adam := screenplay.ActorNamed("Adam")

	t.Run("waits for a given time", func(t *testing.T) {
		start := time.Now()
		require.NoError(t, adam.AttemptsTo(wait.For(20).Milliseconds().Because("it is a test")))
		assert.WithinDuration(t, start.Add(20*time.Millisecond), time.Now(), 5*time.Millisecond)
	})

	t.Run("fails if no reason is provided for the wait", func(t *testing.T) {
		err := adam.AttemptsTo(wait.For(10).Seconds())
		require.ErrorIs(t, err, wait.ErrReasonNotProvided)
		assert.EqualError(t, err, "failed to wait: cannot wait without a reason: you must call the .Because() method")
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		action1 := wait.For(1).Second().Because("it is a test")
		assert.Equal(t, "wait for 1 second because it is a test", action1.String())

		action2 := wait.For(1).Seconds().Because("it is a test")
		assert.Equal(t, "wait for 1 second because it is a test", action2.String())

		action3 := wait.For(20).Milliseconds().Because("it is a test")
		assert.Equal(t, "wait for 20 milliseconds because it is a test", action3.String())

		action4 := wait.For(20).Millisecond().Because("it is a test")
		assert.Equal(t, "wait for 20 milliseconds because it is a test", action4.String())
	})
}
