package proceed

import (
	impl "github.com/grandper/go-screenplay/v2/action/internal/proceed"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// ThroughThe starts an action that goes through the items the answer to a
// question holds. Close it with the body performed for each item:
//
//	proceed.ThroughThe(temporaryFiles).AndForEach(remove.TheFile)
func ThroughThe[T any](question screenplay.Question[[]T]) *impl.ThroughTheBuilder[T] {
	return impl.ThroughThe(question)
}

// Through starts an action that goes through the items written in the
// statement. Close it with the body performed for each item:
//
//	proceed.Through(endpoints...).AndForEach(CheckTheEndpoint)
func Through[T any](items ...T) *impl.ThroughTheBuilder[T] {
	return impl.Through(items...)
}
