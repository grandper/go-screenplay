// Package action holds the errors of the actions that drive REST APIs. The
// actions live in packages named after their verb:
//
//	send.GetRequest().To(url)                    extensions/http/action/send
//	add.TheHeader("Accept", "application/json")  extensions/http/action/add
//	set.TheHeader("Accept", "application/json")  extensions/http/action/set
//
// For example:
//
//	err := anActor.AttemptsTo(
//		add.TheHeader("Accept", "application/json"),
//		send.GetRequest().To("https://api.example.com/birds"),
//	)
package action
