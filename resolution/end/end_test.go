package end_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/v2/resolution/end"
	"github.com/grandper/go-screenplay/v2/resolution/testdata"
)

func TestWith(t *testing.T) {
	t.Parallel()

	resolution := end.With("world")

	testdata.AssertMatch(t, resolution, "hello world")
	testdata.AssertNoMatch(t, resolution, "world hello")
	assert.Equal(t, `end with "world"`, resolution.String())
}
