package proceed_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action"
	"github.com/grandper/go-screenplay/v2/action/proceed"
	"github.com/grandper/go-screenplay/v2/action/see/term"
	"github.com/grandper/go-screenplay/v2/question"
	"github.com/grandper/go-screenplay/v2/resolution/be"
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/grandper/go-screenplay/v2/screenplay/fixture"
)

// UploadThePart is a function of the index that returns a concrete action.
func UploadThePart(index int) *fixture.FakePerformable {
	return fixture.NewFakePerformable(fmt.Sprintf("upload the part %d", index), nil)
}

// SendHTTPRequestNumber is a function of the index whose name holds an acronym.
func SendHTTPRequestNumber(index int) screenplay.Performable {
	return fixture.NewFakePerformable(fmt.Sprintf("send HTTP request number %d", index), nil)
}

func TestRepeatedlyTo(t *testing.T) {
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
	addTheBookToTheCart := createRecordingPerformable("add the book to the cart", nil)

	// queue holds messages that takeAMessage removes, one per pass.
	queue := 0
	queueLength := question.FromFunc("queue length", func(*screenplay.Actor) (int, error) { return queue, nil })
	takeAMessage := action.FromFunc("take a message", func(*screenplay.Actor) error {
		queue--
		record = append(record, "take a message")
		return nil
	})

	t.Run("performs the action exactly the given number of times", func(t *testing.T) {
		reset()
		require.NoError(t, adam.AttemptsTo(proceed.RepeatedlyTo(addTheBookToTheCart).Exactly(3).Times()))
		assert.Equal(
			t,
			[]string{"add the book to the cart", "add the book to the cart", "add the book to the cart"},
			record,
		)

		reset()
		require.NoError(t, adam.AttemptsTo(proceed.RepeatedlyTo(addTheBookToTheCart).Exactly(0).Times()))
		assert.Equal(t, []string{}, record)
	})

	t.Run("gives the index of the pass to a function of the index, starting at 0", func(t *testing.T) {
		reset()
		recordThePart := func(index int) *fixture.FakePerformable {
			record = append(record, fmt.Sprintf("part %d", index))
			return fixture.NewFakePerformable("upload the part", nil)
		}
		require.NoError(t, adam.AttemptsTo(proceed.RepeatedlyTo(recordThePart).Exactly(3).Times()))
		assert.Equal(t, []string{"part 0", "part 1", "part 2"}, record)

		reset()
		recordTheRequest := func(index int) screenplay.Performable {
			record = append(record, fmt.Sprintf("request %d", index))
			return nil
		}
		require.NoError(t, adam.AttemptsTo(proceed.RepeatedlyTo(recordTheRequest).Exactly(2).Times()))
		assert.Equal(t, []string{"request 0", "request 1"}, record)
	})

	t.Run("repeats the action while the condition holds, testing it before each pass", func(t *testing.T) {
		reset()
		queue = 2
		require.NoError(
			t,
			adam.AttemptsTo(proceed.RepeatedlyTo(takeAMessage).WhileThe(queueLength).HappensTo(be.GreaterThan(0))),
		)
		assert.Equal(t, []string{"take a message", "take a message"}, record)

		reset()
		queue = 0
		require.NoError(
			t,
			adam.AttemptsTo(proceed.RepeatedlyTo(takeAMessage).WhileThe(queueLength).HappensTo(be.GreaterThan(0))),
		)
		assert.Equal(t, []string{}, record)
	})

	t.Run("repeats the action until the condition holds", func(t *testing.T) {
		reset()
		queue = 3
		require.NoError(
			t,
			adam.AttemptsTo(proceed.RepeatedlyTo(takeAMessage).UntilThe(queueLength).HappensTo(be.EqualTo(1))),
		)
		assert.Equal(t, []string{"take a message", "take a message"}, record)
	})

	t.Run("performs the action once before the first test with AtLeastOnce", func(t *testing.T) {
		reset()
		queue = 0
		require.NoError(
			t,
			adam.AttemptsTo(
				proceed.RepeatedlyTo(takeAMessage).AtLeastOnce().WhileThe(queueLength).HappensTo(be.GreaterThan(0)),
			),
		)
		assert.Equal(t, []string{"take a message"}, record)

		reset()
		queue = 1
		require.NoError(
			t,
			adam.AttemptsTo(
				proceed.RepeatedlyTo(takeAMessage).AtLeastOnce().UntilThe(queueLength).HappensTo(be.AtMost(1)),
			),
		)
		assert.Equal(t, []string{"take a message"}, record)
	})

	t.Run("names the pass when the action fails", func(t *testing.T) {
		reset()
		passes := 0
		addTheBookUntilItIsOutOfStock := action.FromFunc("add the book to the cart", func(*screenplay.Actor) error {
			passes++
			if passes == 2 {
				return errors.New("the book is out of stock")
			}
			return nil
		})
		err := adam.AttemptsTo(proceed.RepeatedlyTo(addTheBookUntilItIsOutOfStock).Exactly(3).Times())
		require.EqualError(
			t,
			err,
			"Adam failed to proceed repeatedly to add the book to the cart exactly 3 times on the 2nd pass: "+
				"the book is out of stock",
		)
		assert.Equal(t, 2, passes)
	})

	t.Run("writes the rank of the pass in English", func(t *testing.T) {
		for pass, rank := range map[int]string{
			1: "1st", 2: "2nd", 3: "3rd", 4: "4th", 11: "11th", 12: "12th", 13: "13th", 21: "21st", 22: "22nd", 112: "112th",
		} {
			failAtThePass := func(index int) screenplay.Performable {
				if index+1 == pass {
					return fixture.NewFakePerformable("add the book to the cart", assert.AnError)
				}
				return nil
			}
			err := adam.AttemptsTo(proceed.RepeatedlyTo(failAtThePass).Exactly(pass).Times())
			require.ErrorIs(t, err, assert.AnError)
			assert.Contains(t, err.Error(), " on the "+rank+" pass: ")
		}
	})

	t.Run("names the pass when the condition cannot be tested", func(t *testing.T) {
		reset()
		lockedQueue := fixture.NewFailingFakeQuestion[int]("queue length", errors.New("the queue is locked"))
		err := adam.AttemptsTo(proceed.RepeatedlyTo(takeAMessage).WhileThe(lockedQueue).HappensTo(be.GreaterThan(0)))
		require.EqualError(
			t,
			err,
			"Adam failed to proceed repeatedly to take a message while the queue length happens to be greater than 0 "+
				"on the 1st pass: the queue is locked",
		)
		assert.Equal(t, []string{}, record)
	})

	t.Run("fails before the first pass when the condition is built wrongly", func(t *testing.T) {
		for _, loop := range []screenplay.Performable{
			proceed.RepeatedlyTo(takeAMessage).WhileThe(queueLength).HappensTo(),
			proceed.RepeatedlyTo(takeAMessage).AtLeastOnce().UntilThe(queueLength).HappensTo(nil),
			proceed.RepeatedlyTo(takeAMessage).WhileThe[int](nil).HappensTo(be.GreaterThan(0)),
		} {
			reset()
			require.ErrorIs(t, adam.AttemptsTo(loop), term.ErrInvalidTerm)
			assert.Equal(t, []string{}, record)
		}

		reset()
		queue = 1
		err := adam.AttemptsTo(
			proceed.RepeatedlyTo(takeAMessage).AtLeastOnce().WhileThe(queueLength).HappensTo(be.Between(3, 1)),
		)
		require.ErrorIs(t, err, screenplay.ErrInvalidResolution)
	})

	t.Run("fails when the number of times is negative", func(t *testing.T) {
		reset()
		err := adam.AttemptsTo(proceed.RepeatedlyTo(addTheBookToTheCart).Exactly(-1).Times())
		require.ErrorIs(t, err, screenplay.ErrInvalidResolution)
		require.EqualError(
			t,
			err,
			"Adam failed to proceed repeatedly to add the book to the cart exactly -1 times: "+
				"invalid resolution: the number of times -1 is negative",
		)
		assert.Equal(t, []string{}, record)
	})

	t.Run("fails when the body is neither an action nor a function of the index", func(t *testing.T) {
		err := adam.AttemptsTo(proceed.RepeatedlyTo("upload the part").Exactly(3).Times())
		require.ErrorIs(t, err, screenplay.ErrInvalidResolution)
		require.EqualError(
			t,
			err,
			"Adam failed to proceed repeatedly to perform an action exactly 3 times: "+
				"invalid resolution: string is neither an action nor a function of the index",
		)

		var nilAction screenplay.Performable
		require.ErrorIs(
			t,
			adam.AttemptsTo(proceed.RepeatedlyTo(nilAction).Exactly(3).Times()),
			screenplay.ErrInvalidResolution,
		)

		var nilFunction func(index int) *fixture.FakePerformable
		require.ErrorIs(
			t,
			adam.AttemptsTo(proceed.RepeatedlyTo(nilFunction).Exactly(3).Times()),
			screenplay.ErrInvalidResolution,
		)

		require.ErrorIs(
			t,
			adam.AttemptsTo(
				proceed.RepeatedlyTo(func(string) screenplay.Performable { return nil }).Exactly(3).Times(),
			),
			screenplay.ErrInvalidResolution,
		)
	})

	t.Run("narrates the loop, then each pass", func(t *testing.T) {
		recorder := fixture.NewRecorder()
		narrated := narratedActor(recorder)

		require.NoError(t, narrated.AttemptsTo(proceed.RepeatedlyTo(UploadThePart).Exactly(2).Times()))

		messages := recorder.Messages()
		assert.Equal(t, "Adam attempts to proceed repeatedly to upload the part exactly 2 times", messages[0])
		assert.Contains(t, messages, "Adam attempts to upload the part 0")
		assert.Contains(t, messages, "Adam attempts to upload the part 1")
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(
			t,
			"proceed repeatedly to add the book to the cart exactly 3 times",
			proceed.RepeatedlyTo(addTheBookToTheCart).Exactly(3).Times().String(),
		)
		assert.Equal(
			t,
			"proceed repeatedly to add the book to the cart exactly 1 time",
			proceed.RepeatedlyTo(addTheBookToTheCart).Exactly(1).Times().String(),
		)
		assert.Equal(
			t,
			"proceed repeatedly to take a message while the queue length happens to be greater than 0",
			proceed.RepeatedlyTo(takeAMessage).WhileThe(queueLength).HappensTo(be.GreaterThan(0)).String(),
		)
		assert.Equal(
			t,
			"proceed repeatedly to take a message until the queue length happens to be equal to 0",
			proceed.RepeatedlyTo(takeAMessage).UntilThe(queueLength).HappensTo(be.EqualTo(0)).String(),
		)
		assert.Equal(
			t,
			"proceed repeatedly to take a message at least once while the queue length happens to be greater than 0",
			proceed.RepeatedlyTo(takeAMessage).
				AtLeastOnce().
				WhileThe(queueLength).
				HappensTo(be.GreaterThan(0)).
				String(),
		)
		assert.Equal(
			t,
			"proceed repeatedly to take a message at least once until the queue length does not happen to be greater than 0",
			proceed.RepeatedlyTo(takeAMessage).
				AtLeastOnce().
				UntilThe(queueLength).
				DoesNotHappenTo(be.GreaterThan(0)).
				String(),
		)
	})

	t.Run("describes a function of the index with its name", func(t *testing.T) {
		assert.Equal(
			t,
			"proceed repeatedly to upload the part exactly 3 times",
			proceed.RepeatedlyTo(UploadThePart).Exactly(3).Times().String(),
		)
		assert.Equal(
			t,
			"proceed repeatedly to send HTTP request number exactly 3 times",
			proceed.RepeatedlyTo(SendHTTPRequestNumber).Exactly(3).Times().String(),
		)
		assert.Equal(
			t,
			"proceed repeatedly to perform an action exactly 3 times",
			proceed.RepeatedlyTo(func(int) screenplay.Performable { return nil }).Exactly(3).Times().String(),
		)
	})
}
