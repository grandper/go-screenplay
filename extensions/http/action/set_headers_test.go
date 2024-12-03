package action_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/extensions/http/ability"
	"github.com/grandper/go-screenplay/extensions/http/action"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestSetHeaderAction(t *testing.T) {
	t.Run("panics if the header is empty", func(t *testing.T) {
		assert.Panics(t, func() {
			action.SetHeader("", "application/json")
		})
	})

	t.Run("does not panic when the value is empty", func(t *testing.T) {
		assert.NotPanics(t, func() {
			action.SetHeader("Content-Type", "")
		})
	})

	t.Run("panics when a list of headers is empty", func(t *testing.T) {
		assert.Panics(t, func() {
			action.SetHeaders()
		})
	})

	t.Run("panics when the list of header-value pairs is incomplete", func(t *testing.T) {
		assert.Panics(t, func() {
			action.SetHeaders("Content-Type")
		})
		assert.Panics(t, func() {
			action.SetHeaders(
				"Content-Type", "application/json",
				"Authorization")
		})
	})

	t.Run("fails to set headers if the actor does not have the ability MakeHttpRequest", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		require.Error(t, adam.AttemptsTo(action.SetHeader("Content-Type", "application/json")))

		assert.Error(t, adam.AttemptsTo(action.SetHeaders(
			"Content-Type", "application/json",
			"Authorization", "Bearer 84e7750a-582f-4ed7-9510-6e181d530686")))
	})

	t.Run("sets a single header to the session", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())
		assert.NoError(t, adam.AttemptsTo(action.SetHeader("Content-Type", "application/json")))
	})

	t.Run("sets headers to the session", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())
		assert.NoError(t, adam.AttemptsTo(action.SetHeaders(
			"Content-Type", "application/json",
			"Authorization", "Bearer 84e7750a-582f-4ed7-9510-6e181d530686")))
	})

	t.Run("will remove previous headers", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())
		assert.NoError(t, adam.AttemptsTo(action.SetHeaders(
			"Content-Type", "application/json",
			"Authorization", "Bearer 84e7750a-582f-4ed7-9510-6e181d530686")))
		assert.NoError(t, adam.AttemptsTo(action.SetHeader("Content-Type", "application/json")))
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		setHeader := action.SetHeader("Content-Type", "application/json")
		assert.Equal(t, "set the header Content-Type = application/json", setHeader.String())

		setHeaders := action.SetHeaders(
			"Content-Type", "application/json",
			"Authorization", "Bearer 84e7750a-582f-4ed7-9510-6e181d530686")
		assert.Equal(
			t,
			"set the headers Content-Type = application/json, Authorization = Bearer 84e7750a-582f-4ed7-9510-6e181d530686",
			setHeaders.String(),
		)
	})

	t.Run("can be protected from disclosure in logs", func(t *testing.T) {
		setHeader := action.SetHeader("Authorization", "Bearer 84e7750a-582f-4ed7-9510-6e181d530686").
			WhichShouldBeKeptSecret()
		assert.Equal(t, "set the header Authorization = <secret>", setHeader.String())

		setHeader = action.SetHeader("Authorization", "Bearer 84e7750a-582f-4ed7-9510-6e181d530686").Secretly()
		assert.Equal(t, "set the header Authorization = <secret>", setHeader.String())
	})
}
