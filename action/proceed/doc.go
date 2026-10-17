// Package proceed holds the actions that control how other actions are
// performed: on a condition, with a fallback, silently, in parallel or in the
// background.
//
//	proceed.To(action).IfThe(question).HappensTo(...)  performs an action when a condition holds
//	proceed.EitherTo(action).Or(alternative)           falls back when an action fails
//	proceed.SilentlyTo(action)                         mutes the narration of an action
//	proceed.ConcurrentlyTo(actions...)                 performs actions in parallel and waits for them
//	proceed.SimultaneouslyTo(actions...)               behaves like ConcurrentlyTo
//	proceed.AsynchronouslyTo(actions...)               launches actions without waiting for them
//
// A condition states a fact, so the conditions about a question are closed by
// the "happens to" verbs (HappensTo, HappensToBe, HappensToHave, their
// negations DoesNotHappenTo..., HappensToBeAnyOf, and their plural forms
// HappenTo...). They take the same resolutions as the "should" verbs of
// see.The:
//
//	err := actor.AttemptsTo(
//		proceed.To(RetryTheRequest()).IfThe(StatusCode()).HappensTo(be.AtLeast(500)),
//	)
//
// reads "proceed to retry the request if the status code happens to be at
// least 500", which is also what the actor narrates.
package proceed
