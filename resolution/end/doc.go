// Package end holds the resolution that completes "should end ...". It is a
// bare infinitive and takes the verb Should (or ShouldNot):
//
//	err := actor.Should(see.The(fileName).Should(end.With(".json")))
//
// reads "see the file name should end with ".json"".
package end
