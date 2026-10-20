package wait_test

import (
	"errors"
	"log" //nolint:depguard // the example narrates through the standard logger
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action/see/term"
	"github.com/grandper/go-screenplay/v2/action/wait"
	"github.com/grandper/go-screenplay/v2/narrator/logadapter"
	"github.com/grandper/go-screenplay/v2/question"
	"github.com/grandper/go-screenplay/v2/resolution/be"
	"github.com/grandper/go-screenplay/v2/resolution/contain"
	resolutionfixture "github.com/grandper/go-screenplay/v2/resolution/fixture"
	"github.com/grandper/go-screenplay/v2/resolution/length"
	"github.com/grandper/go-screenplay/v2/resolution/match"
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/grandper/go-screenplay/v2/screenplay/fixture"
)

func TestUntilThe(t *testing.T) {
	adam := screenplay.ActorNamed("Adam")
	formField := fixture.NewFakeQuestion("form field", "")
	errFormFieldNotFound := errors.New("failed to find the form field")
	failingFormField := fixture.NewFailingFakeQuestion[string]("form field", errFormFieldNotFound)

	// Bob gives up after 100 ms, so that the tests of a failed wait do not take
	// the whole screenplay.DefaultTimeout.
	production := screenplay.NewProduction(
		screenplay.WithTimeout(100*time.Millisecond),
		screenplay.WithPolling(10*time.Millisecond),
	)
	bob := production.SetTheStage(screenplay.CastOfStandardActors()).ActorNamed("Bob")

	t.Run("waits until a question is answered", func(t *testing.T) {
		formFieldWillChange := fixture.NewFakeQuestion("form field", "")

		go func() {
			time.Sleep(5 * time.Millisecond)

			formFieldWillChange.AnswerWith("hello world")
		}()

		require.NoError(
			t,
			bob.AttemptsTo(wait.UntilThe(formFieldWillChange).HappensTo(be.EqualTo("hello world"))),
		)
	})

	t.Run("returns at once when the condition is already met", func(t *testing.T) {
		started := time.Now()

		require.NoError(t, adam.AttemptsTo(
			wait.UntilThe(fixture.NewFakeQuestion("form field", "ready")).HappensTo(be.EqualTo("ready"))))

		assert.Less(t, time.Since(started), screenplay.DefaultPolling)
	})

	t.Run("fails if the attempt continuously failed", func(t *testing.T) {
		err := bob.AttemptsTo(wait.UntilThe(formField).HappensTo(be.EqualTo("hello world")))
		require.EqualError(t, err,
			`Bob waited for 100ms without seeing that the form field happens to be equal to "hello world": got ""`)
	})

	t.Run("fails when the question fails", func(t *testing.T) {
		err := bob.AttemptsTo(wait.UntilThe(failingFormField).HappensTo(be.EqualTo("hello world")))
		require.ErrorIs(t, err, errFormFieldNotFound)
		require.EqualError(t, err,
			`Bob waited for 100ms without seeing that the form field happens to be equal to "hello world": `+
				"failed to find the form field")
	})

	t.Run("fails without waiting when the condition is built wrongly", func(t *testing.T) {
		started := time.Now()

		require.ErrorIs(t, adam.AttemptsTo(wait.UntilThe(formField).HappensTo()), term.ErrInvalidTerm)
		require.ErrorIs(t,
			adam.AttemptsTo(wait.UntilThe(formField).DoesNotHappenTo(match.PatternString("hello("))),
			screenplay.ErrInvalidResolution,
		)
		assert.Less(t, time.Since(started), screenplay.DefaultTimeout)
	})

	t.Run("fails at once when the polling of the production is longer than its timeout", func(t *testing.T) {
		slowPolling := screenplay.NewProduction(
			screenplay.WithTimeout(100*time.Millisecond),
			screenplay.WithPolling(time.Second),
		)
		bobPollingSlowly := slowPolling.SetTheStage(screenplay.CastOfStandardActors()).ActorNamed("Bob")

		err := bobPollingSlowly.AttemptsTo(wait.UntilThe(formField).HappensTo(be.EqualTo("ready")))

		require.ErrorIs(t, err, wait.ErrPollingPeriodMustBeLessThanOrEqualToTimeout)
		require.EqualError(t, err, `Bob failed to wait until the form field happens to be equal to "ready": `+
			"polling period must be less than or equal to timeout")
	})

	t.Run("fails at once when the timeout of the production is negative", func(t *testing.T) {
		negativeTimeout := screenplay.NewProduction(screenplay.WithTimeout(-1 * time.Second))
		bobWithNoTime := negativeTimeout.SetTheStage(screenplay.CastOfStandardActors()).ActorNamed("Bob")

		err := bobWithNoTime.AttemptsTo(wait.UntilThe(formField).HappensTo(be.EqualTo("ready")))

		require.ErrorIs(t, err, wait.ErrPollingPeriodAndTimeoutMustBePositive)
		require.EqualError(t, err, `Bob failed to wait until the form field happens to be equal to "ready": `+
			"polling period and timeout must be positive")
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		action := wait.UntilThe(formField).HappensTo(be.EqualTo("finished"))
		assert.Equal(t, `wait until the form field happens to be equal to "finished"`, action.String())
	})

	t.Run("waits for a condition stated with every verb, singular and plural", func(t *testing.T) {
		statusCode := fixture.NewFakeQuestion("status code", 503)
		orderList := fixture.NewFakeQuestion("order list", []string{"book", "pen"})
		logs := fixture.NewFakeQuestion("logs", []string{"ready"})
		aServerError := resolutionfixture.NewFakeResolution[int]("a server error", true)
		aSuccess := resolutionfixture.NewFakeResolution[int]("a success", false)
		complete := resolutionfixture.NewFakeResolution[[]string]("complete", true)
		empty := resolutionfixture.NewFakeResolution[[]string]("empty", false)

		for _, testCase := range []struct {
			action screenplay.Performable
			reads  string
		}{
			{wait.UntilThe(statusCode).HappensToBe(aServerError),
				"wait until the status code happens to be a server error"},
			{wait.UntilThe(orderList).HappensToHave(length.Exactly[[]string](2)),
				"wait until the order list happens to have length exactly 2"},
			{wait.UntilThe(statusCode).DoesNotHappenToBe(aSuccess),
				"wait until the status code does not happen to be a success"},
			{wait.UntilThe(orderList).DoesNotHappenToHave(length.Exactly[[]string](0)),
				"wait until the order list does not happen to have length exactly 0"},
			{wait.UntilThe(statusCode).HappensToBeAnyOf(be.EqualTo(502), be.EqualTo(503)),
				"wait until the status code happens to be any of: equal to 502, equal to 503"},
			{wait.UntilThe(logs).HappenTo(contain.Element("ready")),
				`wait until the logs happen to contain element "ready"`},
			{wait.UntilThe(logs).HappenToBe(complete),
				"wait until the logs happen to be complete"},
			{wait.UntilThe(logs).HappenToHave(length.Exactly[[]string](1)),
				"wait until the logs happen to have length exactly 1"},
			{wait.UntilThe(logs).DoNotHappenTo(contain.Element("panic")),
				`wait until the logs do not happen to contain element "panic"`},
			{wait.UntilThe(logs).DoNotHappenToBe(empty),
				"wait until the logs do not happen to be empty"},
			{wait.UntilThe(logs).DoNotHappenToHave(length.Exactly[[]string](0)),
				"wait until the logs do not happen to have length exactly 0"},
			{wait.UntilThe(logs).HappenToBeAnyOf(complete, empty),
				"wait until the logs happen to be any of: complete, empty"},
		} {
			assert.Equal(t, testCase.reads, testCase.action.String())
			require.NoError(t, adam.AttemptsTo(testCase.action))
		}
	})

	t.Run("fails when the condition is met only after the timeout", func(t *testing.T) {
		shortWindow := screenplay.NewProduction(
			screenplay.WithTimeout(30*time.Millisecond),
			screenplay.WithPolling(12*time.Millisecond),
		)
		bobInAHurry := shortWindow.SetTheStage(screenplay.CastOfStandardActors()).ActorNamed("Bob")

		// False during the whole window, true right after it: the last judgment,
		// made after the timeout to explain the failure, sees it met.
		start := time.Now()
		lateCondition := question.FromFunc("late condition", func(_ *screenplay.Actor) (bool, error) {
			return time.Since(start) > 30*time.Millisecond, nil
		})

		err := bobInAHurry.AttemptsTo(wait.UntilThe(lateCondition).HappensTo(be.True()))
		assert.EqualError(
			t,
			err,
			"Bob waited for 30ms without seeing that the late condition happens to be true in time",
		)
	})
}

func ExampleUntilThe() {
	production := screenplay.NewProduction(
		screenplay.WithNarrator(screenplay.NewNarrator(logadapter.New(log.New(os.Stdout, "", 0)))))
	adam := production.SetTheStage(screenplay.CastOfStandardActors()).ActorNamed("Adam")
	logs := fixture.NewFakeQuestion("logs", []string{"ready"})

	_ = adam.AttemptsTo(wait.UntilThe(logs).HappenTo(contain.Element("ready")))

	// Output:
	// Adam attempts to wait until the logs happen to contain element "ready"
}
