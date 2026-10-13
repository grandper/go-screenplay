// Package header holds the question that answers with the Header field of the
// struct another question answers with. Go cannot express "a struct with a
// Header field" as a constraint, so the type of the field is given explicitly
// while the type of the wrapped question is inferred:
//
//	err := actor.Should(see.The(header.Of[string](TheLastResponse)).Should(be.EqualTo("created")))
package header
