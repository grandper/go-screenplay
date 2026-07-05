# go-screenplay

`go-screenplay` is a package that brings the *screenplay* pattern to Golang.

First, let's give credits, where credits is due. This package takes some inspiration from several other libraries written in other languages:
- https://github.com/HamedStack/HamedStack.Playwright.Screenplay
- https://github.com/testla-project/testla-screenplay-playwright-js
- https://github.com/ScreenPyHQ/screenpy
- https://serenity-js.org/handbook/design/screenplay-pattern/

**Main features:**

- Implementation of the *Screenplay* Pattern in Golang.
- Utilities to use *Given-When-Then* instruction in tests.
- Natural-language-enabling syntactic sugar

## Installation

```bash
go get github.com/grandper/go-screenplay
```

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
The actor represents a user of your application. It is the glue to acts on actions, questions, and Resolution.
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

Finally, abilities all implement a `Forget() error` interface that is
used to clean up the of what is inside the ability (e.g., closing
sessions). This method is automatically called when you use
```go
err := neo.Exit()
```

### Actor specific session data
Sometimes we need to store information in a step or rask and reuse it in a subsequent one.
To do so, each actor can call the `remember` and `recall` methods to store and retrieve
data in a key-value store.
```go
actor.Remember("the user id", 1234)
userID := actor.Recall("the user id").(int)
```
You can also hand `Remember` a question directly. The actor answers it and stores
the answer instead of the question. If the question fails to be answered, `nil` is
stored for the key:
```go
actor.Remember("the user id", UserID())
userID := actor.Recall("the user id").(int)
```
An actor can forget the data it has stored using
```go
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
When all actors in a scenario share the same abilities, you can use:
```go
cast := screenplay.CastWhereEveryoneCan(MakeHTTPRequests(), UseDatabase())
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
stage := screenplay.SetTheStage(screenplay.CastWhereEveryoneCan(MakeHTTPRequests()))
```
Use `TheActorCalled` to retrieve an actor by name. If the actor does not exist
yet, the stage creates one and lets the cast prepare it. Calling it a second
time with the same name returns the same actor instance (name matching is
case-insensitive). The actor is also placed in the spotlight.
```go
adam := stage.TheActorCalled("Adam")
sameAdam := stage.TheActorCalled("adam") // returns the same actor
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

A question is simply a struct that implements the following interface
```go
type Question interface {
	AnsweredBy(actor *Actor) (any, error)
	String() string
}
```
By passing the actor to the `AnsweredBy` method, we are able to use the actor `abilities` to answer the question.
The answer can be of any type.

The resolution is a struct that implements the following interface
```go
type Matcher func(obj any) (bool, error)

type Resolution interface {
	Resolve() Matcher
	String() string
}
```
As we can see here the resolution acts as a factory to instantiate a `Matcher`.
`Matcher` receives any type of objects and return `true` if that object
fulfills the expectation(s).

As mentioned before questions and resolutions works together: the actor asks a question to
retrieve an answer (i.e., a value of a given type) which is then asserted using the resolution.

Note that the resolution provides both which assertion is performed and the expected value.
For example, if we have a question `TitleOfThePage` that returns the title of an HTML page,
and a resolution that checks that the value is equal to "Some great title", then we can use
the action `see.The` to perform our test.
```go
err := theActor.Should(see.The(Text.OfThe(TitleOfThePage), is.EqualTo("Some great title")))
```

## Usage
Now that we have seen the foundation of the Screenplay Pattern and that we know how it works,
you probably realize that we need a collection of `abilities`, `actions`, `questions`, and `resolutions`.
In this section we will review what the library as already implemented for you.

#### Pausing, Stopping, and Waiting for Answers to Questions
You may want the test execution to top until the user hit the `Enter` key (mainly for debugging reasons).
In that case you can use the `Stop` action.
```go
err := theActor.Will(Stop())
```
You can also make the execution stop until some resolution is valid:
```go
err := theActor.Will(Stop().UntilThe(HomePage, contains.TheText("Hello World!")))
err := theActor.Will(Stop().UntilThe(CashAccount.Balance(), IsEqualTo(2000)))
```

