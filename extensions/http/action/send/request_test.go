package send_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/extensions/http/ability"
	"github.com/grandper/go-screenplay/v2/extensions/http/action/send"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestRequest(t *testing.T) {
	t.Run("sends DELETE request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodDelete, r.Method)
			assert.Equal(t, http.NoBody, r.Body)
		}))
		defer server.Close()
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())
		assert.NoError(t, adam.AttemptsTo(send.DeleteRequest().To(server.URL)))
	})

	t.Run("sends GET request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, http.NoBody, r.Body)
		}))
		defer server.Close()
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())
		assert.NoError(t, adam.AttemptsTo(send.GetRequest().To(server.URL)))
	})

	t.Run("sends HEAD request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodHead, r.Method)
			assert.Equal(t, http.NoBody, r.Body)
		}))
		defer server.Close()
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())
		assert.NoError(t, adam.AttemptsTo(send.HeadRequest().To(server.URL)))
	})

	t.Run("sends OPTIONS request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodOptions, r.Method)
			assert.Equal(t, http.NoBody, r.Body)
		}))
		defer server.Close()
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())
		assert.NoError(t, adam.AttemptsTo(send.OptionsRequest().To(server.URL)))
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
			adam.AttemptsTo(send.PatchRequest().To(server.URL).WithBody(strings.NewReader("hello world"))),
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
			adam.AttemptsTo(send.PostRequest().To(server.URL).WithBody(strings.NewReader("hello world"))),
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
			adam.AttemptsTo(send.PutRequest().To(server.URL).WithBody(strings.NewReader(`hello world`))),
		)
	})

	t.Run("sends request with basic authentication credential", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			assert.True(t, ok)
			assert.Equal(t, "wanda", username)
			assert.Equal(t, "s3cret", password)
		}))
		defer server.Close()
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())
		assert.NoError(
			t,
			adam.AttemptsTo(send.GetRequest().To(server.URL).WithCredential("wanda", "s3cret")),
		)
	})

	t.Run("sends request with basic authentication using the WithAuth alias", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			assert.True(t, ok)
			assert.Equal(t, "wanda", username)
			assert.Equal(t, "s3cret", password)
		}))
		defer server.Close()
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())
		assert.NoError(
			t,
			adam.AttemptsTo(send.GetRequest().To(server.URL).WithAuth("wanda", "s3cret")),
		)
	})

	t.Run("ignores a nil body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.NoBody, request.Body)
		}))
		defer server.Close()
		wanda := screenplay.ActorNamed("Wanda").WhoCan(ability.MakeHTTPRequests())

		sendTheOrder := send.PostRequest().To(server.URL).WithBody(nil)

		assert.Equal(t, fmt.Sprintf("send a POST request to %q", server.URL), sendTheOrder.String())
		require.NoError(t, wanda.AttemptsTo(sendTheOrder))
	})

	t.Run("fails when the body cannot be read", func(t *testing.T) {
		errDiskGone := errors.New("the disk is gone")
		wanda := screenplay.ActorNamed("Wanda").WhoCan(ability.MakeHTTPRequests())

		sendTheOrder := send.PostRequest().To("https://www.example.com").WithBody(iotest.ErrReader(errDiskGone))

		assert.Equal(t,
			`send a POST request to "https://www.example.com" with a body that could not be read`,
			sendTheOrder.String())
		err := wanda.AttemptsTo(sendTheOrder)
		require.ErrorIs(t, err, errDiskGone)
		require.EqualError(t, err, "failed to read the request body: the disk is gone")
	})

	t.Run("fails when the response body is cut short", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Length", "100")
			_, _ = w.Write([]byte("hello"))
		}))
		defer server.Close()
		wanda := screenplay.ActorNamed("Wanda").WhoCan(ability.MakeHTTPRequests())

		err := wanda.AttemptsTo(send.GetRequest().To(server.URL))

		require.ErrorIs(t, err, io.ErrUnexpectedEOF)
		require.EqualError(t, err, "failed to create a new HTTPResponse: unexpected EOF")
	})

	t.Run("fails to send a request to an invalid URL", func(t *testing.T) {
		wanda := screenplay.ActorNamed("Wanda").WhoCan(ability.MakeHTTPRequests())

		err := wanda.AttemptsTo(send.GetRequest().To("http://exa mple.com"))

		require.EqualError(t, err, `parse "http://exa mple.com": invalid character " " in host name`)
	})

	t.Run("fails to send request if the actor does not have the ability MakeHttpRequest", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		assert.Error(t, adam.AttemptsTo(send.GetRequest().To("https://www.google.com")))
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		sendRequest := send.RequestWithMethod(http.MethodGet).To("https://www.google.com")
		assert.Equal(t, "send a GET request to \"https://www.google.com\"", sendRequest.String())

		sendRequest = send.RequestWithMethod(http.MethodPost).
			To("https://www.google.com").
			WithBody(strings.NewReader("hello world"))
		assert.Equal(
			t,
			"send a POST request to \"https://www.google.com\" with body \"hello world\"",
			sendRequest.String(),
		)
	})

	t.Run("can send request secretly", func(t *testing.T) {
		sendRequest := send.RequestWithMethod(http.MethodGet).To("https://www.google.com").Secretly()
		assert.Equal(t, "send a secret HTTP request", sendRequest.String())

		sendRequest = send.RequestWithMethod(http.MethodGet).To("https://www.google.com").WhichShouldBeKeptSecret()
		assert.Equal(t, "send a secret HTTP request", sendRequest.String())
	})

	t.Run("keeps the URL out of the error of a failed secret request", func(t *testing.T) {
		// Nothing listens on port 1, so the request fails with a transport error
		// that would otherwise embed the URL and its query string.
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())

		err := adam.AttemptsTo(send.PostRequest().To("http://127.0.0.1:1/login?token=SECRET").Secretly())
		require.Error(t, err)
		assert.NotContains(t, err.Error(), "SECRET")
		assert.NotContains(t, err.Error(), "/login")
		assert.Contains(t, err.Error(), "failed to send a secret HTTP request")
	})

	t.Run("reports the cause of a failed secret request that has no URL in its error", func(t *testing.T) {
		wanda := screenplay.ActorNamed("Wanda").WhoCan(ability.MakeHTTPRequests())

		err := wanda.AttemptsTo(
			send.RequestWithMethod("FETCH").To("https://api.example.com/orders?token=s3cr3t").Secretly(),
		)

		require.ErrorIs(t, err, ability.ErrInvalidHTTPMethod)
		require.EqualError(t, err, `failed to send a secret HTTP request: "FETCH" is not a valid HTTP method`)
	})

	t.Run("keeps the URL in the error of a failed request that is not secret", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())

		err := adam.AttemptsTo(send.PostRequest().To("http://127.0.0.1:1/login?token=SECRET"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "http://127.0.0.1:1/login?token=SECRET")
	})

	t.Run("refuses a request body larger than the buffered size", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
		defer server.Close()
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())

		tooLarge := io.LimitReader(endlessZeros{}, screenplay.DefaultMaxBufferedSize+1)
		err := adam.AttemptsTo(send.PostRequest().To(server.URL).WithBody(tooLarge))
		require.EqualError(t, err, "failed to read the request body: the request body exceeds 64 MiB")
	})
}

// endlessZeros is an endless source of zero bytes.
type endlessZeros struct{}

func (endlessZeros) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}
