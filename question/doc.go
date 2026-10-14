// Package question holds the general-purpose questions an actor can ask,
// whatever the application under test, and the tools to create your own.
//
// FromFunc creates a question from a function, without declaring a dedicated
// type. The answer type is inferred from the function. A question describes
// itself without its article: see.The, AsksFor, Log and Remember add "the"
// themselves, so "answer" is narrated "the answer":
//
//	theAnswer := question.FromFunc("answer", func(actor *screenplay.Actor) (int, error) {
//		return 42, nil
//	})
//	answer, err := actor.AsksFor(theAnswer) // answer is an int
//
// Recalled answers with a value the actor remembers, typed as asked. It is the
// counterpart of the action.Remember action:
//
//	err := actor.Should(see.The(question.Recalled[int]("status code")).Should(be.EqualTo(200)))
//
// Silently mutes the narration of a question, for example to keep a secret out
// of the logs:
//
//	token, err := actor.AsksFor(question.Silently(TheAuthToken()))
//
// The sub-packages hold the questions that wrap another question to answer
// with a part of its answer: first and last (an element of a slice), number
// (the count of elements), and body, data and header (a field of a struct).
package question
