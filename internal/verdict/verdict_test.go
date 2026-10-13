package verdict_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/internal/verdict"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestVerdict(t *testing.T) {

	t.Run("describe", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, `"hello"`, verdict.Describe("hello"))
		assert.Equal(t, `"hello"`, verdict.Describe([]byte("hello")))
		assert.Equal(t, "42", verdict.Describe(42))
		assert.Equal(t, "[a b]", verdict.Describe([]string{"a", "b"}))
		assert.Equal(t, "<nil>", verdict.Describe(nil))
	})

	t.Run("of", func(t *testing.T) {
		t.Parallel()

		require.NoError(t, verdict.Of(true, 42))
		require.EqualError(t, verdict.Of(false, 42), "got 42")
		require.EqualError(t, verdict.Of(false, "hello"), `got "hello"`)
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		errBadPattern := errors.New("missing closing )")
		resolution := verdict.Invalid[string]("match pattern `A-(`", errBadPattern)

		assert.Equal(t, "match pattern `A-(`", resolution.String())

		err := resolution.Inspect("A-42")
		require.ErrorIs(t, err, screenplay.ErrInvalidResolution)
		require.ErrorIs(t, err, errBadPattern)
		require.EqualError(t, err, "invalid resolution \"match pattern `A-(`\": missing closing )")
	})
}
