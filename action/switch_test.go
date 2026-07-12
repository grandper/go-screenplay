package action_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/action"
	"github.com/grandper/go-screenplay/fixture"
	"github.com/grandper/go-screenplay/resolution"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestChooseAction(t *testing.T) {
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

	isEqualTo := func(expected any) screenplay.Resolution {
		return resolution.FromFunc("equal to the expected value", func() screenplay.Matcher {
			return func(obj any) (bool, error) {
				return obj == expected, nil
			}
		})
	}
	failingResolution := resolution.FromFunc("valid", func() screenplay.Matcher {
		return func(_ any) (bool, error) {
			return false, assert.AnError
		}
	})

	t.Run("performs the first branch whose boolean condition holds", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(
			action.Choose().
				To(doTask).When(true).
				To(doOtherTask).When(true),
		))
		assert.Equal(t, []string{"do the task"}, record)
	})

	t.Run("performs a later branch when it is the first to hold", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(
			action.Choose().
				To(doTask).When(false).
				To(doOtherTask).When(true),
		))
		assert.Equal(t, []string{"do the other task"}, record)
	})

	t.Run("performs the fallback when no branch holds", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(
			action.Choose().
				To(doTask).When(false).
				To(doTheAlternative).Otherwise(),
		))
		assert.Equal(t, []string{"do the alternative"}, record)
	})

	t.Run("performs nothing when no branch holds and there is no fallback", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(
			action.Choose().
				To(doTask).When(false),
		))
		assert.Equal(t, []string{}, record)
	})

	t.Run("performs the branch when a question matches its resolution", func(t *testing.T) {
		reset()
		question := fixture.NewFakeQuestion("the status code", 200)
		require.NoError(t, adam.AttemptsTo(
			action.Choose().
				To(doTask).WhenThe(question, isEqualTo(200)),
		))
		assert.Equal(t, []string{"do the task"}, record)
	})

	t.Run("chooses the first holding branch whatever the condition kind", func(t *testing.T) {
		reset()
		miss := fixture.NewFakeQuestion("the status code", 500)
		hit := fixture.NewFakeQuestion("the status code", 200)
		require.NoError(t, adam.AttemptsTo(
			action.Choose().
				To(doTask).WhenThe(miss, isEqualTo(200)).
				To(doOtherTask).When(true).
				To(doTheAlternative).WhenThe(hit, isEqualTo(200)),
		))
		assert.Equal(t, []string{"do the other task"}, record)
	})

	t.Run("returns the error when a reached question fails to answer", func(t *testing.T) {
		reset()
		failing := fixture.NewFailingFakeQuestion("the status code", assert.AnError)
		require.ErrorIs(t, adam.AttemptsTo(
			action.Choose().
				To(doTask).WhenThe(failing, isEqualTo(200)).
				To(doTheAlternative).Otherwise(),
		), assert.AnError)
		assert.Equal(t, []string{}, record)
	})

	t.Run("returns the error when a reached resolution matcher fails", func(t *testing.T) {
		reset()
		question := fixture.NewFakeQuestion("the status code", 200)
		require.ErrorIs(t, adam.AttemptsTo(
			action.Choose().
				To(doTask).WhenThe(question, failingResolution).
				To(doTheAlternative).Otherwise(),
		), assert.AnError)
		assert.Equal(t, []string{}, record)
	})

	t.Run("does not evaluate a later branch once an earlier one holds", func(t *testing.T) {
		reset()
		failing := fixture.NewFailingFakeQuestion("the status code", assert.AnError)
		require.NoError(t, adam.AttemptsTo(
			action.Choose().
				To(doTask).When(true).
				To(doOtherTask).WhenThe(failing, isEqualTo(200)),
		))
		assert.Equal(t, []string{"do the task"}, record)
	})

	t.Run("returns the error from the performed branch", func(t *testing.T) {
		reset()
		require.ErrorIs(t, adam.AttemptsTo(
			action.Choose().
				To(doTaskAndFail).When(true),
		), assert.AnError)
		assert.Equal(t, []string{assert.AnError.Error()}, record)
	})

	t.Run("returns the error from the performed fallback", func(t *testing.T) {
		reset()
		require.ErrorIs(t, adam.AttemptsTo(
			action.Choose().
				To(doTask).When(false).
				To(doTheAlternativeAndFail).Otherwise(),
		), assert.AnError)
		assert.Equal(t, []string{assert.AnError.Error()}, record)
	})

	t.Run("performs every performable of a branch in order", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(
			action.Choose().
				To(doTask, doOtherTask).When(true),
		))
		assert.Equal(t, []string{"do the task", "do the other task"}, record)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		question := fixture.NewFakeQuestion("status code", 200)

		booleanChoice := action.Choose().
			To(doTask).When(false).
			To(doOtherTask).When(false).
			To(doTheAlternative).Otherwise()
		assert.Equal(
			t,
			"choose to do the task when the first condition holds, "+
				"to do the other task when the second condition holds, "+
				"or to do the alternative otherwise",
			booleanChoice.String(),
		)

		questionChoice := action.Choose().To(doTask).WhenThe(question, isEqualTo(200))
		assert.Equal(
			t,
			"choose to do the task when the status code is equal to the expected value",
			questionChoice.String(),
		)
	})

	t.Run("supports the Default wording for the fallback", func(t *testing.T) {
		withOtherwise := action.Choose().To(doTask).When(false).To(doTheAlternative).Otherwise()
		withDefault := action.Choose().To(doTask).When(false).To(doTheAlternative).Default()
		assert.Equal(t, withOtherwise.String(), withDefault.String())

		reset()
		require.NoError(t, adam.AttemptsTo(withDefault))
		assert.Equal(t, []string{"do the alternative"}, record)
	})
}