You may also want to ask the execution for a given number of milliseconds, or seconds:
```go
err := theActor.AttemptsTo(PauseFor(10).Seconds().Because("the connection needs time to be setup"))
err := theActor.AttemptsTo(PauseFor(500).Milliseconds().Because("of some obscur reason"))
```
However, if you need to wait until something happens, you may want to use `Stop` instead to save time.

Finally, if you expect an action to fail repeatedly until finally it succeeds, you should use `Eventually`.
```go
err := adam.Should(Eventually(see.The(Text.OfThe(PageTitle), ContainsTheText("Hello World!"))))
err := adam.AttemptsTo(Eventually(Click.OnThe(SaveButton)).TryingEvery(100).Milliseconds())
err := adam.WasAbleTo(Eventually(CancelTheOrder).TryingFor(5).Seconds().PollingEvery(500).Milliseconds())
```
`Eventually` bounds its retry loop with a `RetryWindow`: the total time it keeps
trying and the interval it waits between two tries. Both are configured with a
fluent, natural-language vocabulary. To set how long the actor keeps trying use
`For`, `TryingFor`, `TryingForNoLongerThan`, or `WaitingFor`; to set how often it
tries use `Polling`, `PollingEvery`, or `TryingEvery`. Each of them returns a
`utils.TimeFrameBuilder` so you then pick a unit (`Milliseconds`, `Seconds`, ...):
```go
err := adam.WasAbleTo(Eventually(CancelTheOrder).WaitingFor(5).Seconds().Polling(500).Milliseconds())
```

#### Logging
If you need to log the answer to a question, you can use the `Log` action:
```go
err := theActor.AttemptsTo(Log(HowManyBirdsAreInTheSky()))
err := theActor.AttemptsTo(Log(Number.Of(ItemsInTheList)))
```

