package proceed_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action"
	"github.com/grandper/go-screenplay/v2/action/proceed"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestEitherTo(t *testing.T) {
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
	doTask1 := createRecordingPerformable("do task 1", nil)
	doTask2 := createRecordingPerformable("do task 2", nil)
	doTask1AndFail := createRecordingPerformable("do task 1", errors.New("failed to do task 1"))
	doTask2AndFail := createRecordingPerformable("do task 2", errors.New("failed to do task 2"))
	doTheAlternativeOfTask1 := createRecordingPerformable("do the alternative of task 1", nil)
	doTheAlternativeOfTask2 := createRecordingPerformable("do the alternative of task 2", nil)
	doTheAlternativeOfTask1AndFail := createRecordingPerformable(
		"do the alternative of task 1",
		errors.New("failed to do the alternative of task 1"),
	)
	doTheAlternativeOfTask2AndFail := createRecordingPerformable(
		"do the alternative of task 2",
		errors.New("failed to do the alternative of task 2"),
	)

	t.Run("performs the first task", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(proceed.EitherTo(doTask1).Or(doTheAlternativeOfTask1)))
		assert.Equal(t, []string{"do task 1"}, record)
	})

	t.Run("performs multiple tasks", func(t *testing.T) {
		reset()
		require.NoError(
			t,
			adam.AttemptsTo(proceed.EitherTo(doTask1, doTask2).Or(doTheAlternativeOfTask1, doTheAlternativeOfTask2)),
		)
		assert.Equal(t, []string{"do task 1", "do task 2"}, record)
	})

	t.Run("fails to perform the first task but succeeds to do the alternate task", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(proceed.EitherTo(doTask1AndFail).Or(doTheAlternativeOfTask1)))
		assert.Equal(t, []string{"failed to do task 1", "do the alternative of task 1"}, record)
	})

	t.Run("fails to perform the multiple tasks but succeeds to do the alternate task", func(t *testing.T) {
		reset()
		require.NoError(
			t,
			adam.AttemptsTo(
				proceed.EitherTo(doTask1, doTask2AndFail).Or(doTheAlternativeOfTask1, doTheAlternativeOfTask2),
			),
		)
		assert.Equal(
			t,
			[]string{
				"do task 1",
				"failed to do task 2",
				"do the alternative of task 1",
				"do the alternative of task 2",
			},
			record,
		)
	})

	t.Run("returns an error when all tasks failed", func(t *testing.T) {
		reset()
		require.Error(t, adam.AttemptsTo(proceed.EitherTo(doTask1AndFail).Or(doTheAlternativeOfTask1AndFail)))
		assert.Equal(t, []string{"failed to do task 1", "failed to do the alternative of task 1"}, record)
	})

	t.Run("returns an error when all task groups failed", func(t *testing.T) {
		reset()
		require.Error(
			t,
			adam.AttemptsTo(
				proceed.EitherTo(doTask1, doTask2AndFail).Or(doTheAlternativeOfTask1, doTheAlternativeOfTask2AndFail),
			),
		)
		assert.Equal(
			t,
			[]string{
				"do task 1",
				"failed to do task 2",
				"do the alternative of task 1",
				"failed to do the alternative of task 2",
			},
			record,
		)
	})

	t.Run("returns an error when all task groups failed quickly", func(t *testing.T) {
		reset()
		require.Error(
			t,
			adam.AttemptsTo(
				proceed.EitherTo(doTask1AndFail, doTask2AndFail).
					Or(doTheAlternativeOfTask1AndFail, doTheAlternativeOfTask2AndFail),
			),
		)
		assert.Equal(t, []string{"failed to do task 1", "failed to do the alternative of task 1"}, record)
	})

	t.Run("returns an error when no action is provided", func(t *testing.T) {
		reset()
		require.ErrorIs(t, adam.AttemptsTo(proceed.EitherTo().Or(doTheAlternativeOfTask1AndFail)), proceed.ErrNoPerformableProvided)
		assert.Equal(t, []string{}, record)
	})

	t.Run("returns an error when no alternative action is provided", func(t *testing.T) {
		reset()
		require.ErrorIs(t, adam.AttemptsTo(proceed.EitherTo(doTask1).Or()), proceed.ErrNoAlternateProvided)
		assert.Equal(t, []string{}, record)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		reset()

		action1 := proceed.EitherTo(doTask1).Or(doTheAlternativeOfTask1)
		assert.Equal(t, "proceed either to do task 1 or do the alternative of task 1", action1.String())

		action2 := proceed.EitherTo(doTask1, doTask2).Or(doTheAlternativeOfTask1, doTheAlternativeOfTask2)
		assert.Equal(
			t,
			"proceed either to do task 1, do task 2 or do the alternative of task 1, do the alternative of task 2",
			action2.String(),
		)
	})

}
