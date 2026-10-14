package end_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/v2/resolution/end"
	"github.com/grandper/go-screenplay/v2/resolution/fixture"
)

func TestWith(t *testing.T) {
	t.Parallel()

	resolution := end.With("world")

	fixture.AssertMatch(t, resolution, "hello world")
	fixture.AssertNoMatch(t, resolution, "world hello")
	assert.Equal(t, `end with "world"`, resolution.String())
}
