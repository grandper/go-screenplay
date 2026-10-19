package see_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action/see"
	"github.com/grandper/go-screenplay/v2/action/see/term"
	"github.com/grandper/go-screenplay/v2/resolution/be"
	"github.com/grandper/go-screenplay/v2/resolution/contain"
	resolutionfixture "github.com/grandper/go-screenplay/v2/resolution/fixture"
	"github.com/grandper/go-screenplay/v2/resolution/length"
	"github.com/grandper/go-screenplay/v2/resolution/match"
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/grandper/go-screenplay/v2/screenplay/fixture"
)

func TestSeeTheAction(t *testing.T) {
	t.Parallel()

	adam := screenplay.ActorNamed("Adam")
	formField := fixture.NewFakeQuestion("form field", "hello world")
	errNoField := errors.New("failed to get the field content")
	missingFormField := fixture.NewFailingFakeQuestion[string]("form field", errNoField)

	t.Run("should see something", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, adam.AttemptsTo(see.The(formField).Should(be.EqualTo("hello world"))))
	})

	t.Run("fails when there is nothing to see", func(t *testing.T) {
		t.Parallel()

		err := adam.AttemptsTo(see.The(formField).Should(be.EqualTo("hello everybody")))

		require.EqualError(t, err,
			`Adam failed to see that the form field should be equal to "hello everybody": got "hello world"`)
	})

	t.Run("fails with the error of the question when the actor fails to answer it", func(t *testing.T) {
		t.Parallel()
		require.ErrorIs(t,
			adam.AttemptsTo(see.The(missingFormField).Should(be.EqualTo("hello everybody"))),
			errNoField,
		)
	})

	t.Run("fails a negation when the actor fails to answer the question", func(t *testing.T) {
		t.Parallel()
		require.ErrorIs(t,
			adam.AttemptsTo(see.The(missingFormField).ShouldNot(be.EqualTo("hello everybody"))),
			errNoField,
		)
	})

	t.Run("sees the negation when the resolution is not satisfied", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, adam.AttemptsTo(see.The(formField).ShouldNot(be.EqualTo("hello everybody"))))
	})

	t.Run("fails the negation when the resolution is satisfied", func(t *testing.T) {
		t.Parallel()
		require.Error(t, adam.AttemptsTo(see.The(formField).ShouldNot(be.EqualTo("hello world"))))
	})

	t.Run("sees several resolutions about the same question", func(t *testing.T) {
		t.Parallel()

		require.NoError(t, adam.AttemptsTo(see.The(formField).Should(
			contain.Substring("hello"),
			contain.Substring("world"),
		)))
		require.Error(t, adam.AttemptsTo(see.The(formField).Should(
			contain.Substring("hello"),
			contain.Substring("everybody"),
		)))
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		t.Parallel()

		action := see.The(formField).Should(be.EqualTo("hello world"))
		assert.Equal(t, `see that the form field should be equal to "hello world"`, action.String())

		negated := see.The(formField).ShouldNot(be.EqualTo("hello world"))
		assert.Equal(t, `see that the form field should not be equal to "hello world"`, negated.String())
	})

	t.Run("supports every verb", func(t *testing.T) {
		t.Parallel()

		anything := resolutionfixture.NewFakeResolution[string]("anything", true)

		for _, testCase := range []struct {
			action screenplay.Performable
			reads  string
			passes bool
		}{
			{see.The(formField).Should(contain.Substring("hello")),
				`see that the form field should contain substring "hello"`, true},
			{see.The(formField).ShouldNot(contain.Substring("hello")),
				`see that the form field should not contain substring "hello"`, false},
			{see.The(formField).ShouldBe(anything),
				"see that the form field should be anything", true},
			{see.The(formField).ShouldNotBe(anything),
				"see that the form field should not be anything", false},
			{see.The(formField).ShouldHave(length.Exactly[string](11)),
				"see that the form field should have length exactly 11", true},
			{see.The(formField).ShouldNotHave(length.Exactly[string](11)),
				"see that the form field should not have length exactly 11", false},
			{see.The(formField).ShouldAnyOf(be.EqualTo("hello"), be.EqualTo("hello world")),
				`see that the form field should be any of: equal to "hello", equal to "hello world"`, true},
			{see.The(formField).ShouldAnyOf(be.EqualTo("hello"), be.EqualTo("world")),
				`see that the form field should be any of: equal to "hello", equal to "world"`, false},
		} {
			assert.Equal(t, testCase.reads, testCase.action.String())

			if testCase.passes {
				require.NoError(t, adam.AttemptsTo(testCase.action))
			} else {
				require.Error(t, adam.AttemptsTo(testCase.action))
			}
		}
	})

	t.Run("fails with an invalid term when the sentence is built wrongly", func(t *testing.T) {
		t.Parallel()

		require.ErrorIs(t, adam.AttemptsTo(see.The[string](nil).Should(be.EqualTo("hello"))), term.ErrInvalidTerm)
		require.ErrorIs(t, adam.AttemptsTo(see.The(formField).Should()), term.ErrInvalidTerm)
		require.ErrorIs(t, adam.AttemptsTo(see.The(formField).ShouldNot()), term.ErrInvalidTerm)
		require.ErrorIs(t, adam.AttemptsTo(see.The(formField).ShouldNot(nil)), term.ErrInvalidTerm)
		require.EqualError(t, adam.AttemptsTo(see.The(formField).Should()),
			"invalid term: Should needs at least one resolution")
	})

	t.Run("fails with an invalid resolution when a resolution is built wrongly, whatever the verb", func(t *testing.T) {
		t.Parallel()

		broken := match.PatternString("hello(")

		require.ErrorIs(t, adam.AttemptsTo(see.The(formField).Should(broken)), screenplay.ErrInvalidResolution)
		require.ErrorIs(t, adam.AttemptsTo(see.The(formField).ShouldNot(broken)), screenplay.ErrInvalidResolution)
		require.ErrorContains(t, adam.AttemptsTo(see.The(formField).ShouldNot(broken)),
			"invalid resolution \"match pattern `hello(`\": error parsing regexp")
	})
}
