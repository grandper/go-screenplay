package screenplay_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/fixture"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestActor(t *testing.T) {
	performTesting := performTestingAbility{}
	checkErrors := checkErrorsAbility{}
	flyInTheSky := flyInTheSkyAbility{}

	t.Run("can be created using a name", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		assert.Equal(t, "Adam", adam.Name())
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
