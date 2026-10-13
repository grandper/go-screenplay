package ability_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action"
	"github.com/grandper/go-screenplay/v2/extensions/http/ability"
	httpaction "github.com/grandper/go-screenplay/v2/extensions/http/action"
	"github.com/grandper/go-screenplay/v2/extensions/http/question"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestMakeHTTPRequestsAbility(t *testing.T) {
	t.Run("should add headers", func(t *testing.T) {
		makeHTTPRequests := ability.MakeHTTPRequests()
		headers := makeHTTPRequests.ToRetrieveHeaders()
		assert.Empty(t, headers)

		makeHTTPRequests.ToAddTheHeader("Content-Type", "application/json")
		headers = makeHTTPRequests.ToRetrieveHeaders()
		assert.Equal(t, map[string]string{
			"Content-Type": "application/json",
		}, headers)

		makeHTTPRequests.ToAddTheHeader("Authorization", "bearer 75b0c854-b34a-41de-92ef-2dd0e1baccd2")
		headers = makeHTTPRequests.ToRetrieveHeaders()
		assert.Equal(t, map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "bearer 75b0c854-b34a-41de-92ef-2dd0e1baccd2",
		}, headers)
	})

	t.Run("should reset the headers", func(t *testing.T) {
		makeHTTPRequests := ability.MakeHTTPRequests()
		makeHTTPRequests.ToAddTheHeader("Content-Type", "application/json")
		makeHTTPRequests.ToAddTheHeader("Authorization", "bearer 75b0c854-b34a-41de-92ef-2dd0e1baccd2")
		headers := makeHTTPRequests.ToRetrieveHeaders()
		assert.Equal(t, map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "bearer 75b0c854-b34a-41de-92ef-2dd0e1baccd2",
		}, headers)

		makeHTTPRequests.ResetHeaders()
		headers = makeHTTPRequests.ToRetrieveHeaders()
		assert.Empty(t, headers)
	})

	t.Run("should send a request to a server", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Equal(t, http.MethodPost, r.Method)
			body, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			defer func() {
				assert.NoError(t, r.Body.Close())
			}()
			assert.JSONEq(t, `{"foo":"bar"}`, string(body))
		}))
		defer server.Close()

		makeHTTPRequests := ability.MakeHTTPRequests()
		makeHTTPRequests.ToAddTheHeader("Content-Type", "application/json")

		assert.Empty(t, makeHTTPRequests.ToRetrieveResponses())
		require.NoError(
			t,
			makeHTTPRequests.ToSend(context.Background(), http.MethodPost, server.URL, strings.NewReader(`{"foo":"bar"}`), nil),
		)

		responses := makeHTTPRequests.ToRetrieveResponses()
		assert.Len(t, responses, 1)
	})

	t.Run("should send a request to a server with credentials", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "Basic dXNlcjpwYXNz", r.Header.Get("Authorization"))
			assert.Equal(t, http.MethodGet, r.Method)
		}))
		defer server.Close()
		credential := &ability.Credential{
			Username: "user",
			Password: "pass",
		}
		makeHTTPRequests := ability.MakeHTTPRequests()
		assert.Empty(t, makeHTTPRequests.ToRetrieveResponses())
		require.NoError(t, makeHTTPRequests.ToSend(context.Background(), http.MethodGet, server.URL, nil, credential))
	})

	t.Run("should fail to send request when the method is wrong", func(t *testing.T) {
		makeHTTPRequests := ability.MakeHTTPRequests()
		assert.Error(t, makeHTTPRequests.ToSend(context.Background(), "foobar", "https://www.google.fr", nil, nil))
	})

	t.Run("should fail to send request when the protocol scheme is unsupported", func(t *testing.T) {
		makeHTTPRequests := ability.MakeHTTPRequests()
		assert.Error(t, makeHTTPRequests.Send(context.Background(), http.MethodGet, "foobar", nil, nil))
	})

	t.Run("should be able to forget", func(t *testing.T) {
		makeHTTPRequests := ability.MakeHTTPRequests()
		require.NoError(t, makeHTTPRequests.Forget())

		assert.Nil(t, makeHTTPRequests.ToRetrieveHeaders())
		assert.Nil(t, makeHTTPRequests.ToRetrieveResponses())
	})

	t.Run("should be described as a string", func(t *testing.T) {
		makeHTTPRequests := ability.MakeHTTPRequests()
		assert.Equal(t, "make HTTP requests", makeHTTPRequests.String())
	})

	t.Run("is safe for concurrent use", func(t *testing.T) {
		// Run with -race: headers are added while requests are sent and
		// responses are read, all at the same time.
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())

		err := adam.AttemptsTo(action.Concurrently(
			httpaction.SendGetRequest().To(server.URL+"/orders"),
			httpaction.AddHeader("Authorization", "Bearer tok-123"),
			httpaction.SendGetRequest().To(server.URL+"/invoices"),
			httpaction.AddHeader("Accept", "application/json"),
			httpaction.SendGetRequest().To(server.URL+"/customers"),
			action.FromFunc("read the responses", func(theActor *screenplay.Actor) error {
				_, err := theActor.AsksFor(question.Responses())
				return err
			}),
		))
		require.NoError(t, err)

		makeHTTPRequests, err := screenplay.UseAbilityTo[*ability.MakeHTTPRequestsAbility]().Of(adam)
		require.NoError(t, err)
		assert.Len(t, makeHTTPRequests.ToRetrieveResponses(), 3)
		assert.Len(t, makeHTTPRequests.ToRetrieveHeaders(), 2)
	})

	t.Run("stops a request when the actor's context expires", func(t *testing.T) {
		// The server never answers: only the actor's deadline can end the request.
		released := make(chan struct{})
		server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			<-released
		}))
		defer server.Close()
		defer close(released)

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests()).WithContext(ctx)

		start := time.Now()
		err := adam.AttemptsTo(httpaction.SendGetRequest().To(server.URL))
		require.ErrorIs(t, err, context.DeadlineExceeded)
		assert.Less(t, time.Since(start), 5*time.Second)
	})

	t.Run("does not forward the session headers to another host", func(t *testing.T) {
		var otherHostHeader http.Header
		otherHost := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			otherHostHeader = r.Header.Clone()
			w.WriteHeader(http.StatusOK)
		}))
		defer otherHost.Close()

		var sameHostHeader http.Header
		var server *httptest.Server
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/to-other-host":
				http.Redirect(w, r, otherHost.URL+"/landing", http.StatusFound)
			case "/to-same-host":
				http.Redirect(w, r, server.URL+"/home", http.StatusFound)
			default:
				sameHostHeader = r.Header.Clone()
				w.WriteHeader(http.StatusOK)
			}
		}))
		defer server.Close()

		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())
		require.NoError(t, adam.AttemptsTo(
			httpaction.AddHeader("X-Api-Key", "secret-key"),
			httpaction.SendGetRequest().To(server.URL+"/to-other-host"),
			httpaction.SendGetRequest().To(server.URL+"/to-same-host"),
		))

		assert.Empty(t, otherHostHeader.Get("X-Api-Key"), "the session header must not reach another host")
		assert.Equal(t, "secret-key", sameHostHeader.Get("X-Api-Key"), "the session header is kept on the same host")
	})

	t.Run("stops after ten redirects by default", func(t *testing.T) {
		var server *httptest.Server
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, server.URL+"/again", http.StatusFound)
		}))
		defer server.Close()

		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())
		err := adam.AttemptsTo(httpaction.SendGetRequest().To(server.URL))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "stopped after 10 redirects")
	})

	t.Run("stops after the number of redirects set with WithAtMost", func(t *testing.T) {
		redirects := 0
		var server *httptest.Server
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			redirects++
			http.Redirect(w, r, server.URL+"/again", http.StatusFound)
		}))
		defer server.Close()

		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests().WithAtMost(3).Redirects())
		err := adam.AttemptsTo(httpaction.SendGetRequest().To(server.URL))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "stopped after 3 redirects")
		assert.Equal(t, 3, redirects, "three requests reach the server before the limit stops the chain")
	})

	t.Run("follows no redirect when WithAtMost is zero", func(t *testing.T) {
		var server *httptest.Server
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, server.URL+"/again", http.StatusFound)
		}))
		defer server.Close()

		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests().WithAtMost(0).Redirects())
		err := adam.AttemptsTo(httpaction.SendGetRequest().To(server.URL))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "stopped after 0 redirects")
	})

	t.Run("times out after the duration set with WithTimingOutAfter", func(t *testing.T) {
		// The server never answers: only the ability's timeout can end the request.
		released := make(chan struct{})
		server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			<-released
		}))
		defer server.Close()
		defer close(released)

		adam := screenplay.ActorNamed("Adam").WhoCan(
			ability.MakeHTTPRequests().WithTimingOutAfter(100).Milliseconds(),
		)

		start := time.Now()
		err := adam.AttemptsTo(httpaction.SendGetRequest().To(server.URL))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Client.Timeout exceeded")
		assert.Less(t, time.Since(start), 5*time.Second)
	})

	t.Run("chains the timeout and the redirect limit", func(t *testing.T) {
		var server *httptest.Server
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, server.URL+"/again", http.StatusFound)
		}))
		defer server.Close()

		adam := screenplay.ActorNamed("Adam").WhoCan(
			ability.MakeHTTPRequests().WithTimingOutAfter(5).Seconds().WithAtMost(2).Redirects(),
		)
		err := adam.AttemptsTo(httpaction.SendGetRequest().To(server.URL))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "stopped after 2 redirects")
	})

	t.Run("does not panic when a header is added after Forget", func(t *testing.T) {
		makeHTTPRequests := ability.MakeHTTPRequests()
		require.NoError(t, makeHTTPRequests.Forget())

		makeHTTPRequests.ToAddTheHeader("X-Api-Key", "secret-key")
		assert.Equal(t, map[string]string{"X-Api-Key": "secret-key"}, makeHTTPRequests.ToRetrieveHeaders())
	})
}
