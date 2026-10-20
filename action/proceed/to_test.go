package proceed_test

import (
	"errors"
	"log" //nolint:depguard // the example narrates through the standard logger
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action"
	"github.com/grandper/go-screenplay/v2/action/proceed"
	"github.com/grandper/go-screenplay/v2/action/see/term"
	"github.com/grandper/go-screenplay/v2/narrator/logadapter"
	"github.com/grandper/go-screenplay/v2/resolution/be"
	"github.com/grandper/go-screenplay/v2/resolution/contain"
	resolutionfixture "github.com/grandper/go-screenplay/v2/resolution/fixture"
	"github.com/grandper/go-screenplay/v2/resolution/length"
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/grandper/go-screenplay/v2/screenplay/fixture"
)

func TestTo(t *testing.T) {
	record := []string{}
	createRecordingPerformable := func(description string, err error) screenplay.Performable {
		return action.FromFunc(description, func(_ *screenplay.Actor) error {
			if err != nil {
				record = append(record, err.Error())
				return err
			}
			record = append(record, description)
			return nil
		})
	}
	reset := func() {
		record = []string{}
	}
	adam := screenplay.ActorNamed("Adam")
	doTask := createRecordingPerformable("do the task", nil)
	doOtherTask := createRecordingPerformable("do the other task", nil)
	doTaskAndFail := createRecordingPerformable("do the task", errors.New("failed to do the task"))
	doTheAlternative := createRecordingPerformable("do the alternative", nil)
	doTheAlternativeAndFail := createRecordingPerformable(
		"do the alternative",
		errors.New("failed to do the alternative"),
	)

	t.Run("performs the task when the boolean condition holds", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(proceed.To(doTask).If(true)))
		assert.Equal(t, []string{"do the task"}, record)
	})

	t.Run("performs multiple tasks when the condition holds", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(proceed.To(doTask, doOtherTask).If(true)))
		assert.Equal(t, []string{"do the task", "do the other task"}, record)
	})

	t.Run("performs the alternative when the boolean condition does not hold", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(proceed.To(doTask).If(false).Otherwise(doTheAlternative)))
		assert.Equal(t, []string{"do the alternative"}, record)
	})

	t.Run("performs nothing when the condition does not hold and there is no alternative", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(proceed.To(doTask).If(false)))
		assert.Equal(t, []string{}, record)
	})

	t.Run("negates the boolean condition with Unless", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(proceed.To(doTask).Unless(false)))
		assert.Equal(t, []string{"do the task"}, record)

		reset()
		require.NoError(t, adam.AttemptsTo(proceed.To(doTask).Unless(true).Otherwise(doTheAlternative)))
		assert.Equal(t, []string{"do the alternative"}, record)
	})

	t.Run("performs the task when the question matches the resolution", func(t *testing.T) {
		reset()
		question := fixture.NewFakeQuestion("status code", 200)
		require.NoError(t, adam.AttemptsTo(proceed.To(doTask).IfThe(question).HappensTo(be.EqualTo(200))))
		assert.Equal(t, []string{"do the task"}, record)
	})

	t.Run("performs the alternative when the question does not match the resolution", func(t *testing.T) {
		reset()
		question := fixture.NewFakeQuestion("status code", 500)
		require.NoError(
			t,
			adam.AttemptsTo(
				proceed.To(doTask).IfThe(question).HappensTo(be.EqualTo(200)).Otherwise(doTheAlternative),
			),
		)
		assert.Equal(t, []string{"do the alternative"}, record)
	})

	t.Run("returns the error when the question fails to answer", func(t *testing.T) {
		reset()
		question := fixture.NewFailingFakeQuestion[int]("the status code", assert.AnError)
		require.ErrorIs(
			t,
			adam.AttemptsTo(
				proceed.To(doTask).IfThe(question).HappensTo(be.EqualTo(200)).Otherwise(doTheAlternative),
			),
			assert.AnError,
		)
		assert.Equal(t, []string{}, record)
	})

	t.Run("negates the question/resolution condition with UnlessThe", func(t *testing.T) {
		reset()
		question := fixture.NewFakeQuestion("status code", 500)
		require.NoError(t, adam.AttemptsTo(proceed.To(doTask).UnlessThe(question).HappensTo(be.EqualTo(200))))
		assert.Equal(t, []string{"do the task"}, record)

		reset()
		matching := fixture.NewFakeQuestion("status code", 200)
		require.NoError(
			t,
			adam.AttemptsTo(
				proceed.To(doTask).UnlessThe(matching).HappensTo(be.EqualTo(200)).Otherwise(doTheAlternative),
			),
		)
		assert.Equal(t, []string{"do the alternative"}, record)
	})

	t.Run("fails when the condition is built wrongly, whatever its polarity", func(t *testing.T) {
		question := fixture.NewFakeQuestion("status code", 200)
		broken := be.Between(300, 200)

		for _, conditional := range []screenplay.Performable{
			proceed.To(doTask).IfThe(question).HappensTo().Otherwise(doTheAlternative),
			proceed.To(doTask).UnlessThe(question).HappensTo().Otherwise(doTheAlternative),
			proceed.To(doTask).IfThe(question).DoesNotHappenTo(nil).Otherwise(doTheAlternative),
		} {
			reset()
			require.ErrorIs(t, adam.AttemptsTo(conditional), term.ErrInvalidTerm)
			assert.Equal(t, []string{}, record)
		}

		for _, conditional := range []screenplay.Performable{
			proceed.To(doTask).IfThe(question).HappensTo(broken).Otherwise(doTheAlternative),
			proceed.To(doTask).IfThe(question).DoesNotHappenTo(broken).Otherwise(doTheAlternative),
			proceed.To(doTask).UnlessThe(question).HappensTo(broken).Otherwise(doTheAlternative),
		} {
			reset()
			require.ErrorIs(t, adam.AttemptsTo(conditional), screenplay.ErrInvalidResolution)
			assert.Equal(t, []string{}, record)
		}

		// The plural verbs report the mistake with the words the user wrote.
		logs := fixture.NewFakeQuestion("logs", []string{"ready"})

		reset()
		err := adam.AttemptsTo(proceed.To(doTask).IfThe(logs).HappenTo())
		require.ErrorIs(t, err, term.ErrInvalidTerm)
		require.EqualError(t, err, "invalid term: HappenTo needs at least one resolution")
		assert.Equal(t, []string{}, record)

		reset()
		err = adam.AttemptsTo(proceed.To(doTask).IfThe(logs).DoNotHappenToBe(nil))
		require.ErrorIs(t, err, term.ErrInvalidTerm)
		require.EqualError(t, err, "invalid term: the resolution at position 0 given to DoNotHappenToBe is nil")
		assert.Equal(t, []string{}, record)
	})

	t.Run("returns the error from the performed task", func(t *testing.T) {
		reset()
		require.EqualError(t, adam.AttemptsTo(proceed.To(doTaskAndFail).If(true)), "failed to do the task")
		assert.Equal(t, []string{"failed to do the task"}, record)
	})

	t.Run("returns the error from the performed alternative", func(t *testing.T) {
		reset()
		require.EqualError(
			t,
			adam.AttemptsTo(proceed.To(doTask).If(false).Otherwise(doTheAlternativeAndFail)),
			"failed to do the alternative",
		)
		assert.Equal(t, []string{"failed to do the alternative"}, record)
	})

	t.Run("performs the task when there is no condition", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(proceed.To(doTask)))
		assert.Equal(t, []string{"do the task"}, record)
	})

	t.Run("guards the action with every verb, singular and plural", func(t *testing.T) {
		retryTheRequest := createRecordingPerformable("retry the request", nil)
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
			holds  bool
		}{
			{proceed.To(retryTheRequest).IfThe(statusCode).HappensToBe(aServerError),
				"proceed to retry the request if the status code happens to be a server error", true},
			{proceed.To(retryTheRequest).IfThe(orderList).HappensToHave(length.Exactly[[]string](2)),
				"proceed to retry the request if the order list happens to have length exactly 2", true},
			{proceed.To(retryTheRequest).IfThe(statusCode).DoesNotHappenToBe(aSuccess),
				"proceed to retry the request if the status code does not happen to be a success", true},
			{proceed.To(retryTheRequest).IfThe(orderList).DoesNotHappenToHave(length.Exactly[[]string](2)),
				"proceed to retry the request if the order list does not happen to have length exactly 2", false},
			{proceed.To(retryTheRequest).IfThe(statusCode).HappensToBeAnyOf(be.EqualTo(502), be.EqualTo(503)),
				"proceed to retry the request if the status code happens to be any of: equal to 502, equal to 503", true},
			{proceed.To(retryTheRequest).IfThe(logs).HappenTo(contain.Element("ready")),
				`proceed to retry the request if the logs happen to contain element "ready"`, true},
			{proceed.To(retryTheRequest).IfThe(logs).HappenToBe(complete),
				"proceed to retry the request if the logs happen to be complete", true},
			{proceed.To(retryTheRequest).IfThe(logs).HappenToHave(length.Exactly[[]string](1)),
				"proceed to retry the request if the logs happen to have length exactly 1", true},
			{proceed.To(retryTheRequest).IfThe(logs).DoNotHappenTo(contain.Element("ready")),
				`proceed to retry the request if the logs do not happen to contain element "ready"`, false},
			{proceed.To(retryTheRequest).IfThe(logs).DoNotHappenToBe(empty),
				"proceed to retry the request if the logs do not happen to be empty", true},
			{proceed.To(retryTheRequest).IfThe(logs).DoNotHappenToHave(length.Exactly[[]string](1)),
				"proceed to retry the request if the logs do not happen to have length exactly 1", false},
			{proceed.To(retryTheRequest).IfThe(logs).HappenToBeAnyOf(complete, empty),
				"proceed to retry the request if the logs happen to be any of: complete, empty", true},
		} {
			reset()
			assert.Equal(t, testCase.reads, testCase.action.String())
			require.NoError(t, adam.AttemptsTo(testCase.action))
			if testCase.holds {
				assert.Equal(t, []string{"retry the request"}, record)
			} else {
				assert.Equal(t, []string{}, record)
			}
		}
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		question := fixture.NewFakeQuestion("status code", 200)

		assert.Equal(t, "proceed to do the task", proceed.To(doTask).String())

		booleanAction := proceed.To(doTask).If(true)
		assert.Equal(t, "proceed to do the task if the condition holds", booleanAction.String())

		// Unless is narrated with the word written in the code, whatever the boolean.
		assert.Equal(
			t,
			"proceed to do the task unless the condition holds",
			proceed.To(doTask).Unless(true).String(),
		)
		assert.Equal(
			t,
			"proceed to do the task unless the condition holds",
			proceed.To(doTask).Unless(false).String(),
		)
		assert.Equal(
			t,
			"proceed to do the task unless the condition holds, otherwise do the alternative",
			proceed.To(doTask).Unless(true).Otherwise(doTheAlternative).String(),
		)

		questionAction := proceed.To(doTask).
			IfThe(question).
			HappensTo(be.EqualTo(200)).
			Otherwise(doTheAlternative)
		assert.Equal(
			t,
			"proceed to do the task if the status code happens to be equal to 200, otherwise do the alternative",
			questionAction.String(),
		)

		negatedAction := proceed.To(doTask).UnlessThe(question).HappensTo(be.EqualTo(200))

		withAnotherVerb := proceed.To(doTask).
			IfThe(question).
			DoesNotHappenTo(be.AtLeast(400), be.LessThan(200))
		assert.Equal(
			t,
			"proceed to do the task if the status code does not happen to be at least 400, nor be less than 200",
			withAnotherVerb.String(),
		)

		assert.Equal(
			t,
			"proceed to do the task unless the status code happens to be equal to 200",
			negatedAction.String(),
		)
	})

	t.Run("support alternative wordings", func(t *testing.T) {
		question := fixture.NewFakeQuestion("status code", 200)

		assert.Equal(
			t,
			"proceed to do the task when the condition holds",
			proceed.To(doTask).When(true).String(),
		)
		assert.Equal(
			t,
			"proceed to do the task when the status code happens to be equal to 200",
			proceed.To(doTask).WhenThe(question).HappensTo(be.EqualTo(200)).String(),
		)

		reset()
		require.NoError(t, adam.AttemptsTo(proceed.To(doTask).When(true)))
		assert.Equal(t, []string{"do the task"}, record)
	})
}

func ExampleTo() {
	production := screenplay.NewProduction(
		screenplay.WithNarrator(screenplay.NewNarrator(logadapter.New(log.New(os.Stdout, "", 0)))))
	adam := production.SetTheStage(screenplay.CastOfStandardActors()).ActorNamed("Adam")
	statusCode := fixture.NewFakeQuestion("status code", 503)
	retryTheRequest := fixture.NewFakePerformable("retry the request", nil)

	_ = adam.AttemptsTo(proceed.To(retryTheRequest).IfThe(statusCode).HappensTo(be.AtLeast(500)))

	// Output:
	// Adam attempts to proceed to retry the request if the status code happens to be at least 500
	//     Adam attempts to retry the request
}