func TestChooseBasedOnTheAction(t *testing.T) {
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
	doTheAlternative := createRecordingPerformable("do the alternative", nil)

	isEqualTo := func(expected any) screenplay.Resolution {
		return resolution.FromFunc("equal to the expected value", func() screenplay.Matcher {
			return func(obj any) (bool, error) {
				return obj == expected, nil
			}
		})
	}
	failingResolution := resolution.FromFunc("valid", func() screenplay.Matcher {
		return func(_ any) (bool, error) {
			return false, assert.AnError
		}
	})

	t.Run("performs the branch whose resolution matches the answer", func(t *testing.T) {
		reset()
		question := fixture.NewFakeQuestion("the status code", 200)
		require.NoError(t, adam.AttemptsTo(
			action.ChooseBasedOnThe(question).
				To(doTask).When(isEqualTo(200)).
				To(doOtherTask).When(isEqualTo(404)),
		))
		assert.Equal(t, []string{"do the task"}, record)
	})

	t.Run("performs a later branch when only it matches", func(t *testing.T) {
		reset()
		question := fixture.NewFakeQuestion("the status code", 404)
		require.NoError(t, adam.AttemptsTo(
			action.ChooseBasedOnThe(question).
				To(doTask).When(isEqualTo(200)).
				To(doOtherTask).When(isEqualTo(404)),
		))
		assert.Equal(t, []string{"do the other task"}, record)
	})

	t.Run("performs the fallback when no branch matches", func(t *testing.T) {
		reset()
		question := fixture.NewFakeQuestion("the status code", 500)
		require.NoError(t, adam.AttemptsTo(
			action.ChooseBasedOnThe(question).
				To(doTask).When(isEqualTo(200)).
				To(doTheAlternative).Otherwise(),
		))
		assert.Equal(t, []string{"do the alternative"}, record)
	})

	t.Run("returns the error when the question fails to answer", func(t *testing.T) {
		reset()
		failing := fixture.NewFailingFakeQuestion("the status code", assert.AnError)
		require.ErrorIs(t, adam.AttemptsTo(
			action.ChooseBasedOnThe(failing).
				To(doTask).When(isEqualTo(200)).
				To(doTheAlternative).Otherwise(),
		), assert.AnError)
		assert.Equal(t, []string{}, record)
	})

	t.Run("returns the error when a branch resolution matcher fails", func(t *testing.T) {
		reset()
		question := fixture.NewFakeQuestion("the status code", 200)
		require.ErrorIs(t, adam.AttemptsTo(
			action.ChooseBasedOnThe(question).
				To(doTask).When(failingResolution).
				To(doOtherTask).When(isEqualTo(200)),
		), assert.AnError)
		assert.Equal(t, []string{}, record)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		question := fixture.NewFakeQuestion("status code", 200)

		choice := action.ChooseBasedOnThe(question).
			To(doTask).When(isEqualTo(200)).
			To(doOtherTask).When(isEqualTo(404)).
			To(doTheAlternative).Otherwise()
		assert.Equal(
			t,
			"choose to do the task when the status code is equal to the expected value, "+
				"to do the other task when it is equal to the expected value, "+
				"or to do the alternative otherwise",
			choice.String(),
		)
	})

	t.Run("supports the Default wording for the fallback", func(t *testing.T) {
		question := fixture.NewFakeQuestion("the status code", 500)
		withOtherwise := action.ChooseBasedOnThe(question).
			To(doTask).
			When(isEqualTo(200)).
			To(doTheAlternative).
			Otherwise()
		withDefault := action.ChooseBasedOnThe(question).To(doTask).When(isEqualTo(200)).To(doTheAlternative).Default()
		assert.Equal(t, withOtherwise.String(), withDefault.String())

		reset()
		require.NoError(t, adam.AttemptsTo(withDefault))
		assert.Equal(t, []string{"do the alternative"}, record)
	})
}
