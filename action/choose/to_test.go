package choose_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action"
	"github.com/grandper/go-screenplay/v2/action/choose"
	"github.com/grandper/go-screenplay/v2/fixture"
	"github.com/grandper/go-screenplay/v2/resolution/be"
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/grandper/go-screenplay/v2/term"
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
		assert.Equal(t, "choose to do the task when the first condition holds, or to do the alternative otherwise", withOtherwise.String())
		assert.Equal(t, "choose to do the task when the first condition holds, or to do the alternative default", withDefault.String())

		reset()
		require.NoError(t, adam.AttemptsTo(withDefault))
		assert.Equal(t, []string{"do the alternative"}, record)
	})
}
