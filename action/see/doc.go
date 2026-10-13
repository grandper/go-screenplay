// Package see holds the actions an actor performs to check what the answers to
// questions should be.
//
// One question is checked with The, closed by a verb and one or more
// resolutions:
//
//	err := actor.Should(see.The(statusCode).Should(be.EqualTo(200)))
//	err := actor.Should(see.The(logs).ShouldNot(contain.Substring("panic")))
//	err := actor.Should(see.The(customers).ShouldHave(length.Exactly[[]string](2)))
//
// The verb sets the narration and the polarity: Should, ShouldBe and
// ShouldHave hold when every resolution holds, ShouldNot, ShouldNotBe and
// ShouldNotHave when none does, and ShouldAnyOf when at least one does. The
// resolutions share the answer type of the question, so a mismatch does not
// compile.
//
// Several questions, whatever their answer types, are checked together with
// AllOfTheseMet, AtLeastOneOfTheseMet and NoneOfTheseMet, which take the terms
// built with the term package:
//
//	err := actor.Should(see.AllOfTheseMet(
//		term.That(statusCode).Should(be.EqualTo(200)),
//		term.That(logs).ShouldNot(contain.Substring("panic")),
//	))
//
// Several questions sharing one answer type and one expectation are checked
// with ThatAllOfThe, ThatAnyOfThe and ThatNoneOfThe:
//
//	err := actor.Should(see.ThatAllOfThe(profileWidget, loginForm).Should(contain.Substring("Adam")))
//
// The sentence written in the code is the one the actor narrates and the one a
// failure reports:
//
//	Adam failed to see that the status code should be equal to 200: got 404
package see
