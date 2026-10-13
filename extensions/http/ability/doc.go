// Package ability holds the ability to make HTTP requests:
//
//	anActor := screenplay.ActorNamed("Wanda").WhoCan(ability.MakeHTTPRequests())
//
// The ability is configured fluently, reading "make HTTP requests with timing
// out after 5 seconds with at most 3 redirects":
//
//	anActor := screenplay.ActorNamed("Wanda").WhoCan(
//		ability.MakeHTTPRequests().WithTimingOutAfter(5).Seconds().WithAtMost(3).Redirects(),
//	)
package ability
