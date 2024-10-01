package testdata

import (
	"testing"

	"github.com/grandper/go-screenplay/screenplay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func AssertMatch(t *testing.T, matcher screenplay.Resolution, obj any) {
	t.Helper()

	result, err := matcher.Resolve()(obj)
	require.NoError(t, err)
	assert.True(t, result)
}

func AssertNoMatch(t *testing.T, matcher screenplay.Resolution, obj any) {
	t.Helper()

	result, err := matcher.Resolve()(obj)
	require.NoError(t, err)
	assert.False(t, result)
}

func AssertMatcherFails(t *testing.T, matcher screenplay.Resolution, obj any) {
	t.Helper()

	result, err := matcher.Resolve()(obj)
	require.Error(t, err)
	assert.False(t, result)
}
