// Package ability holds the implementation of the ability to use the file
// system. Only the actions and the questions of the extension use its
// operations; a test gives the ability to an actor through the public package:
//
//	anActor := screenplay.ActorNamed("Boby").WhoCan(ability.UseTheFileSystem())
package ability
