package add

import impl "github.com/grandper/go-screenplay/v2/extensions/http/action/internal/add"

// TheHeader adds a header to the actor's HTTP session:
//
//	add.TheHeader("Accept", "application/json")
func TheHeader(key, value string) *impl.TheHeadersAction {
	return impl.TheHeader(key, value)
}

// TheHeaders adds several headers to the actor's HTTP session.
func TheHeaders(args ...string) *impl.TheHeadersAction {
	return impl.TheHeaders(args...)
}
