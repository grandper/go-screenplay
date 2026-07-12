package see_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/action/see"
	"github.com/grandper/go-screenplay/fixture"
	"github.com/grandper/go-screenplay/resolution/equal"
	"github.com/grandper/go-screenplay/resolution/testdata"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestSeeTheAction(t *testing.T) {
	t.Parallel()

	adam := screenplay.ActorNamed("Adam")
	formField := fixture.NewFakeQuestion("form field", "hello world")
	missingFormField := fixture.NewFailingFakeQuestion("form field", errors.New("failed to get the field content"))
	isEqualButFails := testdata.NewFailingResolution(
		"equal to hello world",
		errors.New("failed to match the content of the field"),
	)

	t.Run("should see something", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, adam.AttemptsTo(see.The(formField).Is(equal.To("hello world"))))
	})

	t.Run("fails when there is nothing to see", func(t *testing.T) {
		t.Parallel()
		require.Error(t, adam.AttemptsTo(see.The(formField).Is(equal.To("hello everybody"))))
	})

	t.Run("fails when the actor fails to answer the question", func(t *testing.T) {
		t.Parallel()
		require.Error(t, adam.AttemptsTo(see.The(missingFormField).Is(equal.To("hello everybody"))))
	})

	t.Run("fails when the resolution fails", func(t *testing.T) {
		t.Parallel()
		require.Error(t, adam.AttemptsTo(see.The(formField).Is(isEqualButFails)))
	})

	t.Run("sees the negation when the resolution does not match", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, adam.AttemptsTo(see.The(formField).IsNot(equal.To("hello everybody"))))
	})

	t.Run("fails the negation when the resolution matches", func(t *testing.T) {
		t.Parallel()
		require.Error(t, adam.AttemptsTo(see.The(formField).IsNot(equal.To("hello world"))))
	})

	t.Run("propagates the resolution error through a negation", func(t *testing.T) {
		t.Parallel()
		require.Error(t, adam.AttemptsTo(see.The(formField).IsNot(isEqualButFails)))
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		t.Parallel()

		action := see.The(formField).Is(equal.To("hello world"))
		assert.Equal(t, "see if the form field is equal to hello world", action.String())

		negated := see.The(formField).IsNot(equal.To("hello world"))
		assert.Equal(t, "see if the form field is not equal to hello world", negated.String())
	})

	t.Run("supports every verb and its negation", func(t *testing.T) {
		t.Parallel()

		resolution := equal.To("hello world")
		for _, testCase := range []struct {
			action *see.TheAction
			reads  string
			passes bool
		}{
			{see.The(formField).Is(resolution), "see if the form field is equal to hello world", true},
			{see.The(formField).IsNot(resolution), "see if the form field is not equal to hello world", false},
			{see.The(formField).Are(resolution), "see if the form field are equal to hello world", true},
			{see.The(formField).AreNot(resolution), "see if the form field are not equal to hello world", false},
			{see.The(formField).Does(resolution), "see if the form field does equal to hello world", true},
			{see.The(formField).DoesNot(resolution), "see if the form field does not equal to hello world", false},
			{see.The(formField).Do(resolution), "see if the form field do equal to hello world", true},
			{see.The(formField).DoNot(resolution), "see if the form field do not equal to hello world", false},
			{see.The(formField).Has(resolution), "see if the form field has equal to hello world", true},
			{see.The(formField).HasNot(resolution), "see if the form field has not equal to hello world", false},
			{see.The(formField).Have(resolution), "see if the form field have equal to hello world", true},
			{see.The(formField).HaveNot(resolution), "see if the form field have not equal to hello world", false},
			{see.The(formField).Had(resolution), "see if the form field had equal to hello world", true},
			{see.The(formField).HadNot(resolution), "see if the form field had not equal to hello world", false},
			{see.The(formField).Was(resolution), "see if the form field was equal to hello world", true},
			{see.The(formField).WasNot(resolution), "see if the form field was not equal to hello world", false},
			{see.The(formField).Were(resolution), "see if the form field were equal to hello world", true},
			{see.The(formField).WereNot(resolution), "see if the form field were not equal to hello world", false},
		} {
			assert.Equal(t, testCase.reads, testCase.action.String())
			if testCase.passes {
				require.NoError(t, adam.AttemptsTo(testCase.action))
			} else {
				require.Error(t, adam.AttemptsTo(testCase.action))
			}
		}
	})
}
