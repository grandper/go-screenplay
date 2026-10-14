package fixture

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/screenplay"
)

// AssertMatch asserts that the value satisfies the resolution.
func AssertMatch[T any](t *testing.T, resolution screenplay.Resolution[T], value T) {
	t.Helper()

	require.NoError(t, resolution.Inspect(value))
}

// AssertNoMatch asserts that the value does not satisfy the resolution, and
// that the resolution was built correctly.
func AssertNoMatch[T any](t *testing.T, resolution screenplay.Resolution[T], value T) {
	t.Helper()

	err := resolution.Inspect(value)
	require.Error(t, err)
	require.NotErrorIs(t, err, screenplay.ErrInvalidResolution)
}

// AssertInvalid asserts that the resolution was built wrongly: inspecting the
// value fails with an error wrapping screenplay.ErrInvalidResolution.
func AssertInvalid[T any](t *testing.T, resolution screenplay.Resolution[T], value T) {
	t.Helper()

	require.ErrorIs(t, resolution.Inspect(value), screenplay.ErrInvalidResolution)
}
