package proceed_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action"
	"github.com/grandper/go-screenplay/v2/action/proceed"
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
	})

	t.Run("returns the error from the performed task", func(t *testing.T) {
		reset()
		require.Error(t, adam.AttemptsTo(proceed.To(doTaskAndFail).If(true)))
		assert.Equal(t, []string{"failed to do the task"}, record)
	})

	t.Run("returns the error from the performed alternative", func(t *testing.T) {
		reset()
		require.Error(t, adam.AttemptsTo(proceed.To(doTask).If(false).Otherwise(doTheAlternativeAndFail)))
		assert.Equal(t, []string{"failed to do the alternative"}, record)
	})

	t.Run("performs the task when there is no condition", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(proceed.To(doTask)))
		assert.Equal(t, []string{"do the task"}, record)
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
