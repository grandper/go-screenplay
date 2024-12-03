package action_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/extensions/http/ability"
	"github.com/grandper/go-screenplay/extensions/http/action"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestSendHTTPRequestAction(t *testing.T) {
	t.Run("sends DELETE request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodDelete, r.Method)
			assert.Equal(t, http.NoBody, r.Body)
		}))
		defer server.Close()
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())
		assert.NoError(t, adam.AttemptsTo(action.SendDeleteRequest().To(server.URL)))
	})

	t.Run("sends GET request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, http.NoBody, r.Body)
		}))
		defer server.Close()
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())
		assert.NoError(t, adam.AttemptsTo(action.SendGetRequest().To(server.URL)))
	})

	t.Run("sends HEAD request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodHead, r.Method)
			assert.Equal(t, http.NoBody, r.Body)
		}))
		defer server.Close()
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())
		assert.NoError(t, adam.AttemptsTo(action.SendHeadRequest().To(server.URL)))
	})

	t.Run("sends OPTIONS request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodOptions, r.Method)
			assert.Equal(t, http.NoBody, r.Body)
		}))
		defer server.Close()
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())
		assert.NoError(t, adam.AttemptsTo(action.SendOptionsRequest().To(server.URL)))
	})

	t.Run("sends PATCH request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPatch, r.Method)
			content, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			assert.Equal(t, "hello world", string(content))
		}))
		defer server.Close()
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())
		assert.NoError(
			t,
			adam.AttemptsTo(action.SendPatchRequest().To(server.URL).WithBody(strings.NewReader("hello world"))),
		)
	})

	t.Run("sends POST request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			content, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			assert.Equal(t, "hello world", string(content))
		}))
		defer server.Close()
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())
		assert.NoError(
			t,
			adam.AttemptsTo(action.SendPostRequest().To(server.URL).WithBody(strings.NewReader("hello world"))),
		)
	})

	t.Run("sends PUT request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPut, r.Method)
			content, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			assert.Equal(t, "hello world", string(content))
		}))
		defer server.Close()
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())
		assert.NoError(
			t,
			adam.AttemptsTo(action.SendPutRequest().To(server.URL).WithBody(strings.NewReader(`hello world`))),
		)
	})

	t.Run("fails to send request if the actor does not have the ability MakeHttpRequest", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		assert.Error(t, adam.AttemptsTo(action.SendGetRequest().To("https://www.google.com")))
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		sendRequest := action.SendHTTPRequest(http.MethodGet).To("https://www.google.com")
		assert.Equal(t, "send a GET request to https://www.google.com", sendRequest.String())

		sendRequest = action.SendHTTPRequest(http.MethodPost).
			To("https://www.google.com").
			WithBody(strings.NewReader("hello world"))
		assert.Equal(t, "send a POST request to https://www.google.com with body 'hello world'", sendRequest.String())
	})

	t.Run("can send request secretly", func(t *testing.T) {
		sendRequest := action.SendHTTPRequest(http.MethodGet).To("https://www.google.com").Secretly()
		assert.Equal(t, "send a secret HTTP request", sendRequest.String())

		sendRequest = action.SendHTTPRequest(http.MethodGet).To("https://www.google.com").WhichShouldBeKeptSecret()
		assert.Equal(t, "send a secret HTTP request", sendRequest.String())
	})
}
