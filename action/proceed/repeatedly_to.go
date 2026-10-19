package proceed

import impl "github.com/grandper/go-screenplay/v2/action/internal/proceed"

// RepeatedlyTo starts an action that performs a body again and again. Close it
// with the clause that says how many passes it makes:
//
//	proceed.RepeatedlyTo(TakeAMessage()).WhileThe(queueLength).HappensTo(be.GreaterThan(0))
//	proceed.RepeatedlyTo(PollTheJob()).UntilThe(jobStatus).HappensTo(be.EqualTo("done"))
//	proceed.RepeatedlyTo(AddTheBookToTheCart()).Exactly(3).Times()
//
// The body is an action, or a function of the index of the pass that returns
// an action, as for index := range 3 does in Go; the index starts at 0:
//
//	proceed.RepeatedlyTo(UploadThePart).Exactly(3).Times()
//
// Any other body builds, and fails with screenplay.ErrInvalidResolution when
// the loop is performed.
func RepeatedlyTo[B any](body B) *impl.RepeatedlyToBuilder[B] {
	return impl.RepeatedlyTo(body)
}
