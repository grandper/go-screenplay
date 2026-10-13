# go-screenplay

[![Test](https://github.com/grandper/go-screenplay/actions/workflows/go-test.yml/badge.svg)](https://github.com/grandper/go-screenplay/actions/workflows/go-test.yml)
[![Lint](https://github.com/grandper/go-screenplay/actions/workflows/go-lint.yml/badge.svg)](https://github.com/grandper/go-screenplay/actions/workflows/go-lint.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/grandper/go-screenplay/v2/screenplay.svg)](https://pkg.go.dev/github.com/grandper/go-screenplay/v2/screenplay)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

`go-screenplay` is a package that brings the *screenplay* pattern to Golang.

**Main features:**

- Implementation of the *Screenplay* Pattern in Golang.
- Utilities to use *Given-When-Then* instruction in tests.
- Natural-language-enabling syntactic sugar
- Questions and resolutions typed with generics: an assertion that cannot hold does not compile.
- Assertions that read as English sentences, and that report failures in the same words.

## Installation

```bash
go get github.com/grandper/go-screenplay/v2
```
The library requires Go 1.27 or later.

## What is the Screenplay Pattern?
The Screenplay Pattern helps
- developing tests that are easy to maintain and that can keep up with a growing code base
that is constantly evolving  to add more features is challenging.
- minimizing the maintenance code, which is useful for developers and testers who try to treat their code for tests as production code.

### The Old Way: The Page Object Pattern
The Page Object Pattern is a simple UI pattern that provide a simple way to test UI code.
The main concept is to create page objects that contains the interaction for a given page.
The problem with the Page Object Pattern is that it lacks scalability because it breaks the Single Responsibility Principle of SOLID.

The basic structure of a test code for page object looks like
- model
  - Domain model classes
- pages
  - Page objects
- Steps
  - Tasks and assertions performed by the user.

### The Screenplay Pattern
This pattern was previously known as "the Journey Pattern" and was originally promoted. It is designed to create
tests from the point of view of the user. It uses a composition of objects that follows the SOLID principles to creates
easy to understand and maintainable tests.

It provides an efficient alternative to the so-called Page Object Pattern.

The pattern organizes the code into several components:
1. Element Locator provide element location details on the UI.
2. Actions present interaction implementation details
3. Tasks provide an object for users.
4. Questions provide information about what the user can see.
5. Actors represents the entity that takes actions and asks questions.

When using the Screenplay Pattern, the test code is usually organized as follows:
- model
  - Domain model classes
- tasks
  - Business tasks
- action
  - UI interactions (sub-tasks)
- pages
  - Page Objects and Elements
- questions
  - Objects used to query the application

## Main Concepts
In this section we go into more details about the different pieces of the Screenplay Pattern.

### Actor
The actor represents a user of your application. It is the glue that acts on actions, questions, and resolutions.
An actor is typically named to give some context about his role. This is also very useful when multiple actors
are expected to interact with the application for a given scenario to differentiate them.

To create an actor use
```go
reviewer := screenplay.ActorNamed("reviewer")
```
An actor can gain additional actions using Abilities.

### Abilities
Abilities brings new feature to your actors.
For example, you can make them use libraries, tools, and resources.

To give your actor an ability you can simply use
```go
neo := screenplay.ActorNamed("Neo").WhoCan(UseTerminalCommands())
```
Ability are then used to perform more actions. For example if you have
an ability that contains database sessions, then you can create an
action `QueryDatabase` that will be able to access the session of an
actor. You can also use ability in Questions.

Inside an action or a question, use the actor's ability in one chained
statement that names the ability, then calls the operation:
```go
func (a *QueryDatabaseAction) PerformAs(actor *screenplay.Actor) error {
	return screenplay.UseAnAbilityOf(actor).That(UseADatabase()).Execute(a.query)
}
```
The argument of `That` is a pointer to an initialized ability, used for its type.
`That` returns the actor's own instance of that type, or a nil pointer when the
actor never learned it. Every method of an ability guards a nil receiver and
returns `screenplay.ErrActorMissingAbility`, so the operation reports the error
and the chain never breaks:
```go
func (db *UseADatabaseAbility) Execute(query string) error {
	if db == nil {
		return screenplay.ErrActorMissingAbility
	}
	_, err := db.session.Exec(query)
	return err
}
```
An actor who never learned the ability fails the action with that error:
```go
eve := screenplay.ActorNamed("Eve")
err := eve.AttemptsTo(QueryDatabase("SELECT 1"))
// actor does not have the required ability
```
`UseAnAbilityOf` is written in `PerformAs` or `AnsweredBy`, never in a test: a
test gives abilities to actors and performs actions; it does not use abilities
itself.

You can also check whether an actor has an ability without using it, with
`actor.HasAbilityTo(UseADatabase())`, and add abilities after creation with
`actor.Can(...)` (an alias of `WhoCan`).

Finally, abilities all implement a `Forget() error` interface that is
used to clean up the of what is inside the ability (e.g., closing
sessions). This method is automatically called when you use
```go
err := neo.Exit()
```

### Actor specific session data
Sometimes we need to store information in a step or task and reuse it in a subsequent one.
To do so, each actor can call the `Remember` and `Recall` methods to store and retrieve
data in a key-value store.
```go
actor.Remember("the user id", 1234)
userID := actor.Recall[int]("the user id").OrElse(0)
```
`Recall` is a generic method: you give it the type of the value you expect, and it hands
the value back typed, wrapped in a `monad.Option` from
[go-monad](https://github.com/grandper/go-monad). The option is empty when nothing is
stored under the key, when `nil` is stored under it, or when the stored value is not of
the type you asked for, so you decide what happens in that case:
```go
userID := actor.Recall[int]("the user id").OrElse(0)                 // fall back to a default value
userID, err := actor.Recall[int]("the user id").OrElseError(ErrNoUser) // or turn the absence into an error
if actor.Recall[int]("the user id").IsPresent() { ... }              // or simply check
```
Use `actor.Recall[any](key)` to retrieve a value whatever its type.

You can also hand `Remember` a question directly, whatever its answer type. The
actor answers it and stores the answer instead of the question. If the question
fails to be answered, `nil` is stored for the key:
```go
actor.Remember("the user id", UserID())
userID := actor.Recall[int]("the user id").OrElse(0)
```
To read a remembered value back as a question, ask the `question.Recalled[T]` question
(see [Remembering answers to questions](#remembering-answers-to-questions)):
```go
userID, err := actor.AsksFor(question.Recalled[int]("the user id"))
```
An actor can share a value with another actor, and forget the data it has stored using
```go
actor.Share("the user id").With(anotherActor)
actor.Forget("the user id")
```
The actor's memory and abilities are safe to access from several goroutines at
once, so `Remember`, `Recall` and `Forget` can be called from actions that run in
parallel. See [Thread safety](#thread-safety) for the details and the caveats.

### Cast
A cast provisions actors with abilities when they enter the stage.
It decouples the creation of actors from the abilities they need, so you only
define the setup once and every actor gets it automatically.

The simplest cast creates actors with no predefined abilities:
```go
cast := screenplay.CastOfStandardActors()
```
When every actor in a scenario needs the same abilities, give the cast their
constructors, without parentheses. Each actor gets its own instances, so what an
actor records through an ability stays its own and its exit forgets its abilities
only:
```go
cast := screenplay.CastWhereEveryoneCan(MakeHTTPRequests).And(UseDatabase)
```
For more control over how each actor is prepared you can supply a function:
```go
cast := screenplay.CastFunc(func(actor *screenplay.Actor) {
	actor.WhoCan(MakeHTTPRequests())
	actor.Remember("role", "reviewer")
})
```

### Stage
The stage manages all actors in a test scenario and tracks which one is
currently in the spotlight. It works together with a cast so that actors
are automatically prepared when they first appear.

To create a stage, provide it with a cast:
```go
stage := screenplay.SetTheStage(screenplay.CastWhereEveryoneCan(MakeHTTPRequests))
```
Use `ActorNamed` to retrieve an actor by name. If the actor does not exist
yet, the stage creates one and lets the cast prepare it. Calling it a second
time with the same name returns the same actor instance (name matching is
case-insensitive). The actor is also placed in the spotlight.
```go
adam := stage.ActorNamed("Adam")
sameAdam := stage.ActorNamed("adam") // returns the same actor
```
You can retrieve the actor currently in the spotlight at any time:
```go
actor, err := stage.TheActorInTheSpotlight()
```
To check whether any actor is currently on stage before accessing the spotlight:
```go
if stage.AnActorIsOnStage() {
	actor, _ := stage.TheActorInTheSpotlight()
	// use actor
}
```
At the end of a scenario, draw the curtain to make every actor exit and clean
up their abilities:
```go
err := stage.DrawTheCurtain()
```

### Actions
An actor will interact with the test under code by performing actions.
Actions can be used to set up your test or execute the operation that
you want to test.

Any action must implement the `Performable` interface:
```go
type Performable interface {
	PerformAs(actor *Actor) error
	String() string
}
```
Usually actions are created using fluent interface pattern to make the
code easier to read. For example, we can have a package `create` that
contains a method `Folder` that will create an action builder that
can be used as follows:
```go
action := create.Folder().Named("my_folder")
```
Now if we want to perform that action we can simply do one of the followings:
```go
bob := screenplay.ActorNamed("Bob")

// All the following functions will give you the name result.
err := bob.AttemptsTo(create.Folder().Named("my_folder"))
err := bob.WasAbleTo(create.Folder().Named("my_folder"))
err := bob.Does(create.Folder().Named("my_folder"))
err := bob.Did(create.Folder().Named("my_folder"))
err := bob.Will(create.Folder().Named("my_folder"))
err := bob.TriesTo(create.Folder().Named("my_folder"))
err := bob.TriedTo(create.Folder().Named("my_folder"))
err := bob.Tries(create.Folder().Named("my_folder"))
err := bob.Tried(create.Folder().Named("my_folder"))
err := bob.Shall(create.Folder().Named("my_folder"))
err := bob.Should(create.Folder().Named("my_folder"))
```
The reason for having all these alias functions available is providing you
the liberty to choose the word that makes the code of your tests easier to read.
The word you choose is also the one the actor narrates (see
[Logging and narration](#logging-and-narration)): `bob.AttemptsTo(...)` is narrated
"Bob attempts to create the folder my_folder", and `bob.WasAbleTo(...)` "Bob was able to
create the folder my_folder".

### Tasks
It happens very often that we need to group several actions together.
We usually do this to avoid repeating ourselves or to give this group
of actions an explicit name.

For example, we may want to create a `Login` task that will group the following actions:
- `NavigateTo.TheHomePage()`
- `Click.On(TheLoginButton)`
- `FillIn.TheField(Username).With(username)`
- `FillIn.TheField(Password).With(password)`

To do so, we can use
```go
login := screenplay.TaskWhere("login to the application", 
	NavigateTo.TheHomePage(), 
	Click.On(TheLoginButton), 
	FillIn.TheField(Username).With(username), 
	FillIn.TheField(Password).With(password))
```
You can also wrap the call in a function to add parameters:
```go
func LoginAs(username, password string) screenplay.Performable {
	return screenplay.TaskWhere("login to the application", 
		NavigateTo.TheHomePage(), 
		Click.On(TheLoginButton), 
		FillIn.TheField(Username).With(username), 
		FillIn.TheField(Password).With(password))
}
```

Tasks implements the `Performable` interface just like any action so that an actor
can use tasks the same way:
```
err := bob.AttemptsTo(LoginAs("bob", "1234"))
```

### Questions and Resolutions
Tests are usually composed of three stages: the setup, the execution of some code,
and finally an assertion.
The counterpart of assertions in the Screenplay Pattern are questions and resolutions.
When an actor asks a question, the answer to this question can be asserted using a resolution.

A question is simply a struct that implements the following generic interface,
where `T` is the type of its answer:
```go
type Question[T any] interface {
	AnsweredBy(actor *Actor) (T, error)
	String() string
}
```
By passing the actor to the `AnsweredBy` method, we are able to use the actor `abilities` to answer the question.
The answer can be of any type, and that type is carried by the question: an actor asking a
`Question[int]` gets an `int` back, without any type assertion, because `AsksFor` (and its alias
`Sees`) is a generic method:
```go
statusCode, err := theActor.AsksFor(StatusCode()) // StatusCode() is a Question[int]
```
A question that is only known at runtime, or that legitimately answers with values of
different types, is simply a `Question[any]`.

The resolution is a struct that implements the following generic interface, where `T` is
the type of the values it inspects:
```go
type Resolution[T any] interface {
	Inspect(value T) error
	String() string
}
```
`Inspect` receives a value of type `T` and returns `nil` if that value fulfills the
expectation, and an error explaining why if it does not. `String` describes the expectation.
It is written to follow the word "should": `be.EqualTo(200)` describes itself as
`be equal to 200`.

As mentioned before questions and resolutions work together: the actor asks a question to
retrieve an answer (i.e., a value of a given type) which is then inspected using the resolution.
The two are joined by a verb such as `Should` or `ShouldNot`, and their types must agree —
`see.The(question)` accepts a `Resolution[T]` when the question is a `Question[T]` — so pairing
a question with a resolution that cannot inspect its answer is a compile-time error rather
than a failed test:
```go
see.The(StatusCode()).Should(be.EqualTo(200))            // Question[int] with Resolution[int]
see.The(StatusCode()).Should(contain.Substring("OK"))    // does not compile: Resolution[string]
```

Note that the resolution provides both which assertion is performed and the expected value.
For example, if we have a question `TitleOfThePage` that returns the title of an HTML page,
and a resolution that checks that the value is equal to "Some great title", then we can use
the action `see.The` to perform our test.
```go
err := theActor.Should(see.The(Text.OfThe(TitleOfThePage)).Should(be.EqualTo("Some great title")))
```
This reads as *"the actor should see the text of the title of the page should be equal to
Some great title"*, and the same sentence is what the actor narrates and what a failure reports:
```
Adam failed to see that the title of the page should be equal to "Some great title": got "Some title"
```

## Usage
Now that we have seen the foundation of the Screenplay Pattern and that we know how it works,
you probably realize that we need a collection of `abilities`, `actions`, `questions`, and `resolutions`.
In this section we will review what the library as already implemented for you.

#### Pausing, Stopping, and Waiting for Answers to Questions
You may want the test execution to stop until the user hit the `Enter` key (mainly for debugging reasons).
In that case you can use the `Stop` action.
```go
err := theActor.Will(Stop())
```
You can also make the execution stop until some resolution is valid:
```go
err := theActor.Will(Stop().UntilThe(HomePage).HappensTo(contain.Substring("Hello World!")))
err := theActor.Will(Stop().UntilThe(CashAccount.Balance()).HappensTo(be.EqualTo(2000)))
```

You may also want to pause the execution for a given number of milliseconds or seconds:
```go
err := theActor.AttemptsTo(PauseFor(10).Seconds().Because("the connection needs time to be setup"))
err := theActor.AttemptsTo(PauseFor(500).Milliseconds().Because("of some obscure reason"))
```
However, if you need to wait until something happens, you may want to use `Stop` instead to save time.

Finally, if you expect an action to fail repeatedly until finally it succeeds, you should use `Eventually`.
```go
err := adam.Should(Eventually(see.The(Text.OfThe(PageTitle)).Should(contain.Substring("Hello World!"))))
err := adam.AttemptsTo(Eventually(Click.OnThe(SaveButton)).TryingEvery(100).Milliseconds())
err := adam.WasAbleTo(Eventually(CancelTheOrder).TryingFor(5).Seconds().PollingEvery(500).Milliseconds())
```
`Eventually` bounds its retry loop with a `timing.Window`: the total time it keeps
trying and the interval it waits between two tries. Both are configured with a
fluent, natural-language vocabulary. To set how long the actor keeps trying use
`For`, `TryingFor`, `TryingForNoLongerThan`, or `WaitingFor`; to set how often it
tries use `Polling`, `PollingEvery`, or `TryingEvery`. Each of them returns a
`timing.DurationBuilder` so you then pick a unit (`Milliseconds`, `Seconds`, ...):
```go
err := adam.WasAbleTo(Eventually(CancelTheOrder).WaitingFor(5).Seconds().Polling(500).Milliseconds())
```

When you do not set the window explicitly, `Eventually` falls back to the actor's
default timeout and polling interval, and `Stop().UntilThe(...)` uses the default
timeout as well. An actor created on its own uses the package defaults
(`screenplay.DefaultTimeout`, 2s, and `screenplay.DefaultPolling`, 500ms). To
change them for a whole scenario, configure the `Production` the actor's stage is
built from with `WithTimeout` and `WithPolling`:
```go
production := screenplay.NewProduction(
    screenplay.WithNarrator(narrator),
    screenplay.WithTimeout(10*time.Second),
    screenplay.WithPolling(250*time.Millisecond),
)
stage := production.SetTheStage(screenplay.CastOfStandardActors())
```
Every actor the production's stage creates then reports those values through
`actor.Timeout()` and `actor.Polling()`, and every `Eventually`/`Stop` they run
inherits them unless the action overrides the window itself.

The production also bounds what one command output, one HTTP body or one file
read may hold in memory. An actor created on its own uses the package default
(`screenplay.DefaultMaxBufferedSize`, 64 MiB). To change it for a whole
scenario, configure the production with `WithMaxBufferedSize`:
```go
production := screenplay.NewProduction(
    screenplay.WithMaxBufferedSize(8 << 20),
)
```
Every actor the production's stage creates then reports the bound through
`actor.MaxBufferedSize()`.

#### Logging and narration
As actors perform actions and ask questions, the library can emit a readable,
hierarchical narrative of what is happening — an idea borrowed from
[ScreenPy's narration model](https://screenpy-docs.readthedocs.io/en/latest/narration.html):
```
Wanda attempts to request the list of birds
    Wanda attempts to set the header 'Accept' to 'application/json'
    Wanda attempts to send a GET request to /birds
Wanda should see that the status code should be equal to 200
Wanda asks for the status code of the last response
    => 200
```
Each line is the name of the actor, the words of the method you called (`AttemptsTo`,
`WasAbleTo`, `Should`, ...), and the description of the action. An action therefore
describes itself in the infinitive, so that it reads after any of them: "send a GET request
to /birds", "see that the status code should be equal to 200".

Narration is **backend independent**: the library does not force a logging
library on you. It exposes a small `screenplay.Adapter` interface and ships two
reference adapters — one over the standard `log` package and one over `log/slog`.

By default **nothing is narrated**: an actor with no adapter attached behaves
exactly as it does without narration, so this feature is entirely opt-in.

##### Attaching a narrator
A `Narrator` is the microphone: it fans every narrated line out to the adapters
attached to it. The recommended way to configure it once, for the whole scenario,
is a `Production` — an explicit configuration root you build at the top of your
program or test:
```go
import (
    "log/slog"

    "github.com/grandper/go-screenplay/v2/narrator/slogadapter"
    "github.com/grandper/go-screenplay/v2/screenplay"
)

narrator := screenplay.NewNarrator(slogadapter.New(slog.Default()))
production := screenplay.NewProduction(screenplay.WithNarrator(narrator))

stage := production.SetTheStage(screenplay.CastWhereEveryoneCan(MakeHTTPRequests))
wanda := stage.ActorNamed("Wanda") // already narrating
```
Every actor the production's stage creates narrates through that same narrator:
each actor keeps a pointer back to the production and reads the narrator from it.
An actor built by hand with `ActorNamed`, outside any production, has no narrator
and stays silent.

##### The reference adapters
The `log` adapter renders an indented, human-readable tree:
```go
import (
    "log"

    "github.com/grandper/go-screenplay/v2/narrator/logadapter"
)

narrator := screenplay.NewNarrator(logadapter.New(log.Default()))
```
The `slog` adapter emits one structured record per step, preserving the actor,
kind, depth and error as attributes:
```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
narrator := screenplay.NewNarrator(slogadapter.New(logger))
```
You can attach several adapters at once; every event is delivered to each of them:
```go
narrator := screenplay.NewNarrator(logadapter.New(log.Default()))
narrator.AttachAdapter(slogadapter.New(slog.Default()))
```

##### Writing a custom adapter
Anything implementing the one-method `screenplay.Adapter` interface works, so you
can render the narrative through `zap`, `zerolog`, an Allure-style reporter, a
`testing.TB` (so the narrative shows up in `go test -v`), an in-memory recorder
for asserting on it in tests, and so on:
```go
type TestingAdapter struct{ t testing.TB }

func (a TestingAdapter) Narrate(event screenplay.Event) {
    if event.Phase == screenplay.PhaseEnd && event.Err == nil && event.Answer == nil {
        return
    }
    a.t.Logf("%s%s", strings.Repeat("  ", event.Depth), event.Message)
}
```
The `fixture` package ships a ready-made `fixture.Recorder` that records every
`screenplay.Event`, which is handy to assert on the narrative a scenario produces.

##### The `Log` action
The `Log` action answers a question and whispers its value as an aside — like
ScreenPy's `Log`, it answers directly (no "asks for" line) and dumps the value
underneath. It honors whatever adapters you configured (with none attached it
stays silent):
```go
err := theActor.AttemptsTo(Log(HowManyBirdsAreInTheSky()))
err := theActor.AttemptsTo(Log(Number.Of(ItemsInTheList)))
```

##### Silencing steps with `Silently`
Wrap a performable with `action.Silently` (or a question with `question.Silently`)
to mute the narration its actions would produce: the nested beats and asides it
would emit never reach the adapters, collapsing a whole sub-tree into a single,
opaque line. It keeps secrets out of the logs — an action that types a password
or sets an authorization header can be silenced so its value never reaches an
adapter — and trims the noise of a verbose task:
```go
err := theActor.AttemptsTo(action.Silently(LogInAs("wanda", secret)))
// narrated: Adam attempts to silently perform an action
answer, err := theActor.AsksFor(question.Silently(TheAuthToken()))
// narrated: Adam asks for a secret
```
Only that one line is narrated: the wrapped action's own description, the
wrapped question's description and its answer are never passed to an adapter.
What remains visible is the error of a failed action, as returned by that
action. `Silently` mutes a private view of the actor rather than the actor
itself, so it is safe under `Concurrently`: a silenced branch never mutes another
branch running at the same time.

**What reaches the adapters.** The narration is the sentence the actor performs,
values included: command arguments, typed text, URLs, request bodies, header
values, logged answers, and the actual value of every unmet expectation
("got ..."). Point the adapters at a sink that may hold those values, or keep a
value out of the narration with `action.Silently`, `question.Silently`, or the
`Secretly()` of header and request actions.

##### Forcing all narration
When debugging, you often want to see every step — including the ones a
`Silently` would hide. Build the production with `WithForceAllNarration` and
every action and question narrates as usual, even the steps a `Silently` would
otherwise mute:
```go
production := screenplay.NewProduction(
    screenplay.WithNarrator(narrator),
    screenplay.WithForceAllNarration(),
)
```
`production.ForceAllNarration()` reports whether it is set.

#### Remembering answers to questions
When you need to ask a question and reuse its answer later, the `Remember` action
stores the answer in the actor's memory under a given key. You can then retrieve it
with `Recall` (see [Actor specific session data](#actor-specific-session-data)):
```go
err := theActor.AttemptsTo(Remember(StatusCode()).As("statusCode"))
statusCode := theActor.Recall[int]("statusCode").OrElse(0)
```
By default the action fails if the question returns a `nil` answer. Use `AllowingNil`
to store it anyway:
```go
err := theActor.AttemptsTo(Remember(OptionalData()).As("data").AllowingNil())
```
Like any other action, it can be combined with `Eventually` to retry until an answer
is available:
```go
err := theActor.AttemptsTo(Eventually(Remember(UserProfile()).As("profile")).TryingFor(10).Seconds())
```

The counterpart of `Remember` is the `question.Recalled[T]` question: it answers with the
value stored under a key, typed as `T` (the zero value of `T` — `nil` for `any` — when
nothing is stored). It fails with `ErrRecalledValueHasWrongType` when the stored value is
not a `T`. Because it is a question, the remembered value can be checked with `see.The`,
waited for with `Eventually`, asked for with a typed answer, or passed anywhere a question
is expected:
```go
err := theActor.Should(see.The(question.Recalled[int]("statusCode")).Should(be.EqualTo(200)))
statusCode, err := theActor.AsksFor(question.Recalled[int]("statusCode"))
```

#### Working with multiple actions
Sometimes your actor needs to do a list of actions. In that case you can simply
provide a list of actions to your actor:
```go
err := theActor.AttemptsTo(DoThis(), DoThat())
```
`AttemptsTo` (and its aliases) ignores `nil` actions. This lets you build the list
of actions conditionally without having to guard against `nil`:
```go
var maybeLogin screenplay.Performable
if needsLogin {
	maybeLogin = Login()
}
err := theActor.AttemptsTo(maybeLogin, OpenTheDashboard())
```

#### Doing actions concurrently
When your actor needs to perform several actions in parallel, use `Concurrently`
(or its alias `Simultaneously`):
```go
err := theActor.AttemptsTo(Concurrently(DoThis(), DoThat(), DoSomethingElse()))
```
By default the actor waits for every action to complete and returns all the errors
that occurred, joined together with `errors.Join`. Calling `WaitingForAll` is optional
and makes this default explicit:
```go
err := theActor.AttemptsTo(Concurrently(DoThis(), DoThat()).WaitingForAll())
```
Use `StoppingOnError` to cancel the remaining actions as soon as one of them fails.
An action still running is canceled through the actor's context: it observes the
cancellation with `actor.Context()`, as described in "Working with contexts". An
action not started yet is skipped:
```go
err := theActor.AttemptsTo(Concurrently(DoThis(), DoThat()).StoppingOnError())
```
Use `IgnoringErrors` to run every action and discard the errors that occurred:
```go
err := theActor.AttemptsTo(Concurrently(DoThis(), DoThat()).IgnoringErrors())
```
Finally, you can cap the number of actions running at the same time with `WithLimit`:
```go
err := theActor.AttemptsTo(Concurrently(DoThis(), DoThat(), DoSomethingElse()).WithLimit(2).WaitingForAll())
```

#### Doing actions asynchronously
When you want the actor to keep going without waiting for the actions to finish,
use `Asynchronously`. Contrary to `Concurrently`, `AttemptsTo` returns immediately
and the actions keep running in the background:
```go
err := theActor.AttemptsTo(Asynchronously(DoThis(), DoThat()))
```
Because the actor does not wait, the errors are not returned by `AttemptsTo`. They
are collected instead, and you can investigate them later with the
`AsynchronousErrors` question. Answering it waits for every pending asynchronous
action to complete and returns the slice of non-nil errors they produced:
```go
theActor.AttemptsTo(Asynchronously(DoThis(), DoThat()))
errs, _ := theActor.AsksFor(AsynchronousErrors()) // errs is a []error
```
Like `Concurrently`, several flavors of execution are available. By default (or
with an explicit `WaitingForAll`) it waits for every action to complete before
collecting the errors. Use `CancelOnError` to cancel the remaining actions as soon
as one of them fails (a running action observes it with `actor.Context()`), or
`IgnoringErrors` to discard the errors altogether:
```go
err := theActor.AttemptsTo(Asynchronously(DoThis(), DoThat()).WaitingForAll())
err := theActor.AttemptsTo(Asynchronously(DoThis(), DoThat()).CancelOnError())
err := theActor.AttemptsTo(Asynchronously(DoThis(), DoThat()).IgnoringErrors())
```
You can also cap the number of actions running at the same time with `WithLimit`:
```go
err := theActor.AttemptsTo(Asynchronously(DoThis(), DoThat(), DoSomethingElse()).WithLimit(2).WaitingForAll())
```

#### Thread safety
When you use `Concurrently` or `Asynchronously`, several actions run in parallel
while sharing the **same** actor. The actor is safe for concurrent use: its
memory (`Remember`, `Recall`, `Forget`) and its abilities (`Can`/`WhoCan`,
`HasAbilityTo`, `UseAnAbilityOf(...).That(...)`) are guarded by an internal
mutex, so those calls will not race even when made from different goroutines at
the same time.

This protects the actor's own state, but it does not automatically make **your**
actions and abilities thread-safe. Keep the following in mind when actions can run
concurrently:

- Any state your actions or abilities own (a shared HTTP client's fields, a
  buffer, a counter, a slice you append to, ...) must be protected by your own
  synchronization if several concurrent actions touch it.
- The value you store with `Remember` is stored as-is. Storing a mutable value
  (a pointer, a map, a slice) and then mutating it from several actions is still a
  data race — the mutex only guards the storing and retrieving, not what you do
  with the value afterwards.
- Reading a value with `Recall` in one action while another action overwrites the
  same key with `Remember` is safe, but the order in which concurrent actions run
  is not guaranteed. Do not rely on one concurrent action seeing the memory
  written by another.

A good way to catch races in your own code is to run your tests with the Go race
detector:
```sh
go test -race ./...
```

#### Trying an action or doing an alternate actions
Sometimes you want the actor to try to do an action, and in case of failure do another one.
You can achieve this using `Either`:
```go
theActor.Will(Either(DoAction()).Or(DoDifferentAction()))
theActor.Will(Either(DoAction()).Otherwise(DoDifferentAction()))
```

#### Performing an action conditionally
Whereas `Either` falls back when an action *fails*, `Conditionally` chooses which
branch to run based on a *condition being true*. Wrap the action (or actions) with
`Conditionally` and attach the condition; only one branch is ever performed:
```go
theActor.AttemptsTo(
	Conditionally(AddHeader("content-type", "application/json")).
		If(theActor.Knows("hostname")).
		Otherwise(GuessTheContentType()),
)
```
The condition can be either a **boolean expression** (`If`, or its alias `When`) or
**what the answer to a question happens to be** (`IfThe`, or its alias `WhenThe`) — the
screenplay-native condition, closed by a verb and one or more resolutions:
```go
theActor.AttemptsTo(
	Conditionally(RetryTheRequest()).
		IfThe(StatusCode).HappensTo(be.AtLeast(500)).
		Otherwise(ReadTheResponse()),
)
```
This reads as *"conditionally retry the request if the status code happens to be at least
500, otherwise read the response"*, which is also what the actor narrates. See
[The verbs of a condition](#the-verbs-of-a-condition) for the full list.
Both forms have a negation, `Unless` and `UnlessThe`:
```go
theActor.AttemptsTo(Conditionally(Log(TheLastResponse)).Unless(response.IsOK()))
theActor.AttemptsTo(Conditionally(LogTheResponse()).UnlessThe(StatusCode).HappensTo(be.EqualTo(200)))
```
The `Otherwise` branch is optional. When it is omitted and the condition does not hold,
the actor does nothing and no error is returned. A boolean condition is evaluated
eagerly (where it is written); use the question form when the condition must
be evaluated at the moment the action is performed.

#### Choosing one action out of many
`Choose` is the multi-way generalization of `Conditionally`: it performs the first branch
whose condition holds. Each branch is a `To(action)` closed by its condition — `When` for
a boolean, `WhenThe` for what the answer to a question happens to be (the multi-way
`If`/`IfThe`) — and the fallback is a `To(action)` closed by `Otherwise()`:
```go
theActor.AttemptsTo(
	Choose().
		To(Log("success")).WhenThe(StatusCode).HappensTo(be.EqualTo(200)).
		To(FollowTheRedirect()).WhenThe(StatusCode).HappensTo(be.EqualTo(404)).
		To(RetryTheRequest()).When(retriesLeft > 0).
		To(AbortTheRequest()).Otherwise(),
)
```
This reads as *"choose to log success when the status code happens to be equal to 200, to
follow the redirect when the status code happens to be equal to 404, to retry the request
when …, or to abort the request otherwise"*. Only the first matching branch runs; there is no fall-through.

When every branch switches on the **same** value, `ChooseBasedOnThe` names the question
once and answers it a single time. Each branch is then closed by `WhenIt()` and a verb:
```go
theActor.AttemptsTo(
	ChooseBasedOnThe(StatusCode).
		To(Log("success")).WhenIt().HappensTo(be.EqualTo(200)).
		To(FollowTheRedirect()).WhenIt().HappensToBeAnyOf(be.EqualTo(301), be.EqualTo(302)).
		To(RetryTheRequest()).WhenIt().HappensTo(be.AtLeast(500)).
		To(AbortTheRequest()).Otherwise(),
)
```
The closing `Otherwise()` is optional; without it, an unmatched choice does nothing and
returns no error. A `When` boolean is evaluated eagerly (where it is written); a `WhenThe`
branch is evaluated when it is reached, so an earlier matching branch means a later
branch's question is never asked.

#### Observing things
The simplest way to ask a question is to use the action `see`.
```go
err := theActor.Should(see.The(Text.OfThe(PageTitle)).Should(start.With("Hello")))
```
A verb takes one or more resolutions, so you can check several things about the same
question in one sentence:
```go
err := theActor.Should(see.The(Text.OfThe(PageTitle)).Should(
			start.With("Hello"),
			end.With("World!")))
```

It is possible to check several questions at once, even when they answer with different
types. Each check is a *term*, built with `term.That` and closed by a verb, and the actor
sees whether all, at least one, or none of the terms are met:
```go
err := theActor.Should(see.AllOfTheseMet(
			term.That(profilWidget).Should(contain.Substring("Adam")),
			term.That(StatusCode()).Should(be.EqualTo(200))))
err := theActor.Should(see.AtLeastOneOfTheseMet(
			term.That(profilWidget).Should(contain.Substring("Adam")),
			term.That(loginForm).Should(contain.Substring("Login"))))
err := theActor.Should(see.NoneOfTheseMet(
			term.That(profilWidget).Should(contain.Substring("error")),
			term.That(StatusCode()).Should(be.EqualTo(500))))
```
A term captures the answer type of its question when it is built, so a resolution of
another type does not compile, and no type assertion is needed anywhere:
```go
term.That(StatusCode()).Should(contain.Substring("OK")) // does not compile: Resolution[string]
```
`see.AllOfTheseMet` checks the terms in order and stops at the first one that is not met.
`see.AtLeastOneOfTheseMet` stops at the first one that is met. `see.NoneOfTheseMet` fails
when a question cannot be answered, because a term that could not be checked is not known
to be unmet.

When the questions share the same answer type and the same expectation, you can check that
all, any, or none of them satisfy it without repeating yourself:
```go
err := theActor.Should(see.ThatAllOfThe(profilWidget, loginForm).Should(contain.Substring("Adam")))
err := theActor.Should(see.ThatAnyOfThe(profilWidget, loginForm).Should(contain.Substring("Adam")))
err := theActor.Should(see.ThatNoneOfThe(profilWidget, loginForm).Should(contain.Substring("Adam")))
```
The questions handed to `ThatAllOfThe`, `ThatAnyOfThe`, and `ThatNoneOfThe` must share
the same answer type (here, both widgets answer with a `string`).

#### Choosing the verb, and negating
`see.The(question)` returns a builder on which you pick the verb that makes the assertion
read as plain English. The verb sets the narration and the polarity, and each family of
resolutions is written to read after one of them:

| Verb | Reads | Takes | Holds when |
| --- | --- | --- | --- |
| `Should` | should … | `be`, `contain`, `match`, `start`, `end` | every resolution holds |
| `ShouldBe` | should be … | complements such as "a valid order id" | every resolution holds |
| `ShouldHave` | should have … | `length` | every resolution holds |
| `ShouldNot` | should not … | same as `Should` | no resolution holds |
| `ShouldNotBe` | should not be … | same as `ShouldBe` | no resolution holds |
| `ShouldNotHave` | should not have … | same as `ShouldHave` | no resolution holds |
| `ShouldAnyOf` | should be any of: … | `be` | at least one resolution holds |

```go
err := theActor.Should(see.The(PageTitle).Should(be.EqualTo("Welcome")))
err := theActor.Should(see.The(TheResponse).ShouldNot(contain.Substring("error")))
err := theActor.Should(see.The(TheItems).ShouldHave(length.AtLeast[[]Item](1)))
err := theActor.Should(see.The(StatusCode()).ShouldAnyOf(be.EqualTo(200), be.EqualTo(201)))
```
The negated verbs simply invert the resolution's outcome, so the same resolution is used
for both the positive and the negative assertion — there are no separate "not" resolutions.
An assertion still fails, whatever its verb, when the question itself cannot be answered.

The same verbs close the terms built with `term.That`.

#### The verbs of a condition
"Should" states an expectation, so it does not read well in a condition, which states a
fact: nobody says "if the status code should be at least 500". The sentences that state a
fact — `Conditionally(...).IfThe`, `WhenThe` and `UnlessThe`, `Choose().WhenThe`,
`ChooseBasedOnThe(...).WhenIt`, and `Stop().UntilThe` — are closed by the "happens to"
verbs. They take the same resolutions as the "should" verbs, and each has a singular and a
plural form so that the sentence agrees with its subject:

| Singular | Plural | Holds when |
| --- | --- | --- |
| `HappensTo` | `HappenTo` | every resolution holds |
| `HappensToBe` | `HappenToBe` | every resolution holds |
| `HappensToHave` | `HappenToHave` | every resolution holds |
| `DoesNotHappenTo` | `DoNotHappenTo` | no resolution holds |
| `DoesNotHappenToBe` | `DoNotHappenToBe` | no resolution holds |
| `DoesNotHappenToHave` | `DoNotHappenToHave` | no resolution holds |
| `HappensToBeAnyOf` | `HappenToBeAnyOf` | at least one resolution holds |

```go
theActor.AttemptsTo(Conditionally(RetryTheRequest()).IfThe(StatusCode).HappensTo(be.AtLeast(500)))
theActor.AttemptsTo(Conditionally(AskForANewPassword()).IfThe(Password).DoesNotHappenToHave(length.AtLeast[string](12)))
theActor.AttemptsTo(Stop().UntilThe(Logs).HappenTo(contain.Substring("ready")))
```
These read *"if the status code happens to be at least 500"*, *"if the password does not
happen to have length at least 12"* and *"stop until the logs happen to contain substring
ready"*. `ChooseBasedOnThe(...).WhenIt()` only offers the singular forms, because its
subject is "it".

#### Checking texts
You can check if a text starts with, ends with, or contains a string:
```go
err := theActor.Should(see.The(Text.OfThe(PageTitle)).Should(start.With("Hello")))
err := theActor.Should(see.The(Text.OfThe(PageTitle)).Should(end.With("World!")))
err := theActor.Should(see.The(Text.OfThe(PageTitle)).Should(contain.Substring("lo Wor")))
```

You can check if a slice of bytes contains another slice of bytes:
```go
err := theActor.Should(see.The(BytesOfThe(LastResponse)).Should(contain.TheBytes([]byte("lo Wor"))))
```

You can match a text exactly:
```go
err := theActor.Should(see.The(Text.OfThe(PageTitle)).Should(be.EqualTo("Hello World!")))
```

You can use regex to match a text:
```go
err := theActor.Should(see.The(Text.OfThe(PageTitle)).Should(match.PatternString(`^Hello \w+`)))
pattern := regexp.MustCompile(`^Hello \w+`)
err := theActor.Should(see.The(Text.OfThe(PageTitle)).Should(match.Pattern(pattern)))
```

#### Checking numbers
You can assess the value of numbers as follows:
```go
err := theActor.Should(see.The(Number.Of(items.InThe(TodoList))).Should(be.EqualTo(0)))
err := theActor.Should(see.The(Number.Of(items.InThe(TodoList))).Should(be.LessThan(1)))
err := theActor.Should(see.The(Number.Of(items.InThe(TodoList))).Should(be.AtMost(1)))
err := theActor.Should(see.The(Number.Of(items.InThe(TodoList))).Should(be.GreaterThan(1)))
err := theActor.Should(see.The(Number.Of(items.InThe(TodoList))).Should(be.AtLeast(1)))
tolerance := 25.0
err := theActor.Should(see.The(Number.Of(items.InThe(TodoList))).Should(be.CloseTo(101, tolerance)))
err := theActor.Should(see.The(Number.Of(items.InThe(TodoList))).Should(be.Between(1, 5)))
```
`be.EqualTo` is exact on the type: a question answering with an `int64` needs
`be.EqualTo(int64(200))`. Use `be.Zero[T]()` to check that a value equals its zero value
(a `false` boolean, the number zero, an empty string, ...), and `be.True()` and `be.False()`
for a question answering with a boolean.

#### Checking collections
There are a couple of collection types in Golang such as slices and maps.

The resolutions that inspect a collection are typed by its elements, so the compiler knows
which questions they pair with and nothing is checked by reflection at run time.

There exists a couple of resolutions to check slices:
```go
err := theActor.Should(see.The(items.InThe(TodoList)).Should(contain.Element("Add tests for the Go package")))
err := theActor.Should(see.The(items.InThe(TodoList)).Should(be.DeepEqualTo([]string{"Write code", "Add tests"})))
err := theActor.Should(see.The(items.InThe(TodoList)).Should(be.Empty[[]string]()))
err := theActor.Should(see.The(items.InThe(TodoList)).ShouldHave(length.Exactly[[]string](5)))
err := theActor.Should(see.The(items.InThe(TodoList)).ShouldHave(length.Between[[]string](1, 5)))
err := theActor.Should(see.The(Text.OfAll(items.InThe(TodoList))).Should(contain.TheItemMatching(pattern)))
```
`contain.Element("...")` is a resolution on a `[]string`: the type of the slice is inferred
from the element. Use `be.DeepEqualTo` to compare slices, maps and the structs holding them,
which Go cannot compare with `==`. The `length` package also provides `length.AtLeast` and
`length.AtMost`, and works on channels, maps, slices, arrays and strings.

Here are a couple of options to test maps:
```go
err := theActor.Should(see.The(HeadersOf(TheLastResponse)).Should(contain.Key[string, string]("Content-Type")))
err := theActor.Should(see.The(HeadersOf(TheLastResponse)).Should(contain.Value[string, string]("application/json")))
err := theActor.Should(see.The(HeadersOf(TheLastResponse)).Should(
			contain.Entries(map[string]string{"Content-Type": "application/json"})))
```
Go only infers type arguments from the arguments of a call, never from the question a
resolution is later paired with. A resolution that does not receive a value of every type
it needs therefore takes its types explicitly: `contain.Key[string, string]("Content-Type")`
names the key and value types of the map, and `length.Exactly[[]string](5)` and
`be.Empty[[]string]()` name the type of the collection. The other resolutions need no type
argument: `be.EqualTo`, `be.GreaterThan`, `be.Between`, `contain.Element` and
`contain.Entries` take their type from the value they are built with, and
`contain.Substring`, `start.With`, `end.With`, and `match.Pattern` inspect a `string`
(`contain.TheBytes` a `[]byte`).

A named slice type such as `type Names []string` is a different type from `[]string`, so a
question that should pair with these resolutions answers with a plain slice or map.

When a question answers with a slice, you may want to assert against a single element of it
rather than the whole collection. The `first` and `last` questions wrap a question answering
with a slice and answer with, respectively, its first and last element — typed as the
element of the slice, so `last.Of(Responses())` answers with a `*HTTPResponse` when
`Responses()` answers with a `[]*HTTPResponse`:
```go
err := theActor.Should(see.The(first.Of(items.InThe(TodoList))).Should(be.EqualTo("Add tests for the Go package")))
err := theActor.Should(see.The(last.Of(items.InThe(TodoList))).Should(be.EqualTo("by the end of the year")))
```
Wrapping a question that does not answer with a slice is a compile-time error. Both fail
with `ErrAnswerIsEmpty` if the wrapped question answers with an empty slice.

When a question answers with a struct, you may want to assert against one of its fields rather
than the whole value. The `body`, `data`, and `header` questions wrap another question and answer
with, respectively, the `Body`, `Data`, and `Header` field of the struct it returns. Go cannot
express "a struct with a `Body` field" as a constraint, so the type of the field is given as the
first type argument while the type of the wrapped question is inferred:
```go
err := theActor.Should(see.The(body.Of[string](TheLastResponse)).Should(be.EqualTo("created")))
err := theActor.Should(see.The(data.Of[string](last.Of(NATSMessages))).Should(be.EqualTo("payload")))
err := theActor.Should(see.The(header.Of[string](last.Of(HTTPResponses))).Should(be.EqualTo("application/json")))
body, err := theActor.AsksFor(body.Of[string](TheLastResponse)) // body is a string
```
Use `body.Of[any](...)` when the field may be of any type. They fail with `ErrAnswerIsNotAStruct`
if the wrapped question does not answer with a struct (or a pointer to one), with
`ErrFieldNotFound` if the struct has no matching field, and with `ErrFieldHasWrongType` if the
field is not of the given type.

When a question answers with a slice or a map, you may want to assert against the number of
elements it contains. The `number` question wraps another question and answers with the count
of elements in the slice or map it returns:
```go
err := theActor.Should(see.The(number.Of(items.InThe(TodoList))).Should(be.EqualTo(2)))
```
It answers with an `int` and fails with `ErrAnswerIsNotACollection` if the wrapped question does
not answer with a slice or a map.

#### Mistakes in a test never panic
A sentence or a resolution can be written wrongly: a verb without a resolution, a pattern
that does not compile, a lower bound greater than the upper one. The library never panics
on such a mistake. Building always succeeds, and the mistake is reported as an error when
the actor performs the action, so it fails the test like any other failure and points at
what is wrong:
```go
err := theActor.Should(see.The(OrderID()).Should(match.PatternString("A-(")))
// invalid resolution "match pattern `A-(`": error parsing regexp: missing closing ): `A-(`

err := theActor.Should(see.The(StatusCode()).Should())
// invalid term: Should needs at least one resolution
```
Two errors tell a mistake in the test apart from an expectation that is not met. You can
check for them with `errors.Is`:

| Error | Reported when |
| --- | --- |
| `screenplay.ErrInvalidResolution` | a resolution was built wrongly: `match.PatternString("A-(")`, `match.Pattern(nil)`, `be.Between(10, 1)`, `be.CloseTo(3.0, -0.1)`, `be.Empty[int]()`, `length.Exactly[int](1)`, `length.Exactly[string](-1)`, `resolution.FromFunc[int]("...", nil)` |
| `term.ErrInvalidTerm` | a sentence was built wrongly: a verb without a resolution, a `nil` resolution, a `nil` question, or a `nil` term |

The extensions follow the same rule. For example, the HTTP actions `AddHeaders` and
`SetHeaders` report an empty key or an incomplete list of key-value pairs with an error
wrapping their `ErrInvalidHeaders`, when the actor performs them.

A mistake fails the action **whatever the verb**. This matters for the negative verbs: an
error returned by a resolution normally means "not satisfied", which `ShouldNot` turns into
a success. A resolution that was built wrongly checked nothing, so it must not pass:
```go
broken := match.PatternString("A-(")

err := theActor.Should(see.The(OrderID()).Should(broken))    // fails: invalid resolution
err := theActor.Should(see.The(OrderID()).ShouldNot(broken)) // fails too: nothing was checked
```
The same holds everywhere a question meets a resolution. `see.AllOfTheseMet`,
`see.AtLeastOneOfTheseMet` and `see.NoneOfTheseMet` report a wrongly built term before
judging any of them, `Conditionally` and `Choose` perform no branch, and `Stop().UntilThe`
fails at once rather than waiting for its timeout.

A resolution of your own reports that it was built wrongly the same way: its `Inspect`
returns an error wrapping `screenplay.ErrInvalidResolution`. Any other error means that the
value does not satisfy the resolution:
```go
func Matching(pattern string) screenplay.Resolution[string] {
	compiled, compileErr := regexp.Compile(pattern)

	return resolution.FromFunc("match pattern "+pattern, func(text string) error {
		if compileErr != nil {
			return fmt.Errorf("%w: %w", screenplay.ErrInvalidResolution, compileErr) // a mistake: every verb fails
		}

		if !compiled.MatchString(text) {
			return fmt.Errorf("got %q", text) // not satisfied: ShouldNot passes
		}

		return nil
	})
}
```

#### Creating an Anonymous Resolution
It happens sometimes that we need a one-off resolution without creating a dedicated struct.
We usually do this for simple or ad-hoc assertions that are not worth extracting into a reusable type.

For example, we may want to assert that a value satisfies a custom condition inline. Use
`resolution.FromFunc` with a description and a function that has the signature of `Inspect`:
it returns `nil` when the value satisfies the resolution, and an error explaining why when it
does not. The type the resolution inspects is inferred from the function:
```go
positive := resolution.FromFunc("be positive", func(n int) error {
	if n <= 0 {
		return fmt.Errorf("got %d", n)
	}

	return nil
})
```
You can also wrap the call in a function to add parameters:
```go
func DivisibleBy(divisor int) screenplay.Resolution[int] {
	return resolution.FromFunc(fmt.Sprintf("be divisible by %d", divisor), func(value int) error {
		if value%divisor != 0 {
			return fmt.Errorf("got %d", value)
		}

		return nil
	})
}
```
Write the description so that it reads after the verb you intend to use: `"be positive"`
after `Should`, `"a valid order id"` after `ShouldBe`, `"three retries left"` after
`ShouldHave`.

Anonymous resolutions implement the `Resolution[T]` interface just like any other resolution so that an actor
can use them the same way:
```go
err := theActor.Should(see.The(CashAccount.Balance()).Should(DivisibleBy(100)))
```

#### Working with contexts
Contexts are hold by actors to share information across different part of a sequence of calls.
```go
ctx := context.Background()
actor := screenplay.ActorNamed("Alice").WithContext(ctx)
```
`WithContext` returns a view of the actor with that context, like
`http.Request.WithContext`: the actor it was called on is untouched, and so is
every other goroutine holding it. The view shares the actor's memory, abilities
and cleanup tasks. The context can later be retrieved using
```go
ctx := actor.Context()
```
When no context has been set, `actor.Context()` returns `context.Background()`, so
it is always safe to call.

##### Context propagation and cancellation
The actor's context is how cancellation and deadlines are propagated to the
actions it performs. Give the actor a cancellable context (for example one derived
from `context.WithTimeout` or `context.WithCancel`, or the one provided by your
test framework), and every action can observe it through `actor.Context()`:
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
actor := screenplay.ActorNamed("Alice").WithContext(ctx)
```
Actions are expected to respect this context. If an action does long-running work
(a network call, a loop, polling, ...), it should check
`actor.Context().Done()` periodically and stop early when the context is
cancelled, returning `actor.Context().Err()`:
```go
func (a *MyAction) PerformAs(actor *screenplay.Actor) error {
    for _, item := range a.items {
        select {
        case <-actor.Context().Done():
            return actor.Context().Err()
        default:
        }

        if err := process(item); err != nil {
            return err
        }
    }

    return nil
}
```
Any operation that already accepts a `context.Context` (HTTP requests, database
queries, ...) should be handed `actor.Context()` so that cancellation flows all
the way down. This matters in particular for `Concurrently` with
`StoppingOnError` and for `Asynchronously` with `CancelOnError`: the remaining
actions are cancelled through the context, so they only actually stop if they are
watching it.

## Using Gherkin-like syntax to write tests
You are maybe familiar with the Gherkin syntax that allow you to write
tests using the _Given-When-Then_ structure. Each keyword as a specific role:
* `Given` indicates that what follows is used to setup the test.
* `When` precedes the actions under test.
* `Then` precedes the assertions of the test.
The language also include the keyword `And` to avoid repeating any of the above
when multiple actions need to happen.

Here's an example:
> Given Adam's account has a balance of $2
> When Adam deposits $100
> Then the account balance is $102

The library provide some functions to reproduce the Gherkin flow in your test:
```go
adam := screenplay.ActorNamed("Adam")
screenplay.Given(adam).WasAbleTo(see.The(AccountBalance).Should(be.EqualTo(2)))
screenplay.When(adam).AttemptsTo(Deposit(100).Dollars())
screenplay.Then(adam).Should(see.The(AccountBalance).Should(be.EqualTo(102)))
```

### Creating your own Performable/Task/Action
The quickest way to create a performable is `action.FromFunc`. It wraps a plain
function of type `func(actor *screenplay.Actor) error` into a `Performable`, using
the `description` you provide as the result of its `String` method:
```go
greet := action.FromFunc("greet the world", func(actor *screenplay.Actor) error {
    fmt.Printf("Hello from %s\n", actor)
    return nil
})
err := adam.AttemptsTo(greet)
```
An action that needs an ability uses it with
`screenplay.UseAnAbilityOf(actor).That(...)` (see [Abilities](#abilities)).

To make the code more fluent, one can wrap `action.FromFunc` (or the task tool)
behind builder methods:
```go
adam := screenplay.ActorNamed("Adam")
err := adam.AttemptsTo(FillIn().TheRegistrationForm().With(adamsData))
```

When the parameters of a task come from questions an actor must ask first, use
`action.FromFuncAndQuestions`. It asks each question, converts the answers to
the parameters of the function, and performs the returned task. The first
argument is the description of the action:
```go
err := actor.AttemptsTo(
    action.FromFuncAndQuestions(
        "connect to the server",
        func(actor *screenplay.Actor, url string, port int) screenplay.Performable {
            return connect.To(url).On(port)
        },
        UrlQuestion{}, PortQuestion{},
    ),
)
```
The first parameter of the function must be `*screenplay.Actor` and the remaining
parameters must match, in order and type, the answers returned by the questions. The
questions may have different answer types, so they are received as plain values and
checked when the action is performed: a value that is not a question fails with
`ErrNotAQuestion`.

### Creating your own Question
The quickest way to create a question is `question.FromFunc`. It wraps a plain
function of type `func(theActor *screenplay.Actor) (T, error)` into a `Question[T]`,
using the `description` you provide as the result of its `String` method. The answer
type is inferred from the function:
```go
theAnswer := question.FromFunc("the answer", func(theActor *screenplay.Actor) (int, error) {
    return 42, nil
})
answer, err := adam.AsksFor(theAnswer) // answer is an int
```
As with actions, you can wrap the call behind builder methods to make the code
more fluent when the answer depends on parameters.

When you write a question as a struct, declare the answer type in its `AnsweredBy`
method and let the compiler check it implements the interface:
```go
type StatusCodeQuestion struct{}

func (q *StatusCodeQuestion) AnsweredBy(actor *screenplay.Actor) (int, error) { ... }
func (q *StatusCodeQuestion) String() string                                 { return "the status code" }

var _ screenplay.Question[int] = (*StatusCodeQuestion)(nil)
```
A question that wraps another one takes it as a `screenplay.Question[T]` — either with the
precise `T` it needs (`status.CodeOf` takes a `Question[*HTTPResponse]`) or with a type
parameter of its own when it works for any answer (`question.Silently[T]`). The same goes
for an action that holds a question: make it generic over the answer type, as
`Remember`, `Log`, and `ChooseBasedOnThe` are, or let a generic function or method capture
the typed question, as `see.The`, `Conditionally(...).IfThe`, `Choose().WhenThe`, and
`Stop().UntilThe` do.

### Creating your own Resolution
The quickest way to create a resolution is `resolution.FromFunc` (see
[Creating an Anonymous Resolution](#creating-an-anonymous-resolution)). When you write a
resolution as a struct, declare the type it inspects in its `Inspect` method and let the
compiler check it implements the interface:
```go
type EndingWithResolution struct {
	suffix string
}

func (r *EndingWithResolution) Inspect(text string) error {
	if !strings.HasSuffix(text, r.suffix) {
		return fmt.Errorf("got %q", text)
	}

	return nil
}

func (r *EndingWithResolution) String() string { return fmt.Sprintf("end with %q", r.suffix) }

var _ screenplay.Resolution[string] = (*EndingWithResolution)(nil)
```
An error returned by `Inspect` means that the value does not satisfy the resolution, which
is what lets `ShouldNot` simply invert the outcome. See
[Mistakes in a test never panic](#mistakes-in-a-test-never-panic) for the one exception: a
resolution that was built wrongly.

### Reusable builders (the `timing` package)
The `timing` package holds small, dependency-free building blocks you can reuse
when writing your own actions, questions, and resolutions.

`DurationBuilder[T]` gives any builder a fluent time API (for example
`.For(100).Milliseconds()` or `.During(5).Seconds()`) without re-implementing the
amount-to-`time.Duration` conversion. It follows a simple pattern — an `amount`,
a `unit`, and a `duration` — and is generic over the parent builder type `T`. You
hand it a pointer to the `time.Duration` field it should fill in; the unit methods
write `amount * unit` into that field and return the parent so the fluent chain
keeps flowing. Pointing several builders at different fields is how a parent can
describe more than one duration (for example how long to keep trying and how
long to wait between tries).

You supply the amount fluently with the ready-made wording vocabulary — `For`,
`During`, `TryingFor`, `TryingForNoLongerThan`, `WaitingFor`, `Every`,
`PollingEvery`, `TryingEvery` — so your API reads naturally without you
re-declaring those aliases:
```go
func (a *MyAction) For(amount int) *timing.DurationBuilder[MyAction] {
    return timing.NewDurationBuilder(a, &a.timeout).For(amount)
}

// enables: MyAction{}.For(30).Seconds()
```
```go
func (a *MyAction) Polling() *timing.DurationBuilder[MyAction] {
    return timing.NewDurationBuilder(a, &a.polling)
}

// enables: MyAction{}.Polling().Every(5).Seconds()
```

Its `String` method describes the duration and matches the unit to the amount,
so it reads `1 second` but `20 seconds`. This lets an action delegate its own
description to the builder — for example `PauseFor(1).Second().Because("...")`
prints `... for 1 second because ...` while `PauseFor(20).Milliseconds()` prints
`... for 20 milliseconds ...`.

`Window` bundles the two durations that bound a retry loop: the `Total` time
during which the actor keeps trying (for how long we repeat) and the `Interval`
it waits between two tries (how much time between every trial). It exposes a
`Valid` method that reports whether the interval is not larger than the total:
```go
window := timing.NewWindow(2*time.Second, 500*time.Millisecond)
if !window.Valid() {
    // the interval between tries is larger than the total time
}
```

`WindowBuilder[T]` layers a fluent, natural-language vocabulary on top of a
`Window`. It is generic over the parent builder type `T` and wires a
`DurationBuilder` to each field of the window: the timeout methods (`For`,
`TryingFor`, `TryingForNoLongerThan`, `WaitingFor`) fill in `Total`, while the
polling methods (`Polling`, `PollingEvery`, `TryingEvery`) fill in `Interval`.
Embed a `*WindowBuilder[T]` into your builder to expose the whole vocabulary
without re-declaring every alias:
```go
type MyAction struct {
    *timing.WindowBuilder[MyAction]
    window timing.Window
}

func NewMyAction() *MyAction {
    a := &MyAction{window: timing.NewWindow(2*time.Second, 500*time.Millisecond)}
    a.WindowBuilder = timing.NewWindowBuilder(a, &a.window)
    return a
}

// enables: NewMyAction().TryingFor(5).Seconds().PollingEvery(500).Milliseconds()
```
This is exactly how `Eventually` configures its timeout and polling period: it
stores them in a `Window` and embeds a `WindowBuilder`, which is how the
same wording drives two different durations.

### Organizing Your Files
The organization of your test code is important. An easy way to organize the code is to group the different concept together.

- tests
  - features
    - feature1.go
      - Test for scenario 1
      - Test for scenario 2
    - feature2.go
  - actions
  - tasks
    - task1.go
    - task2.go
  - abilities
  - questions
  - resolutions

The features folder contains files for the different features. In each file, one can creates tests with different scenarios/test cases.

The tasks are different groups of actions. This way you have a collection of tasks that can be used to develop the tests of your features.

Then ability, actions, questions, and resolutions can also have their own folders.

## Extensions
You can find several extensions in the folder `extensions`. Each one bundles a new
ability together with the actions, questions, and resolutions that go with it, so
you can drive a specific kind of system-under-test with the same screenplay
vocabulary. Every extension has its own README with the full list of what it ships.

### HTTP — [`extensions/http`](extensions/http/README.md)
Drives REST APIs. Give an actor the `MakeHTTPRequests()` ability, then send
requests and ask about the responses:
```go
anActor := screenplay.ActorNamed("Wanda").WhoCan(ability.MakeHTTPRequests())

err := anActor.AttemptsTo(
	action.AddHeader("Accept", "application/json"),
	action.SendGetRequest().To("https://api.example.com/birds"),
)
err = anActor.Should(see.The(status.CodeOf(last.Of(question.Responses()))).Should(be.EqualTo(200)))
```
Highlights: `SendHTTPRequest`/`SendGetRequest`/`SendPostRequest`/… with `.To`,
`.WithBody`, `.WithAuth`/`.WithCredential` and `.Secretly`; the `AddHeader`,
`AddHeaders`, `SetHeader`, `SetHeaders` actions; and the `Responses` and
`status.CodeOf` questions.

### CLI — [`extensions/cli`](extensions/cli/README.md)
Drives command-line programs. Give an actor the `RunCLICommands()` ability, then
run commands, type input, and assert on the result:
```go
anActor := screenplay.ActorNamed("Boby").WhoCan(ability.RunCLICommands())

err := anActor.AttemptsTo(action.RunTheCommand("echo", "Hello World"))
err = anActor.Should(see.The(standard.OutputOf(last.Of(question.Responses()))).Should(contain.TheBytes([]byte("Hello World"))))
```
Highlights: the `RunTheCommand` action (with `.InTheWorkingDirectory` and
`.Interactively`), the `Type` action, and the `EnvironmentVariableNamed`,
`errorcode.Of`, `standard.OutputOf`, and `standard.ErrorOf` questions.

### Filesystem — [`extensions/filesystem`](extensions/filesystem/README.md)
Drives the file system. Give an actor the `UseTheFileSystem()` ability, then
create, append to, overwrite, and remove files and directories, and assert on
their existence and content:
```go
anActor := screenplay.ActorNamed("Boby").WhoCan(ability.UseTheFileSystem())

err := anActor.AttemptsTo(action.Create().TheFile("notes.txt").ContainingTheText("Hello World"))
err = anActor.Should(see.The(question.ContentOfTheFileNamed("notes.txt")).Should(be.DeepEqualTo([]byte("Hello World"))))
```
Highlights: the `Create`, `Remove`, `AppendTheText`/`AppendTheBytes`/`AppendTheContent`,
`OverwriteTo`, and `ChangeDirectoryTo` actions; the `FileNamed`, `DirectoryNamed`,
and `ContentOfTheFileNamed` questions; and the `Exist` resolution.

## References
This package takes some inspiration from several other libraries written in other languages:
- https://github.com/HamedStack/HamedStack.Playwright.Screenplay
- https://github.com/testla-project/testla-screenplay-playwright-js
- https://github.com/ScreenPyHQ/screenpy
- https://serenity-js.org/handbook/design/screenplay-pattern/

## License
Licensed under MIT License.
