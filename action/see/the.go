package see

import (
	impl "github.com/grandper/go-screenplay/v2/action/internal/see"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// The starts an assertion about the answer to a question. Close it with a verb
// and one or more resolutions:
//
//	see.The(statusCode).Should(be.EqualTo(200))
//	see.The(logs).ShouldNot(contain.Substring("panic"))
//	see.The(customers).ShouldHave(length.Exactly[[]string](2))
//
// The resolutions share the answer type of the question, so a mismatch does
// not compile.
func The[T any](question screenplay.Question[T]) *impl.TheSubject[T] {
	return impl.The(question)
}
