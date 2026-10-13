package term_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/fixture"
	"github.com/grandper/go-screenplay/v2/resolution"
	"github.com/grandper/go-screenplay/v2/resolution/be"
	"github.com/grandper/go-screenplay/v2/resolution/contain"
	"github.com/grandper/go-screenplay/v2/resolution/length"
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/grandper/go-screenplay/v2/term"
)

func TestTerm(t *testing.T) {

	t.Run("that", func(t *testing.T) {
		t.Parallel()

		adam := screenplay.ActorNamed("Adam")
		statusCode := fixture.NewFakeQuestion("status code", 200)
		errNoResponse := errors.New("no response received")
		missingStatusCode := fixture.NewFailingFakeQuestion[int]("status code", errNoResponse)

		t.Run("is met when the answer satisfies the clause", func(t *testing.T) {
			t.Parallel()

			met := term.That(statusCode).Should(be.EqualTo(200))

			require.NoError(t, met.MetBy(adam))

			verdict, err := met.JudgedBy(adam)
			require.NoError(t, err)
			assert.True(t, verdict.Met)
		})

		t.Run("is not met when the answer does not satisfy the clause", func(t *testing.T) {
			t.Parallel()

			unmet := term.That(statusCode).Should(be.EqualTo(404))

			err := unmet.MetBy(adam)
			require.EqualError(t, err, "the status code should be equal to 404: got 200")

			verdict, err := unmet.JudgedBy(adam)
			require.NoError(t, err)
			assert.False(t, verdict.Met)
		})

		t.Run("fails with the error of the question when it cannot be answered", func(t *testing.T) {
			t.Parallel()

			unanswered := term.That(missingStatusCode).ShouldNot(be.EqualTo(404))

			err := unanswered.MetBy(adam)
			require.ErrorIs(t, err, errNoResponse)

			_, err = unanswered.JudgedBy(adam)
			require.ErrorIs(t, err, errNoResponse)
		})

		t.Run("describes itself", func(t *testing.T) {
			t.Parallel()
			assert.Equal(t,
				"the status code should be equal to 200",
				term.That(statusCode).Should(be.EqualTo(200)).String(),
			)
		})

		t.Run("is invalid when the question is nil", func(t *testing.T) {
			t.Parallel()

			invalid := term.That[int](nil).Should(be.EqualTo(200))

			require.ErrorIs(t, invalid.MetBy(adam), term.ErrInvalidTerm)
			require.EqualError(t, invalid.MetBy(adam), "invalid term: the question is nil")

			_, err := invalid.JudgedBy(adam)
			require.ErrorIs(t, err, term.ErrInvalidTerm)
		})
	})

	t.Run("verbs", func(t *testing.T) {
		t.Parallel()

		adam := screenplay.ActorNamed("Adam")
		statusCode := fixture.NewFakeQuestion("status code", 200)
		logs := fixture.NewFakeQuestion("logs", "order placed, payment captured")
		customers := fixture.NewFakeQuestion("customers", []string{"Adam", "Eve"})
		aSuccess := resolution.FromFunc("a success", func(code int) error {
			if code < 200 || code >= 300 {
				return fmt.Errorf("%d is not a success", code)
			}

			return nil
		})

		for _, testCase := range []struct {
			name  string
			term  *term.Term
			reads string
			met   bool
		}{
			{
				"Should, met",
				term.That(statusCode).Should(be.EqualTo(200)),
				"the status code should be equal to 200", true,
			},
			{
				"Should, unmet",
				term.That(statusCode).Should(be.EqualTo(404)),
				"the status code should be equal to 404", false,
			},
			{
				"Should with several resolutions, met",
				term.That(logs).Should(contain.Substring("order placed"), contain.Substring("payment captured")),
				`the logs should contain substring "order placed", and contain substring "payment captured"`, true,
			},
			{
				"Should with several resolutions, unmet",
				term.That(logs).Should(contain.Substring("order placed"), contain.Substring("refunded")),
				`the logs should contain substring "order placed", and contain substring "refunded"`, false,
			},
			{
				"ShouldBe, met",
				term.That(statusCode).ShouldBe(aSuccess),
				"the status code should be a success", true,
			},
			{
				"ShouldHave, met",
				term.That(customers).ShouldHave(length.Exactly[[]string](2)),
				"the customers should have length exactly 2", true,
			},
			{
				"ShouldHave, unmet",
				term.That(customers).ShouldHave(length.Exactly[[]string](3)),
				"the customers should have length exactly 3", false,
			},
			{
				"ShouldNot, met",
				term.That(logs).ShouldNot(contain.Substring("panic")),
				`the logs should not contain substring "panic"`, true,
			},
			{
				"ShouldNot, unmet",
				term.That(logs).ShouldNot(contain.Substring("order placed")),
				`the logs should not contain substring "order placed"`, false,
			},
			{
				"ShouldNot with several resolutions, met",
				term.That(logs).ShouldNot(contain.Substring("panic"), contain.Substring("fatal")),
				`the logs should not contain substring "panic", nor contain substring "fatal"`, true,
			},
			{
				"ShouldNot with several resolutions, unmet when one holds",
				term.That(logs).ShouldNot(contain.Substring("panic"), contain.Substring("payment captured")),
				`the logs should not contain substring "panic", nor contain substring "payment captured"`, false,
			},
			{
				"ShouldNotBe, unmet",
				term.That(statusCode).ShouldNotBe(aSuccess),
				"the status code should not be a success", false,
			},
			{
				"ShouldNotHave, met",
				term.That(customers).ShouldNotHave(length.Exactly[[]string](3)),
				"the customers should not have length exactly 3", true,
			},
			{
				"ShouldAnyOf, met",
				term.That(statusCode).ShouldAnyOf(be.EqualTo(200), be.EqualTo(201)),
				"the status code should be any of: equal to 200, equal to 201", true,
			},
			{
				"ShouldAnyOf, unmet",
				term.That(statusCode).ShouldAnyOf(be.EqualTo(201), be.EqualTo(204)),
				"the status code should be any of: equal to 201, equal to 204", false,
			},
		} {
			t.Run(testCase.name, func(t *testing.T) {
				t.Parallel()

				assert.Equal(t, testCase.reads, testCase.term.String())

				verdict, err := testCase.term.JudgedBy(adam)
				require.NoError(t, err)
				assert.Equal(t, testCase.met, verdict.Met)
			})
		}

		t.Run("report exactly the part of the sentence that broke the rule", func(t *testing.T) {
			t.Parallel()

			err := term.That(statusCode).Should(be.AtLeast(200), be.LessThan(200)).MetBy(adam)
			require.EqualError(t, err, "the status code should be less than 200: got 200")

			err = term.That(statusCode).ShouldNot(be.EqualTo(404), be.EqualTo(200)).MetBy(adam)
			require.EqualError(t, err, "the status code should not be equal to 200: got 200")

			err = term.That(statusCode).ShouldAnyOf(be.EqualTo(201), be.EqualTo(204)).MetBy(adam)
			require.EqualError(t, err, "the status code (got 200) should be any of: equal to 201, equal to 204")

			err = term.That(customers).ShouldHave(length.Exactly[[]string](3)).MetBy(adam)
			require.EqualError(t, err, "the customers should have length exactly 3: got [Adam Eve]")

			err = term.That(statusCode).ShouldNotBe(aSuccess).MetBy(adam)
			require.EqualError(t, err, "the status code should not be a success: got 200")
		})

		t.Run("build an invalid term when they are given no resolution or a nil one", func(t *testing.T) {
			t.Parallel()

			require.EqualError(t, term.That(statusCode).Should().MetBy(adam),
				"invalid term: Should needs at least one resolution")
			require.EqualError(t, term.That(statusCode).Should(be.EqualTo(200), nil).MetBy(adam),
				"invalid term: the resolution at position 1 given to Should is nil")

			for _, invalid := range []*term.Term{
				term.That(statusCode).Should(),
				term.That(statusCode).ShouldBe(),
				term.That(statusCode).ShouldHave(),
				term.That(statusCode).ShouldNot(),
				term.That(statusCode).ShouldNotBe(),
				term.That(statusCode).ShouldNotHave(),
				term.That(statusCode).ShouldAnyOf(),
				term.That(statusCode).ShouldNot(nil),
			} {
				require.ErrorIs(t, invalid.MetBy(adam), term.ErrInvalidTerm)

				_, err := invalid.JudgedBy(adam)
				require.ErrorIs(t, err, term.ErrInvalidTerm)
			}
		})

		t.Run("fail on a resolution that was built wrongly, whatever the verb", func(t *testing.T) {
			t.Parallel()

			broken := be.Between(300, 200)

			for _, invalid := range []*term.Term{
				term.That(statusCode).Should(broken),
				term.That(statusCode).ShouldBe(broken),
				term.That(statusCode).ShouldHave(broken),
				term.That(statusCode).ShouldNot(broken),
				term.That(statusCode).ShouldNotBe(broken),
				term.That(statusCode).ShouldNotHave(broken),
				term.That(statusCode).ShouldAnyOf(broken),
				term.That(statusCode).ShouldAnyOf(be.EqualTo(200), broken),
				term.That(statusCode).ShouldNot(be.EqualTo(404), broken),
				term.That(statusCode).Should(be.EqualTo(404), broken),
			} {
				require.ErrorIs(t, invalid.MetBy(adam), screenplay.ErrInvalidResolution)
				require.NotErrorIs(t, invalid.MetBy(adam), term.ErrInvalidTerm)

				_, err := invalid.JudgedBy(adam)
				require.ErrorIs(t, err, screenplay.ErrInvalidResolution)
			}
		})

		t.Run("fail on a custom resolution reporting that it is invalid", func(t *testing.T) {
			t.Parallel()

			custom := resolution.FromFunc("be valid", func(int) error {
				return fmt.Errorf("%w: missing configuration", screenplay.ErrInvalidResolution)
			})

			require.ErrorIs(t, term.That(statusCode).ShouldNot(custom).MetBy(adam), screenplay.ErrInvalidResolution)
		})
	})
}
