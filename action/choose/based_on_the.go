package choose

import (
	impl "github.com/grandper/go-screenplay/v2/action/internal/choose"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// BasedOnThe starts a multi-way choice on the answer of a question; each branch states what
// the answer happens to be:
//
//	choose.BasedOnThe(statusCode).
//		To(ReadTheResponse()).WhenIt().HappensTo(be.EqualTo(200)).
//		To(AbortTheRequest()).Otherwise()
func BasedOnThe[T any](question screenplay.Question[T]) *impl.BasedOnTheAction[T] {
	return impl.BasedOnThe(question)
}
