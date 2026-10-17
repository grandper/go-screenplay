// Package choose holds the actions that perform the first of several branches
// whose condition holds.
//
//	choose.To(action).WhenThe(question).HappensTo(...)...  performs the first branch whose condition holds
//	choose.BasedOnThe(question).To(action).WhenIt()...     switches on the answer to a question
//
// Each branch is closed by a "happens to" verb, and the last one by Otherwise:
//
//	err := actor.AttemptsTo(
//		choose.To(ReadTheResponse()).WhenThe(statusCode).HappensTo(be.EqualTo(200)).
//			To(RetryTheRequest()).WhenThe(statusCode).HappensTo(be.AtLeast(500)).
//			To(AbortTheRequest()).Otherwise(),
//	)
//
// reads "choose to read the response when the status code happens to be equal
// to 200, to retry the request when the status code happens to be at least
// 500, or to abort the request otherwise", which is also what the actor
// narrates.
package choose
