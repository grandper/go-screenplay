// Package start holds the resolution that completes "should start ...". It is
// a bare infinitive and takes the verb Should (or ShouldNot):
//
//	err := actor.Should(see.The(orderID).Should(start.With("A-")))
//
// reads "see the order id should start with "A-"".
package start
