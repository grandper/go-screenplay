package is_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/resolution/is"
	testdata2 "github.com/grandper/go-screenplay/resolution/testdata"
)

func TestIsNotResolution(t *testing.T) {
	t.Parallel()

	const anObject = "some object"
	failingToMatch := testdata2.NewFailingResolution("failing to match", errors.New("a fake error occurred"))
	matchingAnyValue := testdata2.NewFakeResolution("matching any value", true)
	notMatchingAnyValue := testdata2.NewFakeResolution("not matching any value", false)
	isNotFailingMatch := is.Not(failingToMatch)
	isNotMatchingAnyValue := is.Not(matchingAnyValue)
	isNotNotMatchingAnyValue := is.Not(notMatchingAnyValue)

	t.Run("should match when the values of the underlying resolutation fails", func(t *testing.T) {
		t.Parallel()
		testdata2.AssertMatch(t, isNotNotMatchingAnyValue, anObject)
	})

	t.Run("fails when the values of the underlying resolutation succeeds", func(t *testing.T) {
		t.Parallel()
		testdata2.AssertNoMatch(t, isNotMatchingAnyValue, anObject)
	})

	t.Run("returns an error when the underlying resolution fails", func(t *testing.T) {
		t.Parallel()
		testdata2.AssertMatcherFails(t, isNotFailingMatch, anObject)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "not failing to match", isNotFailingMatch.String())
		assert.Equal(t, "not matching any value", isNotMatchingAnyValue.String())
		assert.Equal(t, "not not matching any value", isNotNotMatchingAnyValue.String())
	})
}
