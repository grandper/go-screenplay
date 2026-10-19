package ability

import impl "github.com/grandper/go-screenplay/v2/extensions/http/internal/ability"

var (
	// ErrInvalidHTTPMethod is returned when a request is sent with a method
	// that is not a valid HTTP method.
	ErrInvalidHTTPMethod = impl.ErrInvalidHTTPMethod

	// ErrResponseIsNil is returned when an HTTPResponse is created from a nil
	// http.Response.
	ErrResponseIsNil = impl.ErrResponseIsNil
)
