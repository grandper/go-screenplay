// Package contain holds the resolutions that complete "should contain ...".
// They are bare infinitives and take the verb Should (or ShouldNot):
//
//	err := actor.Should(see.The(logs).Should(contain.Substring("order placed")))
//	err := actor.Should(see.The(customers).Should(contain.Element("Adam")))
//	err := actor.Should(see.The(headers).Should(contain.Key[string, string]("Content-Type")))
//
// The first one reads "see the logs should contain substring "order placed"".
//
// The collection resolutions are typed by their elements: contain.Element("Adam")
// is a resolution on []string. Go only infers type arguments from the arguments
// of a call, so the resolutions that do not receive every type they need take
// them explicitly, as contain.Key[string, string]("Content-Type") does.
package contain