#### Remembering answers to questions
When you need to ask a question and reuse its answer later, the `Remember` action
stores the answer in the actor's memory under a given key. You can then retrieve it
with `Recall` (see [Actor specific session data](#actor-specific-session-data)):
```go
err := theActor.AttemptsTo(Remember(StatusCode()).As("statusCode"))
statusCode := theActor.Recall("statusCode").(int)
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
This mode relies on an `errgroup` and its context cancellation:
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
answer, _ := theActor.AsksFor(AsynchronousErrors())
errs := answer.([]error)
```
Like `Concurrently`, several flavors of execution are available. By default (or
with an explicit `WaitingForAll`) it waits for every action to complete before
collecting the errors. Use `CancelOnError` to cancel the remaining actions as soon
as one of them fails, or `IgnoringErrors` to discard the errors altogether:
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
`HasAbilityTo`, `UseAbilityTo`) are guarded by an internal mutex, so those calls
will not race even when made from different goroutines at the same time.

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
theActor.Will(Either(DoAction()).Or(DoDifferentAction())
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
The condition can be either a **boolean expression** (`If`, or its alias `When`) or a
**question and a resolution** (`IfThe`, or its alias `WhenThe`) — the screenplay-native
condition also used by `see.The`:
```go
theActor.AttemptsTo(
	Conditionally(AddHeader("content-type", "application/json")).
		IfThe(HostName, is.KnownBy(theActor)).
		Otherwise(GuessTheContentType()),
)
```
Both forms have a negation, `Unless` and `UnlessThe`:
```go
theActor.AttemptsTo(Conditionally(Log(TheLastResponse)).Unless(response.IsOK()))
theActor.AttemptsTo(Conditionally(Retry()).UnlessThe(StatusCode, is.EqualTo(200)))
```
The `Otherwise` branch is optional. When it is omitted and the condition does not hold,
the actor does nothing and no error is returned. A boolean condition is evaluated
eagerly (where it is written); use the question/resolution form when the condition must
be evaluated at the moment the action is performed.

#### Observing things
The simplest way to ask a question is tu use the action `see`.
```go
err := theActor.Should(see.The(Text.OfThe(PageTile), StartsWith("Hello")))
```

It is possible to check if the actor sees any or all of a list of question-resolution pair.
```go
err := theActor.Should(see.AnyOf(
			profilWidget, contains.TheText("Adam"),
			loginForm, contains.TheText("Login")))
err := theActor.Should(see.AllOf(
			profilWidget, contains.TheText("Adam"),
			loginForm, contains.TheText("Login")))
```

You can check that all, any, or none of a set of questions match a resolution:
```go
err := theActor.Should(see.ThatAllOfThe(profilWidget, loginForm)(contains.TheText("Adam")))
err := theActor.Should(see.ThatAnyOfThe(profilWidget, loginForm)(contains.TheText("Adam")))
err := theActor.Should(see.ThatNoneOfThe(profilWidget, loginForm)(contains.TheText("Adam")))
```

#### Checking texts
You can check if a text starts with, ends with, or contains a string:
```go
err := theActor.Should(see.The(Text.OfThe(PageTile), StartsWith("Hello")))
err := theActor.Should(see.The(Text.OfThe(PageTile), EndsWith("World!")))
err := theActor.should(see.The(Text.OfThe(PageTitle), contains.TheText("lo Wor")))
```

You can check if a slice of bytes contains another slice of bytes:
```go
err := theActor.Should(see.The(BytesOfThe(LastResponse), contains.TheBytes([]byte("lo Wor"))))
```

You can match a text exactly:
```go
err := theActor.Should(see.The(Text.OfThe(PageTile), ReadsExactly("Hello World!")))
```

You can use regex to match a text:
```go
err := theActor.Should(see.The(Text.OfThe(PageTitle), Matches(`^Hello \w+`)))
pattern := regexp.MustCompile(`^Hello \w+`)
err := theActor.Should(see.The(Text.OfThe(PageTitle), Matches(pattern))
```

#### Checking numbers
You can assess the value of numbers as follows:
```go
err := theActor.Should(see.The(Number.Of(items.InThe(TodoList)), is.EqualTo(0)))
err := theActor.Should(see.The(Number.Of(items.InThe(TodoList)), is.LessThan(1)))
err := theActor.Should(see.The(Number.Of(items.InThe(TodoList)), is.LessThanOrEqualTo(1)))
err := theActor.Should(see.The(Number.Of(items.InThe(TodoList)), is.GreaterThan(1)))
err := theActor.Should(see.The(Number.Of(items.InThe(TodoList)), is.GreaterThanOrEqualTo(1)))
delta := 25
err := theActor.Should(see.The(Number.Of(items.InThe(TodoList)), is.CloseTo(101, delta)))
err := theActor.Should(see.The(Number.Of(items.InThe(TodoList)), is.InRange(1, 5)))
```

#### Checking collections
There are a couple of collection types in Golan such as slices and maps.

There exists a couple of resolutions to check slices:
```go
err := theActor.Should(see.The(List.OfAll(items.InThe(TodoList)), is.Empty()))
err := theActor.Should(see.The(items.InThe(TodoList), contains.TheItem("Add tests for the Go package")))
err := theActor.Should(see.The(items.InThe(TodoList), has.Length(5)))
err := theActor.should(see.The(Text.OfAll(items.InThe(TodoList)), contains.TheItem("by the end of the year")))
```

Here are a couple of options to test maps:
```go
err := theActor.Should(see.The(HeadersOf(TheLastResponse), contains.TheValue("application/json")))
err := theActor.Should(see.The(HeadersOf(TheLastResponse), contains.TheKey("Content-Type")))
err := theActor.Should(see.The(HeadersOf(TheLastResponse), contains.TheEntry("Content-Type", "application/json")))
```

When a question answers with a slice, you may want to assert against a single element of it
rather than the whole collection. The `first` and `last` questions wrap another question and
answer with, respectively, the first and last element of the slice it returns:
```go
err := theActor.Should(see.The(first.Of(items.InThe(TodoList)), ReadsExactly("Add tests for the Go package")))
err := theActor.Should(see.The(last.Of(items.InThe(TodoList)), ReadsExactly("by the end of the year")))
```
Both fail with `ErrAnswerIsNotASlice` if the wrapped question does not answer with a slice,
and with `ErrAnswerIsEmpty` if it answers with an empty one.

When a question answers with a struct, you may want to assert against one of its fields rather
than the whole value. The `body`, `data`, and `header` questions wrap another question and answer
with, respectively, the `Body`, `Data`, and `Header` field of the struct it returns, whatever its
type:
```go
err := theActor.Should(see.The(body.Of(TheLastResponse), is.EqualTo("created")))
err := theActor.Should(see.The(data.Of(last.Of(NATSMessages)), is.EqualTo("payload")))
err := theActor.Should(see.The(header.Of(last.Of(HTTPResponses)), is.EqualTo("application/json")))
```
They fail with `ErrAnswerIsNotAStruct` if the wrapped question does not answer with a struct (or a
pointer to one), and with `ErrFieldNotFound` if the struct has no matching field.

When a question answers with a slice or a map, you may want to assert against the number of
elements it contains. The `number` question wraps another question and answers with the count
of elements in the slice or map it returns:
```go
err := theActor.Should(see.The(number.Of(items.InThe(TodoList)), is.EqualTo(2)))
```
It fails with `ErrAnswerIsNotACollection` if the wrapped question does not answer with a slice or
a map.

#### Logical Operations
We can create the negation of a resolution:
```go
err := theActor.Should(see.The(PageTitle, is.Not(Visible())))
```

#### Creating an Anonymous Resolution
It happens sometimes that we need a one-off resolution without creating a dedicated struct.
We usually do this for simple or ad-hoc assertions that are not worth extracting into a reusable type.

For example, we may want to assert that a value satisfies a custom condition inline:
```go
isPositive := resolution.FromFunc("is positive", func() screenplay.Matcher {
	return func(obj any) (bool, error) {
		n, ok := obj.(int)
		if !ok {
			return false, fmt.Errorf("expected an int, got %T", obj)
		}
		return n > 0, nil
	}
})
```
You can also wrap the call in a function to add parameters:
```go
func IsEqualTo(expected any) screenplay.Resolution {
	return resolution.FromFunc(fmt.Sprintf("is equal to %v", expected), func() screenplay.Matcher {
		return func(obj any) (bool, error) {
			return obj == expected, nil
		}
	})
}
```

Anonymous resolutions implement the `Resolution` interface just like any other resolution so that an actor
can use them the same way:
```go
err := theActor.Should(see.The(CashAccount.Balance(), IsEqualTo(2000)))
```

#### Working with contexts
Contexts are hold by actors to share information across different part of a sequence of calls.
```go
ctx := context.Background()
actor := screenplay.ActorNamed("Alice").WithContext(ctx)
```
The context can later be retrieved using
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
screenplay.Given(adam).WasAbleTo(see.The(AccountBalance, is.Equal(2)))
screenplay.When(adam).AttemptsTo(Deposit(100).Dollars())
screenplay.Then(adam).Should(see.The(AccountBalance), is.EqualTo(102))
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
parameters must match, in order and type, the answers returned by the questions.

### Creating your own Question
The quickest way to create a question is `question.FromFunc`. It wraps a plain
function of type `func(theActor *screenplay.Actor) (any, error)` into a `Question`,
using the `description` you provide as the result of its `String` method:
```go
theAnswer := question.FromFunc("the answer", func(theActor *screenplay.Actor) (any, error) {
    return 42, nil
})
answer, err := theAnswer.AnsweredBy(adam)
```
As with actions, you can wrap the call behind builder methods to make the code
more fluent when the answer depends on parameters.

### Reusable builders (the `utils` package)
The `utils` package holds small, dependency-free building blocks you can reuse
when writing your own actions, questions, and resolutions.

`TimeFrameBuilder[T]` gives any builder a fluent time API (for example
`.For(100).Milliseconds()` or `.During(5).Seconds()`) without re-implementing the
amount-to-`time.Duration` conversion. It follows a simple pattern — an `amount`,
a `unit`, and a `duration` — and is generic over the parent builder type `T`. You
hand it a pointer to the `time.Duration` field it should fill in; the unit methods
write `amount * unit` into that field and return the parent so the fluent chain
keeps flowing. Pointing several builders at different fields is how a parent can
describe more than one time frame (for example how long to keep trying and how
long to wait between tries).

You supply the amount fluently with the ready-made wording vocabulary — `For`,
`During`, `TryingFor`, `TryingForNoLongerThan`, `WaitingFor`, `Every`,
`PollingEvery`, `TryingEvery` — so your API reads naturally without you
re-declaring those aliases:
```go
func (a *MyAction) For(amount int) *utils.TimeFrameBuilder[MyAction] {
    return utils.NewTimeFrameBuilder(a, &a.timeout).For(amount)
}

// enables: MyAction{}.For(30).Seconds()
```
```go
func (a *MyAction) Polling() *utils.TimeFrameBuilder[MyAction] {
    return utils.NewTimeFrameBuilder(a, &a.polling)
}

// enables: MyAction{}.Polling().Every(5).Seconds()
```

Its `String` method describes the time frame and matches the unit to the amount,
so it reads `1 second` but `20 seconds`. This lets an action delegate its own
description to the builder — for example `PauseFor(1).Second().Because("...")`
prints `... for 1 second because ...` while `PauseFor(20).Milliseconds()` prints
`... for 20 milliseconds ...`.

`RetryWindow` bundles the two durations that bound a retry loop: the `Total` time
during which the actor keeps trying (for how long we repeat) and the `Interval`
it waits between two tries (how much time between every trial). It exposes a
`Valid` method that reports whether the interval is not larger than the total:
```go
window := utils.NewRetryWindow(2*time.Second, 500*time.Millisecond)
if !window.Valid() {
    // the interval between tries is larger than the total time
}
```

`RetryWindowBuilder[T]` layers a fluent, natural-language vocabulary on top of a
`RetryWindow`. It is generic over the parent builder type `T` and wires a
`TimeFrameBuilder` to each field of the window: the timeout methods (`For`,
`TryingFor`, `TryingForNoLongerThan`, `WaitingFor`) fill in `Total`, while the
polling methods (`Polling`, `PollingEvery`, `TryingEvery`) fill in `Interval`.
Embed a `*RetryWindowBuilder[T]` into your builder to expose the whole vocabulary
without re-declaring every alias:
```go
type MyAction struct {
    *utils.RetryWindowBuilder[MyAction]
    window utils.RetryWindow
}

func NewMyAction() *MyAction {
    a := &MyAction{window: utils.NewRetryWindow(2*time.Second, 500*time.Millisecond)}
    a.RetryWindowBuilder = utils.NewRetryWindowBuilder(a, &a.window)
    return a
}

// enables: NewMyAction().TryingFor(5).Seconds().PollingEvery(500).Milliseconds()
```
This is exactly how `Eventually` configures its timeout and polling period: it
stores them in a `RetryWindow` and embeds a `RetryWindowBuilder`, which is how the
same wording drives two different time frames.

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
You can find several extensions in the folder `extensions`. Each of them is used to extend the capability of the library to a specific use case.
- `http`: support API testing using REST requests.
- `cli`: support testing CLI applications.
- `filesystem`: support testing file system interactions.

## License
Licensed under MIT License.
