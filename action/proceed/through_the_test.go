package proceed_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action"
	"github.com/grandper/go-screenplay/v2/action/proceed"
	"github.com/grandper/go-screenplay/v2/fixture"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestThroughThe(t *testing.T) {
	record := []string{}
	reset := func() {
		record = []string{}
	}
	adam := screenplay.ActorNamed("Adam")
	removeTheFile := func(fileName string) screenplay.Performable {
		return action.FromFunc("remove the file "+fileName, func(*screenplay.Actor) error {
			if fileName == "b.tmp" {
				return errors.New("file 'b.tmp' does not exist")
			}
			record = append(record, fileName)
			return nil
		})
	}
	checkTheEndpoint := func(endpoint string) *fixture.FakePerformable {
		record = append(record, endpoint)
		return fixture.NewFakePerformable("check the endpoint "+endpoint, nil)
	}

	t.Run("performs the body for each item of the answer, in order", func(t *testing.T) {
		reset()
		temporaryFiles := fixture.NewFakeQuestion("temporary files", []string{"a.tmp", "c.tmp"})
		require.NoError(t, adam.AttemptsTo(proceed.ThroughThe(temporaryFiles).AndForEach(removeTheFile)))
		assert.Equal(t, []string{"a.tmp", "c.tmp"}, record)
	})

	t.Run("performs the body for each item written in the statement", func(t *testing.T) {
		reset()
		endpoints := []string{"/health", "/ready"}
		require.NoError(t, adam.AttemptsTo(proceed.Through(endpoints...).AndForEach(checkTheEndpoint)))
		assert.Equal(t, []string{"/health", "/ready"}, record)
	})

	t.Run("performs nothing when there is no item", func(t *testing.T) {
		reset()
		temporaryFiles := fixture.NewFakeQuestion("temporary files", []string{})
		require.NoError(t, adam.AttemptsTo(proceed.ThroughThe(temporaryFiles).AndForEach(removeTheFile)))
		require.NoError(t, adam.AttemptsTo(proceed.Through[string]().AndForEach(removeTheFile)))
		assert.Equal(t, []string{}, record)
	})

	t.Run("names the pass and the item when the body fails, and stops there", func(t *testing.T) {
		reset()
		temporaryFiles := fixture.NewFakeQuestion("temporary files", []string{"a.tmp", "b.tmp", "c.tmp"})
		err := adam.AttemptsTo(proceed.ThroughThe(temporaryFiles).AndForEach(removeTheFile))
		require.EqualError(
			t,
			err,
			`Adam failed to proceed through the temporary files on the 2nd pass, for "b.tmp": file 'b.tmp' does not exist`,
		)
		assert.Equal(t, []string{"a.tmp"}, record)
	})

	t.Run("returns the error when the question fails to answer", func(t *testing.T) {
		reset()
		temporaryFiles := fixture.NewFailingFakeQuestion[[]string]("temporary files", assert.AnError)
		err := adam.AttemptsTo(proceed.ThroughThe(temporaryFiles).AndForEach(removeTheFile))
		require.ErrorIs(t, err, assert.AnError)
		require.EqualError(t, err, "Adam failed to proceed through the temporary files: "+assert.AnError.Error())
		assert.Equal(t, []string{}, record)
	})

	t.Run("fails when the question or the body is nil", func(t *testing.T) {
		reset()
		require.ErrorIs(
			t,
			adam.AttemptsTo(proceed.ThroughThe[string](nil).AndForEach(removeTheFile)),
			screenplay.ErrInvalidResolution,
		)

		var nilBody func(fileName string) screenplay.Performable
		temporaryFiles := fixture.NewFakeQuestion("temporary files", []string{"a.tmp"})
		err := adam.AttemptsTo(proceed.ThroughThe(temporaryFiles).AndForEach(nilBody))
		require.ErrorIs(t, err, screenplay.ErrInvalidResolution)
		require.EqualError(
			t,
			err,
			"Adam failed to proceed through the temporary files: invalid resolution: the function given to AndForEach is nil",
		)
		assert.Equal(t, []string{}, record)
	})

	t.Run("narrates the loop, then each pass", func(t *testing.T) {
		recorder := fixture.NewRecorder()
		narrated := narratedActor(recorder)

		temporaryFiles := fixture.NewFakeQuestion("temporary files", []string{"a.tmp", "c.tmp"})
		require.NoError(t, narrated.AttemptsTo(proceed.ThroughThe(temporaryFiles).AndForEach(removeTheFile)))

		messages := recorder.Messages()
		assert.Equal(t, "Adam attempts to proceed through the temporary files", messages[0])
		assert.Contains(t, messages, "Adam attempts to remove the file a.tmp")
		assert.Contains(t, messages, "Adam attempts to remove the file c.tmp")
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		temporaryFiles := fixture.NewFakeQuestion("temporary files", []string{})
		assert.Equal(
			t,
			"proceed through the temporary files",
			proceed.ThroughThe(temporaryFiles).AndForEach(removeTheFile).String(),
		)
		assert.Equal(
			t,
			`proceed through "/health", "/ready"`,
			proceed.Through("/health", "/ready").AndForEach(checkTheEndpoint).String(),
		)
		assert.Equal(
			t,
			"proceed through 200, 404",
			proceed.Through(200, 404).AndForEach(func(int) screenplay.Performable { return nil }).String(),
		)
	})
}
