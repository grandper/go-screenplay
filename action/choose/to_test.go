package choose_test

import (
	"log" //nolint:depguard // the example narrates through the standard logger
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action"
	"github.com/grandper/go-screenplay/v2/action/choose"
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
	doTaskAndFail := createRecordingPerformable("do the task", assert.AnError)
	doTheAlternative := createRecordingPerformable("do the alternative", nil)
	doTheAlternativeAndFail := createRecordingPerformable("do the alternative", assert.AnError)

	t.Run("performs the first branch whose boolean condition holds", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(
			choose.To(doTask).When(true).
				To(doOtherTask).When(true),
		))
		assert.Equal(t, []string{"do the task"}, record)
	})

	t.Run("performs a later branch when it is the first to hold", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(
			choose.To(doTask).When(false).
				To(doOtherTask).When(true),
		))
		assert.Equal(t, []string{"do the other task"}, record)
	})

	t.Run("performs the fallback when no branch holds", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(
			choose.To(doTask).When(false).
				To(doTheAlternative).Otherwise(),
		))
		assert.Equal(t, []string{"do the alternative"}, record)
	})

	t.Run("performs nothing when no branch holds and there is no fallback", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(
			choose.To(doTask).When(false),
		))
		assert.Equal(t, []string{}, record)
	})

	t.Run("performs the branch when a question matches its resolution", func(t *testing.T) {
		reset()
		question := fixture.NewFakeQuestion("status code", 200)
		require.NoError(t, adam.AttemptsTo(
			choose.To(doTask).WhenThe(question).HappensTo(be.EqualTo(200)),
		))
		assert.Equal(t, []string{"do the task"}, record)
	})

	t.Run("chooses the first holding branch whatever the condition kind", func(t *testing.T) {
		reset()
		miss := fixture.NewFakeQuestion("status code", 500)
		hit := fixture.NewFakeQuestion("status code", 200)
		require.NoError(t, adam.AttemptsTo(
			choose.To(doTask).WhenThe(miss).HappensTo(be.EqualTo(200)).
				To(doOtherTask).When(true).
				To(doTheAlternative).WhenThe(hit).HappensTo(be.EqualTo(200)),
		))
		assert.Equal(t, []string{"do the other task"}, record)
	})

	t.Run("returns the error when a reached question fails to answer", func(t *testing.T) {
		reset()
		failing := fixture.NewFailingFakeQuestion[int]("the status code", assert.AnError)
		require.ErrorIs(t, adam.AttemptsTo(
			choose.To(doTask).WhenThe(failing).HappensTo(be.EqualTo(200)).
				To(doTheAlternative).Otherwise(),
		), assert.AnError)
		assert.Equal(t, []string{}, record)
	})

	t.Run("fails when a reached branch is built wrongly", func(t *testing.T) {
		question := fixture.NewFakeQuestion("status code", 200)

		reset()
		require.ErrorIs(t, adam.AttemptsTo(
			choose.To(doTask).WhenThe(question).HappensTo().
				To(doTheAlternative).Otherwise(),
		), term.ErrInvalidTerm)
		assert.Equal(t, []string{}, record)

		reset()
		require.ErrorIs(t, adam.AttemptsTo(
			choose.To(doTask).WhenThe(question).DoesNotHappenTo(be.Between(300, 200)).
				To(doTheAlternative).Otherwise(),
		), screenplay.ErrInvalidResolution)
		assert.Equal(t, []string{}, record)
	})

	t.Run("does not evaluate a later branch once an earlier one holds", func(t *testing.T) {
		reset()
		failing := fixture.NewFailingFakeQuestion[int]("the status code", assert.AnError)
		require.NoError(t, adam.AttemptsTo(
			choose.To(doTask).When(true).
				To(doOtherTask).WhenThe(failing).HappensTo(be.EqualTo(200)),
		))
		assert.Equal(t, []string{"do the task"}, record)
	})

	t.Run("returns the error from the performed branch", func(t *testing.T) {
		reset()
		require.ErrorIs(t, adam.AttemptsTo(
			choose.To(doTaskAndFail).When(true),
		), assert.AnError)
		assert.Equal(t, []string{assert.AnError.Error()}, record)
	})

	t.Run("returns the error from the performed fallback", func(t *testing.T) {
		reset()
		require.ErrorIs(t, adam.AttemptsTo(
			choose.To(doTask).When(false).
				To(doTheAlternativeAndFail).Otherwise(),
		), assert.AnError)
		assert.Equal(t, []string{assert.AnError.Error()}, record)
	})

	t.Run("performs every performable of a branch in order", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(
			choose.To(doTask, doOtherTask).When(true),
		))
		assert.Equal(t, []string{"do the task", "do the other task"}, record)
	})

	t.Run("does not fail on a branch built wrongly that the run never reaches", func(t *testing.T) {
		readTheResponse := createRecordingPerformable("read the response", nil)
		abortTheRequest := createRecordingPerformable("abort the request", nil)
		statusCode := fixture.NewFakeQuestion("status code", 200)

		for _, unreached := range []screenplay.Performable{
			choose.To(readTheResponse).When(true).To(abortTheRequest).WhenThe(statusCode).HappensTo(),
			choose.To(readTheResponse).When(true).To(abortTheRequest).WhenThe(statusCode).HappensTo(nil),
			choose.To(readTheResponse).When(true).To(abortTheRequest).WhenThe(statusCode).HappensTo(be.Between(300, 200)),
		} {
			reset()
			require.NoError(t, adam.AttemptsTo(unreached))
			assert.Equal(t, []string{"read the response"}, record)
		}
	})

	t.Run("closes a branch with every verb, singular and plural", func(t *testing.T) {
		readTheResponse := createRecordingPerformable("read the response", nil)
		abortTheRequest := createRecordingPerformable("abort the request", nil)
		statusCode := fixture.NewFakeQuestion("status code", 503)
		orderList := fixture.NewFakeQuestion("order list", []string{"book", "pen"})
		logs := fixture.NewFakeQuestion("logs", []string{"ready"})
		aServerError := resolutionfixture.NewFakeResolution[int]("a server error", true)
		aSuccess := resolutionfixture.NewFakeResolution[int]("a success", false)
		complete := resolutionfixture.NewFakeResolution[[]string]("complete", true)
		empty := resolutionfixture.NewFakeResolution[[]string]("empty", false)

		for _, testCase := range []struct {
			action    screenplay.Performable
			reads     string
			performed string
		}{
			{choose.To(readTheResponse).WhenThe(statusCode).HappensToBe(aServerError).To(abortTheRequest).Otherwise(),
				"choose to read the response when the status code happens to be a server error, " +
					"or to abort the request otherwise", "read the response"},
			{choose.To(readTheResponse).WhenThe(orderList).HappensToHave(length.Exactly[[]string](2)).To(abortTheRequest).Otherwise(),
				"choose to read the response when the order list happens to have length exactly 2, " +
					"or to abort the request otherwise", "read the response"},
			{choose.To(readTheResponse).WhenThe(statusCode).DoesNotHappenToBe(aSuccess).To(abortTheRequest).Otherwise(),
				"choose to read the response when the status code does not happen to be a success, " +
					"or to abort the request otherwise", "read the response"},
			{choose.To(readTheResponse).WhenThe(orderList).DoesNotHappenToHave(length.Exactly[[]string](2)).To(abortTheRequest).Otherwise(),
				"choose to read the response when the order list does not happen to have length exactly 2, " +
					"or to abort the request otherwise", "abort the request"},
			{choose.To(readTheResponse).WhenThe(statusCode).HappensToBeAnyOf(be.EqualTo(502), be.EqualTo(503)).To(abortTheRequest).Otherwise(),
				"choose to read the response when the status code happens to be any of: equal to 502, equal to 503, " +
					"or to abort the request otherwise", "read the response"},
			{choose.To(readTheResponse).WhenThe(logs).HappenTo(contain.Element("ready")).To(abortTheRequest).Otherwise(),
				`choose to read the response when the logs happen to contain element "ready", ` +
					"or to abort the request otherwise", "read the response"},
			{choose.To(readTheResponse).WhenThe(logs).HappenToBe(complete).To(abortTheRequest).Otherwise(),
				"choose to read the response when the logs happen to be complete, " +
					"or to abort the request otherwise", "read the response"},
			{choose.To(readTheResponse).WhenThe(logs).HappenToHave(length.Exactly[[]string](1)).To(abortTheRequest).Otherwise(),
				"choose to read the response when the logs happen to have length exactly 1, " +
					"or to abort the request otherwise", "read the response"},
			{choose.To(readTheResponse).WhenThe(logs).DoNotHappenTo(contain.Element("ready")).To(abortTheRequest).Otherwise(),
				`choose to read the response when the logs do not happen to contain element "ready", ` +
					"or to abort the request otherwise", "abort the request"},
			{choose.To(readTheResponse).WhenThe(logs).DoNotHappenToBe(empty).To(abortTheRequest).Otherwise(),
				"choose to read the response when the logs do not happen to be empty, " +
					"or to abort the request otherwise", "read the response"},
			{choose.To(readTheResponse).WhenThe(logs).DoNotHappenToHave(length.Exactly[[]string](0)).To(abortTheRequest).Otherwise(),
				"choose to read the response when the logs do not happen to have length exactly 0, " +
					"or to abort the request otherwise", "read the response"},
			{choose.To(readTheResponse).WhenThe(logs).HappenToBeAnyOf(complete, empty).To(abortTheRequest).Otherwise(),
				"choose to read the response when the logs happen to be any of: complete, empty, " +
					"or to abort the request otherwise", "read the response"},
		} {
			reset()
			assert.Equal(t, testCase.reads, testCase.action.String())
			require.NoError(t, adam.AttemptsTo(testCase.action))
			assert.Equal(t, []string{testCase.performed}, record)
		}
	})

	t.Run("names the sixth and the twenty-first boolean conditions in English", func(t *testing.T) {
		readTheResponse := createRecordingPerformable("read the response", nil)
		abortTheRequest := createRecordingPerformable("abort the request", nil)

		choice := choose.To(readTheResponse).When(false)
		for range 20 {
			choice = choice.To(readTheResponse).When(false)
		}
		choice = choice.To(abortTheRequest).Otherwise()

		assert.Contains(t, choice.String(), "to read the response when the fifth condition holds, ")
		assert.Contains(t, choice.String(), "to read the response when the 6th condition holds, ")
		assert.Contains(t, choice.String(), "to read the response when the 21st condition holds, ")
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		question := fixture.NewFakeQuestion("status code", 200)

		booleanChoice := choose.To(doTask).When(false).
			To(doOtherTask).When(false).
			To(doTheAlternative).Otherwise()
		assert.Equal(
			t,
			"choose to do the task when the first condition holds, "+
				"to do the other task when the second condition holds, "+
				"or to do the alternative otherwise",
			booleanChoice.String(),
		)

		questionChoice := choose.To(doTask).WhenThe(question).HappensTo(be.EqualTo(200))
		assert.Equal(
			t,
			"choose to do the task when the status code happens to be equal to 200",
			questionChoice.String(),
		)
	})

	t.Run("supports the Default wording for the fallback", func(t *testing.T) {
		withOtherwise := choose.To(doTask).When(false).To(doTheAlternative).Otherwise()
		withDefault := choose.To(doTask).When(false).To(doTheAlternative).Default()
		assert.Equal(
			t,
			"choose to do the task when the first condition holds, or to do the alternative otherwise",
			withOtherwise.String(),
		)
		assert.Equal(
			t,
			"choose to do the task when the first condition holds, or to do the alternative default",
			withDefault.String(),
		)

		reset()
		require.NoError(t, adam.AttemptsTo(withDefault))
		assert.Equal(t, []string{"do the alternative"}, record)
	})
}

func ExampleTo() {
	production := screenplay.NewProduction(
		screenplay.WithNarrator(screenplay.NewNarrator(logadapter.New(log.New(os.Stdout, "", 0)))))
	adam := production.SetTheStage(screenplay.CastOfStandardActors()).ActorNamed("Adam")
	statusCode := fixture.NewFakeQuestion("status code", 200)
	readTheResponse := fixture.NewFakePerformable("read the response", nil)
	abortTheRequest := fixture.NewFakePerformable("abort the request", nil)

	_ = adam.AttemptsTo(choose.To(readTheResponse).WhenThe(statusCode).HappensTo(be.EqualTo(200)).
		To(abortTheRequest).Otherwise())

	// Output:
	// Adam attempts to choose to read the response when the status code happens to be equal to 200, or to abort the request otherwise
	//     Adam attempts to read the response
}
