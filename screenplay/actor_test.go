package screenplay_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/fixture"
	"github.com/grandper/go-screenplay/screenplay"
)

// contextKey is a dedicated type for context keys, avoiding collisions with keys
// defined in other packages.
type contextKey string

// countingPerformable is a performable that counts how many times it is performed.
type countingPerformable struct {
	count int
}

// PerformAs records that the performable has been performed.
func (c *countingPerformable) PerformAs(_ *screenplay.Actor) error {
	c.count++

	return nil
}

// String describes the performable.
func (c *countingPerformable) String() string {
	return "counting performable"
}

// countingPerformable implements the screenplay.Performable interface.
var _ screenplay.Performable = (*countingPerformable)(nil)

func TestActor(t *testing.T) {
	performTesting := performTestingAbility{}
	checkErrors := checkErrorsAbility{}
	flyInTheSky := flyInTheSkyAbility{}

	t.Run("can be created using a name", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		assert.Equal(t, "Adam", adam.Name())
	})

	t.Run("has a default context", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		assert.Equal(t, context.Background(), adam.Context())
	})

	t.Run("can set a custom context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), contextKey("key"), "value")
		adam := screenplay.ActorNamed("Adam").WithContext(ctx)
		assert.Equal(t, ctx, adam.Context())
		assert.NotEqual(t, context.Background(), adam.Context())
	})

	t.Run("can remember and recall information", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		adam.Remember("username", "adam123")
		adam.Remember("is_logged_in", true)

		username, ok := adam.Recall("username").(string)
		require.True(t, ok)
		assert.Equal(t, "adam123", username)

		isLoggedIn, ok := adam.Recall("is_logged_in").(bool)
		require.True(t, ok)
		assert.True(t, isLoggedIn)

		assert.Nil(t, adam.Recall("non_existent_key"))
	})

	t.Run("remembers the answer when the value is a question", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		adam.Remember("status_code", fixture.NewFakeQuestion("status code", 200))

		assert.Equal(t, 200, adam.Recall("status_code"))
	})

	t.Run("remembers nil when the question fails to be answered", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		adam.Remember("status_code", fixture.NewFailingFakeQuestion(
			"status code",
			errors.New("failed to get the status code"),
		))

		assert.Nil(t, adam.Recall("status_code"))
	})

	t.Run("can share information with another actor", func(t *testing.T) {
		t.Run("copies the value to the target actor's memory", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam")
			bob := screenplay.ActorNamed("Bob")

			adam.Remember("session_token", "abc123")
			adam.Share("session_token").With(bob)

			assert.Equal(t, "abc123", bob.Recall("session_token"))
		})

		t.Run("does not remove the value from the source actor", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam")
			bob := screenplay.ActorNamed("Bob")

			adam.Remember("session_token", "abc123")
			adam.Share("session_token").With(bob)

			assert.Equal(t, "abc123", adam.Recall("session_token"))
		})

		t.Run("does nothing when the key does not exist in the source actor", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam")
			bob := screenplay.ActorNamed("Bob")

			adam.Share("unknown_key").With(bob)

			assert.Nil(t, bob.Recall("unknown_key"))
		})

		t.Run("overwrites the value in the target actor if the key already exists", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam")
			bob := screenplay.ActorNamed("Bob")

			bob.Remember("session_token", "old_value")
			adam.Remember("session_token", "new_value")
			adam.Share("session_token").With(bob)

			assert.Equal(t, "new_value", bob.Recall("session_token"))
		})
	})

	t.Run("can forget information", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		adam.Remember("session_token", "abc123")
		assert.Equal(t, "abc123", adam.Recall("session_token"))

		adam.Forget("session_token")
		assert.Nil(t, adam.Recall("session_token"))
	})

	t.Run("does not have ability at first", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		assert.Equal(t, 0, adam.NumAbilities())
	})

	t.Run("has multiple abilities", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		adam.WhoCan(performTesting, checkErrors)
		assert.True(t, adam.HasAbilityTo(performTesting))
		assert.True(t, adam.HasAbilityTo(checkErrors))
		assert.False(t, adam.HasAbilityTo(flyInTheSky))
		assert.Equal(t, 2, adam.NumAbilities())
	})

	t.Run("can use an ability", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		adam.WhoCan(performTesting)
		ability, err := screenplay.UseAbilityTo[performTestingAbility]().Of(adam)
		require.NoError(t, err)
		assert.Equal(t, performTesting, ability)
	})

	t.Run("cannot use a missing ability", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		_, err := screenplay.UseAbilityTo[performTestingAbility]().Of(adam)
		assert.Error(t, err)
	})

	t.Run("forget an ability on exit", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		assert.Equal(t, 0, adam.NumAbilities())
		adam.WhoCan(performTesting)
		assert.Equal(t, 1, adam.NumAbilities())
		require.NoError(t, adam.Exit())
		assert.Equal(t, 0, adam.NumAbilities())
	})

	t.Run("should perform ordered cleanup tasks in order", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		var record []int
		task1 := testOrderedTask{
			id:     1,
			record: &record,
			err:    nil,
		}
		task2 := testOrderedTask{
			id:     2,
			record: &record,
			err:    nil,
		}
		task3 := testOrderedTask{
			id:     3,
			record: &record,
			err:    nil,
		}
		adam.HasOrderedCleanupTasks(task1, task2)
		adam.WithOrderedCleanupTasks(task3)
		require.NoError(t, adam.Exit())
		assert.Equal(t, []int{1, 2, 3}, record)
	})

	t.Run("should stop cleaning if one ordred cleanup task failed", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		var record []int
		task1 := testOrderedTask{
			id:     1,
			record: &record,
			err:    nil,
		}
		task2 := testOrderedTask{
			id:     2,
			record: &record,
			err:    errors.New("failed to perform task 2"),
		}
		task3 := testOrderedTask{
			id:     3,
			record: &record,
			err:    nil,
		}
		adam.HasOrderedCleanupTasks(task1, task2)
		adam.WithOrderedCleanupTasks(task3)
		require.Error(t, adam.Exit())
		assert.Equal(t, []int{1}, record)
	})

	t.Run("should perform independent cleanup tasks in order", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		var record []int
		task1 := testOrderedTask{
			id:     1,
			record: &record,
			err:    nil,
		}
		task2 := testOrderedTask{
			id:     2,
			record: &record,
			err:    nil,
		}
		task3 := testOrderedTask{
			id:     3,
			record: &record,
			err:    nil,
		}
		adam.HasIndependentCleanupTasks(task1, task2)
		adam.WithIndependentCleanupTasks(task3)
		require.NoError(t, adam.Exit())
		assert.Equal(t, []int{1, 2, 3}, record)
	})

	t.Run("should stop continue if one independent cleanup task failed", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		var record []int
		task1 := testOrderedTask{
			id:     1,
			record: &record,
			err:    nil,
		}
		task2 := testOrderedTask{
			id:     2,
			record: &record,
			err:    errors.New("failed to perform task 2"),
		}
		task3 := testOrderedTask{
			id:     3,
			record: &record,
			err:    nil,
		}
		adam.HasIndependentCleanupTasks(task1, task2)
		adam.WithIndependentCleanupTasks(task3)
		require.Error(t, adam.Exit())
		assert.Equal(t, []int{1, 3}, record)
	})

	t.Run("should still run ordered cleanup tasks when an independent cleanup task fails", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		var record []int
		independent := testOrderedTask{id: 1, record: &record, err: errors.New("independent task failed")}
		ordered := testOrderedTask{id: 2, record: &record, err: nil}

		adam.HasIndependentCleanupTasks(independent)
		adam.HasOrderedCleanupTasks(ordered)

		err := adam.Exit()
		require.Error(t, err)
		assert.Equal(t, []int{2}, record, "ordered task must run even though independent task failed")
	})

	t.Run(
		"should not re-run independent cleanup tasks on a second Exit call after partial failure",
		func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam")
			var record []int
			task1 := testOrderedTask{id: 1, record: &record, err: nil}
			task2 := testOrderedTask{id: 2, record: &record, err: errors.New("task 2 failed")}
			task3 := testOrderedTask{id: 3, record: &record, err: nil}

			adam.HasIndependentCleanupTasks(task1, task2, task3)

			require.Error(t, adam.Exit())
			assert.Equal(t, []int{1, 3}, record)

			// Second Exit must not replay any task — the list was cleared despite the error.
			require.NoError(t, adam.Exit())
			assert.Equal(t, []int{1, 3}, record)
		},
	)

	openTheHomePage := fixture.NewFakePerformable("open the home page", nil)
	openTheHomePageButFailed := fixture.NewFakePerformable(
		"open the home page",
		errors.New("the actor failed to perform the task"),
	)

	t.Run("should attempt to perform a task", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		require.NoError(t, adam.AttemptsTo(openTheHomePage))
	})

	t.Run("should return an error when he failed to perform a task", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		require.Error(t, adam.AttemptsTo(openTheHomePageButFailed))
	})

	t.Run("should ignore nil actions", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		require.NoError(t, adam.AttemptsTo(nil))
	})

	t.Run("should perform non-nil actions and ignore the nil ones", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		firstAction := &countingPerformable{}
		secondAction := &countingPerformable{}

		require.NoError(t, adam.AttemptsTo(nil, firstAction, nil, secondAction, nil))
		assert.Equal(t, 1, firstAction.count)
		assert.Equal(t, 1, secondAction.count)
	})

	t.Run("should ignore nil actions but still report the error of a failing action", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		require.Error(t, adam.AttemptsTo(nil, openTheHomePageButFailed))
	})

	thePhoneNumber := fixture.NewFakeQuestion("phone number", "0123456789")
	thePhoneNumberButAnErrorOccurred := fixture.NewFailingFakeQuestion(
		"phone number",
		errors.New("cannot find the phone number"),
	)

	t.Run("should asks a question", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		answer, err := adam.AsksFor(thePhoneNumber)
		require.NoError(t, err)
		assert.Equal(t, "0123456789", answer)
	})

	t.Run("should asks a question and fail to get an answer", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		answer, err := adam.AsksFor(thePhoneNumberButAnErrorOccurred)
		require.Error(t, err)
		assert.Nil(t, answer)
	})

	t.Run("should see something", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		answer, err := adam.Sees(thePhoneNumber)
		require.NoError(t, err)
		assert.Equal(t, "0123456789", answer)
	})

	t.Run("should fail to see something", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		answer, err := adam.Sees(thePhoneNumberButAnErrorOccurred)
		require.Error(t, err)
		assert.Nil(t, answer)
	})
}

