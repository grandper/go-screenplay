// Package screenplay brings the Screenplay Pattern to Go: tests are written
// from the point of view of the users of an application, as actors who perform
// actions and ask questions.
//
// # Actors and abilities
//
// An Actor represents a user of the application. Abilities give an actor the
// means to interact with it (a browser, an HTTP client, a terminal, ...):
//
//	neo := screenplay.ActorNamed("Neo").WhoCan(UseTerminalCommands())
//
// Inside an action or a question, the actor's ability is used with
// UseAnAbilityOf, chaining the operation:
//
//	err := screenplay.UseAnAbilityOf(actor).That(UseTerminalCommands()).Run("ls")
//
// The argument of That carries the type; That returns the actor's own instance
// of it, or a nil pointer when the actor has none, and the operation then
// reports ErrActorMissingAbility. UseAnAbilityOf is written in PerformAs or
// AnsweredBy, never in a test.
//
// # Actions and tasks
//
// An actor interacts with the application by performing actions. An action is
// anything that implements Performable, and a task groups several of them
// under one name with TaskWhere:
//
//	err := neo.AttemptsTo(create.Folder().Named("my_folder"))
//
// AttemptsTo has many aliases (WasAbleTo, Does, Should, ...) so that a test
// reads as plain English.
//
// # Questions and resolutions
//
// The counterpart of assertions are questions and resolutions. A Question[T]
// is answered by an actor with a value of type T, so asking one yields a typed
// answer without any assertion:
//
//	statusCode, err := neo.AsksFor(StatusCode()) // StatusCode() is a Question[int]
//
// A Resolution[T] is the expectation a value of type T should satisfy. A
// question and a resolution are paired by the actions of the action/see
// package, which only accept a resolution of the answer type of the question:
//
//	err := neo.Should(see.The(StatusCode()).Should(be.EqualTo(200)))
//
// The built-in resolutions live under the resolution packages (be, contain,
// length, match, start, end).
//
// # Memory
//
// An actor remembers values under a key and recalls them typed. Recall returns
// an option, which is empty when nothing of the asked type is stored under the
// key:
//
//	neo.Remember("the user id", 1234)
//	userID := neo.Recall[int]("the user id").OrElse(0)
//
// # Cast, stage and production
//
// A Cast provisions actors with abilities, a Stage manages the actors of a
// scenario, and a Production configures the narration, the timeout and the
// polling interval they share:
//
//	production := screenplay.NewProduction(screenplay.WithNarrator(narrator))
//	stage := production.SetTheStage(screenplay.CastWhereEveryoneCan(MakeHTTPRequests))
//	adam := stage.ActorNamed("Adam")
//
// # Given, When, Then
//
// Given, When, Then and And reproduce the Gherkin flow in a test:
//
//	screenplay.Given(adam).WasAbleTo(see.The(AccountBalance).Should(be.EqualTo(2)))
//	screenplay.When(adam).AttemptsTo(Deposit(100).Dollars())
//	screenplay.Then(adam).Should(see.The(AccountBalance).Should(be.EqualTo(102)))
package screenplay
