package see_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action/see"
	"github.com/grandper/go-screenplay/v2/fixture"
	"github.com/grandper/go-screenplay/v2/resolution/be"
	"github.com/grandper/go-screenplay/v2/resolution/contain"
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/grandper/go-screenplay/v2/term"
)

func TestThese(t *testing.T) {

	t.Run("see all of these met", func(t *testing.T) {
		t.Parallel()

		adam := screenplay.ActorNamed("Adam")
		statusCode := fixture.NewFakeQuestion("status code", 200)
		logs := fixture.NewFakeQuestion("logs", "order placed")
		customers := fixture.NewFakeQuestion("customers", []string{"Adam", "Eve"})
		errNoLogs := errors.New("failed to read the logs")
		missingLogs := fixture.NewFailingFakeQuestion[string]("logs", errNoLogs)

		t.Run("succeeds when every term is met, whatever the answer types", func(t *testing.T) {
			t.Parallel()
			require.NoError(t, adam.AttemptsTo(see.AllOfTheseMet(
				term.That(statusCode).Should(be.EqualTo(200)),
				term.That(logs).Should(contain.Substring("order placed")),
				term.That(customers).Should(contain.Element("Adam")),
			)))
		})

		t.Run("fails when one term is not met", func(t *testing.T) {
			t.Parallel()

			err := adam.AttemptsTo(see.AllOfTheseMet(
				term.That(statusCode).Should(be.EqualTo(200)),
				term.That(logs).ShouldNot(contain.Substring("order placed")),
			))

			require.EqualError(t, err,
				`Adam failed to see that the logs should not contain substring "order placed": got "order placed"`)
		})

		t.Run("stops at the first term that is not met", func(t *testing.T) {
			t.Parallel()

			err := adam.AttemptsTo(see.AllOfTheseMet(
				term.That(statusCode).Should(be.EqualTo(404)),
				term.That(missingLogs).Should(contain.Substring("order placed")),
			))

			require.Error(t, err)
			require.NotErrorIs(t, err, errNoLogs)
		})

		t.Run("fails when the actor fails to answer one of the questions", func(t *testing.T) {
			t.Parallel()
			require.ErrorIs(t, adam.AttemptsTo(see.AllOfTheseMet(
				term.That(statusCode).Should(be.EqualTo(200)),
				term.That(missingLogs).Should(contain.Substring("order placed")),
			)), errNoLogs)
		})

		t.Run("implements the stringer interface", func(t *testing.T) {
			t.Parallel()
			assert.Equal(t,
				"see that all of these terms are met",
				see.AllOfTheseMet(
					term.That(statusCode).Should(be.EqualTo(200)),
					term.That(logs).Should(contain.Substring("order placed")),
				).String(),
			)
		})

		t.Run("fails with an invalid term when a term is nil or built wrongly", func(t *testing.T) {
			t.Parallel()

			err := adam.AttemptsTo(see.AllOfTheseMet(term.That(statusCode).Should(be.EqualTo(200)), nil))
			require.ErrorIs(t, err, term.ErrInvalidTerm)
			require.EqualError(t, err, "invalid term: the term at position 1 given to see.AllOfTheseMet is nil")

			require.ErrorIs(t, adam.AttemptsTo(see.AllOfTheseMet(term.That(statusCode).Should())), term.ErrInvalidTerm)
		})

		t.Run("fails with an invalid resolution when a resolution is built wrongly", func(t *testing.T) {
			t.Parallel()
			require.ErrorIs(t,
				adam.AttemptsTo(see.AllOfTheseMet(term.That(statusCode).ShouldNot(be.Between(300, 200)))),
				screenplay.ErrInvalidResolution,
			)
		})
	})

	t.Run("see at least one of these met", func(t *testing.T) {
		t.Parallel()

		adam := screenplay.ActorNamed("Adam")
		statusCode := fixture.NewFakeQuestion("status code", 200)
		logs := fixture.NewFakeQuestion("logs", "order placed")
		errNoLogs := errors.New("failed to read the logs")
		missingLogs := fixture.NewFailingFakeQuestion[string]("logs", errNoLogs)

		t.Run("succeeds when one term is met", func(t *testing.T) {
			t.Parallel()
			require.NoError(t, adam.AttemptsTo(see.AtLeastOneOfTheseMet(
				term.That(statusCode).Should(be.EqualTo(404)),
				term.That(logs).Should(contain.Substring("order placed")),
			)))
		})

		t.Run("succeeds when a term is met although another question cannot be answered", func(t *testing.T) {
			t.Parallel()
			require.NoError(t, adam.AttemptsTo(see.AtLeastOneOfTheseMet(
				term.That(missingLogs).Should(contain.Substring("order placed")),
				term.That(statusCode).Should(be.EqualTo(200)),
			)))
		})

		t.Run("fails when no term is met", func(t *testing.T) {
			t.Parallel()

			err := adam.AttemptsTo(see.AtLeastOneOfTheseMet(
				term.That(statusCode).Should(be.EqualTo(404)),
				term.That(missingLogs).Should(contain.Substring("order placed")),
			))

			require.EqualError(t, err, "Adam failed to see that at least one of these terms was met:\n"+
				"- the status code should be equal to 404: got 200\n"+
				"- the logs should contain substring \"order placed\": failed to read the logs")
			require.ErrorIs(t, err, errNoLogs)
		})

		t.Run("fails when there is no term", func(t *testing.T) {
			t.Parallel()
			require.Error(t, adam.AttemptsTo(see.AtLeastOneOfTheseMet()))
		})

		t.Run("implements the stringer interface", func(t *testing.T) {
			t.Parallel()
			assert.Equal(t,
				"see that at least one of these terms is met",
				see.AtLeastOneOfTheseMet(
					term.That(statusCode).Should(be.EqualTo(404)),
					term.That(logs).Should(contain.Substring("order placed")),
				).String(),
			)
		})

		t.Run("fails on a mistake even when another term is met", func(t *testing.T) {
			t.Parallel()

			met := term.That(statusCode).Should(be.EqualTo(200))

			require.ErrorIs(t, adam.AttemptsTo(see.AtLeastOneOfTheseMet(met, nil)), term.ErrInvalidTerm)
			require.ErrorIs(t,
				adam.AttemptsTo(see.AtLeastOneOfTheseMet(term.That(statusCode).Should(), met)),
				term.ErrInvalidTerm,
			)
			require.ErrorIs(t,
				adam.AttemptsTo(see.AtLeastOneOfTheseMet(term.That(statusCode).Should(be.Between(300, 200)), met)),
				screenplay.ErrInvalidResolution,
			)
		})
	})

	t.Run("see none of these met", func(t *testing.T) {
		t.Parallel()

		adam := screenplay.ActorNamed("Adam")
		statusCode := fixture.NewFakeQuestion("status code", 200)
		logs := fixture.NewFakeQuestion("logs", "order placed")
		errNoLogs := errors.New("failed to read the logs")
		missingLogs := fixture.NewFailingFakeQuestion[string]("logs", errNoLogs)

		t.Run("succeeds when no term is met", func(t *testing.T) {
			t.Parallel()
			require.NoError(t, adam.AttemptsTo(see.NoneOfTheseMet(
				term.That(statusCode).Should(be.EqualTo(500)),
				term.That(logs).Should(contain.Substring("panic")),
			)))
		})

		t.Run("fails when one term is met", func(t *testing.T) {
			t.Parallel()

			err := adam.AttemptsTo(see.NoneOfTheseMet(
				term.That(statusCode).Should(be.EqualTo(500)),
				term.That(logs).Should(contain.Substring("order placed")),
			))

			require.EqualError(t, err, `Adam failed to see that none of these terms was met: `+
				`"the logs should contain substring "order placed"" is met: got "order placed"`)
		})

		t.Run("fails when a question cannot be answered: the term is not known to be unmet", func(t *testing.T) {
			t.Parallel()
			require.ErrorIs(t, adam.AttemptsTo(see.NoneOfTheseMet(
				term.That(statusCode).Should(be.EqualTo(500)),
				term.That(missingLogs).Should(contain.Substring("panic")),
			)), errNoLogs)
		})

		t.Run("implements the stringer interface", func(t *testing.T) {
			t.Parallel()
			assert.Equal(t,
				"see that none of these terms is met",
				see.NoneOfTheseMet(
					term.That(statusCode).Should(be.EqualTo(500)),
					term.That(logs).Should(contain.Substring("panic")),
				).String(),
			)
		})

		t.Run("fails on a mistake rather than counting it as a term that is not met", func(t *testing.T) {
			t.Parallel()

			require.ErrorIs(t, adam.AttemptsTo(see.NoneOfTheseMet(nil)), term.ErrInvalidTerm)
			require.ErrorIs(t,
				adam.AttemptsTo(see.NoneOfTheseMet(term.That(statusCode).Should())),
				term.ErrInvalidTerm,
			)
			require.ErrorIs(t,
				adam.AttemptsTo(see.NoneOfTheseMet(term.That(statusCode).Should(be.Between(300, 200)))),
				screenplay.ErrInvalidResolution,
			)
		})
	})

	t.Run("see that all any none of the", func(t *testing.T) {
		t.Parallel()

		adam := screenplay.ActorNamed("Adam")
		helloAdam := fixture.NewFakeQuestion("profile widget", "hello Adam")
		helloEve := fixture.NewFakeQuestion("login form", "hello Eve")
		goodbyeEve := fixture.NewFakeQuestion("login form", "goodbye Eve")
		errNoWidget := errors.New("failed to get the widget")
		missingWidget := fixture.NewFailingFakeQuestion[string]("widget", errNoWidget)

		t.Run("ThatAllOfThe", func(t *testing.T) {
			t.Parallel()

			require.NoError(t, adam.AttemptsTo(
				see.ThatAllOfThe(helloAdam, helloEve).Should(contain.Substring("hello"))))
			require.Error(t, adam.AttemptsTo(
				see.ThatAllOfThe(helloAdam, goodbyeEve).Should(contain.Substring("hello"))))
			require.ErrorIs(t, adam.AttemptsTo(
				see.ThatAllOfThe(helloAdam, missingWidget).Should(contain.Substring("hello"))), errNoWidget)

			assert.Equal(t,
				"see that all of these terms are met",
				see.ThatAllOfThe(helloAdam, helloEve).Should(contain.Substring("hello")).String(),
			)
		})

		t.Run("ThatAnyOfThe", func(t *testing.T) {
			t.Parallel()

			require.NoError(t, adam.AttemptsTo(
				see.ThatAnyOfThe(helloAdam, goodbyeEve).Should(contain.Substring("goodbye"))))
			require.NoError(t, adam.AttemptsTo(
				see.ThatAnyOfThe(missingWidget, helloAdam).Should(contain.Substring("hello"))))
			require.Error(t, adam.AttemptsTo(
				see.ThatAnyOfThe(helloAdam, helloEve).Should(contain.Substring("goodbye"))))

			assert.Equal(t,
				"see that at least one of these terms is met",
				see.ThatAnyOfThe(helloAdam, goodbyeEve).Should(contain.Substring("goodbye")).String(),
			)
		})

		t.Run("ThatNoneOfThe", func(t *testing.T) {
			t.Parallel()

			require.NoError(t, adam.AttemptsTo(
				see.ThatNoneOfThe(helloAdam, helloEve).Should(contain.Substring("error"))))
			require.Error(t, adam.AttemptsTo(
				see.ThatNoneOfThe(helloAdam, goodbyeEve).Should(contain.Substring("goodbye"))))
			require.ErrorIs(t, adam.AttemptsTo(
				see.ThatNoneOfThe(helloAdam, missingWidget).Should(contain.Substring("error"))), errNoWidget)

			assert.Equal(t,
				"see that none of these terms is met",
				see.ThatNoneOfThe(helloAdam, helloEve).Should(contain.Substring("error")).String(),
			)
		})

		t.Run("fail with an invalid term when a question is nil", func(t *testing.T) {
			t.Parallel()

			hello := contain.Substring("hello")

			require.ErrorIs(t, adam.AttemptsTo(see.ThatAllOfThe(helloAdam, nil).Should(hello)), term.ErrInvalidTerm)
			require.ErrorIs(t, adam.AttemptsTo(see.ThatAnyOfThe(helloAdam, nil).Should(hello)), term.ErrInvalidTerm)
			require.ErrorIs(t, adam.AttemptsTo(see.ThatNoneOfThe[string](nil).Should(hello)), term.ErrInvalidTerm)
		})
	})

	t.Run("see these narration", func(t *testing.T) {
		t.Parallel()

		statusCode := fixture.NewFakeQuestion("status code", 200)
		logs := fixture.NewFakeQuestion("logs", "order placed")
		profileWidget := fixture.NewFakeQuestion("profile widget", "hello Adam")
		loginForm := fixture.NewFakeQuestion("login form", "hello Eve")

		narrated := func(action screenplay.Performable) []string {
			recorder := fixture.NewRecorder()
			production := screenplay.NewProduction(screenplay.WithNarrator(screenplay.NewNarrator(recorder)))
			adam := production.SetTheStage(screenplay.CastOfStandardActors()).ActorNamed("Adam")

			require.NoError(t, adam.AttemptsTo(action))

			messages := []string{}
			for _, message := range recorder.Messages() {
				if len(messages) == 0 || messages[len(messages)-1] != message {
					messages = append(messages, message)
				}
			}

			return messages
		}

		t.Run("all: the lead-in, then one line per term", func(t *testing.T) {
			t.Parallel()

			expected := []string{
				"Adam attempts to see that all of these terms are met",
				"Adam attempts to see that the status code should be equal to 200",
				`Adam attempts to see that the logs should contain substring "order placed"`,
				"Adam attempts to see that all of these terms are met",
			}

			assert.Equal(t, expected, narrated(see.AllOfTheseMet(
				term.That(statusCode).Should(be.EqualTo(200)),
				term.That(logs).Should(contain.Substring("order placed")),
			)))
		})

		t.Run("none: the lead-in, then one line per term that should not be met", func(t *testing.T) {
			t.Parallel()

			expected := []string{
				"Adam attempts to see that none of these terms is met",
				`Adam attempts to see that "the status code should be equal to 500" is not met`,
				`Adam attempts to see that "the logs should contain substring "panic"" is not met`,
				"Adam attempts to see that none of these terms is met",
			}

			assert.Equal(t, expected, narrated(see.NoneOfTheseMet(
				term.That(statusCode).Should(be.EqualTo(500)),
				term.That(logs).Should(contain.Substring("panic")),
			)))
		})

		t.Run("ThatNoneOfThe follows the structure of NoneOfTheseMet", func(t *testing.T) {
			t.Parallel()

			expected := []string{
				"Adam attempts to see that none of these terms is met",
				`Adam attempts to see that "the profile widget should contain substring "error"" is not met`,
				`Adam attempts to see that "the login form should contain substring "error"" is not met`,
				"Adam attempts to see that none of these terms is met",
			}

			assert.Equal(t, expected, narrated(
				see.ThatNoneOfThe(profileWidget, loginForm).Should(contain.Substring("error"))))
		})

		t.Run("ThatAllOfThe follows the structure of AllOfTheseMet", func(t *testing.T) {
			t.Parallel()

			expected := []string{
				"Adam attempts to see that all of these terms are met",
				`Adam attempts to see that the profile widget should contain substring "hello"`,
				`Adam attempts to see that the login form should contain substring "hello"`,
				"Adam attempts to see that all of these terms are met",
			}

			assert.Equal(t, expected, narrated(
				see.ThatAllOfThe(profileWidget, loginForm).Should(contain.Substring("hello"))))
		})
	})
}
