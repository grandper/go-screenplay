package ability

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// HTTPResponse represents a HTTP response.
type HTTPResponse struct {
	body       string
	headers    map[string]string
	statusCode int
}

// NewHTTPResponseFrom creates a new HTTPResponse from a http.Response.
func NewHTTPResponseFrom(response *http.Response) (*HTTPResponse, error) {
	if response == nil {
		return nil, errors.New("failed to create a new HTTPResponse: the http.Response is nil")
	}
	body, err := io.ReadAll(response.Body)
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
