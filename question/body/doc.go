// Package body holds the question that answers with the Body field of the
// struct another question answers with. Go cannot express "a struct with a
// Body field" as a constraint, so the type of the field is given explicitly
// while the type of the wrapped question is inferred:
//
//	err := actor.Should(see.The(body.Of[string](TheLastResponse)).Should(be.EqualTo("created")))
package body
