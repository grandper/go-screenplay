package ability

import "errors"

var (
	// ErrInvalidHTTPMethod is returned when a request is sent with a method
	// that is not a valid HTTP method.
	ErrInvalidHTTPMethod = errors.New("is not a valid HTTP method")

	// ErrResponseIsNil is returned when an HTTPResponse is created from a nil
	// http.Response.
	ErrResponseIsNil = errors.New("the http.Response is nil")
)
