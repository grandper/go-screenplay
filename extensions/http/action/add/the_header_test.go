package add_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/extensions/http/ability"
	"github.com/grandper/go-screenplay/v2/extensions/http/action"
	"github.com/grandper/go-screenplay/v2/extensions/http/action/add"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestTheHeader(t *testing.T) {
	t.Run("fails with invalid headers when the key is empty", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())

		err := adam.AttemptsTo(add.TheHeader("", "application/json"))
		require.ErrorIs(t, err, action.ErrInvalidHeaders)
		require.EqualError(t, err, "invalid headers: a key cannot be empty")
	})

	t.Run("accepts an empty value", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())

		require.NoError(t, adam.AttemptsTo(add.TheHeader("Content-Type", "")))
	})

	t.Run("fails with invalid headers when the list of headers is empty", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())

		err := adam.AttemptsTo(add.TheHeaders())
		require.ErrorIs(t, err, action.ErrInvalidHeaders)
		require.EqualError(t, err, "invalid headers: add.TheHeaders should receive at least one key and its value")
	})

	t.Run("fails with invalid headers when the list of header-value pairs is incomplete", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())

		err := adam.AttemptsTo(add.TheHeaders("Content-Type"))
		require.ErrorIs(t, err, action.ErrInvalidHeaders)
		require.EqualError(t, err, "invalid headers: add.TheHeaders should receive a list of key-value pairs")

		err = adam.AttemptsTo(add.TheHeaders("Content-Type", "application/json", "Authorization"))
		require.ErrorIs(t, err, action.ErrInvalidHeaders)
	})

	t.Run("never panics when it is built wrongly", func(t *testing.T) {
		assert.NotPanics(t, func() {
			add.TheHeader("", "application/json")
			add.TheHeaders()
		})
	})

	t.Run("fails to add headers if the actor does not have the ability MakeHttpRequest", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		require.Error(t, adam.AttemptsTo(add.TheHeader("Content-Type", "application/json")))

		assert.Error(t, adam.AttemptsTo(add.TheHeaders(
			"Content-Type", "application/json",
			"Authorization", "Bearer 84e7750a-582f-4ed7-9510-6e181d530686")))
	})

	t.Run("adds a single header to the session", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())
		assert.NoError(t, adam.AttemptsTo(add.TheHeader("Content-Type", "application/json")))
	})

	t.Run("adds headers to the session", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())
		assert.NoError(t, adam.AttemptsTo(add.TheHeaders(
			"Content-Type", "application/json",
			"Authorization", "Bearer 84e7750a-582f-4ed7-9510-6e181d530686")))
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		addHeader := add.TheHeader("Content-Type", "application/json")
		assert.Equal(t, "add the header \"Content-Type\" = \"application/json\"", addHeader.String())

		addHeaders := add.TheHeaders(
			"Content-Type", "application/json",
			"Authorization", "Bearer 84e7750a-582f-4ed7-9510-6e181d530686")
		assert.Equal(
			t,
			"add the headers \"Authorization\" = \"Bearer 84e7750a-582f-4ed7-9510-6e181d530686\", \"Content-Type\" = \"application/json\"",
			addHeaders.String(),
		)
	})

	t.Run("can be protected from disclosure in logs", func(t *testing.T) {
		addHeader := add.TheHeader("Authorization", "Bearer 84e7750a-582f-4ed7-9510-6e181d530686").
			WhichShouldBeKeptSecret()
		assert.Equal(t, "add the header \"Authorization\" = <secret>", addHeader.String())

		addHeaders := add.TheHeader("Authorization", "Bearer 84e7750a-582f-4ed7-9510-6e181d530686").Secretly()
		assert.Equal(t, "add the header \"Authorization\" = <secret>", addHeaders.String())
	})
}