func TestActorAttemptsToAliases(t *testing.T) {
	openTheHomePage := fixture.NewFakePerformable("open the home page", nil)
	openTheHomePageButFailed := fixture.NewFakePerformable(
		"open the home page",
		errors.New("the actor failed to perform the task"),
	)

	aliases := map[string]func(*screenplay.Actor, ...screenplay.Performable) error{
		"WasAbleTo": (*screenplay.Actor).WasAbleTo,
		"Does":      (*screenplay.Actor).Does,
		"Did":       (*screenplay.Actor).Did,
		"Will":      (*screenplay.Actor).Will,
		"TriesTo":   (*screenplay.Actor).TriesTo,
		"TriedTo":   (*screenplay.Actor).TriedTo,
		"Tries":     (*screenplay.Actor).Tries,
		"Tried":     (*screenplay.Actor).Tried,
		"Shall":     (*screenplay.Actor).Shall,
		"Should":    (*screenplay.Actor).Should,
	}

	for name, perform := range aliases {
		t.Run(name+" performs the actions like AttemptsTo", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam")
			action := &countingPerformable{}

			require.NoError(t, perform(adam, action))
			assert.Equal(t, 1, action.count)
		})

		t.Run(name+" returns the error when an action fails", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam")
			require.NoError(t, perform(adam, openTheHomePage))
			require.Error(t, perform(adam, openTheHomePageButFailed))
		})
	}
}

