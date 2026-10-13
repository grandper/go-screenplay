package start_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/v2/resolution/start"
	"github.com/grandper/go-screenplay/v2/resolution/testdata"
)

func TestWith(t *testing.T) {
	t.Parallel()

	resolution := start.With("hello")

	testdata.AssertMatch(t, resolution, "hello world")
	testdata.AssertNoMatch(t, resolution, "world hello")
	assert.Equal(t, `start with "hello"`, resolution.String())
}
