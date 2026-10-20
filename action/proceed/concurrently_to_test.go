package proceed_test

import (
	"errors"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action"
	impl "github.com/grandper/go-screenplay/v2/action/internal/proceed"
	"github.com/grandper/go-screenplay/v2/action/proceed"
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/grandper/go-screenplay/v2/screenplay/fixture"
)

func TestConcurrentlyTo(t *testing.T) {
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
		require.NoError(t, adam.AttemptsTo(proceed.ConcurrentlyTo(doTask1, doTask2, doTask3)))
		assert.Equal(t, []string{"do task 1", "do task 2", "do task 3"}, sortedRecord())
	})

	t.Run("waiting for all is the implicit default", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(proceed.ConcurrentlyTo(doTask1, doTask2).WaitingForAll()))
		assert.Equal(t, []string{"do task 1", "do task 2"}, sortedRecord())
	})

	t.Run("waiting for all joins every error", func(t *testing.T) {
		reset()
		err := adam.AttemptsTo(proceed.ConcurrentlyTo(doTask1AndFail, doTask2AndFail).WaitingForAll())
		require.EqualError(t, err, "do task 1: failed to do task 1\ndo task 2: failed to do task 2")
		assert.Equal(t, []string{"failed to do task 1", "failed to do task 2"}, sortedRecord())
	})

	t.Run("waiting for all runs every task even when one fails", func(t *testing.T) {
		reset()
		err := adam.AttemptsTo(proceed.ConcurrentlyTo(doTask1, doTask2AndFail, doTask3).WaitingForAll())
		require.EqualError(t, err, "do task 2: failed to do task 2")
		assert.Equal(t, []string{"do task 1", "do task 3", "failed to do task 2"}, sortedRecord())
	})

	t.Run("stopping on error returns the error that occurred", func(t *testing.T) {
		reset()
		err := adam.AttemptsTo(proceed.ConcurrentlyTo(doTask1AndFail).StoppingOnError())
		require.EqualError(t, err, "failed to do task 1")
	})

	t.Run("stopping on error skips the actions that have not started", func(t *testing.T) {
		errPaymentServiceDown := errors.New("the payment service is down")
		var started atomic.Int32
		checkTheStock := action.FromFunc("check the stock", func(*screenplay.Actor) error {
			started.Add(1)
			time.Sleep(50 * time.Millisecond)
			return nil
		})
		chargeTheCard := action.FromFunc("charge the card", func(*screenplay.Actor) error {
			return errPaymentServiceDown
		})

		err := adam.AttemptsTo(proceed.ConcurrentlyTo(chargeTheCard, checkTheStock, checkTheStock, checkTheStock).
			StoppingOnError().WithLimit(1))

		require.ErrorIs(t, err, errPaymentServiceDown)
		assert.Equal(t, int32(0), started.Load())
	})

	t.Run("stopping on first error behaves like stopping on error", func(t *testing.T) {
		reset()
		err := adam.AttemptsTo(proceed.ConcurrentlyTo(doTask1AndFail).StoppingOnFirstError())
		require.EqualError(t, err, "failed to do task 1")
	})

	t.Run("stopping on error succeeds when no task fails", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(proceed.ConcurrentlyTo(doTask1, doTask2).StoppingOnError()))
		assert.Equal(t, []string{"do task 1", "do task 2"}, sortedRecord())
	})

	t.Run("ignoring errors never returns an error", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(proceed.ConcurrentlyTo(doTask1AndFail, doTask2).IgnoringErrors()))
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
			proceed.ConcurrentlyTo(
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
			"proceed concurrently to do task 1, do task 2",
			proceed.ConcurrentlyTo(doTask1, doTask2).String(),
		)
		assert.Equal(
			t,
			"proceed concurrently to do task 1, do task 2, waiting for all",
			proceed.ConcurrentlyTo(doTask1, doTask2).WaitingForAll().String(),
		)
		assert.Equal(
			t,
			"proceed concurrently to do task 1, do task 2, stopping on error",
			proceed.ConcurrentlyTo(doTask1, doTask2).StoppingOnError().String(),
		)
		assert.Equal(
			t,
			"proceed concurrently to do task 1, do task 2, ignoring errors",
			proceed.ConcurrentlyTo(doTask1, doTask2).IgnoringErrors().String(),
		)
		assert.Equal(
			t,
			"proceed concurrently to do task 1, do task 2, with limit 3",
			proceed.ConcurrentlyTo(doTask1, doTask2).WithLimit(3).String(),
		)
		assert.Equal(
			t,
			"proceed concurrently to do task 1, do task 2, stopping on first error",
			proceed.ConcurrentlyTo(doTask1, doTask2).StoppingOnFirstError().String(),
		)
		// The options are narrated as they were written, in call order.
		assert.Equal(
			t,
			"proceed concurrently to do task 1, do task 2, stopping on error, with limit 2",
			proceed.ConcurrentlyTo(doTask1, doTask2).StoppingOnError().WithLimit(2).String(),
		)
		assert.Equal(
			t,
			"proceed concurrently to do task 1, do task 2, with limit 2, stopping on error",
			proceed.ConcurrentlyTo(doTask1, doTask2).WithLimit(2).StoppingOnError().String(),
		)
	})

	t.Run("simultaneously behaves like concurrently and is narrated simultaneously", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(proceed.SimultaneouslyTo(doTask1, doTask2)))
		assert.Equal(t, []string{"do task 1", "do task 2"}, sortedRecord())
		assert.Equal(
			t,
			"proceed simultaneously to do task 1, do task 2",
			proceed.SimultaneouslyTo(doTask1, doTask2).String(),
		)
	})

	t.Run("ignores nil performables instead of panicking", func(t *testing.T) {
		recorder := fixture.NewRecorder()
		narratedAdam := narratedActor(recorder)

		require.NoError(
			t,
			narratedAdam.AttemptsTo(proceed.ConcurrentlyTo(fixture.NewFakePerformable("retry", nil), nil)),
		)
		assert.Contains(t, recorder.Messages(), "Adam attempts to proceed concurrently to retry")

		require.NoError(
			t,
			narratedAdam.AttemptsTo(proceed.AsynchronouslyTo(nil, fixture.NewFakePerformable("retry", nil))),
		)
		require.NoError(
			t,
			narratedAdam.AttemptsTo(
				proceed.EitherTo(nil, fixture.NewFakePerformable("retry", nil)).
					Or(nil, fixture.NewFakePerformable("wait", nil)),
			),
		)
		// A list made only of nil is empty: the mistake fails when performed, without a panic.
		require.ErrorIs(
			t,
			narratedAdam.AttemptsTo(proceed.EitherTo(fixture.NewFakePerformable("retry", nil)).Or(nil)),
			proceed.ErrNoAlternateProvided,
		)
		require.NoError(
			t,
			narratedAdam.AttemptsTo(proceed.To(nil, fixture.NewFakePerformable("retry", nil)).Otherwise(nil)),
		)
	})

	t.Run("narrates each of its performables", func(t *testing.T) {
		for name, mode := range map[string]func(*impl.ConcurrentlyToAction) *impl.ConcurrentlyToAction{
			"waiting for all":   (*impl.ConcurrentlyToAction).WaitingForAll,
			"ignoring errors":   (*impl.ConcurrentlyToAction).IgnoringErrors,
			"stopping on error": (*impl.ConcurrentlyToAction).StoppingOnError,
		} {
			t.Run(name, func(t *testing.T) {
				recorder := fixture.NewRecorder()
				narratedAdam := narratedActor(recorder)

				concurrently := proceed.ConcurrentlyTo(
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
		err := adam.AttemptsTo(proceed.ConcurrentlyTo(
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
