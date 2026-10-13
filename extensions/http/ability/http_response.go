package ability

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/grandper/go-screenplay/v2/screenplay"
)

// HTTPResponse represents an HTTP response.
type HTTPResponse struct {
	body       string
	headers    map[string]string
	statusCode int
}

// NewHTTPResponseFrom creates a new HTTPResponse from an http.Response.
func NewHTTPResponseFrom(response *http.Response) (*HTTPResponse, error) {
	if response == nil {
		return nil, errors.New("failed to create a new HTTPResponse: the http.Response is nil")
	}
	body, err := readTheBody(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new HTTPResponse: %w", err)
	}
	headers := map[string]string{}
	for key, value := range response.Header {
		headers[key] = strings.Join(value, ",")
	}
	return &HTTPResponse{
		body:       string(body),
		headers:    headers,
		statusCode: response.StatusCode,
	}, nil
}

// readTheBody reads the body whole, up to screenplay.DefaultMaxBufferedSize, and
// fails when the body is larger, so a large download never exhausts the memory
// of the test process.
func readTheBody(body io.Reader) ([]byte, error) {
	content, err := io.ReadAll(io.LimitReader(body, screenplay.DefaultMaxBufferedSize+1))
	if err != nil {
		return nil, err
	}
	if len(content) > screenplay.DefaultMaxBufferedSize {
		return nil, fmt.Errorf("the response body exceeds %d MiB", screenplay.DefaultMaxBufferedSize>>20)
	}
	return content, nil
}

// Body returns the body.
func (hr *HTTPResponse) Body() string {
	return hr.body
}

// Headers returns the headers.
func (hr *HTTPResponse) Headers() map[string]string {
	return hr.headers
}

// StatusCode returns the status code.
func (hr *HTTPResponse) StatusCode() int {
	return hr.statusCode
}
