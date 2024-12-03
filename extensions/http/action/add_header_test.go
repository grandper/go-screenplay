package action_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/extensions/http/ability"
	"github.com/grandper/go-screenplay/extensions/http/action"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestAddHeaderAction(t *testing.T) {
	t.Run("panics if the header is empty", func(t *testing.T) {
		assert.Panics(t, func() {
			action.AddHeader("", "application/json")
		})
	})

	t.Run("does not panic when the value is empty", func(t *testing.T) {
		assert.NotPanics(t, func() {
			action.AddHeader("Content-Type", "")
		})
	})

	t.Run("panics when a list of headers is empty", func(t *testing.T) {
		assert.Panics(t, func() {
			action.AddHeaders()
		})
	})

	t.Run("panics when the list of header-value pairs is incomplete", func(t *testing.T) {
		assert.Panics(t, func() {
			action.AddHeaders("Content-Type")
		})
		assert.Panics(t, func() {
			action.AddHeaders(
				"Content-Type", "application/json",
				"Authorization")
		})
	})

	t.Run("fails to add headers if the actor does not have the ability MakeHttpRequest", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		require.Error(t, adam.AttemptsTo(action.AddHeader("Content-Type", "application/json")))

		assert.Error(t, adam.AttemptsTo(action.AddHeaders(
			"Content-Type", "application/json",
			"Authorization", "Bearer 84e7750a-582f-4ed7-9510-6e181d530686")))
	})

	t.Run("adds a single header to the session", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())
		assert.NoError(t, adam.AttemptsTo(action.AddHeader("Content-Type", "application/json")))
	})

	t.Run("adds headers to the session", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())
		assert.NoError(t, adam.AttemptsTo(action.AddHeaders(
			"Content-Type", "application/json",
			"Authorization", "Bearer 84e7750a-582f-4ed7-9510-6e181d530686")))
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		addHeader := action.AddHeader("Content-Type", "application/json")
		assert.Equal(t, "add the header Content-Type = application/json", addHeader.String())

		addHeaders := action.AddHeaders(
			"Content-Type", "application/json",
			"Authorization", "Bearer 84e7750a-582f-4ed7-9510-6e181d530686")
		assert.Equal(
			t,
			"add the headers Content-Type = application/json, Authorization = Bearer 84e7750a-582f-4ed7-9510-6e181d530686",
			addHeaders.String(),
		)
	})

	t.Run("can be protected from disclosure in logs", func(t *testing.T) {
		addHeader := action.AddHeader("Authorization", "Bearer 84e7750a-582f-4ed7-9510-6e181d530686").
			WhichShouldBeKeptSecret()
		assert.Equal(t, "add the header Authorization = <secret>", addHeader.String())

		addHeaders := action.AddHeader("Authorization", "Bearer 84e7750a-582f-4ed7-9510-6e181d530686").Secretly()
		assert.Equal(t, "add the header Authorization = <secret>", addHeaders.String())
	})
}
