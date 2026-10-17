// Package remember holds the action that stores the answer to a question in
// the actor's memory.
//
//	err := actor.AttemptsTo(remember.The(bearerToken).As("session token"))
//
// reads "remember the bearer token as "session token"", which is also what the
// actor narrates. The answer is recalled later with
// actor.Recall[string]("session token").
package remember
