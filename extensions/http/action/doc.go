// Package action holds the actions to drive REST APIs: sending requests and setting their headers.
//
//	err := anActor.AttemptsTo(
//		action.AddHeader("Accept", "application/json"),
//		action.SendGetRequest().To("https://api.example.com/birds"),
//	)
package action
