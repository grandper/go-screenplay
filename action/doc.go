// Package action holds the general-purpose actions an actor can perform,
// whatever the application under test.
//
// Some actions control the flow of a scenario:
//
//	Eventually(action)                        retries an action until it succeeds
//	Either(action).Or(alternative)            falls back when an action fails
//	Conditionally(action).If(condition)       performs an action when a condition holds
//	Choose().To(action).When(condition)       performs the first branch whose condition holds
//	ChooseBasedOnThe(question).To(action)...  switches on the answer to a question
//	Concurrently(actions...)                  performs actions in parallel and waits for them
//	Asynchronously(actions...)                launches actions without waiting for them
//	Stop().UntilThe(question).HappensTo(...)  waits until an answer happens to be as stated
//	PauseFor(10).Seconds().Because(reason)    pauses the execution
//
// A condition states a fact, so the conditions about a question are closed by
// the "happens to" verbs (HappensTo, HappensToBe, HappensToHave, their
// negations DoesNotHappenTo..., HappensToBeAnyOf, and their plural forms
// HappenTo...). They take the same resolutions as the "should" verbs of
// see.The:
//
//	err := actor.AttemptsTo(
//		action.Conditionally(RetryTheRequest()).IfThe(StatusCode()).HappensTo(be.AtLeast(500)),
//	)
//
// reads "conditionally retry the request if the status code happens to be at
// least 500", which is also what the actor narrates.
//
// Other actions deal with the answers to questions and with narration:
//
//	Remember(question).As(key)   stores the answer to a question in the actor's memory
//	Silently(action)             mutes the narration of an action
//
// Finally, FromFunc and FromFuncAndQuestions create an action from a function,
// without declaring a dedicated type:
//
//	greet := action.FromFunc("greet the world", func(actor *screenplay.Actor) error {
//		fmt.Printf("Hello from %s\n", actor.Name())
//		return nil
//	})
//
// The actions that check what the answers to questions should be live in the
// action/see package.
package action
