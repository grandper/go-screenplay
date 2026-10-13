// Package fixture holds the fakes to test code built on go-screenplay without a
// real application: a fake performable, a fake question, a fake ability, and a
// recorder of the narration.
//
//	statusCode := fixture.NewFakeQuestion("status code", 200)
//	err := actor.Should(see.The(statusCode).Should(be.EqualTo(200)))
package fixture
