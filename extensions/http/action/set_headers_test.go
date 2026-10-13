package action_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/extensions/http/ability"
	"github.com/grandper/go-screenplay/v2/extensions/http/action"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestSetHeaderAction(t *testing.T) {
	t.Run("fails with invalid headers when the key is empty", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())

		err := adam.AttemptsTo(action.SetHeader("", "application/json"))
		require.ErrorIs(t, err, action.ErrInvalidHeaders)
		require.EqualError(t, err, "invalid headers: a key cannot be empty")
	})

	t.Run("accepts an empty value", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())

		require.NoError(t, adam.AttemptsTo(action.SetHeader("Content-Type", "")))
	})

	t.Run("fails with invalid headers when the list of headers is empty", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())

		err := adam.AttemptsTo(action.SetHeaders())
		require.ErrorIs(t, err, action.ErrInvalidHeaders)
		require.EqualError(t, err, "invalid headers: SetHeaders should receive at least one key and its value")
	})

	t.Run("fails with invalid headers when the list of header-value pairs is incomplete", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())

		err := adam.AttemptsTo(action.SetHeaders("Content-Type"))
		require.ErrorIs(t, err, action.ErrInvalidHeaders)
		require.EqualError(t, err, "invalid headers: SetHeaders should receive a list of key-value pairs")

		err = adam.AttemptsTo(action.SetHeaders("Content-Type", "application/json", "Authorization"))
		require.ErrorIs(t, err, action.ErrInvalidHeaders)
	})

	t.Run("never panics when it is built wrongly", func(t *testing.T) {
		assert.NotPanics(t, func() {
			action.SetHeader("", "application/json")
			action.SetHeaders()
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