func TestActorTimeoutAndPolling(t *testing.T) {
	t.Run("falls back to the default timeout and polling without a production", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		assert.Equal(t, screenplay.DefaultTimeout, adam.Timeout())
		assert.Equal(t, screenplay.DefaultPolling, adam.Polling())
	})

	t.Run("reads the timeout and polling from the production", func(t *testing.T) {
		production := screenplay.NewProduction(
			screenplay.WithTimeout(10*time.Second),
			screenplay.WithPolling(250*time.Millisecond),
		)
		adam := production.SetTheStage(screenplay.CastOfStandardActors()).ActorNamed("Adam")

		assert.Equal(t, 10*time.Second, adam.Timeout())
		assert.Equal(t, 250*time.Millisecond, adam.Polling())
	})
}

func TestActorMuted(t *testing.T) {
	t.Run("shares the memory of the actor it is derived from", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		adam.Remember("greeting", "hello world")

		assert.Equal(t, "hello world", adam.Muted().Recall("greeting"))
	})

	t.Run("still performs the actions it is asked to", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		action := &countingPerformable{}

		require.NoError(t, adam.Muted().AttemptsTo(action))
		assert.Equal(t, 1, action.count)
	})

	t.Run("suppresses narration", func(t *testing.T) {
		recorder := fixture.NewRecorder()
		adam := narratedActor(recorder)

		require.NoError(t, adam.Muted().AttemptsTo(fixture.NewFakePerformable("waves", nil)))
		assert.Empty(t, recorder.Events())
	})
}

// TestActorIsThreadSafe exercises the actor from many goroutines at once. It is
// meant to be run with the race detector (go test -race) to make sure the
// concurrent access to the actor's memory and abilities is properly guarded.
func TestActorIsThreadSafe(_ *testing.T) {
	adam := screenplay.ActorNamed("Adam")

	const goroutines = 50

	var waitGroup sync.WaitGroup
	waitGroup.Add(goroutines)

	for i := range goroutines {
		go func(i int) {
			defer waitGroup.Done()

			key := fmt.Sprintf("key-%d", i)

			adam.Remember(key, i)
			_ = adam.Recall(key)
			adam.Can(flyInTheSkyAbility{})
			_ = adam.HasAbilityTo(flyInTheSkyAbility{})
			_ = adam.NumAbilities()
			adam.Forget(key)
		}(i)
	}

	waitGroup.Wait()
}

type testOrderedTask struct {
	id     int
	record *[]int
	err    error
}

func (tt testOrderedTask) PerformAs(_ *screenplay.Actor) error {
	if tt.err != nil {
		return tt.err
	}
	*tt.record = append(*tt.record, tt.id)
	return nil
}

func (tt testOrderedTask) String() string {
	return "test an ordered task execution"
}

type performTestingAbility struct {
	err error
}

func (tpt performTestingAbility) Forget() error {
	return tpt.err
}

type checkErrorsAbility struct {
	err error
}

func (tce checkErrorsAbility) Forget() error {
	return tce.err
}

type flyInTheSkyAbility struct {
	err error
}

func (fits flyInTheSkyAbility) Forget() error {
	return fits.err
}
