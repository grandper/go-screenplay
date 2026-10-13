package action_test

import (
	"errors"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action"
	"github.com/grandper/go-screenplay/v2/fixture"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestConcurrentlyAction(t *testing.T) {
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
	sortedRecord := func() []string {
		mutex.Lock()
		defer mutex.Unlock()

		sorted := append([]string{}, record...)
		sort.Strings(sorted)

		return sorted
	}
	reset := func() {
		mutex.Lock()
		defer mutex.Unlock()

		record = []string{}
	}

	adam := screenplay.ActorNamed("Adam")
	doTask1 := createRecordingPerformable("do task 1", nil)
	doTask2 := createRecordingPerformable("do task 2", nil)
	doTask3 := createRecordingPerformable("do task 3", nil)
	doTask1AndFail := createRecordingPerformable("do task 1", errors.New("failed to do task 1"))
	doTask2AndFail := createRecordingPerformable("do task 2", errors.New("failed to do task 2"))

	t.Run("performs every task in parallel", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(action.Concurrently(doTask1, doTask2, doTask3)))
		assert.Equal(t, []string{"do task 1", "do task 2", "do task 3"}, sortedRecord())
	})

	t.Run("waiting for all is the implicit default", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(action.Concurrently(doTask1, doTask2).WaitingForAll()))
		assert.Equal(t, []string{"do task 1", "do task 2"}, sortedRecord())
	})

	t.Run("waiting for all joins every error", func(t *testing.T) {
		reset()
		err := adam.AttemptsTo(action.Concurrently(doTask1AndFail, doTask2AndFail).WaitingForAll())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to do task 1")
		assert.Contains(t, err.Error(), "failed to do task 2")
		assert.Equal(t, []string{"failed to do task 1", "failed to do task 2"}, sortedRecord())
	})

	t.Run("waiting for all runs every task even when one fails", func(t *testing.T) {
		reset()
		err := adam.AttemptsTo(action.Concurrently(doTask1, doTask2AndFail, doTask3).WaitingForAll())
		require.Error(t, err)
		assert.Equal(t, []string{"do task 1", "do task 3", "failed to do task 2"}, sortedRecord())
	})

	t.Run("stopping on error returns the error that occurred", func(t *testing.T) {
		reset()
		err := adam.AttemptsTo(action.Concurrently(doTask1AndFail).StoppingOnError())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to do task 1")
	})

	t.Run("stopping on first error is an alias for stopping on error", func(t *testing.T) {
		reset()
		err := adam.AttemptsTo(action.Concurrently(doTask1AndFail).StoppingOnFirstError())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to do task 1")
	})

	t.Run("stopping on error succeeds when no task fails", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(action.Concurrently(doTask1, doTask2).StoppingOnError()))
		assert.Equal(t, []string{"do task 1", "do task 2"}, sortedRecord())
	})

	t.Run("ignoring errors never returns an error", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(action.Concurrently(doTask1AndFail, doTask2).IgnoringErrors()))
		assert.Equal(t, []string{"do task 2", "failed to do task 1"}, sortedRecord())
	})

	t.Run("with a limit never runs more tasks than allowed at the same time", func(t *testing.T) {
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
			action.Concurrently(
				countingPerformable("do task 1"),
				countingPerformable("do task 2"),
				countingPerformable("do task 3"),
				countingPerformable("do task 4"),
				countingPerformable("do task 5"),
			).WithLimit(2).WaitingForAll(),
		))

		concurrencyMutex.Lock()
		defer concurrencyMutex.Unlock()
		assert.LessOrEqual(t, peak, 2)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(
			t,
			"concurrently do task 1, do task 2",
			action.Concurrently(doTask1, doTask2).String(),
		)
		assert.Equal(
			t,
			"concurrently do task 1, do task 2",
			action.Concurrently(doTask1, doTask2).WaitingForAll().String(),
		)
		assert.Equal(
			t,
			"concurrently (stopping on error) do task 1, do task 2",
			action.Concurrently(doTask1, doTask2).StoppingOnError().String(),
		)
		assert.Equal(
			t,
			"concurrently (ignoring errors) do task 1, do task 2",
			action.Concurrently(doTask1, doTask2).IgnoringErrors().String(),
		)
		assert.Equal(
			t,
			"concurrently (limited to 3 at a time) do task 1, do task 2",
			action.Concurrently(doTask1, doTask2).WithLimit(3).String(),
		)
	})

	t.Run("simultaneously is an alias for concurrently", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(action.Simultaneously(doTask1, doTask2)))
		assert.Equal(t, []string{"do task 1", "do task 2"}, sortedRecord())
		assert.Equal(t, action.Concurrently(doTask1), action.Simultaneously(doTask1))
	})

	t.Run("ignores nil performables instead of panicking", func(t *testing.T) {
		recorder := fixture.NewRecorder()
		narratedAdam := narratedActor(recorder)

		require.NoError(t, narratedAdam.AttemptsTo(action.Concurrently(fixture.NewFakePerformable("retry", nil), nil)))
		assert.Contains(t, recorder.Messages(), "Adam attempts to concurrently retry")

		require.NoError(
			t,
			narratedAdam.AttemptsTo(action.Asynchronously(nil, fixture.NewFakePerformable("retry", nil))),
		)
		require.NoError(
			t,
			narratedAdam.AttemptsTo(
				action.Either(nil, fixture.NewFakePerformable("retry", nil)).
					Or(nil, fixture.NewFakePerformable("wait", nil)),
			),
		)
		// A list made only of nil is empty: the mistake fails when performed, without a panic.
		require.Error(t, narratedAdam.AttemptsTo(action.Either(fixture.NewFakePerformable("retry", nil)).Or(nil)))
		require.NoError(
			t,
			narratedAdam.AttemptsTo(action.Conditionally(nil, fixture.NewFakePerformable("retry", nil)).Otherwise(nil)),
		)
	})

	t.Run("narrates each of its performables", func(t *testing.T) {
		for name, mode := range map[string]func(*action.ConcurrentlyAction) *action.ConcurrentlyAction{
			"waiting for all":   (*action.ConcurrentlyAction).WaitingForAll,
			"ignoring errors":   (*action.ConcurrentlyAction).IgnoringErrors,
			"stopping on error": (*action.ConcurrentlyAction).StoppingOnError,
		} {
			t.Run(name, func(t *testing.T) {
				recorder := fixture.NewRecorder()
				narratedAdam := narratedActor(recorder)

				concurrently := action.Concurrently(
					fixture.NewFakePerformable("open the door", nil),
					fixture.NewFakePerformable("close the window", nil),
				)
				require.NoError(t, narratedAdam.AttemptsTo(mode(concurrently)))

				messages := recorder.Messages()
				assert.Contains(t, messages, "Adam attempts to open the door")
				assert.Contains(t, messages, "Adam attempts to close the window")
			})
		}
	})

	t.Run("stopping on error cancels the actions still running through the actor's context", func(t *testing.T) {
		boom := errors.New("boom")
		failImmediately := action.FromFunc("fail immediately", func(_ *screenplay.Actor) error { return boom })
		waitUnlessCanceled := func(name string) screenplay.Performable {
			return action.FromFunc(name, func(actor *screenplay.Actor) error {
				select {
				case <-actor.Context().Done():
					return actor.Context().Err()
				case <-time.After(2 * time.Second):
					return errors.New(name + " was not canceled")
				}
			})
		}

		start := time.Now()
		err := adam.AttemptsTo(action.Concurrently(
			failImmediately,
			waitUnlessCanceled("wait for the first order"),
			waitUnlessCanceled("wait for the second order"),
		).StoppingOnError())
		require.ErrorIs(t, err, boom)
		assert.Less(t, time.Since(start), time.Second)

		// The actor itself, and whoever else holds it, never saw the run's context.
		assert.NoError(t, adam.Context().Err())
	})
}
