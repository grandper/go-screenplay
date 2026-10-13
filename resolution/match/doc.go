// Package match holds the resolutions that complete "should match ...". They
// are bare infinitives and take the verb Should (or ShouldNot):
//
//	err := actor.Should(see.The(orderID).Should(match.PatternString(`^A-\d+$`)))
//
// reads "see the order id should match pattern `^A-\d+$`".
package match
