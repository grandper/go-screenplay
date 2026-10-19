package set

import impl "github.com/grandper/go-screenplay/v2/extensions/http/action/internal/set"

// TheHeader sets a header to the actor's HTTP session:
//
//	set.TheHeader("Authorization", "Bearer "+token).Secretly()
//
// This action will remove any headers previously set in the session.
func TheHeader(key, value string) *impl.TheHeadersAction {
	return impl.TheHeader(key, value)
}

// TheHeaders sets several headers to the actor's HTTP session.
// This action will remove any headers previously set in the session.
func TheHeaders(args ...string) *impl.TheHeadersAction {
	return impl.TheHeaders(args...)
}
