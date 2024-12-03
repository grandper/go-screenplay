package ability_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/extensions/http/ability"
)

func TestHTTPResponse(t *testing.T) {
	t.Run("can be created from http.Response", func(t *testing.T) {
		response := &http.Response{
			Status:     "200 OK",
			StatusCode: http.StatusOK,
			Proto:      "HTTP/1.0",
			ProtoMajor: 1,
			ProtoMinor: 0,
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
			Body: io.NopCloser(strings.NewReader(`{"message": "hello world!"}`)),
		}

		httpResponse, err := ability.NewHTTPResponseFrom(response)
		require.NoError(t, err)
		assert.Equal(t, 200, httpResponse.StatusCode())
		assert.Equal(t, map[string]string{"Content-Type": "application/json"}, httpResponse.Headers())
		assert.JSONEq(t, `{"message": "hello world!"}`, httpResponse.Body())
	})

	t.Run("cannot be created from a nil http.Response", func(t *testing.T) {
		httpResponse, err := ability.NewHTTPResponseFrom(nil)
		require.Error(t, err)
		assert.Nil(t, httpResponse)
	})

	t.Run("cannot be created from a http.Response when the body cannot be read", func(t *testing.T) {
		response := &http.Response{
			Status:     "200 OK",
			StatusCode: http.StatusOK,
			Proto:      "HTTP/1.0",
			ProtoMajor: 1,
			ProtoMinor: 0,
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
			Body: io.NopCloser(&fakeReader{}),
		}

		httpResponse, err := ability.NewHTTPResponseFrom(response)
		require.Error(t, err)
		assert.Nil(t, httpResponse)
	})
}

type fakeReader struct{}

// Read reads up to len(p) bytes into p. It returns the number
// of bytes read (0 <= n <= len(p)) and any error encountered.
func (fr *fakeReader) Read(_ []byte) (int, error) {
	return 0, errors.New("failed to read the body")
}
