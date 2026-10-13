package resolution_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/resolution"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestFromFunc(t *testing.T) {
	t.Parallel()

	t.Run("creates a resolution with a description and an inspect function", func(t *testing.T) {
		t.Parallel()

		errNotTheAnswer := errors.New("not the answer")
		theAnswer := resolution.FromFunc("be the answer", func(value int) error {
			if value != 42 {
				return errNotTheAnswer
			}

			return nil
		})

		assert.Implements(t, (*screenplay.Resolution[int])(nil), theAnswer)
		assert.Equal(t, "be the answer", theAnswer.String())
		require.NoError(t, theAnswer.Inspect(42))
		require.ErrorIs(t, theAnswer.Inspect(0), errNotTheAnswer)
	})

	t.Run("is an invalid resolution when the inspect function is nil", func(t *testing.T) {
		t.Parallel()

		broken := resolution.FromFunc[int]("be anything", nil)

		assert.Equal(t, "be anything", broken.String())
		require.ErrorIs(t, broken.Inspect(42), screenplay.ErrInvalidResolution)
	})
}
