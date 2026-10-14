package start_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/v2/resolution/fixture"
	"github.com/grandper/go-screenplay/v2/resolution/start"
)

func TestWith(t *testing.T) {
	t.Parallel()

	resolution := start.With("hello")

	fixture.AssertMatch(t, resolution, "hello world")
	fixture.AssertNoMatch(t, resolution, "world hello")
	assert.Equal(t, `start with "hello"`, resolution.String())
}
