// Package term states what the answer to a question should be, as one typed
// sentence:
//
//	term.That(statusCode).Should(be.EqualTo(200))
//
// reads "term that status code should be equal to 200".
//
// A term captures the answer type of its question when it is built, so a
// resolution of another type does not compile, and terms about questions of
// different types can be handed together to the actions that judge several of
// them: see.AllOfTheseMet, see.AtLeastOneOfTheseMet and see.NoneOfTheseMet.
//
//	err := actor.Should(see.AllOfTheseMet(
//		term.That(statusCode).Should(be.EqualTo(200)),
//		term.That(customers).Should(contain.Element("Adam")),
//	))
//
// # Verbs
//
// A term is closed by a verb, the same ones as see.The. A verb takes one or
// more resolutions and sets the narration and the polarity of the sentence:
//
//	Should        "should ..."            every resolution holds
//	ShouldBe      "should be ..."         every resolution holds
//	ShouldHave    "should have ..."       every resolution holds
//	ShouldNot     "should not ..."        no resolution holds
//	ShouldNotBe   "should not be ..."     no resolution holds
//	ShouldNotHave "should not have ..."   no resolution holds
//	ShouldAnyOf   "should be any of: ..." at least one resolution holds
//
// # Mistakes
//
// Nothing panics when a term is written wrongly. A verb without a resolution,
// a nil resolution or a nil question builds a term that fails with an error
// wrapping ErrInvalidTerm when the actor performs the action it is handed to.
package term
