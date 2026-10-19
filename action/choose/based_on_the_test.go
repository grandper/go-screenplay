package choose_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action"
	"github.com/grandper/go-screenplay/v2/action/choose"
	"github.com/grandper/go-screenplay/v2/action/see/term"
	"github.com/grandper/go-screenplay/v2/resolution/be"
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/grandper/go-screenplay/v2/screenplay/fixture"
)

func TestBasedOnThe(t *testing.T) {
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

	t.Run("performs the branch whose resolution matches the answer", func(t *testing.T) {
		reset()
		question := fixture.NewFakeQuestion("status code", 200)
		require.NoError(t, adam.AttemptsTo(
			choose.BasedOnThe(question).
				To(doTask).WhenIt().HappensTo(be.EqualTo(200)).
				To(doOtherTask).WhenIt().HappensTo(be.EqualTo(404)),
		))
		assert.Equal(t, []string{"do the task"}, record)
	})

	t.Run("performs a later branch when only it matches", func(t *testing.T) {
		reset()
		question := fixture.NewFakeQuestion("status code", 404)
		require.NoError(t, adam.AttemptsTo(
			choose.BasedOnThe(question).
				To(doTask).WhenIt().HappensTo(be.EqualTo(200)).
				To(doOtherTask).WhenIt().HappensTo(be.EqualTo(404)),
		))
		assert.Equal(t, []string{"do the other task"}, record)
	})

	t.Run("fails when a reached branch is built wrongly", func(t *testing.T) {
		question := fixture.NewFakeQuestion("status code", 200)

		reset()
		require.ErrorIs(t, adam.AttemptsTo(
			choose.BasedOnThe(question).
				To(doTask).WhenIt().DoesNotHappenTo().
				To(doTheAlternative).Otherwise(),
		), term.ErrInvalidTerm)
		assert.Equal(t, []string{}, record)

		reset()
		require.ErrorIs(t, adam.AttemptsTo(
			choose.BasedOnThe(question).
				To(doTask).WhenIt().DoesNotHappenTo(be.Between(300, 200)).
				To(doTheAlternative).Otherwise(),
		), screenplay.ErrInvalidResolution)
		assert.Equal(t, []string{}, record)
	})

	t.Run("supports every verb to close a branch", func(t *testing.T) {
		reset()
		question := fixture.NewFakeQuestion("status code", 503)
		require.NoError(t, adam.AttemptsTo(
			choose.BasedOnThe(question).
				To(doTask).WhenIt().HappensToBeAnyOf(be.EqualTo(200), be.EqualTo(201)).
				To(doOtherTask).WhenIt().DoesNotHappenTo(be.LessThan(500)),
		))
		assert.Equal(t, []string{"do the other task"}, record)
	})

	t.Run("performs the fallback when no branch matches", func(t *testing.T) {
		reset()
		question := fixture.NewFakeQuestion("status code", 500)
		require.NoError(t, adam.AttemptsTo(
			choose.BasedOnThe(question).
				To(doTask).WhenIt().HappensTo(be.EqualTo(200)).
				To(doTheAlternative).Otherwise(),
		))
		assert.Equal(t, []string{"do the alternative"}, record)
	})

	t.Run("returns the error when the question fails to answer", func(t *testing.T) {
		reset()
		failing := fixture.NewFailingFakeQuestion[int]("the status code", assert.AnError)
		require.ErrorIs(t, adam.AttemptsTo(
			choose.BasedOnThe(failing).
				To(doTask).WhenIt().HappensTo(be.EqualTo(200)).
				To(doTheAlternative).Otherwise(),
		), assert.AnError)
		assert.Equal(t, []string{}, record)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		question := fixture.NewFakeQuestion("status code", 200)

		choice := choose.BasedOnThe(question).
			To(doTask).WhenIt().HappensTo(be.EqualTo(200)).
			To(doOtherTask).WhenIt().HappensTo(be.EqualTo(404)).
			To(doTheAlternative).Otherwise()
		assert.Equal(
			t,
			"choose based on the status code to do the task when it happens to be equal to 200, "+
				"to do the other task when it happens to be equal to 404, "+
				"or to do the alternative otherwise",
			choice.String(),
		)
	})

	t.Run("supports the Default wording for the fallback", func(t *testing.T) {
		question := fixture.NewFakeQuestion("status code", 500)
		withOtherwise := choose.BasedOnThe(question).
			To(doTask).
			WhenIt().HappensTo(be.EqualTo(200)).
			To(doTheAlternative).
			Otherwise()
		withDefault := choose.BasedOnThe(question).
			To(doTask).WhenIt().HappensTo(be.EqualTo(200)).
			To(doTheAlternative).Default()
		assert.Equal(
			t,
			"choose based on the status code to do the task when it happens to be equal to 200, or to do the alternative otherwise",
			withOtherwise.String(),
		)
		assert.Equal(
			t,
			"choose based on the status code to do the task when it happens to be equal to 200, or to do the alternative default",
			withDefault.String(),
		)

		reset()
		require.NoError(t, adam.AttemptsTo(withDefault))
		assert.Equal(t, []string{"do the alternative"}, record)
	})
}
