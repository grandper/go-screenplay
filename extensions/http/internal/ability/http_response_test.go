package ability_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/extensions/http/internal/ability"
	"github.com/grandper/go-screenplay/v2/screenplay"
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
		require.ErrorIs(t, err, ability.ErrResponseIsNil)
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

	t.Run("refuses a body larger than the buffered size", func(t *testing.T) {
		tooLarge := io.LimitReader(zeroReader{}, screenplay.DefaultMaxBufferedSize+1)
		response := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(tooLarge)}

		_, err := ability.NewHTTPResponseFrom(response)
		require.EqualError(t, err, "failed to create a new HTTPResponse: the response body exceeds 64 MiB")
	})

	t.Run("accepts a body of exactly the buffered size", func(t *testing.T) {
		exact := io.LimitReader(zeroReader{}, screenplay.DefaultMaxBufferedSize)
		response := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(exact)}

		httpResponse, err := ability.NewHTTPResponseFrom(response)
		require.NoError(t, err)
		assert.Len(t, httpResponse.Body(), screenplay.DefaultMaxBufferedSize)
	})
}

type fakeReader struct{}

// Read reads up to len(p) bytes into p. It returns the number
// of bytes read (0 <= n <= len(p)) and any error encountered.
func (fr *fakeReader) Read(_ []byte) (int, error) {
	return 0, errors.New("failed to read the body")
}

// zeroReader is an endless source of zero bytes.
type zeroReader struct{}

func (zeroReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}
