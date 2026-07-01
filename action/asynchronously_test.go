package action_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/action"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestAsynchronouslyAction(t *testing.T) {
	var (
		mutex  sync.Mutex
		record []string
	)

	createRecordingPerformable := func(description string, err error) screenplay.Performable {
		return action.FromFunc(description, func(_ *screenplay.Actor) error {
			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				record = append(record, err.Error())
				return err
			}

			record = append(record, description)
			return nil
		})
	}
	recorded := func() []string {
		mutex.Lock()
		defer mutex.Unlock()

		return append([]string{}, record...)
	}
	reset := func() {
		mutex.Lock()
		defer mutex.Unlock()

		record = []string{}
	}

	doTask1 := createRecordingPerformable("do task 1", nil)
	doTask2 := createRecordingPerformable("do task 2", nil)
	doTask1AndFail := createRecordingPerformable("do task 1", errors.New("failed to do task 1"))
	doTask2AndFail := createRecordingPerformable("do task 2", errors.New("failed to do task 2"))

	asynchronousErrors := func(t *testing.T, actor *screenplay.Actor) []error {
		t.Helper()

		answer, err := action.AsynchronousErrors().AnsweredBy(actor)
		require.NoError(t, err)

		errs, ok := answer.([]error)
		require.True(t, ok)

		return errs
	}

	t.Run("returns immediately and keeps running in the background", func(t *testing.T) {
		reset()
		adam := screenplay.ActorNamed("Adam")

		started := make(chan struct{})
		release := make(chan struct{})
		blocking := action.FromFunc("blocking task", func(_ *screenplay.Actor) error {
			close(started)
			<-release
			return nil
		})

		require.NoError(t, adam.AttemptsTo(action.Asynchronously(blocking)))

		// AttemptsTo has already returned while the task is still blocked.
		<-started
		close(release)

		// Answering the question waits for the background action to complete.
		assert.Empty(t, asynchronousErrors(t, adam))
	})

	t.Run("waiting for all collects every error", func(t *testing.T) {
		reset()
		adam := screenplay.ActorNamed("Adam")

		require.NoError(t, adam.AttemptsTo(action.Asynchronously(doTask1AndFail, doTask2AndFail).WaitingForAll()))

		errs := asynchronousErrors(t, adam)
		require.Len(t, errs, 2)
		messages := []string{errs[0].Error(), errs[1].Error()}
		assert.Contains(t, messages[0]+messages[1], "failed to do task 1")
		assert.Contains(t, messages[0]+messages[1], "failed to do task 2")
	})

	t.Run("waiting for all is the implicit default", func(t *testing.T) {
		reset()
		adam := screenplay.ActorNamed("Adam")

		require.NoError(t, adam.AttemptsTo(action.Asynchronously(doTask1, doTask2)))

		assert.Empty(t, asynchronousErrors(t, adam))
		assert.ElementsMatch(t, []string{"do task 1", "do task 2"}, recorded())
	})

	t.Run("cancel on error collects the error", func(t *testing.T) {
		reset()
		adam := screenplay.ActorNamed("Adam")

		require.NoError(t, adam.AttemptsTo(action.Asynchronously(doTask1AndFail).CancelOnError()))

		errs := asynchronousErrors(t, adam)
		require.Len(t, errs, 1)
		assert.Contains(t, errs[0].Error(), "failed to do task 1")
	})

	t.Run("ignoring errors collects nothing", func(t *testing.T) {
		reset()
		adam := screenplay.ActorNamed("Adam")

		require.NoError(t, adam.AttemptsTo(action.Asynchronously(doTask1AndFail, doTask2).IgnoringErrors()))

		assert.Empty(t, asynchronousErrors(t, adam))
		assert.ElementsMatch(t, []string{"do task 2", "failed to do task 1"}, recorded())
	})

	t.Run("collects the errors of several asynchronous actions", func(t *testing.T) {
		reset()
		adam := screenplay.ActorNamed("Adam")

		require.NoError(t, adam.AttemptsTo(action.Asynchronously(doTask1AndFail)))
		require.NoError(t, adam.AttemptsTo(action.Asynchronously(doTask2AndFail)))

		assert.Len(t, asynchronousErrors(t, adam), 2)
	})

	t.Run("has no errors when nothing was launched asynchronously", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")

		assert.Empty(t, asynchronousErrors(t, adam))
	})

	t.Run("with a limit never runs more tasks than allowed at the same time", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")

		var (
			concurrencyMutex sync.Mutex
			current          int
			peak             int
		)

		countingPerformable := func(description string) screenplay.Performable {
			return action.FromFunc(description, func(_ *screenplay.Actor) error {
				concurrencyMutex.Lock()
				current++
				if current > peak {
					peak = current
				}
				concurrencyMutex.Unlock()

				time.Sleep(20 * time.Millisecond)

				concurrencyMutex.Lock()
				current--
				concurrencyMutex.Unlock()

				return nil
			})
		}

		require.NoError(t, adam.AttemptsTo(
			action.Asynchronously(
				countingPerformable("do task 1"),
				countingPerformable("do task 2"),
				countingPerformable("do task 3"),
				countingPerformable("do task 4"),
				countingPerformable("do task 5"),
			).WithLimit(2).WaitingForAll(),
		))

		// Wait for the background actions to complete.
		assert.Empty(t, asynchronousErrors(t, adam))

		concurrencyMutex.Lock()
		defer concurrencyMutex.Unlock()
		assert.LessOrEqual(t, peak, 2)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(
			t,
			"asynchronously do task 1, do task 2",
			action.Asynchronously(doTask1, doTask2).String(),
		)
		assert.Equal(
			t,
			"asynchronously do task 1, do task 2",
			action.Asynchronously(doTask1, doTask2).WaitingForAll().String(),
		)
		assert.Equal(
			t,
			"asynchronously (cancelling on error) do task 1, do task 2",
			action.Asynchronously(doTask1, doTask2).CancelOnError().String(),
		)
		assert.Equal(
			t,
			"asynchronously (ignoring errors) do task 1, do task 2",
			action.Asynchronously(doTask1, doTask2).IgnoringErrors().String(),
		)
		assert.Equal(
			t,
			"asynchronously (limited to 3 at a time) do task 1, do task 2",
			action.Asynchronously(doTask1, doTask2).WithLimit(3).String(),
		)
	})

	t.Run("the asynchronous errors question is a stringer", func(t *testing.T) {
		assert.Equal(t, "the asynchronous errors", action.AsynchronousErrors().String())
	})
}
