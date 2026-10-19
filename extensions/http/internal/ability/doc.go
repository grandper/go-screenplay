// Package ability holds the implementation of the ability to make HTTP
// requests. Only the actions and the questions of the extension use its
// operations, through TheAbilityOf; a test gives the ability to an actor
// through the public package:
//
//	anActor := screenplay.ActorNamed("Wanda").WhoCan(ability.MakeHTTPRequests())
package ability
