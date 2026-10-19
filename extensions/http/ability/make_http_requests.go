package ability

import impl "github.com/grandper/go-screenplay/v2/extensions/http/internal/ability"

// MakeHTTPRequests enables the actor to make HTTP requests. The ability is
// configured fluently, reading "make HTTP requests with ...":
//
//	MakeHTTPRequests().WithTimingOutAfter(5).Seconds().WithAtMost(3).Redirects()
func MakeHTTPRequests() *impl.MakeHTTPRequestsAbility {
	return impl.MakeHTTPRequests()
}
