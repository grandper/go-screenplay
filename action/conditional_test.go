package action_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/action"
	"github.com/grandper/go-screenplay/fixture"
	"github.com/grandper/go-screenplay/resolution"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestConditionalAction(t *testing.T) {
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

	isEqualTo := func(expected any) screenplay.Resolution {
		return resolution.FromFunc("equal to the expected value", func() screenplay.Matcher {
			return func(obj any) (bool, error) {
				return obj == expected, nil
			}
		})
	}

	t.Run("performs the task when the boolean condition holds", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(action.Conditionally(doTask).If(true)))
		assert.Equal(t, []string{"do the task"}, record)
	})

	t.Run("performs multiple tasks when the condition holds", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(action.Conditionally(doTask, doOtherTask).If(true)))
		assert.Equal(t, []string{"do the task", "do the other task"}, record)
	})

	t.Run("performs the alternative when the boolean condition does not hold", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(action.Conditionally(doTask).If(false).Otherwise(doTheAlternative)))
		assert.Equal(t, []string{"do the alternative"}, record)
	})

	t.Run("performs nothing when the condition does not hold and there is no alternative", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(action.Conditionally(doTask).If(false)))
		assert.Equal(t, []string{}, record)
	})

	t.Run("negates the boolean condition with Unless", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(action.Conditionally(doTask).Unless(false)))
		assert.Equal(t, []string{"do the task"}, record)

		reset()
		require.NoError(t, adam.AttemptsTo(action.Conditionally(doTask).Unless(true).Otherwise(doTheAlternative)))
		assert.Equal(t, []string{"do the alternative"}, record)
	})

	t.Run("performs the task when the question matches the resolution", func(t *testing.T) {
		reset()
		question := fixture.NewFakeQuestion("the status code", 200)
		require.NoError(t, adam.AttemptsTo(action.Conditionally(doTask).IfThe(question, isEqualTo(200))))
		assert.Equal(t, []string{"do the task"}, record)
	})

	t.Run("performs the alternative when the question does not match the resolution", func(t *testing.T) {
		reset()
		question := fixture.NewFakeQuestion("the status code", 500)
		require.NoError(
			t,
			adam.AttemptsTo(action.Conditionally(doTask).IfThe(question, isEqualTo(200)).Otherwise(doTheAlternative)),
		)
		assert.Equal(t, []string{"do the alternative"}, record)
	})

	t.Run("returns the error when the question fails to answer", func(t *testing.T) {
		reset()
		question := fixture.NewFailingFakeQuestion("the status code", assert.AnError)
		require.ErrorIs(
			t,
			adam.AttemptsTo(action.Conditionally(doTask).IfThe(question, isEqualTo(200)).Otherwise(doTheAlternative)),
			assert.AnError,
		)
		assert.Equal(t, []string{}, record)
	})

	t.Run("returns the error when the resolution matcher fails", func(t *testing.T) {
		reset()
		question := fixture.NewFakeQuestion("the status code", 200)
		failingResolution := resolution.FromFunc("valid", func() screenplay.Matcher {
			return func(_ any) (bool, error) {
				return false, assert.AnError
			}
		})
		require.ErrorIs(
			t,
			adam.AttemptsTo(
				action.Conditionally(doTask).IfThe(question, failingResolution).Otherwise(doTheAlternative),
			),
			assert.AnError,
		)
		assert.Equal(t, []string{}, record)
	})

	t.Run("negates the question/resolution condition with UnlessThe", func(t *testing.T) {
		reset()
		question := fixture.NewFakeQuestion("the status code", 500)
		require.NoError(t, adam.AttemptsTo(action.Conditionally(doTask).UnlessThe(question, isEqualTo(200))))
		assert.Equal(t, []string{"do the task"}, record)

		reset()
		matching := fixture.NewFakeQuestion("the status code", 200)
		require.NoError(
			t,
			adam.AttemptsTo(
				action.Conditionally(doTask).UnlessThe(matching, isEqualTo(200)).Otherwise(doTheAlternative),
			),
		)
		assert.Equal(t, []string{"do the alternative"}, record)
	})

	t.Run("returns the error from the performed task", func(t *testing.T) {
		reset()
		require.Error(t, adam.AttemptsTo(action.Conditionally(doTaskAndFail).If(true)))
		assert.Equal(t, []string{"failed to do the task"}, record)
	})

	t.Run("returns the error from the performed alternative", func(t *testing.T) {
		reset()
		require.Error(t, adam.AttemptsTo(action.Conditionally(doTask).If(false).Otherwise(doTheAlternativeAndFail)))
		assert.Equal(t, []string{"failed to do the alternative"}, record)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		question := fixture.NewFakeQuestion("status code", 200)

		booleanAction := action.Conditionally(doTask).If(true)
		assert.Equal(t, "do the task if the condition holds", booleanAction.String())

		questionAction := action.Conditionally(doTask).IfThe(question, isEqualTo(200)).Otherwise(doTheAlternative)
		assert.Equal(
			t,
			"do the task if the status code is equal to the expected value, otherwise do the alternative",
			questionAction.String(),
		)

		negatedAction := action.Conditionally(doTask).UnlessThe(question, isEqualTo(200))
		assert.Equal(t, "do the task if the status code is not equal to the expected value", negatedAction.String())
	})

	t.Run("support alternative wordings", func(t *testing.T) {
		question := fixture.NewFakeQuestion("status code", 200)

		assert.Equal(
			t,
			action.Conditionally(doTask).If(true).String(),
			action.Conditionally(doTask).When(true).String(),
		)
		assert.Equal(
			t,
			action.Conditionally(doTask).IfThe(question, isEqualTo(200)).String(),
			action.Conditionally(doTask).WhenThe(question, isEqualTo(200)).String(),
		)

		reset()
		require.NoError(t, adam.AttemptsTo(action.Conditionally(doTask).When(true)))
		assert.Equal(t, []string{"do the task"}, record)
	})
}
