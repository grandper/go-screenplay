package log_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action/log"
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/grandper/go-screenplay/v2/screenplay/fixture"
)

func TestThe(t *testing.T) {
	adam := screenplay.ActorNamed("Adam")

	t.Run("answers the question", func(t *testing.T) {
		require.NoError(t, adam.AttemptsTo(log.The(fixture.NewFakeQuestion("status code", 200))))
	})

	t.Run("returns the error of the question", func(t *testing.T) {
		missingStatusCode := fixture.NewFailingFakeQuestion[int]("status code", errors.New("no response"))

		require.EqualError(t, adam.AttemptsTo(log.The(missingStatusCode)), "no response")
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "log the status code", log.The(fixture.NewFakeQuestion("status code", 200)).String())
	})
}
